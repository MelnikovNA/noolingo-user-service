package mysql

import (
	"database/sql"
	"strings"

	"github.com/noolingo/user-service/internal/domain"

	_ "github.com/go-sql-driver/mysql"
)

func New(cfg *domain.Mysql) (*sql.DB, error) {

	dsn := strings.TrimPrefix(cfg.DSN, "mysql://")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return db, err
}
