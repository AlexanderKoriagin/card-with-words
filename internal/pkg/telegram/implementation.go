package telegram

import (
	"fmt"
	"sync"

	"cardWithWords/internal/pkg/base"
	"cardWithWords/internal/pkg/groq"
	"cardWithWords/internal/services"

	tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	butRus = "Русский"
	butEng = "English"

	butDictRus = "Русский / Словарь"
	butGroqRus = "Русский / ИИ"
	butDictEng = "English / Dictionary"
	butGroqEng = "English / AI"

	butDictWordsRus = "Русский / Словарь / Получить карту"
	butDictWordsEng = "English / Dictionary / Get a card"

	butGroqChildRus = "Русский / ИИ / Для детей"
	butGroqTeenRus  = "Русский / ИИ / Для подростков"
	butGroqAdultRus = "Русский / ИИ / Для взрослых"
	butGroqChildEng = "English / AI / For kids"
	butGroqTeenEng  = "English / AI / For teens"
	butGroqAdultEng = "English / AI / For adults"

	butGroqChildWordsRus = "Русский / ИИ / Для детей / Получить карту"
	butGroqTeenWordsRus  = "Русский / ИИ / Для подростков / Получить карту"
	butGroqAdultWordsRus = "Русский / ИИ / Для взрослых / Получить карту"
	butGroqChildWordsEng = "English / AI / For kids / Get a card"
	butGroqTeenWordsEng  = "English / AI / For teens / Get a card"
	butGroqAdultWordsEng = "English / AI / For adults / Get a card"

	butBackToMainRus   = "⬅️ К выбору языка"
	butBackToMainEng   = "⬅️ Back to languages"
	butBackToModesRus  = "⬅️ К выбору режима"
	butBackToModesEng  = "⬅️ Back to modes"
	butBackToAiDiffRus = "⬅️ К выбору сложности"
	butBackToAiDiffEng = "⬅️ Back to difficulty"

	msgDefault = "нажми на кнопку / press the button"
)

type Channels struct {
	Done   chan struct{}
	Errors chan error
}

type Bot struct {
	Api        *tba.BotAPI
	WordsLocal services.Data
	WordsGroq  groq.Words
	Wg         *sync.WaitGroup
	Channels   Channels
}

type Features interface {
	Open(token string) error
	Worker()
}

func (b *Bot) Open(token string) (err error) {
	b.Api, err = tba.NewBotAPI(token)
	if err != nil {
		return err
	}

	return nil
}

// Worker generate words by request from user
func (b *Bot) Worker() {
	var (
		keyboardMain = tba.NewReplyKeyboard(
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butRus),
				tba.NewKeyboardButton(butEng),
			),
		)

		keyboardRusLayer0 = tba.NewReplyKeyboard(
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butDictRus),
				tba.NewKeyboardButton(butGroqRus),
			),
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butBackToMainRus),
			),
		)

		keyboardEngLayer0 = tba.NewReplyKeyboard(
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butDictEng),
				tba.NewKeyboardButton(butGroqEng),
			),
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butBackToMainEng),
			),
		)

		keyboardRusDict = tba.NewReplyKeyboard(
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butDictWordsRus),
			),
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butBackToModesRus),
			),
		)

		keyboardEngDict = tba.NewReplyKeyboard(
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butDictWordsEng),
			),
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butBackToModesEng),
			),
		)

		keyboardRusAiLayer1 = tba.NewReplyKeyboard(
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butGroqChildRus),
			),
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butGroqTeenRus),
			),
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butGroqAdultRus),
			),
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butBackToModesRus),
			),
		)

		keyboardEngAiLayer1 = tba.NewReplyKeyboard(
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butGroqChildEng),
			),
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butGroqTeenEng),
			),
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butGroqAdultEng),
			),
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butBackToModesEng),
			),
		)

		keyboardRusAiChild = tba.NewReplyKeyboard(
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butGroqChildWordsRus),
			),
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butBackToAiDiffRus),
			),
		)

		keyboardRusAiTeen = tba.NewReplyKeyboard(
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butGroqTeenWordsRus),
			),
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butBackToAiDiffRus),
			),
		)

		keyboardRusAiAdult = tba.NewReplyKeyboard(
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butGroqAdultWordsRus),
			),
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butBackToAiDiffRus),
			),
		)

		keyboardEngAiChild = tba.NewReplyKeyboard(
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butGroqChildWordsEng),
			),
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butBackToAiDiffEng),
			),
		)

		keyboardEngAiTeen = tba.NewReplyKeyboard(
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butGroqTeenWordsEng),
			),
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butBackToAiDiffEng),
			),
		)

		keyboardEngAiAdult = tba.NewReplyKeyboard(
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butGroqAdultWordsEng),
			),
			tba.NewKeyboardButtonRow(
				tba.NewKeyboardButton(butBackToAiDiffEng),
			),
		)
	)

	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values, and we don't
	// need them repeated.
	uc := tba.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	uc.Timeout = 30

	// Start polling Telegram for updates
	updates := b.Api.GetUpdatesChan(uc)

	for {
		select {
		case <-b.Channels.Done:
			b.Wg.Done()
			return
		case u := <-updates:
			var (
				output   = msgDefault
				keyboard tba.ReplyKeyboardMarkup
			)

			if u.Message == nil {
				continue
			} else {
				switch u.Message.Text {
				case butRus, butBackToModesRus:
					keyboard = keyboardRusLayer0
				case butEng, butBackToModesEng:
					keyboard = keyboardEngLayer0
				case butBackToMainRus, butBackToMainEng:
					keyboard = keyboardMain
				case butDictRus:
					keyboard = keyboardRusDict
				case butDictWordsRus:
					output = b.WordsLocal.GetRussian()
					keyboard = keyboardRusDict
				case butDictEng:
					keyboard = keyboardEngDict
				case butDictWordsEng:
					output = b.WordsLocal.GetEnglish()
					keyboard = keyboardEngDict
				case butGroqRus, butBackToAiDiffRus:
					keyboard = keyboardRusAiLayer1
				case butGroqEng, butBackToAiDiffEng:
					keyboard = keyboardEngAiLayer1
				case butGroqChildRus:
					keyboard = keyboardRusAiChild
				case butGroqTeenRus:
					keyboard = keyboardRusAiTeen
				case butGroqAdultRus:
					keyboard = keyboardRusAiAdult
				case butGroqChildEng:
					keyboard = keyboardEngAiChild
				case butGroqTeenEng:
					keyboard = keyboardEngAiTeen
				case butGroqAdultEng:
					keyboard = keyboardEngAiAdult
				case butGroqChildWordsRus:
					groqCard, err := b.WordsGroq.Card8Words(base.Russian, base.Child)
					if err != nil {
						b.Channels.Errors <- fmt.Errorf("[Worker] couldn't get card ru-child from Groq: %v\n", err)
					}

					output = *groqCard
					keyboard = keyboardRusAiChild
				case butGroqTeenWordsRus:
					groqCard, err := b.WordsGroq.Card8Words(base.Russian, base.Teen)
					if err != nil {
						b.Channels.Errors <- fmt.Errorf("[Worker] couldn't get card ru-teen from Groq: %v\n", err)
					}

					output = *groqCard
					keyboard = keyboardRusAiTeen
				case butGroqAdultWordsRus:
					groqCard, err := b.WordsGroq.Card8Words(base.Russian, base.Adult)
					if err != nil {
						b.Channels.Errors <- fmt.Errorf("[Worker] couldn't get card ru-adult from Groq: %v\n", err)
					}

					output = *groqCard
					keyboard = keyboardRusAiAdult
				case butGroqChildWordsEng:
					groqCard, err := b.WordsGroq.Card8Words(base.English, base.Child)
					if err != nil {
						b.Channels.Errors <- fmt.Errorf("[Worker] couldn't get card en-child from Groq: %v\n", err)
					}

					output = *groqCard
					keyboard = keyboardEngAiChild
				case butGroqTeenWordsEng:
					groqCard, err := b.WordsGroq.Card8Words(base.English, base.Teen)
					if err != nil {
						b.Channels.Errors <- fmt.Errorf("[Worker] couldn't get card en-teen from Groq: %v\n", err)
					}

					output = *groqCard
					keyboard = keyboardEngAiTeen
				case butGroqAdultWordsEng:
					groqCard, err := b.WordsGroq.Card8Words(base.English, base.Adult)
					if err != nil {
						b.Channels.Errors <- fmt.Errorf("[Worker] couldn't get card en-adult from Groq: %v\n", err)
					}

					output = *groqCard
					keyboard = keyboardEngAiAdult
				default:
					keyboard = keyboardMain
				}
			}

			msg := tba.NewMessage(u.Message.Chat.ID, output)
			msg.ReplyMarkup = keyboard

			_, err := b.Api.Send(msg)
			if err != nil {
				b.Channels.Errors <- fmt.Errorf("[Worker] couldn't send msg to %s: %v\n", u.Message.From.UserName, err)
			}
		}
	}
}
