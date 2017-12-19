{ PubSub } = require 'graphql-subscriptions'
{SubscriptionServer} = require 'subscriptions-transport-ws'
{ execute, subscribe } = require 'graphql'

module.exports = (server)->
  schema = require('./makeSchema.coffee')()

  pubsub = new PubSub()
  db = 
    messages: [
      (id:'1', text: "yay", room: "default")
    ]
    online: []
    users: []
    sessions: []

  new SubscriptionServer {
    execute
    subscribe
    schema
    onConnect: ({token}, webSocket)->
      return {
        currentUser:{id:"1",name:"mike"}
        token
      }
    onDisconnect: (ws)->
      sessions = ws._sessions or []
      
      for sIdDel in sessions
        for i,s in db.sessions 
          if s.id == sIdDel
            db.sessions[i].deleted = true
            pubsub.publish("updatedSession_#{s.room}",s)

    onOperation: (message, params,ws)-> 
      
      { ...params, context:{...params.context, db,pubsub,ws}}
    
  }, {
    server
    path: '/ws'
  }
