package ws

import (
	graphql "github.com/neelance/graphql-go"
	queriesSelect "github.com/neelance/graphql-go/ws/queries"
	"github.com/gorilla/websocket"
	//"github.com/davecgh/go-spew/spew"

	"net/http"
	"encoding/json"
  "context"
	"log"
)

type Message struct {
  Type   string          `json:"type"`
	Id     string          `json:"id,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type Handler struct {
  Schema    *graphql.Schema
  OnConnect func(json.RawMessage, *http.Request) (context.Context,error)
  Upgrader  websocket.Upgrader
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	headers := http.Header{}
	headers.Set("sec-websocket-protocol","graphql-ws")
	conn, err := h.Upgrader.Upgrade(w, request, headers)
	if err != nil {
		log.Printf("cant connect websockets: %v",err)
		return
	}
	defer conn.Close()
  var ctx context.Context
  
	//////////////////
	// auth: run OnConnect, can update ctx
  for { 
    var msg Message
    err := conn.ReadJSON(&msg)
    if err != nil {
			 return // disconnect
		}
    if msg.Type == GQL_CONNECTION_INIT {
			var errAuth error
			ctx, errAuth = h.OnConnect(msg.Payload,request)
      if errAuth != nil {
         errW := conn.WriteJSON(Message{
            Type: GQL_CONNECTION_ERROR,
            Payload: json.RawMessage(err.Error()),
				 })
         if errW != nil {
           panic(errW)
				 }
         // XXX: Is the message send synchronously?
         return // disconnect
			}
      errW := conn.WriteJSON(Message{ Type: GQL_CONNECTION_ACK })
			if errW != nil {
				 panic(errW)
			}
      break // OnConnect ok
		}
	}    

	subscribeFn := func(ctx context.Context, data json.RawMessage) (<-chan json.RawMessage, error) {
		var params struct {
			Query         string                 `json:"query"`
			OperationName string                 `json:"operationName"`
			Variables     map[string]interface{} `json:"variables"`
		}
		err := json.Unmarshal(data, &params)
		if err != nil {
			return nil, err
		}
    return h.Schema.Subscribe(ctx, params.Query, params.OperationName, params.Variables)
	}

	o := ProtoObserver{conn}
  queries := queriesSelect.NewQueries(ctx, &o, subscribeFn) 

	//  conn.WriteMessage allowed in one gorutine
	go queries.RunSelectLoop()
  defer queries.RemoveAll()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				// TODO: log staus
				/*
				CloseNormalClosure           = 1000
				CloseGoingAway               = 1001
				CloseProtocolError           = 1002
				CloseUnsupportedData         = 1003
				CloseNoStatusReceived        = 1005
				CloseAbnormalClosure         = 1006
				CloseInvalidFramePayloadData = 1007
				ClosePolicyViolation         = 1008
				CloseMessageTooBig           = 1009
				CloseMandatoryExtension      = 1010
				CloseInternalServerErr       = 1011
				CloseServiceRestart          = 1012
				CloseTryAgainLater           = 1013
				CloseTLSHandshake            = 1015
				*/
				return
			}
			panic(err)
		}
		switch msg.Type {
		case GQL_START:
			queries.Create(msg.Id, msg.Payload)
		case GQL_STOP:
			queries.Remove(msg.Id)
		case GQL_CONNECTION_TERMINATE:
			return
		}
	}
}
