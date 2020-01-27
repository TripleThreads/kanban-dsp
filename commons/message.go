package commons

type Message struct {
	RequestType string
	Port        string
	Operations  []Operation
}
