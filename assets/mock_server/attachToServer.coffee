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

  new SubscriptionServer {
    execute
    subscribe
    schema
    onConnect: ({token}, webSocket)->
      return {
        currentUser:{id:"1",name:"mike"}
        token
      }

    onOperation: (message, params, webSocket)-> 
      { ...params, context:{...params.context, db,pubsub}}
    
  }, {
    server
    path: '/ws'
  }
