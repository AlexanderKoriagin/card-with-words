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
                 4. Adhere to the constraints on minimum word length and age category.
                 5. Use only single words. Do not provide phrases or compound words consisting of multiple separate words.
                 6. If you cannot generate enough unique words due to the STOP-LIST, generate the closest possible synonyms that are NOT in the list.`

	childPrompt = `Language: %s. 
                   Difficulty: %s.
                   Requirements: Concrete, high-imagery nouns only. Focus on toys, common animals, household items, food, and characters found in children's books or cartoons. Use words that a child can describe by their appearance, sound, or function. Strictly avoid abstract concepts, metaphors, and complex emotions.
                   Minimum word length: 3 letters.
                   
                   ### FORBIDDEN WORDS (STOP-LIST):
                   [%s] 
                   
                   ### TASK: 
                   1. CRITICAL: Ensure NONE of the generated words match any words in the STOP-LIST above.
                   2. Check for semantic and root-word duplicates in the STOP-LIST.
                   3. Make the list diverse and unpredictable.
                   4. Generate %d words in JSON.`

	teenPrompt = `Language: %s. 
                  Difficulty: %s.
                  Requirements: Intermediate-level nouns. Focus on themes: education, technology, emotions, and urban life. Requirements: words that a teenager uses or hears daily but are not part of a 'basic objects' list. Avoid words with obvious 3nd-grade definitions. Exclude highly technical or professional jargon.
                  Minimum word length: 4 letters.

                  ### FORBIDDEN WORDS (STOP-LIST):
                  [%s] 

                  ### TASK: 
                  1. CRITICAL: Ensure NONE of the generated words match any words in the STOP-LIST above.
                  2. Check for semantic and root-word duplicates in the STOP-LIST.
                  3. Make the list diverse and unpredictable.
                  4. Generate %d words in JSON.`

	adultPrompt = `Language: %s.
				   Difficulty: %s.
				   Requirements: A diverse mix of intermediate and advanced nouns. Include broad abstract concepts, professional fields, modern social phenomena, and cultural idioms. Focus on words that are well-known but require a creative explanation. Avoid basic everyday objects and highly technical jargon.
                   Minimum word length: 5 letters.

                   ### FORBIDDEN WORDS (STOP-LIST):
                   [%s] 
                   
                   ### TASK: 
                   1. CRITICAL: Ensure NONE of the generated words match any words in the STOP-LIST above.
                   2. Check for semantic and root-word duplicates in the STOP-LIST.
                   3. Make the list diverse and unpredictable.
                   4. Generate %d words in JSON.`
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
