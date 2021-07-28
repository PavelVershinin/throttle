package rolling_window

import (
	"container/list"
	"time"
)

type RollingWindow interface {
	Add(int64)
	Count() int64
}

type rollingWindow struct {
	atomDuration       time.Duration
	windowSize         int
	initTime           time.Time
	queue              *list.List
	currentAtomElement *list.Element
	currentAtomNumber  int64
	count              int64
}

// New Вернёт инициализированный счётчик
// atomDuration Минимальная единица времени
// windowSize Ширина окна в atomDuration
func New(atomDuration time.Duration, windowSize int) RollingWindow {
	rw := &rollingWindow{
		atomDuration:      atomDuration,
		windowSize:        windowSize,
		initTime:          time.Now(),
		queue:             list.New(),
		currentAtomNumber: 0,
		count:             0,
	}
	return rw
}

// Count Вернет количество событий в окне на текущий момент
func (rw *rollingWindow) Count() int64 {
	rw.moveWindow()
	return rw.count
}

// Add Добавит count событий в окно
func (rw *rollingWindow) Add(count int64) {
	rw.moveWindow()
	rw.currentAtomElement.Value = rw.currentAtomElement.Value.(int64) + count
	rw.count += count
}

func (rw *rollingWindow) moveWindow() {
	atomsAfterInit := (time.Now().UnixNano() - rw.initTime.UnixNano()) / rw.atomDuration.Nanoseconds()
	for ; rw.currentAtomNumber <= atomsAfterInit; rw.currentAtomNumber++ {
		rw.currentAtomElement = rw.queue.PushFront(int64(0))
	}
	for rw.windowSize < rw.queue.Len() {
		el := rw.queue.Back()
		rw.count -= el.Value.(int64)
		rw.queue.Remove(el)
	}
}
