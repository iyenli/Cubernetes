package faker

import "time"

type Recorder interface {
	WatchCenter()
	QueryRecord(fn string, period time.Duration) int
	WatchRequest() <-chan string
}

func NewRecorder() (Recorder, error) {
	return &fakeRecorder{
		begin:       time.Now(),
		requestChan: make(chan string),
	}, nil
}

type fakeRecorder struct {
	begin       time.Time
	requestChan chan string
}

func (fr *fakeRecorder) QueryRecord(fn string, period time.Duration) int {

	since := time.Since(fr.begin)
	if since < time.Second*30 {
		return 0
	} else if since < time.Second*60 {
		return 5
	} else if since < time.Second*90 {
		return 20
	} else {
		return 0
	}
}

func (fr *fakeRecorder) WatchRequest() <-chan string {
	return fr.requestChan
}

func (fr *fakeRecorder) WatchCenter() {
	defer close(fr.requestChan)
	for {
		time.Sleep(time.Second * 5)
		since := time.Since(fr.begin)
		if since > time.Second*30 && since < time.Second*90 {
			fr.requestChan <- "fake-func"
		}
	}
}
