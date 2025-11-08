package tools

// Params is a struct that contains the parameters to connect to the RabbitMQ
type Params struct {
	Host     string
	Port     string
	User     string
	Password string
	Vhost    string
}
