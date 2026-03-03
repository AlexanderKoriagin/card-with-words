package base

const (
	DefaultQty       = 8
	DefaultCacheSize = 256

	SeparatorComma   = ", "
	SeparatorNextRow = "\n"
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
