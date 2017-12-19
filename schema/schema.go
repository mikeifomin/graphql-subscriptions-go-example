package schema

var Schema = `
# start_used_for_js_mock_server
schema {
	query: Query
  subscription: Subscription
	mutation: Mutation
}

type Query { 
	messages: [Message!]! 
	#sessions: [Session!]!
	#online: [User!]!
}

type Subscription {
	addedMessage(room: String!): Message!
	#updatedSession(room: String!): Session!
	#sessionUpdated(id:ID!, x:Int!, y:Int!, selected:String!, zoom
	#sessionCreated(room: String!): Session!
	#sessionRemoved(room: String!): ID!
}

type Mutation {
	newMessage(id: ID!, text: String!, room: String!): Message!

	#newSession(room: String!): Session!
	#updateSession(id: ID!, state: InputState): Session!
	#removeSession(id: ID!): Boolean
}

type Message { 
	id: ID! 
	text: String!
	room: String!
}

##input InputState {

##	x: Int 
##	y: Int
##	selected: String
##	zoom: Float
##	offsetX: Int
##	offsetY: Int
##}

##type Session {
##	id: ID!
##	room: String!
##	color: String!
##	user: User!

##	x: Int!
##	y: Int!
##	selected: String!
##	zoom: String!
##	offsetX: Int!
##	offsetY: Int!
##}
##type SessionState {
##	x: Int 
##	y: Int
##	selected: String
##	zoom: Float
##	offsetX: Int
##	offsetY: Int
##}

# end_used_for_js_mock_server
`

