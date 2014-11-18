package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/crowdmob/goamz/sqs"
	"github.com/gchaincl/goku"
	. "github.com/gchaincl/goku/sqs"
)

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
	sqs_provider := New(aws_queue)

	q := goku.NewQueue(sqs_provider, sqs_provider)

	// Send Messages (4/sec)
	go func() {
		for i := 0; ; i++ {
			q.Receiver() <- fmt.Sprintf("Hello #%d", i)
			time.Sleep(250 * time.Millisecond)
		}
	}()

	// Receive Messages
	for i := 0; ; i++ {
		goku_msg := <-q.Sender()

		// Cast to sqs.Message so we can operate with goamz
		sqs_msg := goku_msg.(sqs.Message)
		log.Println(sqs_msg.Body)

		// Delete message
		_, err := aws_queue.DeleteMessage(&sqs_msg)
		if err != nil {
			panic(err)
		}
	}
}
