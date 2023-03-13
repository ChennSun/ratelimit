package bucket

import "time"

func Init(interval time.Duration, times uint) (chan struct{}, chan struct{}) {
	var bucketChan = make(chan struct{}, times)
	var stopChan = make(chan struct{})
	go func() {
		t := time.NewTicker(interval / time.Duration(times))
		defer t.Stop()
		for {
			select {
			case <-t.C:
				bucketChan <- struct{}{}
			case <-stopChan:
				return
			}
		}
	}()
	return bucketChan, stopChan
}
