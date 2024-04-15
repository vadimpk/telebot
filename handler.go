package telebot

// Handler is a struct that holds all the information about a handler.
type Handler struct {
	synchronous bool
	verbose     bool
	parseMode   ParseMode

	onError func(error, Context)

	// handlers is a map of all the handlers.
	handlers map[string]HandlerFunc
	// middleware is a main chain of middleware functions.
	middleware []MiddlewareFunc
}

func NewHandler(settings HandlerSettings) *Handler {
	return &Handler{
		synchronous: settings.Synchronous,
		verbose:     settings.Verbose,
		parseMode:   settings.ParseMode,
		onError:     settings.OnError,

		handlers: make(map[string]HandlerFunc),
	}
}

// HandlerSettings represents a utility struct for passing certain
// properties of a handler around and is required to make handlers.
type HandlerSettings struct {
	// Synchronous prevents handlers from running in parallel.
	// It makes ProcessUpdate return after the handler is finished.
	Synchronous bool

	// Verbose forces bot to log all upcoming requests.
	// Use for debugging purposes only.
	Verbose bool

	// ParseMode used to set default parse mode of all sent messages.
	// It attaches to every send, edit or whatever method. You also
	// will be able to override the default mode by passing a new one.
	ParseMode ParseMode

	// OnError is a callback function that will get called on errors
	// resulted from the handler. It is used as post-middleware function.
	// Notice that context can be nil.
	OnError func(error, Context)
}

// Handle lets you set the handler for some command name or
// one of the supported endpoints. It also applies middleware
// if such passed to the function.
//
// Example:
//
//	b.Handle("/start", func (c tele.Context) error {
//		return c.Reply("Hello!")
//	})
//
//	b.Handle(&inlineButton, func (c tele.Context) error {
//		return c.Respond(&tele.CallbackResponse{Text: "Hello!"})
//	})
//
// Middleware usage:
//
//	b.Handle("/ban", onBan, middleware.Whitelist(ids...))
func (h *Handler) Handle(endpoint interface{}, hf HandlerFunc, m ...MiddlewareFunc) {
	// append main middleware to provided middleware.
	if len(h.middleware) > 0 {
		m = appendMiddleware(h.middleware, m)
	}

	handler := func(c Context) error {
		return applyMiddleware(hf, m...)(c)
	}

	switch end := endpoint.(type) {
	case string:
		h.handlers[end] = handler
	case CallbackEndpoint:
		h.handlers[end.CallbackUnique()] = handler
	default:
		panic("telebot: unsupported endpoint")
	}
}

// Use adds middleware to the chain.
func (h *Handler) Use(middleware ...MiddlewareFunc) {
	h.middleware = append(h.middleware, middleware...)
}

// Group returns a new group.
func (h *Handler) Group() *Group {
	return &Group{h: h}
}

// Group represents a group of handlers. Can be used to apply specific middleware
// to a group of handlers.
type Group struct {
	h          *Handler
	middleware []MiddlewareFunc
}

// Use adds middleware to the group chain.
func (g *Group) Use(middleware ...MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
}

// Handle adds endpoint handler to the Handler, combining group's middleware
// with the optional given middleware.
func (g *Group) Handle(endpoint interface{}, hf HandlerFunc, m ...MiddlewareFunc) {
	g.h.Handle(endpoint, hf, appendMiddleware(g.middleware, m)...)
}
