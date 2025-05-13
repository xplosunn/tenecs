package async_test

import (
	"errors"
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/typer/async"
	"testing"
	"time"
)

func TestRunAsyncSuccess(t *testing.T) {
	fn := func() (int, error) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	}

	async := async.RunAsync(fn)

	start := time.Now()
	result, err := async.Await()
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != 42 {
		t.Errorf("Expected result 42, got %d", result)
	}

	if duration < 100*time.Millisecond {
		t.Error("Await returned before the async function completed")
	}
}

func TestRunAsyncError(t *testing.T) {
	fn := func() (string, error) {
		return "", errors.New("something went wrong")
	}

	async := async.RunAsync(fn)

	_, err := async.Await()
	assert.EqualError(t, err, "something went wrong")
}

func TestRunAsyncConcurrent(t *testing.T) {
	slowFn := func() (int, error) {
		time.Sleep(200 * time.Millisecond)
		return 1, nil
	}

	fastFn := func() (int, error) {
		time.Sleep(50 * time.Millisecond)
		return 2, nil
	}

	slowAsync := async.RunAsync(slowFn)
	fastAsync := async.RunAsync(fastFn)

	start := time.Now()
	fastResult, fastErr := fastAsync.Await()
	fastDuration := time.Since(start)

	if fastErr != nil {
		t.Errorf("Fast function error: %v", fastErr)
	}
	if fastResult != 2 {
		t.Errorf("Expected fast result 2, got %d", fastResult)
	}
	if fastDuration < 50*time.Millisecond || fastDuration > 100*time.Millisecond {
		t.Errorf("Fast function took unexpected time: %v", fastDuration)
	}

	slowResult, slowErr := slowAsync.Await()
	slowDuration := time.Since(start)

	if slowErr != nil {
		t.Errorf("Slow function error: %v", slowErr)
	}
	if slowResult != 1 {
		t.Errorf("Expected slow result 1, got %d", slowResult)
	}
	if slowDuration < 200*time.Millisecond {
		t.Errorf("Slow function completed too quickly: %v", slowDuration)
	}
}

func TestRunAsyncImmediate(t *testing.T) {
	fn := func() (float64, error) {
		return 3.14, nil
	}

	async := async.RunAsync(fn)

	result, err := async.Await()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 3.14 {
		t.Errorf("Expected 3.14, got %f", result)
	}
}
