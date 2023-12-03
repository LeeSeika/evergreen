package mq

import (
	"context"
	"encoding/json"
	"evergreen/model"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func PublishComment(comment *model.Comment) error {
	b, err := json.Marshal(comment)
	if err != nil {
		zap.L().Error("json marshal comment failed", zap.Error(err))
		return err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	err = producerChComment.PublishWithContext(ctx, CommentExchange, VoteRoutingKey, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        b,
	})
	if err != nil {
		zap.L().Error("[RabbitMQ] publish post failed", zap.Error(err))
		return err
	}

	return nil
}
