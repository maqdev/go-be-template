package config

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/kkyr/fig"
	"github.com/stretchr/testify/require"
)

func LoadConfig(envPrefix, configPath string) (*AppConfig, error) {
	var cfg AppConfig
	dir := fig.Dirs(".")
	if path.IsAbs(configPath) {
		dir = fig.Dirs(path.Dir(configPath))
		configPath = path.Base(configPath)
	}
	err := fig.Load(&cfg, dir, fig.File(configPath), fig.UseEnv(envPrefix))
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// nolint: gochecknoglobals // test helper and needs to be stored in a global variable
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
	cfg, err := LoadConfig("TEST", filepath.Join(GetSourcesRootPath(), "infra/config.yaml"))
	require.NoError(t, err)
	return cfg
}
