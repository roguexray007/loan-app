package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/pressly/goose"

	"github.com/roguexray007/loan-app/internal/boot"
	_ "github.com/roguexray007/loan-app/internal/database/migrations"
	"github.com/roguexray007/loan-app/internal/provider"
)

var (
	flags   = flag.NewFlagSet("goose", flag.ExitOnError)
	dir     = flags.String("dir", "internal/database/migrations", "Directory with migration files")
	verbose = flags.Bool("v", false, "Enable verbose mode")
)

func main() {
	// Initialize context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// initialize app with all dependencies required
	(&boot.Migrations{}).Init(ctx)

	run(ctx)
}

func run(ctx context.Context) {

	flags.Usage = usage
	flags.Parse(os.Args[1:])
	args := flags.Args()
	if *verbose {
		goose.SetVerbose(true)
	}

	// I.e. no command provided, hence print usage and return.
	if len(args) < 1 {
		flags.Usage()
		return
	}

	// Prepares command and arguments for goose's run.
	command := args[0]
	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	// If command is create or fix, no need to connect to db and hence the
	// specific case handling.
	switch command {
	case "create":
		if err := goose.Run("create", nil, *dir, arguments...); err != nil {
			fmt.Println("failed to run command: " + err.Error())
		}
		return
	case "fix":
		if err := goose.Run("fix", nil, *dir); err != nil {
			fmt.Println("failed to run command: " + err.Error())
		}
		return
	}

	dialect := provider.GetConfig(ctx).Db.Master.Dialect
	if err := goose.SetDialect(dialect); err != nil {
		fmt.Println("failed to run command: " + err.Error())
	}

	provider.GetDatabase(nil).GetConnection(ctx, "master").LogMode(true)
	db := provider.GetDatabase(nil).GetConnection(ctx, "master").DB()

	// Finally, executes the goose's command.
	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		fmt.Println("failed to run command: " + err.Error())
	}

}

func usage() {
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usageCommands = `
Commands:
	up                   Migrate the DB to the most recent version available
	up-to VERSION        Migrate the DB to a specific VERSION
	down                 Roll back the version by 1
	down-to VERSION      Roll back to a specific VERSION
	redo                 Re-run the latest migration
	reset                Roll back all migrations
	status               Dump the migration status for the current DB
	version              Print the current version of the database
	create NAME          Creates new migration file with the current timestamp
	fix                  Apply sequential ordering to migrations
`
)
