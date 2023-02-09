package mock

import "errors"

type TelegramFailed struct{}

func (tf *TelegramFailed) Open(_ string) error {
	return errors.New("error")
}

func (tf *TelegramFailed) Worker() {

}
