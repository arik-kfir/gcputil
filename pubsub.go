package gcputil

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/secureworks/errors"
)

func CreateTopic(ctx context.Context, client *pubsub.Client, topicName string, topicConfig *pubsub.TopicConfig) (*pubsub.Topic, error) {
	topic := client.Topic(topicName)
	if exists, err := topic.Exists(ctx); err != nil {
		return nil, errors.Chain(err, "failed checking if topic '%s' exists", topicName)
	} else if exists {
		return topic, nil
	} else {
		topic, err := client.CreateTopicWithConfig(ctx, topicName, topicConfig)
		if err != nil {
			return nil, errors.Chain(err, "failed creating topic '%s'", topicName)
		} else {
			return topic, nil
		}
	}
}

func CreateSubscription(ctx context.Context, client *pubsub.Client, subscriptionName string, subOpts pubsub.SubscriptionConfig) (*pubsub.Subscription, error) {
	sub := client.Subscription(subscriptionName)
	if exists, err := sub.Exists(ctx); err != nil {
		return nil, errors.Chain(err, "failed checking if subscription '%s' exists", subscriptionName)
	} else if exists {
		return sub, nil
	}

	sub, err := client.CreateSubscription(ctx, subscriptionName, subOpts)
	if err != nil {
		return nil, errors.Chain(err, "failed creating subscription '%s'", subscriptionName)
	} else {
		return sub, nil
	}
}
