package main

import (
	"flag"

	"github.com/crowdmob/goamz/sqs"
	"github.com/gchaincl/goku"
	. "github.com/gchaincl/goku/sqs"
)

type writer struct{}

func (self writer) Write(msgs []goku.Message) error {
	for i, msg := range msgs {
		println(i, msg.(string))
	}
	return nil
}

func getAWSQueue(access, secret, region, name string) *sqs.Queue {
	aws_sqs, err := sqs.NewFrom(access, secret, region)
	if err != nil {
		panic(err)
	}

	queue, err := aws_sqs.GetQueue(name)
	if err != nil {
		panic(err)
	}

	return queue
}

func main() {
	var access, secret, region, queueName string
	flag.StringVar(&access, "access", "", "AWS Access Key")
	flag.StringVar(&secret, "secret", "", "AWS Secret Key")
	flag.StringVar(&region, "region", "us-east-1", "AWS Region")
	flag.StringVar(&queueName, "queue", "", "SQS Queue Name")
	flag.Parse()

	if flag.NFlag() < 3 {
		println("Missing flags, usage:")
		flag.PrintDefaults()
		return
	}

	aws_queue := getAWSQueue(access, secret, region, queueName)
	reader := New(aws_queue)

	q := goku.NewQueue(reader, &writer{})

	for {
		// Receive messages from goku
		goku_msg := <-q.Sender()

		// Cast to sqs.Message
		sqs_msg := goku_msg.(sqs.Message)

		_, err := aws_queue.DeleteMessage(&sqs_msg)
		if err != nil {
			panic(err)
		}

		// Send to Writer
		q.Receiver() <- goku.Message(sqs_msg.Body)
	}
}
