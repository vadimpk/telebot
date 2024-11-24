package telebot

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Handler manages handlers that can be shared between multiple bots
type Handler struct {
	// shared options
	URL         string
	synchronous bool
	verbose     bool
	offline     bool
	parseMode   ParseMode
	onError     func(error, Context)
	client      *http.Client

	updates chan Update

	stop chan struct{}

	handlers   map[string]HandlerFunc
	middleware []MiddlewareFunc
}

// NewHandler creates a new handler manager
func NewHandler(pref HandlerSettings) *Handler {
	client := pref.Client
	if client == nil {
		client = &http.Client{Timeout: time.Minute}
	}

	if pref.URL == "" {
		pref.URL = DefaultApiURL
	}
	if pref.OnError == nil {
		pref.OnError = defaultOnError
	}

	hm := &Handler{
		URL:         pref.URL,
		onError:     pref.OnError,
		synchronous: pref.Synchronous,
		verbose:     pref.Verbose,
		offline:     pref.Offline,
		parseMode:   pref.ParseMode,
		client:      client,

		handlers:   make(map[string]HandlerFunc),
		middleware: make([]MiddlewareFunc, 0),
	}

	return hm
}

// HandlerSettings represents a utility struct for passing certain
// properties of a handler around and is required to make a handler manager.
type HandlerSettings struct {
	URL string

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

	// HTTP Client used to make requests to telegram api
	Client *http.Client

	// Offline allows to create a bot without network for testing purposes.
	Offline bool
}

var defaultOnError = func(err error, c Context) {
	if c != nil {
		log.Println(c.Update().ID, err)
	} else {
		log.Println(err)
	}
}

// Handle lets you set the handler for some command name or
// one of the supported endpoints. It also applies middleware
// if such passed to the function.
//
// Example:
//
//	hm.Handle("/start", func (c tele.Context) error {
//		return c.Reply("Hello!")
//	})
//
//	hm.Handle(&inlineButton, func (c tele.Context) error {
//		return c.Respond(&tele.CallbackResponse{Text: "Hello!"})
//	})
//
// Middleware usage:
//
//	hm.Handle("/ban", onBan, middleware.Whitelist(ids...))
func (hm *Handler) Handle(endpoint interface{}, h HandlerFunc, m ...MiddlewareFunc) {
	end := extractEndpoint(endpoint)
	if end == "" {
		panic("telebot: unsupported endpoint")
	}

	if len(hm.middleware) > 0 {
		m = appendMiddleware(hm.middleware, m)
	}

	hm.handlers[end] = func(c Context) error {
		return applyMiddleware(h, m...)(c)
	}
}

// Trigger executes the registered handler by the endpoint.
func (hm *Handler) Trigger(endpoint interface{}, c Context) error {
	end := extractEndpoint(endpoint)
	if end == "" {
		return fmt.Errorf("telebot: unsupported endpoint")
	}

	handler, ok := hm.handlers[end]
	if !ok {
		return fmt.Errorf("telebot: no handler found for given endpoint")
	}

	return handler(c)
}

// Use adds middleware to the global chain
func (hm *Handler) Use(middleware ...MiddlewareFunc) {
	hm.middleware = append(hm.middleware, middleware...)
}

// Group creates a new group of handlers with shared middleware
func (hm *Handler) Group() *Group {
	return &Group{
		hm:         hm,
		middleware: make([]MiddlewareFunc, 0),
	}
}
