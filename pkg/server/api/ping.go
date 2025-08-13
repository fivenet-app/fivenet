package api

// PingResponse is the "pre-created" response returned for ping requests.
var PingResponse = Ping{
	Message: "Pong",
}

// Ping represents the structure of a ping response.
type Ping struct {
	Message string `json:"message"`
}
