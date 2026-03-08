package subscription

import (
	"sync"

	"go.vervstack.ru/matreshka/internal/domain"
	"go.vervstack.ru/matreshka/internal/service"
)

type PubSubService struct {
	m                        sync.RWMutex
	serviceNameToSubscribers map[string]map[service.Subscriber]struct{}
}

func New() *PubSubService {
	return &PubSubService{
		serviceNameToSubscribers: make(map[string]map[service.Subscriber]struct{}),
	}
}

func (s *PubSubService) Publish(patch domain.PatchConfigRequest) {
	s.m.RLock()
	defer s.m.RUnlock()

	for c := range s.serviceNameToSubscribers[patch.ConfigName] {
		c.Consume(patch)
	}
}

func (s *PubSubService) Subscribe(sub service.Subscriber, serviceNames ...string) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, serviceName := range serviceNames {
		topicSubsMap, ok := s.serviceNameToSubscribers[serviceName]
		if !ok {
			topicSubsMap = make(map[service.Subscriber]struct{})
			s.serviceNameToSubscribers[serviceName] = topicSubsMap
		}

		topicSubsMap[sub] = struct{}{}
	}
}

func (s *PubSubService) Unsubscribe(sub service.Subscriber, serviceNames ...string) {
	s.m.Lock()
	defer s.m.Unlock()

	for _, serviceName := range serviceNames {
		delete(s.serviceNameToSubscribers[serviceName], sub)
	}
}

func (s *PubSubService) Shutdown() {
	s.m.Lock()
	defer s.m.Unlock()

	subs := make(map[service.Subscriber]struct{})

	for serviceName, publishers := range s.serviceNameToSubscribers {
		for sub := range publishers {
			subs[sub] = struct{}{}
		}

		delete(s.serviceNameToSubscribers, serviceName)
	}

	for sub := range subs {
		sub.Stop()
	}
}

func (s *PubSubService) StopSubscription(c service.Subscriber) {
	s.m.Lock()
	defer s.m.Unlock()

	c.Stop()

	for _, topic := range s.serviceNameToSubscribers {
		delete(topic, c)
	}

}
