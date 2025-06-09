package services

import (
	"errors"

	"washo.com/main/repositories"
)

type ProcessVehicles struct {
	Process   string                 // The status name (e.g., "Python", "Go").
	Vehicles []repositories.Vehicle // The list of vehicles for this status.
}

type VehicleService struct {
	repo *repositories.VehicleRepository
}

func NewVehicleService(repo *repositories.VehicleRepository) *VehicleService {
	return &VehicleService{repo: repo}
}

func (s *VehicleService) CreateVehicle(input *repositories.CreateVehicleRequest) (string, error) {
	if s.repo == nil {
		return "", errors.New("repository is nil")
	}
	id, err := s.repo.Create(input)
	return id, err
}

func (s *VehicleService) GetVehicleByID(id string) (*repositories.Vehicle, error) {
	if s.repo == nil {
		return nil, errors.New("repository is nil")
	}
	return s.repo.FindByID(id)
}

func (s *VehicleService) UpdateVehicle(id string, input repositories.CreateVehicleRequest) error {
	if s.repo == nil {
		return errors.New("repository is nil")
	}
	return s.repo.Update(id, &input)
}

func (s *VehicleService) GetVehiclesByProcess(processes []string) ([]ProcessVehicles, error) {
	groupedVehicles := []ProcessVehicles{}

	for _, process := range processes {
		vehicles, err := s.repo.FindByProcess(process)
		if err != nil {
			return nil, err
		}
		groupedVehicles = append(groupedVehicles, ProcessVehicles{
			Process:   process,
			Vehicles: vehicles,
		})
	}

	return groupedVehicles, nil
}

func (s *VehicleService) DeleteVehicle(id string) error {
	if s.repo == nil {
		return errors.New("repository is nil")
	}
	return s.repo.Delete(id)
}
