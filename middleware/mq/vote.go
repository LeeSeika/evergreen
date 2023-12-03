package mq

import "github.com/rabbitmq/amqp091-go"

func GetVoteConsumerMsg() (<-chan amqp091.Delivery, error) {
	return consumerChVote.Consume(VoteQueue, "", false, false, false, false, nil)
}
