package goku

import (
	"runtime"
	"testing"
	"time"
)

type TestReader struct{}

func (self TestReader) Read() ([]Message, error) {
	time.Sleep(10 * time.Millisecond)
	return []Message{"Hello"}, nil
}

type TestWriter struct {
	msgs []Message
}

func (self *TestWriter) Write(msgs []Message) error {
	for _, msg := range msgs {
		self.msgs = append(self.msgs, msg)
	}
	return nil
}

func TestNewQueueSetupReader(t *testing.T) {
	t.Parallel()

	q := NewQueue(&TestReader{}, nil)

	select {
	case <-q.Sender():
		break
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout expired")
	}
}

func TestNewQueueSetsupWriter(t *testing.T) {
	t.Parallel()

	writer := &TestWriter{}
	q := NewQueue(nil, writer)

	q.Receiver() <- "Hello"
	q.Receiver() <- "World"

	// Ensure writer gorouting run
	runtime.Gosched()

	if count := len(writer.msgs); count != 2 {
		t.Errorf("messages written == %d, want 2", count)
	}
}
