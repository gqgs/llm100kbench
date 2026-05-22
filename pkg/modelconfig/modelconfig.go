package modelconfig

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Models []Model `json:"models"`
}

type Model struct {
	Alias    string `json:"alias"`
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Enabled  bool   `json:"enabled"`
	Archived bool   `json:"archived,omitempty"`
	Env      string `json:"env,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c Config) EnabledModels() []Model {
	models := make([]Model, 0, len(c.Models))
	for _, model := range c.Models {
		if model.Enabled && !model.Archived {
			models = append(models, model)
		}
	}
	return models
}

func (c Config) Validate() error {
	aliases := map[string]struct{}{}
	for _, model := range c.Models {
		if model.Alias == "" {
			return fmt.Errorf("model alias is required")
		}
		if model.Provider == "" {
			return fmt.Errorf("provider is required for %s", model.Alias)
		}
		if model.Model == "" {
			return fmt.Errorf("model id is required for %s", model.Alias)
		}
		if model.Enabled && model.Env == "" {
			return fmt.Errorf("env is required for enabled model %s", model.Alias)
		}
		if _, ok := aliases[model.Alias]; ok {
			return fmt.Errorf("duplicate model alias %s", model.Alias)
		}
		aliases[model.Alias] = struct{}{}
	}
	return nil
}
