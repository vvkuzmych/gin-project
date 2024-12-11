package repository

import (
	"database/sql"
	"gorm.io/gorm"
)

type Storage interface {
	WithTx(f func(tx *gorm.DB) error) error
	Begin(opts ...*sql.TxOptions) *gorm.DB
	Commit(db *gorm.DB) *gorm.DB
	Rollback(db *gorm.DB) *gorm.DB
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
	IsError(result *gorm.DB) bool
}

type storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) Storage {
	return &storage{db: db}
}

func (s *storage) WithTx(f func(tx *gorm.DB) error) error {
	tx := s.db.Begin()

	err := f(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (s *storage) Begin(opts ...*sql.TxOptions) *gorm.DB {
	return s.db.Begin(opts...)
}

func (s *storage) Commit(db *gorm.DB) *gorm.DB {
	return db.Commit()
}

func (s *storage) Rollback(db *gorm.DB) *gorm.DB {
	return db.Rollback()
}

func (s *storage) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	return s.db.Transaction(fc, opts...)
}

func (s *storage) IsError(result *gorm.DB) bool {
	return result.Error != nil
}
