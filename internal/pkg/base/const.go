package base

const (
	DefaultQty          = 8
	DefaultCacheSize    = 256
	DefaultGroqAttempts = 3
	DefaultGroqReqPause = 100 // in milliseconds

	SeparatorComma   = ", "
	SeparatorNextRow = "\n"

	MsgGroqProblemsEng = "Sorry, I couldn't get the card from Groq. Please, try again later."
	MsgGroqProblemsRus = "Извините, не удалось получить карту от Groq. Пожалуйста, попробуйте позже."
)

type Language string

const (
	Russian Language = "russian"
	English Language = "english"
)

type Difficulty string

const (
	Child Difficulty = "child"
	Teen  Difficulty = "teen"
	Adult Difficulty = "adult"
)
