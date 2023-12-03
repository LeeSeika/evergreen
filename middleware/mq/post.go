package mq

import (
	"context"
	"encoding/json"
	"evergreen/model"
	"time"

	"github.com/rabbitmq/amqp091-go"

	"go.uber.org/zap"
)

func PublishPost(post *model.Post) error {
	b, err := json.Marshal(post)
	if err != nil {
		zap.L().Error("json marshal post failed", zap.Error(err))
		return err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	err = producerChPost.PublishWithContext(ctx, PostExchange, PostRoutingKey, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        b,
	})
	if err != nil {
		zap.L().Error("[RabbitMQ] publish post error failed", zap.Error(err))
		return err
	}

	return nil
}

func GetPostConsumerMsg() (<-chan amqp091.Delivery, error) {
	return consumerChPost.Consume(PostQueue, "", false, false, false, false, nil)
}
