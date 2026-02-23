package scheduler

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/service"
)

type Scheduler struct {
	cron               *cron.Cron
	buildingService    *service.BuildingService
	researchService    *service.ResearchService
	fleetService       *service.FleetService
	productionService  *service.ProductionService
	planetRepo        repository.PlanetRepository
	buildingQueueRepo repository.BuildingQueueRepository
	researchQueueRepo repository.ResearchQueueRepository
}

func NewScheduler(
	buildingService *service.BuildingService,
	researchService *service.ResearchService,
	fleetService *service.FleetService,
	productionService *service.ProductionService,
	planetRepo repository.PlanetRepository,
	buildingQueueRepo repository.BuildingQueueRepository,
	researchQueueRepo repository.ResearchQueueRepository,
) *Scheduler {
	return &Scheduler{
		cron:                cron.New(),
		buildingService:     buildingService,
		researchService:     researchService,
		fleetService:        fleetService,
		productionService:   productionService,
		planetRepo:         planetRepo,
		buildingQueueRepo:   buildingQueueRepo,
		researchQueueRepo:   researchQueueRepo,
	}
}

func (s *Scheduler) Start() {
	s.cron.AddFunc("@every 1m", s.processQueues)
	s.cron.AddFunc("@every 1m", s.processFleets)
	s.cron.AddFunc("@every 5m", s.updateProduction)
	
	s.cron.Start()
	log.Info().Msg("Scheduler started")
}

func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Info().Msg("Scheduler stopped")
}

func (s *Scheduler) processQueues() {
	ctx := context.Background()
	
	s.processBuildingQueues(ctx)
	s.processResearchQueues(ctx)
}

func (s *Scheduler) processBuildingQueues(ctx context.Context) {
	log.Debug().Msg("Processing building queues")
	
	queues, err := s.buildingQueueRepo.GetAllWithActive(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get building queues")
		return
	}
	
	for _, queue := range queues {
		if time.Now().Before(queue.EndTime) {
			continue
		}
		
		err = s.buildingService.CompleteBuilding(ctx, queue.ID)
		if err != nil {
			log.Debug().Err(err).Uint("queue_id", queue.ID).Msg("Failed to complete building")
		} else {
			log.Info().Uint("planet_id", queue.PlanetID).Int("building_id", queue.BuildingID).Int("level", queue.Level).Msg("Building completed")
		}
	}
}

func (s *Scheduler) processResearchQueues(ctx context.Context) {
	log.Debug().Msg("Processing research queues")
	
	queues, err := s.researchQueueRepo.GetAllWithActive(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get research queues")
		return
	}
	
	for _, queue := range queues {
		if time.Now().Before(queue.EndTime) {
			continue
		}
		
		err = s.researchService.CompleteResearch(ctx, queue.ID)
		if err != nil {
			log.Debug().Err(err).Uint("queue_id", queue.ID).Msg("Failed to complete research")
		} else {
			log.Info().Uint("user_id", queue.UserID).Int("research_id", queue.ResearchID).Int("level", queue.Level).Msg("Research completed")
		}
	}
}

func (s *Scheduler) processFleets() {
	ctx := context.Background()
	
	err := s.fleetService.ProcessArrivingMissions(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error processing fleet missions")
	}
}

func (s *Scheduler) updateProduction() {
	ctx := context.Background()
	
	err := s.productionService.UpdateAllPlanets(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error updating production")
	} else {
		log.Debug().Msg("Production updated for all planets")
	}
}
