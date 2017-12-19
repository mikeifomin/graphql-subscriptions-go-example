<template lang="pug">
.plot
  .messages
    input(
      @keyup.enter="send" 
      v-model="text" 
      placeholder="type message"
    )
    p(v-for="m in messages") {{m.text}} 
  .pointer(
    v-for="p in sessions" 
    :style="{top:p.y+'px',left:p.x+'px'}"
    )
</template>

<script lang="coffee">
import gql from 'graphql-tag'

export default {
  name: 'Plot'
  data: ->
    text: ""
    messages: []
    sessions: []
    sessionId: null
  props:
    room: default: "default"

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
        variables: -> {@room}
        updateQuery: (prev, {subscriptionData })->
          messages: [
            ...prev.messages
            subscriptionData.data.addedMessage
          ]
    sessions:
      query: gql """
        query sessions {
          sessions {id x y agent}
        }
      """
      subscribeToMore:
        document: gql """
          subscription updatedSession($room: String!) {
            updatedSession(room: $room) { id x y deleted}
          }"""
        variables: -> {@room}
        updateQuery: (prev, {subscriptionData})->
          prev = prev.sessions

          upd = subscriptionData.data.updatedSession
          for i,s in prev
            if s.id == upd.id
              if upd.deleted
                return sessions:[
                  ...prev.slice(0,i)
                  ...prev.slice(i+1)
                ]
              sNew = Object.assign {},s,upd
              return sessions:[
                ...prev.slice(0,i)
                sNew   
                ...prev.slice(i+1)
              ]
          return sessions: [ ...prev,upd]
  created: -> document.addEventListener 'mousemove',@mousemove
  beforeDestroy: -> document.removeEventListener 'mousemove',@mousemove

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

    mousemove: (e)->
      x = e.x
      y = e.y
      if not @sessionId
        @newSession(x,y)
        return
      @updateSession(x,y)

    newSession: (x,y)->
      agent = navigator.userAgent 
      {data} = await @$apollo.mutate 
        mutation: gql """
          mutation newSession($x: Int!, $y: Int!, $agent: String!, $room: String!) {
            newSession(x:$x, y:$y, agent:$agent, room:$room) {
              id x y agent 
            } }
        """
        variables: {x,y, agent,@room}
      @sessionId = data.newSession.id

    updateSession: (x,y)->
      @$apollo.mutate 
        mutation: gql """
          mutation newSession($x: Int!, $y: Int!, $id: ID!) {
            updateSession(x:$x, y:$y, id:$id ) { id x y }
          } """
        variables: {x,y,id:@sessionId}

}
</script>

<style lang="stylus">
.plot
 width 100%
 height 100vh
.pointer
  position absolute
  background red
  width 10px
  height 10px
</style>
