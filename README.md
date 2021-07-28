# throttle


[![Go Report Card](https://goreportcard.com/badge/github.com/PavelVershinin/throttle)](https://goreportcard.com/report/github.com/PavelVershinin/throttle)

Ограничивает количество выполняемых заданий за единицу времени. Если добавить заданий больше установленного лимита, они будут добавлены в очередь и распределены по времени запуска, так чтобы не превысить лимит.

# Использование
```
// Создаём дроссель с пропускной способностью 50 заданий в секунду
th := throttle.New(50, time.Second)

// Добавляем задания
th.Push(func() {
    // Какое-то очень полезное действие
})
...

// Дождёмся полного завершения всех заданий
th.Wait()
```