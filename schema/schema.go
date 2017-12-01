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
	#online: [User!]!
}

type Subscription {
	addedMessage(room: String!): Message!
	#joinUser: User!
	#leaveUser: User!
}

type Mutation {
	newMessage(id: ID!, text: String!, room: String!): Message!
}

type Message { 
	id: ID! 
	text: String!
	room: String!
}
# end_used_for_js_mock_server
`

