package main

import (
	"log"
	"sync"
	"time"
)

var winMutex map[string]*sync.RWMutex

func init() {
	winMutex = make(map[string]*sync.RWMutex)
}

type timeSlot struct {
	timeStamp time.Time
	count     int
}

type SlidingWindowLimiter struct {
	SlotDuration time.Duration
	WinDuration  time.Duration
	numSlots     int
	windows      map[string][]*timeSlot
	maxReq       int
}

func countReq(win []*timeSlot) int {
	count := 0
	for _, w := range win {
		count += w.count
	}
	return count
}

func NewSlidingWindow(slotDuration, winDuration time.Duration, maxReq int) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		SlotDuration: slotDuration,
		WinDuration:  winDuration,
		numSlots:     int(winDuration / slotDuration),
		windows:      make(map[string][]*timeSlot),
		maxReq:       maxReq,
	}
}

func (swl *SlidingWindowLimiter) getWindow(uid string) []*timeSlot {
	win, ok := swl.windows[uid]
	if !ok {
		win = make([]*timeSlot, 0, swl.numSlots)
	}
	return win
}

func (swl *SlidingWindowLimiter) storeWindow(uid string, win []*timeSlot) {
	swl.windows[uid] = win
}

func (swl *SlidingWindowLimiter) validate(uid string) bool {
	mutex, ok := winMutex[uid]
	if !ok {
		var m sync.RWMutex
		mutex = &m
		winMutex[uid] = mutex
	}
	mutex.Lock()
	defer mutex.Unlock()

	win := swl.getWindow(uid)
	now := time.Now()
	timeoutOffset := -1
	for index, ts := range win {
		if ts.timeStamp.Add(swl.WinDuration).After(now) {
			break
		}
		timeoutOffset = index
	}
	if timeoutOffset != -1 {
		win = win[timeoutOffset+1:]
	}

	var result bool
	if countReq(win) < swl.maxReq {
		result = true
	}
	var lastSlot *timeSlot
	if len(win) > 0 {
		lastSlot = win[len(win)-1]
		if lastSlot.timeStamp.Add(swl.SlotDuration).Before(now) {
			lastSlot = &timeSlot{timeStamp: now, count: 1}
			win = append(win, lastSlot)
		} else {
			lastSlot.count++
		}
	} else {
		lastSlot = &timeSlot{timeStamp: now, count: 1}
		win = append(win, lastSlot)
	}
	swl.storeWindow(uid, win)
	return result
}

func (swl *SlidingWindowLimiter) getUid() string {
	return "127.0.0.1"
}

func (swl *SlidingWindowLimiter) IsLimited() bool {
	return !swl.validate(swl.getUid())
}

func main() {
	limiter := NewSlidingWindow(100*time.Millisecond, time.Second, 10)
	for i := 0; i < 5; i++ {
		log.Println(limiter.IsLimited())
	}
	time.Sleep(100 * time.Millisecond)
	for i := 0; i < 5; i++ {
		log.Println(limiter.IsLimited())
	}
	log.Println("=====1=====")
	log.Println(limiter.IsLimited())
	for _, v := range limiter.windows[limiter.getUid()] {
		log.Println(v.timeStamp, v.count)
	}
	log.Println("=====2=====")
	time.Sleep(time.Second)
	for i := 0; i < 20; i++ {
		log.Println(limiter.IsLimited())
	}
	for _, v := range limiter.windows[limiter.getUid()] {
		log.Println(v.timeStamp, v.count)
	}
}
