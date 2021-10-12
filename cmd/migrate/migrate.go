package migrate

import (
	"errors"

	"addysnip.dev/api/pkg/database"
	"addysnip.dev/api/pkg/utils"
	"addysnip.dev/emailer/pkg/logger"
	models "addysnip.dev/types/database"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:   "migrate",
		Usage:  "Run database migration",
		Action: Run,
	}
}

var log = logger.Category("cmd/migrate")

func Run(c *cli.Context) error {
	if database.DB == nil {
		if utils.Getenv("DATABASE_DSN", "") == "" {
			log.Error("DATABASE_DSN is not set and is required")
			return errors.New("database_dsn is not set and is required")
		}

		log.Info("Connecting to database")
		err := database.Connect(utils.Getenv("DATABASE_DSN", ""), database.DBOptions{})
		if err != nil {
			log.Error("Error connecting to database: %s", err.Error())
			return err
		}
	}

	log.Info("Migrating database tables")
	err := database.DB.AutoMigrate(
		&models.Template{},
	)
	if err != nil {
		log.Error("Error migrating database: %s", err.Error())
		return err
	}

	log.Info("Migration complete")

	return nil
}
