package core

import (
	"sync"
	"sync/atomic"
	"time"
)

type slot struct {
	mu    sync.Mutex
	conns map[*Connx]struct{}
}

// TimingWheel 基于布尔标记位的单层时间轮
type TimingWheel struct {
	interval time.Duration
	ticker   *time.Ticker
	slots    []*slot
	current  int32

	deadReqCh chan *Connx
	stopCh    chan struct{}
}

func NewTimingWheel(interval time.Duration, slotNum int) *TimingWheel {
	tw := &TimingWheel{
		interval:  interval,
		slots:     make([]*slot, slotNum),
		current:   0,
		deadReqCh: make(chan *Connx, 1024), // 缓冲池，存放待清理连接
		stopCh:    make(chan struct{}),
	}
	for i := range slotNum {
		tw.slots[i] = &slot{
			conns: make(map[*Connx]struct{}),
		}
	}
	// 启动专门负责关闭死连接的异步协程
	go tw.asyncCloser()

	return tw
}

func (tw *TimingWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.run()
}

func (tw *TimingWheel) Stop() {
	close(tw.stopCh)
	if tw.ticker != nil {
		tw.ticker.Stop()
	}
}

// addConn 添加连接到当前的槽位。整个生命周期内连接不再更换槽位。
func (tw *TimingWheel) addConn(conn *Connx) {
	idx := atomic.LoadInt32(&tw.current)
	s := tw.slots[idx]

	s.mu.Lock()
	s.conns[conn] = struct{}{}
	s.mu.Unlock()
}

func (tw *TimingWheel) run() {
	for {
		select {
		case <-tw.ticker.C:
			tw.tick()
		case <-tw.stopCh:
			return
		}
	}
}

func (tw *TimingWheel) tick() {
	// 指针向前推进一步
	next := (atomic.LoadInt32(&tw.current) + 1) % int32(len(tw.slots))
	atomic.StoreInt32(&tw.current, next)

	s := tw.slots[next]

	s.mu.Lock()
	defer s.mu.Unlock()

	for conn := range s.conns {
		// 检查布尔活跃标记
		alive := atomic.LoadInt32(&conn.isAlive)
		if alive == 1 {
			// 如果在过去的一个周期内存活过，重置回 0
			atomic.StoreInt32(&conn.isAlive, 0)
		} else {
			// 一整个周期内没有动静，宣告死亡
			delete(s.conns, conn) // 立即从槽里的链表(Map)中断开

			// 扔给本地缓冲池（如果满了，开启临时协程以保证时间轮不被阻塞）
			select {
			case tw.deadReqCh <- conn:
			default:
				go conn.s.Close(conn)
			}
		}
	}
}

func (tw *TimingWheel) asyncCloser() {
	// 异步线程负责逐个关闭连接，免受网络IO影响阻塞时间轮滴答
	for conn := range tw.deadReqCh {
		conn.s.Close(conn)
	}
}
