fs = require 'fs'
{ merge, property } = require 'lodash'
{ makeExecutableSchema } = require 'graphql-tools'
{ graphqlExpress, graphiqlExpress } = require 'apollo-server-express'
{ PubSub } = require 'graphql-subscriptions'

{ SubscriptionServer } = require 'subscriptions-transport-ws'
{ execute, subscribe } = require 'graphql'

pubsub = new PubSub()
 

fullSchema = ->
  schema = makeExecutableSchema {
    typeDefs: loadSchema()
    resolvers
  }

dynamic = ->
  db = require './db.coffee'
  schema = fullSchema()
  graphqlExpress (req)->
    token = req.get('token')
    return {
      schema
      tracing: true
      context: {token,db,req, pubsub}
    }

serverReady = (server)->
  schema = fullSchema()
  new SubscriptionServer {
    execute
    subscribe
    schema
    onConnect: ({token}, webSocket)->
      new Promise (resolve,reject)->
        try
          user = await userFromToken(token) 
          resolve {
            currentUser:user
            token
          }
        catch e
          reject(e)

    onOperation: (message, params, webSocket)-> 
      db = require './db.coffee'
      { ...params, context:{...params.context, db,pubsub}}
    
  }, {
    server
    path: '/ws'
  }

module.exports = {middleware,serverReady}
