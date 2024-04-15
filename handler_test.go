package telebot

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerHandle(t *testing.T) {
	h := NewHandler(HandlerSettings{})

	h.Handle("/start", func(c Context) error { return nil })
	assert.Contains(t, h.handlers, "/start")

	reply := ReplyButton{Text: "reply"}
	h.Handle(&reply, func(c Context) error { return nil })

	inline := InlineButton{Unique: "inline"}
	h.Handle(&inline, func(c Context) error { return nil })

	btnReply := (&ReplyMarkup{}).Text("btnReply")
	h.Handle(&btnReply, func(c Context) error { return nil })

	btnInline := (&ReplyMarkup{}).Data("", "btnInline")
	h.Handle(&btnInline, func(c Context) error { return nil })

	assert.Contains(t, h.handlers, btnReply.CallbackUnique())
	assert.Contains(t, h.handlers, btnInline.CallbackUnique())
	assert.Contains(t, h.handlers, reply.CallbackUnique())
	assert.Contains(t, h.handlers, inline.CallbackUnique())
}

func TestHandlerOnError(t *testing.T) {
	b, err := NewBot(Settings{
		Handler: NewHandler(HandlerSettings{
			Synchronous: true,
		}),
		Offline: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	var ok bool
	b.handler.onError = func(err error, c Context) {
		assert.Equal(t, b, c.(*nativeContext).b)
		assert.NotNil(t, err)
		ok = true
	}

	b.runHandler(func(c Context) error {
		return errors.New("not nil")
	}, &nativeContext{b: b})

	assert.True(t, ok)
}
