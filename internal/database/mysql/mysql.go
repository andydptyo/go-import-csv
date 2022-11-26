package mysql

import (
	"database/sql"
	"embed"
	"time"

	"github.com/Boostport/migration"
	"github.com/Boostport/migration/driver/mysql"
	"github.com/andydptyo/go-import-csv/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

//go:embed migrations
var embedFS embed.FS

type Mysql struct {
	DB *sql.DB
}

func MigrationSource() *migration.EmbedMigrationSource {
	return &migration.EmbedMigrationSource{
		EmbedFS: embedFS,
		Dir:     "migrations",
	}
}

func RunMigration(dsn string, direction int) (int, error) {
	var maxApplied int
	driver, err := mysql.New(dsn)

	if err != nil {
		return 0, err
	}

	if direction == 1 {
		maxApplied = 1
	}

	applied, err := migration.Migrate(driver, MigrationSource(), migration.Direction(direction), maxApplied)
	if err != nil {
		return applied, err
	}

	return applied, nil
}

func New(c *config.Mysql) (*Mysql, error) {
	db, err := sql.Open("mysql", c.GetDsn())
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * time.Duration(c.MaxLifeTime))
	db.SetMaxOpenConns(c.MaxOpenConnection)
	db.SetMaxIdleConns(c.MaxIdleConnection)

	conn := &Mysql{
		DB: db,
	}

	return conn, nil
}

func (c *Mysql) Close() {
	c.DB.Close()
}
