package database

import (
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	log "github.com/arthur-trt/bechdelproxy/log"
)

var Conn *gorm.DB

func init() {
	var err error
	defaultPort := "5432"
	defaultTZ := "Europe/Paris"
	defaultBatchSize := 1000

	pgsqlConnection := map[string]*string{
		"PGSQL_HOST": nil,
		"PGSQL_USER": nil,
		"PGSQL_PASS": nil,
		"PGSQL_DB":   nil,
		"PGSQL_PORT": &defaultPort,
		"PGSQL_TZ":   &defaultTZ,
	}

	for key := range pgsqlConnection {
		if value, found := os.LookupEnv(key); found {
			pgsqlConnection[key] = &value
		} else if pgsqlConnection[key] != nil {
			log.Warn(key + " is not defined, using default value: " + *pgsqlConnection[key])
		} else {
			log.Fatal(key + " is not defined")
		}
	}

	dsn := "host=" + *pgsqlConnection["PGSQL_HOST"] +
		" user=" + *pgsqlConnection["PGSQL_USER"] +
		" password=" + *pgsqlConnection["PGSQL_PASS"] +
		" dbname=" + *pgsqlConnection["PGSQL_DB"] +
		" port=" + *pgsqlConnection["PGSQL_PORT"] +
		" TimeZone=" + *pgsqlConnection["PGSQL_TZ"]

	log.Debug("Connecting to database with DSN: " + dsn)

	Conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		CreateBatchSize: defaultBatchSize,
	})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Info("Connection to database succeeded")

	if batchSizeEnv, found := os.LookupEnv("PGSQL_BATCH_SIZE"); found {
		if batchSize, err := strconv.Atoi(batchSizeEnv); err == nil && batchSize > 0 {
			Conn = Conn.Session(&gorm.Session{CreateBatchSize: batchSize})
			log.Info("Batch size set to: ", batchSize)
		} else {
			log.Warn("Invalid PGSQL_BATCH_SIZE, using default: ", defaultBatchSize)
		}
	}

	if err := migrateDatabase(Conn); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	log.Info("Database migration completed successfully")
}

func migrateDatabase(db *gorm.DB) error {
	log.Info("Starting database migration...")

	models := []interface{}{
		&Movie{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}

	return nil
}

func InsertOrUpdateMovies(movies []Movie) error {
	if len(movies) == 0 {
		log.Info("No movies to insert or update")
		return nil
	}

	log.Info("Inserting or updating movies in batch...")

	if err := Conn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "imdb_id"}},
		UpdateAll: true,
	}).Create(&movies).Error; err != nil {
		log.Error("Failed to insert or update movies: ", err)
		return err
	}

	log.Info("Added " + strconv.FormatInt(int64(len(movies)), 10) + " movies")
	log.Info("Batch insert or update completed successfully")
	return nil
}
