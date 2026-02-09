package worker

import (
	"log"
	"reconTool/internal/model"
	"reconTool/internal/service"
)

type SyncSimProcessor struct{}

func (p *SyncSimProcessor) Process(id int, sim model.SimModel) {
	err := service.SyncSim(sim)

	if err != nil {
		log.Printf("[Worker %d][SKIP] %v", id, err)
	} else {
		log.Printf("[Worker %d][SUCCESS] Process done SIM %s", id, sim.Isdn)
	}
}
