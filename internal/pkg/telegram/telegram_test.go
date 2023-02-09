package telegram

import (
	"cardWithWords/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTelegram_PlayCards(t *testing.T) {
	testData := []struct {
		source Telegram
		isErr  bool
	}{
		{Telegram{Features: &mock.TelegramFailed{}}, true},
		{Telegram{Features: &mock.TelegramSuccess{}}, false},
	}

	for _, td := range testData {
		if td.isErr {
			assert.NotNil(t, td.source.PlayCards("token"))
			continue
		}

		assert.Nil(t, td.source.PlayCards("token"))
	}
}
