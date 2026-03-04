package groq

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

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
                   4. Generate exactly %d words in JSON format.`

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
                  4. Generate exactly %d words in JSON format.`

	adultPrompt = `Language: %s.
                   Difficulty: %s.
                   Requirements: A balanced mix of nouns for a mature audience. 

                   ### WORD SELECTION RULES:
                   1. BALANCE: Generate %d words total.
                   2. "HARD" LEVEL (Exactly 1 word): Intellectually stimulating or rare concepts (e.g., philosophical terms, specific professional roles, or complex social phenomena like "existentialism" or "stagnation").
                   3. "EASY-ADULT" LEVEL (Exactly 2 words): Common nouns that are simple to explain but belong to the adult world (e.g., "vacation", "insurance", "freelance", "colleague", "gym"). No 3rd-grade objects like "table" or "cat".
                   4. "MEDIUM-ADULT" LEVEL (All remaining words): Sophisticated but common nouns. Distribute them across these categories:
                       - 1. Work & Career (deadline, promotion, burnout, expertise).
                       - 2. Relationships & Psychology (empathy, commitment, boundary, nostalgia).
                       - 3. Finance & Economy (mortgage, investment, inflation, budget).
                       - 4. Lifestyle & Urban Life (infrastructure, sustainability, trendsetter).
                       - 5. Health & Wellbeing (prevention, metabolism, mindfulness, resilience).
                       - 6. Society & Politics (legislation, advocacy, consensus, bureaucracy).
                       - 7. Media & Technology (algorithm, privacy, authenticity, integration).
                       - 8. Travel & Culture (heritage, hospitality, destination, itinerary).
                       - 9. Personal Growth & Education (mentorship, discipline, perspective).
                       - 10. Legal & Formalities (agreement, liability, entitlement, procedure).

                   ### FORMAT & STRUCTURE RULES (STRICT):
                   1. SINGLE WORDS ONLY: Each entry must be exactly ONE word. 
                   2. NO PHRASES: Absolutely no word combinations (e.g., NO "social security", NO "means of prevention").
                   3. NO HYPHENS: Do not use words with dashes or hyphens (e.g., NO "socio-economic").
                   4. NO PLURALS: Use singular form where appropriate.
                   5. EXCLUSIONS: Avoid basic objects (chair, bread) and overly obscure jargon.
                   6. Minimum word length: 5 letters.

                   ### FORBIDDEN WORDS (STOP-LIST):
                   [%s] 

                   ### TASK: 
                   1. CRITICAL: Ensure NONE of the generated words match any words in the STOP-LIST.
                   2. Check for semantic and root-word duplicates in the STOP-LIST.
                   3. Maintain the 1-Hard / 2-Easy / Rest-Medium ratio strictly.
                   4. Ensure all outputs are strictly single-word nouns.
                   5. Generate exactly %d words in JSON format.`
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

	wordsCache := w.cache.get(language, difficulty)
	switch difficulty {
	case base.Child:
		userPrompt = fmt.Sprintf(childPrompt, string(language), string(difficulty), wordsCache, base.DefaultQty)
	case base.Teen:
		userPrompt = fmt.Sprintf(teenPrompt, string(language), string(difficulty), wordsCache, base.DefaultQty)
	case base.Adult:
		userPrompt = fmt.Sprintf(adultPrompt, string(language), string(difficulty), base.DefaultQty, wordsCache, base.DefaultQty)
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

	var parsedResult dto.Card8WordsResult
	for i := 0; i < base.DefaultGroqAttempts; i++ {
		groqResult, err := w.client.CreateChatCompletion(params)
		if err != nil {
			return nil, fmt.Errorf("could not create completion from groq: %w", err)
		}

		if groqResult == nil || len(groqResult.Choices) == 0 {
			return nil, fmt.Errorf("groq returned no choices")
		}

		err = json.Unmarshal([]byte(groqResult.Choices[0].Message.Content), &parsedResult)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal choices - %s - from groq: %w", groqResult.Choices[0].Message.Content, err)
		}

		if len(parsedResult.Result) == base.DefaultQty {
			break
		}

		// if groq didn't return exactly 8 words, try again with the same prompt
		parsedResult = dto.Card8WordsResult{}
		time.Sleep(base.DefaultGroqReqPause * time.Millisecond)
	}

	if len(parsedResult.Result) != base.DefaultQty {
		card = base.MsgGroqProblemsRus
		if language == base.English {
			card = base.MsgGroqProblemsEng
		}

		return &card, nil
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
