package main

import (
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
)

type configFile string

func (c configFile) path(s string) (string, error) {
	configDir := configdir.LocalConfig("gitemoji")
	err := os.MkdirAll(configDir, 0o755)
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, s), nil
}

func (c configFile) Write(s []byte) error {
	p, err := c.path(string(c))
	if err != nil {
		return err
	}
	return os.WriteFile(p, s, 0o644)
}

func (c configFile) Read() ([]byte, error) {
	p, err := c.path(string(c))
	if err != nil {
		return nil, err
	}
	return os.ReadFile(p)
}
