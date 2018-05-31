package master

import (
	"errors"
	"log"
)

type serverPhaseFailure struct {
	genericServerPhase
	cause error
}

func (server *server) enterFailureServerPhase(cause error) {
	server.currentPhase = &serverPhaseFailure{genericServerPhase{Name: "Fail"}, cause}
	log.Printf("Server entered failure state due to %v\n", cause)
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
