package mock

type TelegramSuccess struct{}

func (ts *TelegramSuccess) Open(_ string) error {
	return nil
}

func (ts *TelegramSuccess) Worker() {

}
