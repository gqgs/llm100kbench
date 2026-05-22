package modelconfig

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "models.json")
	require.NoError(t, os.WriteFile(path, []byte(`{
		"models": [
			{"alias":"chatgpt","provider":"github","model":"openai/gpt-4.1","enabled":true,"env":"GITHUB_TOKEN"},
			{"alias":"perplexity","provider":"perplexity","model":"sonar","enabled":false,"archived":true}
		]
	}`), 0o644))

	cfg, err := Load(path)
	require.NoError(t, err)
	require.Len(t, cfg.Models, 2)
	require.Len(t, cfg.EnabledModels(), 1)
	require.Equal(t, "chatgpt", cfg.EnabledModels()[0].Alias)
}

func TestValidateRequiresEnvForEnabledModel(t *testing.T) {
	cfg := Config{
		Models: []Model{
			{Alias: "chatgpt", Provider: "github", Model: "openai/gpt-4.1", Enabled: true},
		},
	}

	require.Error(t, cfg.Validate())
}
