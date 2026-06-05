package repository

import (
	"oral_practice/internal/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}

func (r *Repository) Create(value interface{}) *gorm.DB {
	return r.db.Create(value)
}

func (r *Repository) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return r.db.First(dest, conds...)
}

func (r *Repository) Where(query interface{}, args ...interface{}) *gorm.DB {
	return r.db.Where(query, args...)
}

func (r *Repository) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return r.db.Find(dest, conds...)
}

func (r *Repository) Save(value interface{}) *gorm.DB {
	return r.db.Save(value)
}

func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(&model.User{}, &model.Session{}, &model.Message{})
}
