package schema

import (
	"testing"
	"context"
	"reflect"
	//"sync"
	"time"
	graphql "github.com/neelance/graphql-go"
)

func TestOneMessage(t *testing.T) {
	
  r := NewResolver() 
	ctx := context.Background()
	ch, err := r.AddedMessage(ctx, struct{Room string}{"foo"})
	if err != nil {
		t.Fatal(err)
	}
  r.NewMessage(struct {
		Id graphql.ID
		Text string
		Room string
	}{"1", "hi foo", "foo"})
	select {
		case msg, ok := <-ch:
			if !ok {
				t.Error("channel closed, before first message")
			}
			expected := Message{id:"1",room:"foo",text:"hi foo"}
			if !reflect.DeepEqual(msg, expected) {
				t.Fatalf("messages %v not received. got: %v",msg,expected)
			}
		case <- time.After(time.Second):
			t.Error("timeout")
		}

}
