package pubsub

import (
	"sync"
	"time"
	"context"
)

type Pubsub struct {
	capacity   int
	timeout    time.Duration
	rooms      map[string][]chan interface{}
	mutex      sync.Mutex
}
func New() *Pubsub {
	p := Pubsub{
		capacity: 10,
		timeout: time.Second,
		rooms: make(map[string][](chan interface{})),
	}
	return &p
}

func (p *Pubsub) Sub(ctx context.Context, room string) <-chan interface{} {
	p.mutex.Lock(); defer p.mutex.Unlock()
	_, ok := p.rooms[room] 
	if !ok {
		p.rooms[room] = make([]chan interface{},10)
	}
	s := make(chan interface{}, p.capacity)
	p.rooms[room] = append(p.rooms[room], s)

	go func(){
		<- ctx.Done()
		close(s)
	}()
	return s
}

func (p *Pubsub) Pub(room string, data interface{}) {
	p.mutex.Lock(); defer p.mutex.Unlock()
	list, ok := p.rooms[room]
	if !ok {
		return 
	}

	go func() {
		for _,ch := range list {
			select {
			case ch <- data:
			default:
			   // oops need to close	
			}
		}
	}()
}

