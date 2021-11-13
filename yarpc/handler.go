package yarpc

// Handler defines an interface that can be used for handling requests.
type Handler interface {
	Handle(Stream) error
}

// HandlerFunc provides users with a simple functional interface for a Handler.
type HandlerFunc func(Stream) error

func (fn HandlerFunc) Handle(stream Stream) error {
	if fn == nil {
		return nil
	}

	return fn(stream)
}

// DefaultServer is a global server definition that can be leveraged by hosting program.
var DefaultServer = &Server{}

// Handle adds the provided handler to the default server.
func Handle(pattern string, handler Handler) {
	DefaultServer.Handle(pattern, handler)
}

// HandleFunc adds the provided handler function to the default server.
func HandleFunc(pattern string, handler func(Stream) error) {
	Handle(pattern, HandlerFunc(handler))
}

// ListenAndServe starts the default server on the provided network and address.
func ListenAndServe(network, address string, opts ...Option) error {
	return DefaultServer.ListenAndServe(network, address, opts...)
}
