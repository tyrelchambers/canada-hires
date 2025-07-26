package controllers

import (
	"canada-hires/services"
	"net/http"
)

type BusinessController interface {
	GetBusinesses(w http.ResponseWriter, r *http.Request)
	CreateBusiness(w http.ResponseWriter, r *http.Request)
	GetBusiness(w http.ResponseWriter, r *http.Request)
	UpdateBusiness(w http.ResponseWriter, r *http.Request)
}

type businessController struct{
	service services.BusinessService
}

func NewBusinessController(service services.BusinessService) BusinessController {
	return &businessController{service: service}
}

func (c *businessController) GetBusinesses(w http.ResponseWriter, r *http.Request) {}
func (c *businessController) CreateBusiness(w http.ResponseWriter, r *http.Request) {}
func (c *businessController) GetBusiness(w http.ResponseWriter, r *http.Request)    {}
func (c *businessController) UpdateBusiness(w http.ResponseWriter, r *http.Request) {}
