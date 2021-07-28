package throttle

import (
	"container/list"
	"sync"
	"time"

	"github.com/PavelVershinin/throttle/rolling_window"
)

type throttle struct {
	rollingWindow rolling_window.RollingWindow
	requestsLimit int64
	tickCh        chan struct{}
	waitCh        chan struct{}
	mu            sync.Mutex
	queue         *list.List
}

// New Вернёт инициализированный шлюз
// limit максимальное количество заданий обрабатываемых за period времени
func New(limit int64, period time.Duration) *throttle {
	t := &throttle{
		rollingWindow: rolling_window.New(time.Millisecond, int(period.Milliseconds())),
		requestsLimit: limit,
		tickCh:        make(chan struct{}, 1),
		waitCh:        make(chan struct{}, 1),
		queue:         list.New(),
	}
	go func() {
		for {
			select {
			case <-time.Tick(time.Microsecond):
				t.tick()
			case <-t.tickCh:
				t.call()
			}
		}
	}()
	return t
}

// Push Постановка задания в очередь
func (t *throttle) Push(request func()) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.queue.PushBack(request)
	t.tick()
}

// QueueIsFree Вернёт true если следующее задание, может быть обработано без очереди
func (t *throttle) QueueIsFree() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.rollingWindow.Count()+int64(t.queue.Len()) < t.requestsLimit
}

// QueueLength Вернёт количество заданий в очереди
func (t *throttle) QueueLength() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.queue.Len()
}

// Wait Дождётся завершения всех заданий в очереди
func (t *throttle) Wait() {
	<-t.waitCh
}

func (t *throttle) call() {
	t.mu.Lock()
	defer t.mu.Unlock()

	for t.queue.Len() > 0 && t.rollingWindow.Count() < t.requestsLimit {
		t.rollingWindow.Add(1)
		el := t.queue.Front()
		el.Value.(func())()
		t.queue.Remove(el)
	}
	if t.queue.Len() == 0 {
		select {
		case t.waitCh <- struct{}{}:
		default:
		}
	}
}

func (t *throttle) tick() {
	select {
	case t.tickCh <- struct{}{}:
	default:
	}
}
