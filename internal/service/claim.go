package service

import "go.uber.org/dig"

type ClaimService struct {
}

type InputClaimService struct {
	dig.In
}

func NewClaimService(param InputClaimService) *ClaimService {
	return &ClaimService{

	}
}
