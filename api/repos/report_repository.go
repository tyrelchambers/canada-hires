package repos

import "github.com/jmoiron/sqlx"

type ReportRepository interface{}

type reportRepository struct{
	db *sqlx.DB
}

func NewReportRepository(db *sqlx.DB) ReportRepository {
	return &reportRepository{db: db}
}
