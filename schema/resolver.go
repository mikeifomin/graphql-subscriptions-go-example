package schema 

import (
	graphql "github.com/neelance/graphql-go"
	"sync"
	"context"
	"github.com/mikeifomin/graphql-subscriptions-go-example/schema/pubsub"
)

type Message struct{
	id  string
	text string
	room string
}

func NewResolver() *Resolver{
	return &Resolver{
		storage: []Message{},
		pubsub: pubsub.New(),
	}
}

type Resolver struct {
	pubsub   *pubsub.Pubsub
  storage []Message
	mutex   sync.Mutex
}

func (r *Resolver) Messages() []Message {
	return []Message{}
}

func (r *Resolver) NewMessage(arg struct{Id graphql.ID; Text string; Room string}) Message {
  r.mutex.Lock(); defer r.mutex.Unlock()
	m := Message{id:string(arg.Id), text:arg.Text, room:arg.Room}
	r.storage = append(r.storage,m)
	r.pubsub.Pub(m.room, m)
  return m
}

// AddedMessage 
func (r *Resolver) AddedMessage(ctx context.Context, arg struct{Room string}) (<-chan Message, error) {
	ch := make(chan Message, 1)
	dataChannel := r.pubsub.Sub(ctx, arg.Room)
	go func(){
		for {
			data, ok := <- dataChannel
			if !ok {
				return
			}
			if msg, ok := data.(Message); ok {
				ch <- msg
			}
		}
	}()
  return ch, nil
}

func (m Message) Id() graphql.ID {
	return graphql.ID(m.id)
}

func (m Message) Text() string {
	return m.text
}

func (m Message) Room() string {
	return m.room
} 
