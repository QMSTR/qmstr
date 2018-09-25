package master

import (
	"errors"
	"log"

	"github.com/QMSTR/qmstr/pkg/qmstr/service"
)

type serverPhaseFailure struct {
	genericServerPhase
	cause error
}

func (server *server) enterFailureServerPhase(cause error) {
	server.publishEvent(&service.Event{Class: string(EventPhase), Message: "Entering failure phase"})
	server.currentPhase = &serverPhaseFailure{genericServerPhase{Name: "Fail"}, cause}
	log.Printf("Server entered failure phase due to %v\n", cause)
}

func (phase *serverPhaseFailure) Activate() error {
	log.Println("server in failure phase")
	return nil
}

func (phase *serverPhaseFailure) Shutdown() error {
	return errors.New("shutdown not possible failure phase is terminal")
}

func (phase *serverPhaseFailure) GetPhaseID() service.Phase {
	return service.Phase_FAIL
}
