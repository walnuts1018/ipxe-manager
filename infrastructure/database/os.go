package database

import (
	"context"
	"fmt"

	"go.etcd.io/bbolt"
)

const (
	bucketNameOS = "os"
	currentKey   = "current"
)

func (db *DB) GetOS(ctx context.Context) (string, error) {
	var osstr string
	if err := db.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketNameOS))
		if bucket == nil {
			return nil
		}

		value := bucket.Get([]byte(currentKey))
		if value == nil {
			return fmt.Errorf("no OS set")
		}

		osstr = string(value)
		return nil
	}); err != nil {
		return "", fmt.Errorf("failed to get OS from bbolt: %w", err)
	}

	return osstr, nil
}

func (db *DB) SetOS(ctx context.Context, osstr string) error {
	if err := db.db.Batch(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketNameOS))
		if err != nil {
			return err
		}

		if err := bucket.Put([]byte(currentKey), []byte(osstr)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to set OS in bbolt: %w", err)
	}

	return nil
}
