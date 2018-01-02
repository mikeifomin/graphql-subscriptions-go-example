package queries

import (
 "encoding/json"
 "reflect"
 "context"
)

type subscribeFn = func(context.Context, json.RawMessage) (<-chan json.RawMessage, error)

type Queries struct {
  o observable
	subscribe  subscribeFn
	ctx        context.Context

	// the first element of chanList for refreshing reflect.Select
	// signal refresh reflect.Select whan new query comes
	// close for shutdown select loop
	sigCreated    chan json.RawMessage
	idList    []string
	list  []reflect.SelectCase
}

type observable interface {
	Next(id string, data json.RawMessage)
	Error(id string, err error)
	Complete(id string)
}

func NewQueries(ctx context.Context,  o observable, subscribe subscribeFn) *Queries {
	q := Queries{  
		ctx:   ctx,
		o:     o,
		subscribe: subscribe,
	}
  q.sigCreated = make(chan json.RawMessage, 1)	
	q.idList = []string{""}
	q.list = []reflect.SelectCase{
		reflect.SelectCase{
			Dir: reflect.SelectRecv,
			Chan: reflect.ValueOf(q.sigCreated),
		},
	}
	return &q
}

func (q *Queries) Create(id string, payload json.RawMessage) {
  dataChan, err := q.subscribe(q.ctx,payload)
	if err != nil {
		q.o.Error(id,err)
		q.o.Complete(id)
    return
	}
	sel := reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(dataChan),
	}   
	q.list = append(q.list, sel)
	q.idList = append(q.idList, id)
}

func (q *Queries) Remove(id string) {

}

func (q *Queries) RemoveAll() {
 close(q.sigCreated)
}


func (q *Queries) RunSelectLoop() {
  for {
	  index,v,more := reflect.Select(q.list)
    if index == 0 { // sigCreated weaked up
			if !more { 
				break 
			}
      // chanList changed
      continue
		}
		// XXX: does we have garantee that chanList didn't change?
		id := q.idList[index]
		if !more {
			q.o.Complete(id)
			q.remove(index)
			continue
		}
    data := v.Interface().(json.RawMessage)
		q.o.Next(id,data)
	}
}

func (q *Queries) add(<-chan json.RawMessage) {
   
}
func (q *Queries) remove(index int){
	q.list = append(q.list[:index],q.list[index+1:]...)
}

