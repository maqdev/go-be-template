package config

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/kkyr/fig"
	"github.com/stretchr/testify/require"
)

func LoadConfig(envPrefix, configPath string) (*AppConfig, error) {
	var cfg AppConfig
	err := fig.Load(&cfg, fig.File(configPath), fig.UseEnv(envPrefix))
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

//nolint: gochecknoglobals // test helper and needs to be stored in a global variable
var sourcesRoot = findSourcesRoot()

func findSourcesRoot() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Can't get sources root path")
	}
	return filepath.Join(filepath.Dir(filename), "..")
}

func GetSourcesRootPath() string {
	return sourcesRoot
}

func TestConfig(t testing.TB) *AppConfig {
	cfg, err := LoadConfig(filepath.Join(GetSourcesRootPath(), "test.yaml"), "TEST")
	require.NoError(t, err)
	return cfg
}
