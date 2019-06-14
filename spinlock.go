package spinlock

import (
	"runtime"
	"sync/atomic"
)

const (
	// same as go mutex active_spin
	activeSpin = 4
)

type SpinLock struct {
	f uint32
}

func New() *SpinLock {
	return &SpinLock{}
}

func (sl *SpinLock) LockSched() {
	for !sl.TryLock() {
		runtime.Gosched() //allow other goroutines to do stuff.
	}
}

func (sl *SpinLock) Lock() {
	for index := 0; index < activeSpin; index++ {
		if sl.TryLock() {
			return
		}
		continue
	}

	sl.LockSched()
}

func (sl *SpinLock) Unlock() {
	atomic.StoreUint32(&sl.f, 0)
}

func (sl *SpinLock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&sl.f, 0, 1)
}

func (sl *SpinLock) String() string {
	if atomic.LoadUint32(&sl.f) == 1 {
		return "Locked"
	}
	return "Unlocked"
}
