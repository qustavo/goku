package goku

type Message interface{}

type Reader interface {
	Read() (msgs []Message, err error)
}

type Writer interface {
	Write([]Message) (err error)
}

type Queue struct {
	in, out chan Message
}

func (self Queue) Sender() <-chan Message {
	return self.out
}

func (self Queue) Receiver() chan<- Message {
	return self.in
}

func NewQueue(r Reader, w Writer) *Queue {
	q := &Queue{
		in:  make(chan Message),
		out: make(chan Message),
	}

	// Setup Reader
	go func() {
		if r == nil {
			return
		}

		for {
			msgs, err := r.Read()
			if err != nil {
				panic(err)
			}

			for _, msg := range msgs {
				q.out <- msg
			}
		}

	}()

	// Setup Writer
	go func() {
		for {
			msg := <-q.in
			err := w.Write([]Message{msg})
			if err != nil {
				panic(err)
			}
		}
	}()

	return q
}