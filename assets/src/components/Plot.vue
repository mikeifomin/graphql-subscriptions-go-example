<template lang="pug">
.plot
  
  input(
    @keyup.enter="send" 
    v-model="text" 
    placeholder="type message"
  )
  p(v-for="m in messages") {{m.text}} 
  
</template>

<script lang="coffee">
import gql from 'graphql-tag'

export default {
  name: 'Plot'
  data: ->
    text: ""
    messages: []
  apollo:
    messages:
      query: gql """
        query messages {
          messages {id text room}
        }
      """
      subscribeToMore:
        document: gql """
          subscription addedMessage($room: String!) {
            addedMessage(room: $room) { id text room }
          }"""
        variables: room: "default"
        updateQuery: (prev, {subscriptionData })->
          messages: [
            ...prev.messages
            subscriptionData.data.addedMessage
          ]

  methods:
    send: ->
      await @$apollo.mutate
        mutation: gql """
          mutation newMessage($id: ID!, $text: String!, $room: String!) {
            newMessage(id:$id, text:$text, room:$room) {id}
          }"""
        variables:
          id: Math.random().toString(32)
          text: @text
          room: 'default'
      @text = ""

}
</script>

<style lang="stylus">
</style>
