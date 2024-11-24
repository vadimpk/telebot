package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tele "github.com/vadimpk/telebot"
)

var b, _ = tele.NewBot(tele.Settings{Handler: tele.NewHandler(tele.HandlerSettings{Offline: true})})

func TestRecover(t *testing.T) {
	onError := func(err error, c tele.Context) {
		require.Error(t, err, "recover test")
	}

	h := func(c tele.Context) error {
		panic("recover test")
	}

	assert.Panics(t, func() {
		h(nil)
	})

	assert.NotPanics(t, func() {
		Recover(onError)(h)(nil)
	})
}
