package scheduler

import (
	"context"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"

	"ogamex-go/internal/service"
)

type Scheduler struct {
	cron               *cron.Cron
	buildingService    *service.BuildingService
	researchService    *service.ResearchService
	fleetService       *service.FleetService
	productionService  *service.ProductionService
}

func NewScheduler(
	buildingService *service.BuildingService,
	researchService *service.ResearchService,
	fleetService *service.FleetService,
	productionService *service.ProductionService,
) *Scheduler {
	return &Scheduler{
		cron:              cron.New(),
		buildingService:   buildingService,
		researchService:   researchService,
		fleetService:      fleetService,
		productionService: productionService,
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
}

func (s *Scheduler) processResearchQueues(ctx context.Context) {
	log.Debug().Msg("Processing research queues")
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
