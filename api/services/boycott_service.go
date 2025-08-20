package services

import (
	"canada-hires/models"
	"canada-hires/repos"
	"fmt"
	"strings"
)

type ToggleBoycottRequest struct {
	UserID          string `json:"user_id"`
	BusinessName    string `json:"business_name"`
	BusinessAddress string `json:"business_address"`
}

type BoycottService interface {
	ToggleBoycott(req *ToggleBoycottRequest) (*models.Boycott, error)
	GetUserBoycotts(userID string, limit, offset int) ([]*models.Boycott, error)
	GetTopBoycottedBusinesses(limit int) ([]*models.BoycottStats, error)
	GetBoycottCount(businessName, businessAddress string) (int, error)
	IsBoycottedByUser(userID, businessName, businessAddress string) (bool, error)
}

type boycottService struct {
	repo repos.BoycottRepository
}

func NewBoycottService(repo repos.BoycottRepository) BoycottService {
	return &boycottService{repo: repo}
}

func (s *boycottService) ToggleBoycott(req *ToggleBoycottRequest) (*models.Boycott, error) {
	if err := s.validateToggleRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	businessName := strings.TrimSpace(req.BusinessName)
	businessAddress := strings.TrimSpace(req.BusinessAddress)

	boycott, err := s.repo.Toggle(req.UserID, businessName, businessAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to toggle boycott: %w", err)
	}

	return boycott, nil
}

func (s *boycottService) validateToggleRequest(req *ToggleBoycottRequest) error {
	if req.UserID == "" {
		return fmt.Errorf("user ID is required")
	}
	if strings.TrimSpace(req.BusinessName) == "" {
		return fmt.Errorf("business name is required")
	}
	return nil
}

func (s *boycottService) GetUserBoycotts(userID string, limit, offset int) ([]*models.Boycott, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	boycotts, err := s.repo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user boycotts: %w", err)
	}

	return boycotts, nil
}

func (s *boycottService) GetTopBoycottedBusinesses(limit int) ([]*models.BoycottStats, error) {
	if limit <= 0 {
		limit = 3 // Default to top 3
	}
	if limit > 10 {
		limit = 10 // Max limit
	}

	stats, err := s.repo.GetTopBoycotted(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top boycotted businesses: %w", err)
	}

	return stats, nil
}

func (s *boycottService) GetBoycottCount(businessName, businessAddress string) (int, error) {
	if strings.TrimSpace(businessName) == "" {
		return 0, fmt.Errorf("business name is required")
	}

	count, err := s.repo.GetBoycottCount(strings.TrimSpace(businessName), strings.TrimSpace(businessAddress))
	if err != nil {
		return 0, fmt.Errorf("failed to get boycott count: %w", err)
	}

	return count, nil
}

func (s *boycottService) IsBoycottedByUser(userID, businessName, businessAddress string) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("user ID is required")
	}
	if strings.TrimSpace(businessName) == "" {
		return false, fmt.Errorf("business name is required")
	}

	isBoycotted, err := s.repo.IsBoycottedByUser(userID, strings.TrimSpace(businessName), strings.TrimSpace(businessAddress))
	if err != nil {
		return false, fmt.Errorf("failed to check boycott status: %w", err)
	}

	return isBoycotted, nil
}