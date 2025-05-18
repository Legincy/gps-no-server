package mqtt

import (
	"github.com/rs/zerolog"
	"gps-no-server/internal/common/logger"
	"gps-no-server/internal/infrastructure/mqtt/interfaces"
	"strings"
)

type Registry struct {
	subscriptions []interfaces.Subscription
	log           zerolog.Logger
}

func NewSubscriptionRegistry() *Registry {
	return &Registry{
		subscriptions: make([]interfaces.Subscription, 0),
		log:           logger.GetLogger("mqtt"),
	}
}

func (r *Registry) Register(subscription interfaces.Subscription) {
	r.subscriptions = append(r.subscriptions, subscription)
}

func (r *Registry) GetAllTopics() []string {
	topics := make([]string, 0)
	for _, subscription := range r.subscriptions {
		topics = append(topics, subscription.GetTopics()...)
	}
	return topics
}

func (r *Registry) GetAllSubscriptions(topic string) []interfaces.Subscription {
	handlers := make([]interfaces.Subscription, 0)

	for _, subscription := range r.subscriptions {
		for _, pattern := range subscription.GetTopics() {
			if TopicMatches(pattern, topic) {
				handlers = append(handlers, subscription)
				break
			}
		}
	}

	return handlers
}

func TopicMatches(pattern, topic string) bool {
	if pattern == topic {
		return true
	}

	if strings.HasSuffix(pattern, "#") {
		prefix := pattern[:len(pattern)-1]

		if !strings.HasPrefix(prefix, "/") {
			return false
		}

		return strings.HasPrefix(topic, prefix)
	}

	patternSegments := strings.Split(pattern, "/")
	topicSegments := strings.Split(topic, "/")

	if len(patternSegments) != len(topicSegments) {
		return false
	}

	for i := 0; i < len(patternSegments); i++ {
		if patternSegments[i] == "+" {
			continue
		}

		if patternSegments[i] != topicSegments[i] && patternSegments[i] != "+" {
			return false
		}
	}

	return true
}
