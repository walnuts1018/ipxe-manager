package database

import (
	"os"
	"path/filepath"

	"github.com/walnuts1018/ipxe-manager/config"
	"go.etcd.io/bbolt"
)

type DB struct {
	db *bbolt.DB
}

func NewBoltDB(cfg config.DBConfig) (*DB, error) {
	dir, err := filepath.Abs(filepath.Dir(cfg.DBPath))
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := bbolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}
