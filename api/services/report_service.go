package services

import "canada-hires/repos"

type ReportService interface{}

type reportService struct{
	repo repos.ReportRepository
}

func NewReportService(repo repos.ReportRepository) ReportService {
	return &reportService{repo: repo}
}
