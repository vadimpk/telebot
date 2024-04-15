package telebot

// MiddlewareFunc represents a middleware processing function,
// which get called before the endpoint group or specific handler.
type MiddlewareFunc func(HandlerFunc) HandlerFunc

func appendMiddleware(a, b []MiddlewareFunc) []MiddlewareFunc {
	if len(a) == 0 {
		return b
	}

	m := make([]MiddlewareFunc, 0, len(a)+len(b))
	return append(m, append(a, b...)...)
}

func applyMiddleware(h HandlerFunc, m ...MiddlewareFunc) HandlerFunc {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}
