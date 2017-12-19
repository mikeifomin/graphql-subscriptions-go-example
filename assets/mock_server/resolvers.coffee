module.exports = 
  Query:
    messages: (_,__,{db})->
      db.messages
    #sessions: (_,__,{db})->
      #db.sessions.filter (s)->
        #!s.deleted

  Subscription:
    addedMessage:
      resolve: (payload, args, context, info)->
        return payload
      subscribe: (_,{room},{db,pubsub})->
        pubsub.asyncIterator("addedMessage_#{room}")
    #updatedSession: 
      #subscribe: (_,{room},{db,pubsub})->
        #pubsub.asyncIterator("updatedSession_#{room}")
      #resolve: (payload, args, context, info)->
        #return payload
      
  Mutation:
    newMessage: (_,{id,text,room},{db,pubsub})->
      msg = {id, text, room}
      db.messages.push(msg)
      pubsub.publish("addedMessage_#{room}", msg)
      msg
    #newSession: (_,{x,y,agent,room},{ws,db,pubsub})->
      #id = Math.random().toString(36)
      #doc = {id,x,y,agent,room}
      #db.sessions.push(doc)
      #pubsub.publish("updatedSession_#{room}",doc)
      #ws._sessions ?= []
      #ws._sessions.push(id)
      #doc

    #updateSession: (_,{id,x,y},{db,pubsub})->
      #for s in db.sessions
        #if s.id == id
          #s.x = x
          #s.y = y
          #pubsub.publish("updatedSession_#{s.room}",s)
          #return s
      #console.error "not found"
  #Session:
    #deleted: (obj)->
      #!!obj.deleted

  #CurrentUser:
    #orgs: (root,_,{db})-> 
      #db("orgs")
      

