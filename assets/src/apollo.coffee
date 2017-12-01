{ ApolloClient } = require 'apollo-client'
{ InMemoryCache } = require 'apollo-cache-inmemory'
{ WebSocketLink } = require "apollo-link-ws";
{ SubscriptionClient } = require 'subscriptions-transport-ws'

ws = new SubscriptionClient "ws://#{window.location.host}/ws",
  lazy: false
  reconnect: true

  connectionParams: ->
    token: Math.random().toString(30)
  connectionCallback: (err)->
    if err
      # close connection?
      console.error err

window.WS = ws

link = new WebSocketLink(ws);
cache = new InMemoryCache()

apolloClient = new ApolloClient {
  link
  cache
  connectToDevTools: true
}

module.exports = {
  apolloClient
}
