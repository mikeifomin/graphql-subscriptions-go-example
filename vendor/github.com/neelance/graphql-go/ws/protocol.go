package ws

// https://github.com/apollographql/subscriptions-transport-ws/blob/a56491c6feacd96cab47b7a3df8c2cb1b6a96e36/src/message-types.ts
const (
	GQL_CONNECTION_INIT = "connection_init" // Client -> Server
	GQL_CONNECTION_ACK = "connection_ack" // Server -> Client
	GQL_CONNECTION_ERROR = "connection_error" // Server -> Client
	// NOTE: The keep alive message type does not follow the standard due to connection optimizations
	GQL_CONNECTION_KEEP_ALIVE = "ka" // Server -> Client
	GQL_CONNECTION_TERMINATE = "connection_terminate" // Client -> Server
	GQL_START = "start" // Client -> Server
	GQL_DATA = "data" // Server -> Client
	GQL_ERROR = "error" // Server -> Client
	GQL_COMPLETE = "complete" // Server -> Client
	GQL_STOP = "stop" // Client -> Server
)
