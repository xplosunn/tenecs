package async

type Async[R any] struct {
	Await func() (R, error)
}

func RunAsync[R any](f func() (R, error)) Async[R] {
	sem := make(semaphore)
	var result R
	var err error
	go func() {
		result, err = f()
		sem <- empty{}
	}()
	return Async[R]{
		Await: func() (R, error) {
			<-sem
			return result, err
		},
	}
}

type empty struct{}
type semaphore chan empty
