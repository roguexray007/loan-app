package base

import (
	"context"
	goErr "errors"
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/roguexray007/loan-app/internal/routing/tenant"
	"github.com/roguexray007/loan-app/pkg/db"
)

type Repo struct {
	*db.Connections
}

// NewRepo creates new repo instance for the model
func NewRepo(dbConnections *db.Connections) Repo {
	return Repo{
		dbConnections,
	}
}

// FindByID fetches the record which matches the ID provided from the entity defined by receiver
// and the result will be loaded into receiver
func (repo Repo) FindByID(ctx context.Context, receiver IModel, id int64) error {
	q := repo.GetConnection(ctx).Where("id = ?", id)
	tnt := tenant.From(ctx)
	if tnt.IsUser() {
		q.Where("user_id = ?", tnt.User().GetID())
	}

	q = q.First(receiver)

	return GetDBError(q)
}

// FindByIDs fetches the all the records which matches the IDs provided from the entity defined by receivers
// and the result will be loaded into receivers
func (repo Repo) FindByIDs(ctx context.Context, receivers interface{}, ids []int64) error {
	q := repo.GetConnection(ctx).Where(AttributeID+" in (?)", ids).Find(receivers)

	return GetDBError(q)
}

// Create inserts a new record in the entity defined by the receiver
// all data filled in the receiver will inserted
func (repo Repo) Create(ctx context.Context, receiver IModel) error {
	q := repo.GetConnection(ctx, db.Master, false).Create(receiver)
	return GetDBError(q)
}

// CreateOrIgnoreIfDuplicate inserts a new record in the entity defined by the receiver
// all data filled in the receiver will inserted
func (repo Repo) CreateOrIgnoreIfDuplicate(ctx context.Context, receiver IModel) error {
	q := repo.GetConnection(ctx, db.Master, false).
		Set("gorm:insert_modifier", "IGNORE").
		Create(receiver)

	return GetDBError(q)
}

// Update will update the given model with respect to primary key / id available in it.
// if the selective list if passed then only defined attributes will be updated along with updated_at time
func (repo Repo) Update(ctx context.Context, receiver IModel, selectiveList ...string) error {
	q := repo.GetConnection(ctx, db.Master, false).Model(receiver)

	if len(selectiveList) > 0 {
		q = q.Select(selectiveList)
	}

	q = q.Update(receiver)

	if q.RowsAffected == 0 {
		return goErr.New(errNoRowAffected)
	}

	return GetDBError(q)
}

// Delete deletes the given model
// Soft or hard delete of model depends on the models implementation
// if the model composites SoftDeletableModel then it'll be soft deleted
func (repo Repo) Delete(ctx context.Context, receiver IModel) error {
	q := repo.GetConnection(ctx, db.Master, false).Delete(receiver)

	return GetDBError(q)
}

// FindMany will fetch multiple records from the entity defined by receiver which matched the condition provided
// note: this won't work for in clause. can be used only for `=` conditions
func (repo Repo) FindMany(
	ctx context.Context,
	receivers interface{},
	condition interface{}) error {

	q := repo.GetConnection(ctx).Where(condition).Find(receivers)

	return GetDBError(q)
}

// Reload refresh the model using its id
func (repo Repo) Reload(ctx context.Context, model IModel) error {
	return repo.FindByID(ctx, model, model.GetID())
}

// Transaction will manage the execution inside a transactions
// adds the txn db in the context for downstream use case
func (repo Repo) Transaction(ctx context.Context, fc func(ctx context.Context) error) (err error) {
	// If there is an active transaction then do not create a new transactions
	// use the same and continue
	// *Note: make sure there error occurred in nested transaction should be
	// propagated to outer transaction in this approach
	// as currently we dont have support for savepoint we wont be rollback to particular point
	if _, ok := ctx.Value(db.ContextKeyDatabase).(*gorm.DB); ok {
		return fc(ctx)
	}

	panicked := true
	// start transactions
	tx := repo.GetConnection(ctx, db.Master, false).Begin()

	// post operation check if the operation was successfully completed
	// if failed to complete, then rollback the transaction
	defer func() {
		// Make sure to rollback when panic, Block error or Commit error
		if panicked || err != nil {
			tx.Rollback()
		}
	}()

	// call the transaction handled with tx added in the context key
	err = fc(context.WithValue(ctx, db.ContextKeyDatabase, tx))
	if err == nil {
		err = GetDBError(tx.Commit())
	}

	// if there we no panics then this will be set to false
	// this will be further checked in defer
	panicked = false
	return
}

// IsTransactionActive returns true if a transaction is active
func (repo Repo) IsTransactionActive(ctx context.Context) bool {
	_, ok := ctx.Value(db.ContextKeyDatabase).(*gorm.DB)
	return ok
}

/*
findByIDForUpdate fetches the record which matches the ID provided from the
entity defined by receiver and the result will be loaded into receiver. This should
be called only within a transaction and lock will be taken on the corresponding
row until transaction either commits or rolled back
*/
func (repo Repo) findByIDForUpdate(ctx context.Context, receiver IModel, id int64) error {
	q := repo.GetConnection(ctx).
		Set("gorm:query_option", "FOR UPDATE").
		Where("id = ?", id).First(receiver)

	return GetDBError(q)
}

func (r Repo) LoadAssociation(ctx context.Context, receiver interface{}, association interface{}) error {
	q := r.GetConnection(ctx).Model(receiver).Related(association)

	err := GetDBError(q)

	if err != nil {
		var receiverEntityName, associationEntityName, receiverID string

		if entity, ok := receiver.(interface{ EntityName() string }); ok {
			receiverEntityName = entity.EntityName()
		}

		if entity, ok := receiver.(interface{ GetID() string }); ok {
			receiverID = entity.GetID()
		}

		if entity, ok := association.(interface{ EntityName() string }); ok {
			associationEntityName = entity.EntityName()
		}

		fmt.Println(map[string]interface{}{
			"error":                   err.Error(),
			"receiver_entity_name":    receiverEntityName,
			"receiver_id":             receiverID,
			"association_entity_name": associationEntityName,
		})
	}

	return err
}
