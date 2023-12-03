package mq

import (
	"evergreen/settings"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	CommentExchange = "comment_exchange"
	VoteQueue       = "vote_queue"
	VoteRoutingKey  = "vote_key"

	PostExchange   = "post_exchange"
	PostQueue      = "post_queue"
	PostRoutingKey = "post_key"

	DeadLetterExchange = "dlx_exchange"
	DeadLetterQueue    = "dlx_queue"
)

var (
	conn              *amqp091.Connection
	producerChComment *amqp091.Channel
	producerChPost    *amqp091.Channel
	consumerChPost    *amqp091.Channel
	consumerChVote    *amqp091.Channel
)

func Init() error {
	user := settings.Conf.RabbitMQConfig.User
	password := settings.Conf.RabbitMQConfig.Password
	host := settings.Conf.RabbitMQConfig.Host
	port := settings.Conf.RabbitMQConfig.Port

	var err error
	var url string

	url = fmt.Sprintf("amqp://%s:%s@%s:%d/", user, password, host, port)
	conn, err = amqp091.Dial(url)
	if err != nil {
		zap.L().Error("[RabbitMQ] dial failed", zap.Error(err))
		return err
	}

	err = initPostChannel()
	if err != nil {
		return err
	}

	err = initVoteChannel()
	if err != nil {
		return err
	}

	err = initCommentChannel()
	if err != nil {
		return err
	}

	return nil
}

func Close() error {
	err := conn.Close()
	if err != nil {
		zap.L().Error("[RabbitMQ] rabbitmq close error", zap.Error(err))
		return err
	}
	return nil
}

func initPostChannel() error {
	var err error
	// producer
	producerChPost, err = conn.Channel()
	if err != nil {
		zap.L().Error("[RabbitMQ] get post producer channel from connection failed", zap.Error(err))
		return err
	}

	err = producerChPost.ExchangeDeclare(
		PostExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		zap.L().Error("[RabbitMQ] declare post exchange failed", zap.Error(err))
		return err
	}

	_, err = producerChPost.QueueDeclare(
		PostQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Error("[RabbitMQ] declare post queue failed", zap.Error(err))
		return err
	}

	err = producerChPost.QueueBind(PostQueue, PostRoutingKey, PostExchange, false, nil)
	if err != nil {
		zap.L().Error("[RabbitMQ] bind post queue failed", zap.Error(err))
		return err
	}

	// dlx
	err = initDLX(producerChPost, PostRoutingKey)
	if err != nil {
		return err
	}

	// consumer
	consumerChPost, err = conn.Channel()
	if err != nil {
		zap.L().Error("[RabbitMQ] get post consumer channel from connection failed", zap.Error(err))
		return err
	}
	_, err = consumerChPost.QueueDeclare(
		PostQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	return nil
}

func initCommentChannel() error {
	var err error
	// producer
	producerChComment, err = conn.Channel()
	if err != nil {
		zap.L().Error("[RabbitMQ] get comment producer channel from connection failed", zap.Error(err))
		return err
	}

	err = producerChComment.ExchangeDeclare(
		CommentExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		zap.L().Error("[RabbitMQ] declare comment exchange failed", zap.Error(err))
		return err
	}

	_, err = producerChComment.QueueDeclare(
		VoteQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		zap.L().Error("[RabbitMQ] declare vote queue failed", zap.Error(err))
		return err
	}

	err = producerChComment.QueueBind(VoteQueue, VoteRoutingKey, CommentExchange, false, nil)
	if err != nil {
		zap.L().Error("[RabbitMQ] bind vote queue failed", zap.Error(err))
		return err
	}

	// dlx
	err = initDLX(producerChComment, VoteRoutingKey)
	if err != nil {
		return err
	}

	return nil
}

func initVoteChannel() error {
	// consumer
	var err error
	consumerChVote, err = conn.Channel()
	if err != nil {
		zap.L().Error("[RabbitMQ] get vote consumer channel from connection failed", zap.Error(err))
		return err
	}
	_, err = consumerChVote.QueueDeclare(
		VoteQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	return nil
}

func initDLX(ch *amqp091.Channel, routingKey string) error {
	var err error
	// dead letter
	err = ch.ExchangeDeclare(
		DeadLetterExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		zap.L().Error("[RabbitMQ] declare dlx exchange failed", zap.Error(err))
		return err
	}
	args := map[string]interface{}{}
	args["x-dead-letter-exchange"] = DeadLetterExchange
	_, err = ch.QueueDeclare(
		DeadLetterQueue,
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		zap.L().Error("[RabbitMQ] declare post dlx queue failed", zap.Error(err))
		return err
	}
	err = ch.QueueBind(DeadLetterQueue, routingKey, DeadLetterExchange, false, nil)
	if err != nil {
		zap.L().Error("[RabbitMQ] bind post dlx queue failed", zap.Error(err))
		return err
	}

	return nil
}
