package produceriface

type Producer interface {
	Put([]byte, string) error
}

type Mock struct {
	putCh chan put
}

func NewMock() *Mock {
	return &Mock{
		putCh: make(chan put),
	}
}

func (mock *Mock) DrainReads() {
	go func() {
		for {
			_, _, reply := mock.ExpectRead()
			reply(nil)
		}
	}()
}

type putArgs struct {
	message   []byte
	partition string
}

type put struct {
	req      putArgs
	replyErr chan<- error
}

func (mock *Mock) Put(message []byte, partition string) error {
	replyErr := make(chan error)

	mock.putCh <- put{
		req:      putArgs{message, partition},
		replyErr: replyErr,
	}
	return <-replyErr
}

func (mock *Mock) ExpectRead() ([]byte, string, func(error)) {
	put := <-mock.putCh
	return put.req.message, put.req.partition, func(e error) { put.replyErr <- e }
}
