package log

const (
	warnSize = 500
)

var logChan chan string

func InitAsync() {
	logChan = make(chan string, 1000)
	go func() {
		for {
			if data, ok := <-logChan; ok {
				if len(logChan) > warnSize {
					Warnf("async log size warning, size:%d", len(logChan))
				}
				Debug(data)
			} else {
				close(logChan)
				break
			}
		}
	}()
}

func Log(msg string) {
	logChan <- msg
}
