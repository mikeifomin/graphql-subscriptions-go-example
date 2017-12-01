module.exports = 
  Query:
    messages: (_,__,{db})->
      console.log db
      db.messages

  Subscription:
    addedMessage:
      resolve: (payload, args, context, info)->
        return payload
      subscribe: (_,{room},{db,pubsub})->
        pubsub.asyncIterator("addedMessage_#{room}")
      
  Mutation:
    newMessage: (_,{id,text,room},{db,pubsub})->
      msg = {id, text, room}
      db.messages.push(msg)
      console.log db
      pubsub.publish("addedMessage_#{room}", msg)
      msg

  #CurrentUser:
    #orgs: (root,_,{db})-> 
      #db("orgs")
      

