package exercises

type SupportFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type Validator struct {
	Type           string `json:"type"`
	TestCode       string `json:"test_code"`
	TimeoutSeconds int    `json:"timeout_seconds,omitempty"`
}

type Exercise struct {
	ID               string        `json:"id"`
	Language         string        `json:"language"`
	Category         string        `json:"category"`
	Title            string        `json:"title"`
	Topic            string        `json:"topic"`
	Objective        string        `json:"objective"`
	Description      string        `json:"description"`
	Requirements     []string      `json:"requirements"`
	Example          string        `json:"example,omitempty"`
	ExampleOutput    string        `json:"example_output"`
	Hint             string        `json:"hint,omitempty"`
	StarterCode      string        `json:"starter_code"`
	EstimatedMinutes int           `json:"estimated_minutes"`
	Difficulty       int           `json:"difficulty"`
	TimeoutSeconds   int           `json:"timeout_seconds,omitempty"`
	SupportFiles     []SupportFile `json:"support_files,omitempty"`
	Validator        *Validator    `json:"validator,omitempty"`
}

type Track struct {
	ID               string     `json:"id"`
	Language         string     `json:"language"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	EstimatedMinutes int        `json:"estimated_minutes"`
	Exercises        []Exercise `json:"exercises"`
}

// TranslationCatalog mantém texto localizado separado do conteúdo executável.
// Starter code, arquivos de apoio e validadores continuam definidos uma única vez.
type TranslationCatalog struct {
	TrackID     string                         `json:"track_id"`
	Locale      string                         `json:"locale"`
	Title       string                         `json:"title"`
	Description string                         `json:"description"`
	Exercises   map[string]ExerciseTranslation `json:"exercises"`
}

type ExerciseTranslation struct {
	Category      string   `json:"category"`
	Title         string   `json:"title"`
	Topic         string   `json:"topic"`
	Objective     string   `json:"objective"`
	Description   string   `json:"description"`
	Requirements  []string `json:"requirements"`
	Example       string   `json:"example,omitempty"`
	ExampleOutput string   `json:"example_output"`
	Hint          string   `json:"hint,omitempty"`
}
