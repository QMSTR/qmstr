package master

import (
	"fmt"
	"log"
	"time"

	"github.com/QMSTR/go-qmstr/service"
	"golang.org/x/net/context"
)

func (s *server) SubscribeEvents(in *service.EventMessage, stream service.ControlService_SubscribeEventsServer) error {
	eventC := make(chan *service.Event)
	err := s.registerEventChannel(in.Class, eventC)
	if err != nil {
		return err
	}

	for event := range eventC {
		stream.Send(event)
	}
	return nil
}

func (s *server) registerEventChannel(eventClass service.EventClass, eventChannel chan *service.Event) error {
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
	if _, ok := s.eventChannels[service.EventClass(event.Class)]; !ok {
		return fmt.Errorf("No such event class %v", event.Class)
	}
	for _, sub := range []service.EventClass{service.EventClass(event.Class), service.EventClass_ALL} {
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
