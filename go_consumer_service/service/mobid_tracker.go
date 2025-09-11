package service

import (
	"context"
	"go_consumer_service/repositories"
	"go_consumer_service/types"
)

type MobidTrackerService struct {
	repo repositories.MobidTrackerRepo
}

func NewMobidTrackerService(repo repositories.MobidTrackerRepo) *MobidTrackerService {
	return &MobidTrackerService{repo: repo}
}

// get all mobid tracker
func (s *MobidTrackerService) GetAllMobidTracker(ctx context.Context, limit, offset int) (*[]types.MobidTracker, int, error) {
	return s.repo.GetAll(ctx, limit, offset)
}
