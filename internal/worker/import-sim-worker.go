package worker

import (
	"log"
	"reconTool/internal/model"
	"reconTool/internal/service"
)

type ImportSimProcessor struct{}

func (p *ImportSimProcessor) Process(id int, sim model.SimModel) {
	err := service.CheckAndInsertSim(sim)

	if err != nil {
		log.Printf("[Worker %d][SKIP] %v", id, err)
	} else {
		log.Printf("[Worker %d][SUCCESS] Process done SIM %s", id, sim.Isdn)
	}
}
