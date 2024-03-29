package dbmigrat

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

//go:embed testdata
var exampleFixture embed.FS

func Example() {
	// resetDB only for testing purposes - you may ignore it
	err := th.resetDB()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := sqlx.Open("postgres", os.Getenv("DBMIGRAT_TEST_DB_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	pgStore := &PostgresStore{DB: db}
	err = pgStore.CreateLogTable()
	if err != nil {
		log.Fatalln(err)
	}

	authMigrations, err := ReadDir(exampleFixture, "testdata/auth")
	if err != nil {
		log.Fatalln(err)
	}
	billingMigrations, err := ReadDir(exampleFixture, "testdata/billing")
	if err != nil {
		log.Fatalln(err)
	}
	migrations := Migrations{
		"auth":    authMigrations,
		"billing": billingMigrations,
	}

	checkRes, err := CheckLogTableIntegrity(pgStore, migrations)
	if err != nil {
		log.Fatalln(err)
	}
	if checkRes.IsCorrupted {
		log.Fatalf("Db migrations are corrupted: %+v\n", checkRes)
	}

	logsCount, err := Migrate(pgStore, migrations, RepoOrder{"auth", "billing"})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("[dbmigrat] applied %d migrations\n", logsCount)

	// Rollback migrations
	logsCount, err = Rollback(pgStore, migrations, RepoOrder{"billing", "auth"}, -1)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("[dbmigrat] rolled back %d migrations\n", logsCount)

	// Output:
	// [dbmigrat] applied 3 migrations
	// [dbmigrat] rolled back 3 migrations
}
