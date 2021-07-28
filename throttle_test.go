package throttle_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/PavelVershinin/throttle"
)

func TestThrottle(t *testing.T) {
	executedTasks := int32(0)
	th := throttle.New(5, time.Millisecond)

	// Ни одного задания не добавлено, задания будут обработаны без очереди
	assert.Equal(t, true, th.QueueIsFree())

	for i := 0; i < 5; i++ {
		th.Push(func() {
			atomic.AddInt32(&executedTasks, 1)
		})
	}

	// Добавлено 5 заданий, очередь занята
	assert.Equal(t, false, th.QueueIsFree())

	for i := 0; i < 15; i++ {
		th.Push(func() {
			atomic.AddInt32(&executedTasks, 1)
		})
	}

	th.Wait()

	// После Wait очередь будет пуста
	require.Equal(t, 0, th.QueueLength())

	// Все задания будут выполнены
	require.Equal(t, int32(20), executedTasks)
}
