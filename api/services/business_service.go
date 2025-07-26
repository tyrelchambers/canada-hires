package services

import "canada-hires/repos"

type BusinessService interface{}

type businessService struct{
	repo repos.BusinessRepository
}

func NewBusinessService(repo repos.BusinessRepository) BusinessService {
	return &businessService{repo: repo}
}
