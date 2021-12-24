package headers

// New constructs a Header for use.
func New() Header {
	return make(Header)
}

// Header defines an abstract definition of a header.
type Header map[string][]string

// SetAll sets the values for the provides key.
func (h Header) SetAll(key string, values []string) {
	h[key] = values
}

// Set sets a single value for the provided key.
func (h Header) Set(key, value string) {
	h.SetAll(key, []string{value})
}

// GetAll returns all possible values for a key.
func (h Header) GetAll(key string) []string {
	return h[key]
}

// Get returns the first possible header value for a key (if present).
func (h Header) Get(key string) string {
	all := h.GetAll(key)
	if len(all) > 0 {
		return all[0]
	}
	return ""
}
