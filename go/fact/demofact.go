package fact

import (
	"context"
	"errors"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	ErrFactorizationCancelled = errors.New("cancelled")
	ErrWriterInteraction      = errors.New("writer interaction")
)

type Factorizer interface {
	Factorize(ctx context.Context, numbers []int, writer io.Writer) error
}

type factorizerImpl struct {
	WithFactorizationWorkers int
	WithWriteWorkers         int
}

var ErrInvalidFactWorkersAmount = errors.New("unable to complete factorization")
var ErrInvalidWriteWorkersAmount = errors.New("unable to write")

func New(opts ...FactorizeOption) (*factorizerImpl, error) {
	fi := &factorizerImpl{
		WithFactorizationWorkers: runtime.GOMAXPROCS(0),
		WithWriteWorkers:         runtime.GOMAXPROCS(0),
	}

	for _, opt := range opts {
		opt(fi)
	}

	if fi.WithFactorizationWorkers < 1 {
		return nil, fmt.Errorf(
			"%w: factorization workers %d",
			ErrInvalidFactWorkersAmount,
			fi.WithFactorizationWorkers,
		)
	}

	if fi.WithWriteWorkers < 1 {
		return nil, fmt.Errorf(
			"%w: write workers %d",
			ErrInvalidWriteWorkersAmount,
			fi.WithWriteWorkers,
		)
	}

	return fi, nil
}

type FactorizeOption func(*factorizerImpl)

func WithFactorizationWorkers(workers int) FactorizeOption {
	return func(fi *factorizerImpl) {
		fi.WithFactorizationWorkers = workers
	}
}

func WithWriteWorkers(workers int) FactorizeOption {
	return func(fi *factorizerImpl) {
		fi.WithWriteWorkers = workers
	}
}

func (f *factorizerImpl) Factorize(
	ctx context.Context,
	numbers []int,
	writer io.Writer,
) error {
	mainWg := sync.WaitGroup{}
	defer mainWg.Wait()

	numbersToFactorize := make(chan int)
	factorizationsToWrite := make(chan string)
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	var (
		firstErr error
		errOnce  sync.Once
	)

	fail := func(err error) {
		errOnce.Do(func() {
			firstErr = err
			cancel(err)
		})
	}

	mainWg.Add(1)
	go func() {
		defer close(numbersToFactorize)
		defer mainWg.Done()
		for _, num := range numbers {
			select {
			case numbersToFactorize <- num:
			case <-ctx.Done():
				fail(errors.Join(context.Cause(ctx), ErrFactorizationCancelled))
				return
			}
		}
	}()

	wgFW := new(sync.WaitGroup)
	for i := 0; i < f.WithFactorizationWorkers; i++ {
		wgFW.Add(1)
		mainWg.Add(1)
		go func() {
			defer wgFW.Done()
			defer mainWg.Done()
			for {
				select {
				case <-ctx.Done():
					fail(errors.Join(context.Cause(ctx), ErrFactorizationCancelled))
					return
				case n, ok := <-numbersToFactorize:
					if !ok {
						return
					}

					select {
					case <-ctx.Done():
						fail(errors.Join(context.Cause(ctx), ErrFactorizationCancelled))
						return
					case factorizationsToWrite <- f.FactorizeSingle(n):
					}
				}
			}
		}()
	}

	mainWg.Add(1)
	go func() {
		defer close(factorizationsToWrite)
		defer mainWg.Done()
		wgFW.Wait()
	}()

	wgWW := new(sync.WaitGroup)
	for i := 0; i < f.WithWriteWorkers; i++ {
		wgWW.Add(1)
		mainWg.Add(1)
		go func() {
			defer wgWW.Done()
			defer mainWg.Done()
			for {
				select {
				case <-ctx.Done():
					fail(errors.Join(context.Cause(ctx), ErrFactorizationCancelled))
					return
				case fact, ok := <-factorizationsToWrite:
					if !ok {
						return
					}

					if _, err := writer.Write([]byte(fact)); err != nil {
						fail(fmt.Errorf("%w: %w", ErrWriterInteraction, err))
						return
					}
				}
			}
		}()
	}

	wgWW.Wait()

	if firstErr != nil {
		return firstErr
	}

	if cause := context.Cause(ctx); cause != nil {
		if errors.Is(cause, context.Canceled) ||
			errors.Is(cause, context.DeadlineExceeded) {
			return ErrFactorizationCancelled
		}
		return cause
	}

	return nil
}

func (f *factorizerImpl) FactorizeSingle(cur int) string {
	var divisors []int
	if cur == 1 {
		return "1 = 1\n"
	}

	if cur == 0 {
		return "0 = 0\n"
	}

	orig := cur
	if cur < 0 {
		divisors = append(divisors, -1)
		cur *= -1
	}

	divider := 2
	for cur%divider == 0 {
		cur /= divider
		divisors = append(divisors, divider)
	}

	for divider = 3; divider*divider <= cur && divider*divider > 0; divider += 2 {
		for cur%divider == 0 {
			divisors = append(divisors, divider)
			cur /= divider
		}
	}

	if cur > 1 {
		divisors = append(divisors, cur)
	}

	var sb strings.Builder
	for i, div := range divisors {
		if i > 0 {
			sb.WriteString(" * ")
		}
		sb.WriteString(strconv.Itoa(div))
	}

	return strconv.Itoa(orig) + " = " + sb.String() + "\n"
}
