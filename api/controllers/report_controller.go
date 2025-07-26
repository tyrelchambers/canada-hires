package controllers

import (
	"canada-hires/services"
	"net/http"
)

type ReportController interface {
	CreateReport(w http.ResponseWriter, r *http.Request)
	GetReportsByBusiness(w http.ResponseWriter, r *http.Request)
	UpdateReport(w http.ResponseWriter, r *http.Request)
	DeleteReport(w http.ResponseWriter, r *http.Request)
}

type reportController struct{
	service services.ReportService
}

func NewReportController(service services.ReportService) ReportController {
	return &reportController{service: service}
}

func (c *reportController) CreateReport(w http.ResponseWriter, r *http.Request)         {}
func (c *reportController) GetReportsByBusiness(w http.ResponseWriter, r *http.Request) {}
func (c *reportController) UpdateReport(w http.ResponseWriter, r *http.Request)         {}
func (c *reportController) DeleteReport(w http.ResponseWriter, r *http.Request)         {}
