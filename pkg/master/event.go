package master

import (
	"fmt"
	"log"
	"time"

	"github.com/QMSTR/qmstr/pkg/service"
	"golang.org/x/net/context"
)

type EventClass string

const (
	EventAll    EventClass = "all"
	EventPhase  EventClass = "phase"
	EventModule EventClass = "module"
)

var events = map[string]EventClass{string(EventAll): EventAll, string(EventModule): EventModule, string(EventPhase): EventPhase}

func (s *server) SubscribeEvents(in *service.EventMessage, stream service.ControlService_SubscribeEventsServer) error {
	if _, ok := events[in.Class]; !ok {
		return fmt.Errorf("No such event %s", in.Class)
	}
	eventKey := events[in.Class]
	eventC := make(chan *service.Event)
	err := s.registerEventChannel(eventKey, eventC)
	if err != nil {
		return err
	}

	for event := range eventC {
		stream.Send(event)
	}
	return nil
}

func (s *server) registerEventChannel(eventClass EventClass, eventChannel chan *service.Event) error {
	s.eventMutex.RLock()
	if _, ok := s.eventChannels[eventClass]; !ok {
		return fmt.Errorf("No such event class %v", eventClass)
	}
	s.eventMutex.RUnlock()

	s.eventMutex.Lock()
	s.eventChannels[eventClass] = append(s.eventChannels[eventClass], eventChannel)
	s.eventMutex.Unlock()
	return nil
}

func (s *server) publishEvent(event *service.Event) error {
	s.eventMutex.RLock()
	defer s.eventMutex.RUnlock()
	if _, ok := s.eventChannels[EventClass(event.Class)]; !ok {
		return fmt.Errorf("No such event class %v", event.Class)
	}
	for _, sub := range []EventClass{EventClass(event.Class), EventAll} {
		for _, channel := range s.eventChannels[sub] {
			go func(c chan *service.Event) {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				select {
				case c <- event:
					return
				case <-ctx.Done():
					log.Printf("Could not send event %v to channel", event)
				}
			}(channel)
		}
	}

	return nil
}
