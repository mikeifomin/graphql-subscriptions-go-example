package graphql

import (
	"testing"
	"context"
	"reflect"
	"encoding/json"
	"fmt"
)

var testSchema = `
# start_used_for_js_mock_server
schema {
	query: Query
  subscription: Subscription
	mutation: Mutation
}

type Query { 
	messages: [Message!]! 
	#online: [User!]!
}

type Subscription {
	addedMessage(room: String!): Message!
	#joinUser: User!
	#leaveUser: User!
}

type Mutation {
	newMessage(id: ID!, text: String!, room: String!): Message!
}

type Message { 
	id: ID! 
	text: String!
	room: String!
}
# end_used_for_js_mock_server
`

type Resolver struct {}
func (r *Resolver) Messages() []Message {
	return []Message{}
}
func (r *Resolver) NewMessage(arg struct{Id ID; Text string; Room string}) Message {
	return Message{}
}

type Message struct{
	id  string
	text string
	room string
}

func (m Message) Id() ID {
	return ID(m.id)
}

func (m Message) Text() string {
	return m.text
}

func (m Message) Room() string {
	return m.room
} 

func (r *Resolver) AddedMessage(arg struct{Room string}) (<-chan Message, error) {
	ch := make(chan Message, 1)
	count := 0
	max := 3

	go func(){
		for {
			ch <- Message{id:fmt.Sprint(count), room:arg.Room}
			count++
			if count > max {
				close(ch)
				return
			}
		}
	}()
  return ch, nil
}

func TestSubs(t *testing.T) {
	schemaCompiled := MustParseSchema(testSchema, &Resolver{})

  ctx := context.Background()
	queryString := `
	subscription addedMessage($room: String!) {
		addedMessage(room: $room) { id room }
	}`
	variables := map[string]interface{}{"room":"default"}
	opName := "addedMessage"
	ch := schemaCompiled.Subscribe(ctx,queryString,opName,variables)

	type event struct {
		Id   string `json:"id"`
		Room string `json:"room"`
	}
  expected := []event{
		event{"0","default"},
		event{"1","default"},
		event{"2","default"},
		event{"3","default"},
	}

  actual := []event{} 
	for {
		data, ok := <- ch
		if !ok {
			break
		}
		if data.Errors != nil &&  len(data.Errors) > 0 {
			t.Errorf("errors comes from Subscription: %v",data.Errors)
		}
		if data.Data == nil {
			t.Errorf("no data found from Subs: %v", data)
			continue
		}
		var e event
		err := json.Unmarshal(data.Data, &e)
		if err != nil {
			t.Error(err)
		}
		actual = append(actual, e)
	}
	if !reflect.DeepEqual(actual,expected) {
    t.Fatalf("actual data \n %v \n not equal expected: \n %v",actual,expected)
	}
}
