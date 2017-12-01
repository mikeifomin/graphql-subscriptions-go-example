package schema 

import (
	graphql "github.com/neelance/graphql-go"
)

type Resolver struct {}
func (r *Resolver) Messages() []Message {
	return []Message{}
}

func (r *Resolver) AddedMessage(arg struct{Room string}) Message {
  return Message{}
}

func (r *Resolver) NewMessage(arg struct{Id graphql.ID; Text string; Room string}) Message {
	return Message{}
}

type Message struct{
	id  string
	text string
	room string
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

