package pubsub

import (
	"testing"
	"context"
	"time"
)
type testType struct {
	id  string
}

func TestSimple(t *testing.T){
	p := New()
	ctx, stopSub := context.WithCancel(context.Background())
	defer stopSub()

	ch := p.Sub(ctx,"foo")
	payload := testType{"id"}
	p.Pub("foo",payload)

	select {
	case data, has := <- ch:
		if !has {
			t.Error("channel closed unexpected")
		  return
		}
		actual, ok := data.(testType)
		if !ok {
     t.Error("subscription receive wrong type data")     
		 return
		}
		if actual != payload {
			t.Errorf("received not expected data: %v",actual)
		}
	case <- time.After(time.Millisecond*200):
		t.Error("timeout channel read")
	}
}

