package sqs

import (
	"github.com/crowdmob/goamz/sqs"
	"github.com/gchaincl/goku"
)

type SQSProvider struct {
	queue *sqs.Queue
}

func (self SQSProvider) Read() ([]goku.Message, error) {
	resp, err := self.queue.ReceiveMessage(10)
	if err != nil {
		return nil, err
	}

	msgs := make([]goku.Message, len(resp.Messages))
	for i, msg := range resp.Messages {
		msgs[i] = goku.Message(msg)
	}

	return msgs, nil
}

func New(q *sqs.Queue) *SQSProvider {
	provider := &SQSProvider{q}
	return provider
}
