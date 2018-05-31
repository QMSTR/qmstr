package master

import (
	"errors"
	"log"
)

type serverPhaseFailure struct {
	genericServerPhase
}

func (server *server) enterFailureServerPhase() {
	server.currentPhase = &serverPhaseFailure{genericServerPhase{Name: "Fail"}}
	log.Println("Server entered failure state.")
}

func (phase *serverPhaseFailure) Activate() error {
	log.Println("server in failure phase")
	return nil
}

func (phase *serverPhaseFailure) Shutdown() error {
	return errors.New("shutdown not possible failure phase is terminal")
}

func (phase *serverPhaseFailure) GetPhaseID() int32 {
	return PhaseIDFailure
}
