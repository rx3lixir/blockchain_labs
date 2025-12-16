package fork

import (
	"encoding/json"
	"os"
	"time"
)

type ChainInfo struct {
	File      string  `json:"file"`
	CreatedAt int64   `json:"created_at"`
	ForkFrom  *string `json:"fork_from"`
	ForkPoint *int    `json:"fork_point"`
}

type Config struct {
	Chains map[string]*ChainInfo `json:"chains"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Chains: make(map[string]*ChainInfo)}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) Save(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0o644)
}

func (c *Config) AddChain(name, file string, forkFrom *string, forkPoint *int) {
	c.Chains[name] = &ChainInfo{
		File:      file,
		CreatedAt: time.Now().Unix(),
		ForkFrom:  forkFrom,
		ForkPoint: forkPoint,
	}
}

func (c *Config) GetChain(name string) (*ChainInfo, bool) {
	info, ok := c.Chains[name]
	return info, ok
}
