package log

var logChan chan string

func InitAsync() {
	logChan = make(chan string, 1000)
	go func() {
		for {
			data, ok := <-logChan
			if ok {
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
