package rolling_window_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"bou.ke/monkey"
	"github.com/PavelVershinin/throttle/rolling_window"
)

func TestRollingWindow(t *testing.T) {
	pathTime := func(initTime time.Time, timePassed time.Duration) {
		monkey.Patch(time.Now, func() time.Time {
			return initTime.Add(timePassed)
		})
	}

	for _, atomDuration := range []time.Duration{time.Nanosecond, time.Microsecond, time.Millisecond, time.Second, time.Minute, time.Hour, time.Nanosecond * 6} {
		t.Run(atomDuration.String(), func(t *testing.T) {
			initTime := time.Now()
			rw := rolling_window.New(atomDuration*1, 2)

			// Добавим 1 и сразу получим, должно получится 1
			pathTime(initTime, atomDuration*0)
			rw.Add(1)
			cnt := rw.Count()
			require.Equal(t, int64(1), cnt)

			// Подождём 2ед и получим, окно должно успеть съехать и счётчик обнулиться
			pathTime(initTime, atomDuration*2)
			cnt = rw.Count()
			require.Equal(t, int64(0), cnt)

			// Добавим 1, подождём 1ед и добавим ещё 1, окно уехать не успеет, результат будет 2
			rw.Add(1)
			pathTime(initTime, atomDuration*3)
			rw.Add(1)
			cnt = rw.Count()
			require.Equal(t, int64(2), cnt)

			// Подождём ещё 1ед, окно должно уехать наполовину, результат будет 1
			pathTime(initTime, atomDuration*4)
			require.Equal(t, int64(1), rw.Count())
		})
	}
}
