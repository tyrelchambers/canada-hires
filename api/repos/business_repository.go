package repos

import "github.com/jmoiron/sqlx"

type BusinessRepository interface{}

type businessRepository struct{
	db *sqlx.DB
}

func NewBusinessRepository(db *sqlx.DB) BusinessRepository {
	return &businessRepository{db: db}
}
