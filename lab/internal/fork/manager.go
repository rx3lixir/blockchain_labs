package fork

import (
	"fmt"
	"os"

	"github.com/rx3lixir/lab_bc/internal/blockchain"
	"github.com/rx3lixir/lab_bc/internal/storage"
)

type Manager struct {
	Config     *Config
	configFile string
}

func NewManager(configFile string) (*Manager, error) {
	cfg, err := LoadConfig(configFile)
	if err != nil {
		return nil, err
	}
	return &Manager{
		Config:     cfg,
		configFile: configFile,
	}, nil
}

func (m *Manager) CreateFork(sourceName, targetName string) error {
	sourceInfo, ok := m.Config.GetChain(sourceName)
	if !ok {
		return fmt.Errorf("source chain '%s' not found", sourceName)
	}

	if _, exists := m.Config.GetChain(targetName); exists {
		return fmt.Errorf("target chain '%s' already exists", targetName)
	}

	sourceStorage := storage.NewJSONStorage(sourceInfo.File)
	sourceBC, err := sourceStorage.Load()
	if err != nil {
		return fmt.Errorf("failed to load source chain: %w", err)
	}
	if sourceBC == nil {
		return fmt.Errorf("source chain is empty")
	}

	targetFile := fmt.Sprintf("blockchain_%s.json", targetName)

	srcData, err := os.ReadFile(sourceInfo.File)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	if err := os.WriteFile(targetFile, srcData, 0o644); err != nil {
		return fmt.Errorf("failed to create fork file: %w", err)
	}

	forkPoint := sourceBC.Length() - 1
	m.Config.AddChain(targetName, targetFile, &sourceName, &forkPoint)

	if err := m.Config.Save(m.configFile); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

func (m *Manager) GetChainFile(name string) (string, error) {
	info, ok := m.Config.GetChain(name)
	if !ok {
		return "", fmt.Errorf("chain '%s' not found in config", name)
	}
	return info.File, nil
}

func (m *Manager) RegisterChain(name string) error {
	if _, exists := m.Config.GetChain(name); exists {
		return nil
	}

	file := fmt.Sprintf("blockchain_%s.json", name)
	m.Config.AddChain(name, file, nil, nil)
	return m.Config.Save(m.configFile)
}

func (m *Manager) FindCommonAncestor(chain1, chain2 *blockchain.Blockchain) int {
	blocks1 := chain1.Blocks()
	blocks2 := chain2.Blocks()

	minLen := len(blocks1)
	if len(blocks2) < minLen {
		minLen = len(blocks2)
	}

	lastCommon := -1
	for i := 0; i < minLen; i++ {
		if blocks1[i].Hash == blocks2[i].Hash {
			lastCommon = i
		} else {
			break
		}
	}

	return lastCommon
}
