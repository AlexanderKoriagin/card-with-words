package groq

import (
	"encoding/json"
	"fmt"
	"strings"

	"cardWithWords/internal/pkg/base"
	"cardWithWords/internal/pkg/groq/dto"

	jgroq "github.com/jpoz/groq"
)

const (
	baseModel = "llama-3.1-8b-instant"

	sysPrompt = `You are a word generator for the Alias game. Your task is to provide exactly 8 words in JSON format.
                 Rules:
                 1. The answer must be strictly in JSON format: {"result": ["word1", ..., "word8"]}. Strictly follow this format and do not add any additional words or characters.
                 2. Use only nominative case, singular form (where applicable).
                 3. Do not use rare proper names or highly specialized terms unless specifically instructed.
                 4. Adhere to the constraints on minimum word length and age category.`

	childPrompt = `Language: %s. 
                   Difficulty: %s.
                   Requirements: Very simple everyday objects and animals that a 6-year-old can easily explain.
                   Minimum word length: 3 letters.
                   Make the list diverse and unpredictable.
                   Do not use these words: %s.
                   Generate %d words in JSON.`

	teenPrompt = `Language: %s. 
                  Difficulty: %s.
                  Requirements: Common nouns, abstract concepts known to students, objects from modern life. Avoid overly primitive or overly archaic words.
                  Minimum word length: 4 letters.
                  Make the list diverse and unpredictable.
                  Do not use these words: %s.
                  Generate %d words in JSON.`

	adultPrompt = `Language: %s.
				   Difficulty: %s.
				   Requirements: A mix of common and less common nouns, including some abstract concepts, cultural references, and modern terms. Avoid overly simple or overly obscure words.
				   Minimum word length: 5 letters.
				   Make the list diverse and unpredictable.
				   Do not use these words: %s.
				   Generate %d words in JSON.`
)

type words struct {
	client *jgroq.Client
	cache  *cache
}

func New(token string) (Words, error) {
	client := jgroq.NewClient(jgroq.WithAPIKey(token))

	c, err := newCache(base.DefaultCacheSize)
	if err != nil {
		return nil, err
	}

	return &words{client: client, cache: c}, nil
}

func (w *words) Card8Words(language base.Language, difficulty base.Difficulty) (*string, error) {
	var (
		userPrompt string
		card       string
	)

	switch difficulty {
	case base.Child:
		userPrompt = fmt.Sprintf(childPrompt, string(language), string(difficulty), w.cache.get(language, difficulty), base.DefaultQty)
	case base.Teen:
		userPrompt = fmt.Sprintf(teenPrompt, string(language), string(difficulty), w.cache.get(language, difficulty), base.DefaultQty)
	case base.Adult:
		userPrompt = fmt.Sprintf(adultPrompt, string(language), string(difficulty), w.cache.get(language, difficulty), base.DefaultQty)
	}

	params := jgroq.CompletionCreateParams{
		Model: baseModel,
		Messages: []jgroq.Message{
			{
				Role:    "system",
				Content: sysPrompt,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
		ResponseFormat: jgroq.ResponseFormat{Type: "json_object"},
	}

	groqResult, err := w.client.CreateChatCompletion(params)
	if err != nil {
		return nil, fmt.Errorf("could not create completion from groq: %w", err)
	}

	if groqResult == nil || len(groqResult.Choices) == 0 {
		return nil, fmt.Errorf("groq returned no choices")
	}

	var parsedResult dto.Card8WordsResult
	err = json.Unmarshal([]byte(groqResult.Choices[0].Message.Content), &parsedResult)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal choices - %s - from groq: %w", groqResult.Choices[0].Message.Content, err)
	}

	w.cache.add(language, difficulty, parsedResult.Result)
	for i := range parsedResult.Result {
		card += strings.ToUpper(parsedResult.Result[i])
		if i != len(parsedResult.Result)-1 {
			card += base.SeparatorNextRow
		}
	}

	return &card, nil
}
