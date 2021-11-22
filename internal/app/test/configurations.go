package test

import (
	"github.com/khaiql/dbcleaner/engine"
	"github.com/qor/validations"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/migrations"
	dbcleaner "gopkg.in/khaiql/dbcleaner.v2"
)

var Cleaner = dbcleaner.New()

func init() {
	validations.RegisterCallbacks(database.DB())
	mysql := engine.NewMySQLEngine(test_utils.ConfToString())
	Cleaner.SetEngine(mysql)
	migrations.Migrate()
}
