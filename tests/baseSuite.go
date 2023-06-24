package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"

	"github.com/nsf/jsondiff"

	"github.com/roguexray007/loan-app/internal/app/loans"
	"github.com/roguexray007/loan-app/internal/app/loans/payments"
	"github.com/roguexray007/loan-app/internal/app/users"
	"github.com/roguexray007/loan-app/internal/boot"
	"github.com/roguexray007/loan-app/internal/provider"
	"github.com/roguexray007/loan-app/internal/routing/tenant"
	"github.com/roguexray007/loan-app/pkg/container"
	"github.com/roguexray007/loan-app/pkg/db"
	"github.com/roguexray007/loan-app/tests/fixtures"
)

type Single struct {
	ctx   context.Context
	Suite string
	Name  string
}

func (single *Single) SetTester(tnt *tenant.Tenant) *Single {
	single.ctx, _ = tenant.Attach(single.ctx, tnt)
	return single
}

type BaseSuite struct {
	suite.Suite
	Current       *Single
	DBConnections *db.Connections
	Server        *gin.Engine
	Recorder      *httptest.ResponseRecorder
}

func (s *BaseSuite) SetupSuite() {
	ctx := context.Background()

	container.Init(ctx, "dev", provider.GetManager())
	s.Server = (&boot.FunctionalTest{}).Init(ctx)

	s.Server.Use(func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(),
			db.ContextKeyDatabaseConnection, db.Master))
		ctx.Next()
	})

	boot.RegisterApplicationHandler(s.Server)
}

func (s *BaseSuite) BeforeTest(suiteName, testName string) {
	// Create a new test object
	s.Current = &Single{
		ctx:   context.Background(),
		Suite: suiteName,
		Name:  testName,
	}
	s.DBConnections = provider.GetDatabase(nil)

	fixtures.Init()
}

func (s *BaseSuite) AfterTest(suiteName, testName string) {
	// Verify if the after test is called for current test only
	if s.Current.Suite != suiteName || s.Current.Name != testName {
		s.T().Fatalf("test: Before test and after test are not called sequencially")
	}
}

func (s *BaseSuite) TearDownTest() {

	gormDb := s.DBConnections.GetConnection(s.Current.ctx)
	s.Current = nil
	for _, t := range []string{users.TableUser, loans.TableLoan, payments.TableLoanPayment} {
		gormDb.Exec(fmt.Sprintf("delete from %s", t))

	}

}

func (s *BaseSuite) User(tester *tenant.Tenant) {
	s.Current.SetTester(tester)
}

func AssertJSONSelectiveMatch(t *testing.T, actual []byte, expected interface{}) {
	options := jsondiff.DefaultJSONOptions()

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Errorf("failed to marshal expected response: %s", err)
	}

	difference, differenceHumanReadable := jsondiff.Compare(actual, expectedJSON, &options)
	if difference != jsondiff.SupersetMatch && difference != jsondiff.FullMatch {
		t.Errorf("response do not match:\ndifference: %s\nactual: %s\nexpected: %s\n",
			differenceHumanReadable, actual, expectedJSON)
	}
}
