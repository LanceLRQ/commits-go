package config

type Config struct {
	Locale string          `yaml:"locale"`
	Proxy  ProxyConfig     `yaml:"proxy"`
	LLM    LLMConfig       `yaml:"llm"`
	Server OpenAIConfig    `yaml:"server"`
	Git    GitCommitConfig `yaml:"git"`
}

type ProxyConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Url      string `yaml:"url"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

type OpenAIConfig struct {
	ApiUrl  string `yaml:"api_url"`
	ApiKey  string `yaml:"api_key"`
	Model   string `yaml:"model"`
	Timeout int    `yaml:"timeout"`
}

type LLMConfig struct {
	Temperature      float64 `yaml:"temperature"`
	TopP             float64 `yaml:"top_p"`
	MaxTokens        int     `yaml:"max_tokens"`
	FrequencyPenalty float64 `yaml:"frequency_penalty"`
	PresencePenalty  float64 `yaml:"presence_penalty"`
}

type GitCommitConfig struct {
	MaxLength     int    `yaml:"max_length"`
	SuggestLength int    `yaml:"suggest_length"`
	CommitStyle   string `yaml:"commit_style"`
}
