package windows

import (
	"sync"
	"time"
)

// 滑动窗口实现
// 需要：限流主体、限流频率、访问记录、滑动窗口实现
// 数据结构：需要支持：倒序访问、尾插入、每次访问删除过期的访问记录
// 数据结构：链表：插入性能好，但是并不需要中间插入，只需要尾部插入；查询性能差，无序；内存不连续，每个节点需要单独分配内存，加剧内存碎片化；
// 切片：只需要append; 可以排序，有序；连续的内存空间

type LimitWindow struct {
	lock     sync.Mutex
	body     int64
	interval time.Duration
	times    uint
	windows  []int64
}

type Option func(*LimitWindow)

// 设置间隔
func SetInterval(interval time.Duration) Option {
	return func(l *LimitWindow) {
		l.interval = interval
	}
}

// 设置间隔内的次数
func SetTimes(times uint) Option {
	return func(l *LimitWindow) {
		l.times = times
	}
}

func Init(body int64, option ...Option) *LimitWindow {
	l := &LimitWindow{
		body:     body,
		interval: 60 * time.Second,
		times:    10,
	}
	for _, o := range option {
		o(l)
	}
	return l
}

func (l *LimitWindow) Slide() bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	now := time.Now().Unix()
	if len(l.windows) == 0 {
		l.windows = append(l.windows, now)
		return true
	}
	// 丢弃 当前时间 - interval 之前的窗口数据
	discardTime := time.Now().Add(-l.interval).Unix()
	lastTime := l.windows[len(l.windows)-1]
	if lastTime < discardTime {
		l.windows = l.windows[0:0]
		l.windows = append(l.windows, now)
		return true
	}
	discardIndex := 0
	for k, v := range l.windows {
		if v >= discardTime {
			discardIndex = k
			break
		}
	}
	if discardIndex > 0 {
		l.windows = l.windows[discardIndex:]
	}
	// 是否超频
	if len(l.windows) >= int(l.times) {
		return false
	}
	// 加入本次访问数据, 串行访问，不需要排序
	l.windows = append(l.windows, now)
	return true
}
