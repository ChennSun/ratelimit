package windows

import (
	"testing"
	"time"
)

func TestLimitWindow_Slide(t *testing.T) {
	rateLimit := Init(SetTimes(1000), SetInterval(time.Second*10))
	t.Run("test1", func(t *testing.T) {
		if got := rateLimit.Slide(); got != true {
			t.Errorf("Slide() = %v, want %v", got, true)
		}
	})
}
