package postgres

import "database/sql"

type DBEngine interface {
	GetDB() *sql.DB
	GetDBRead() *sql.DB
	Configure(...Option) DBEngine
	Close()
}
