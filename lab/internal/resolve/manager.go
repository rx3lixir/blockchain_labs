package resolve

import (
	"fmt"

	"github.com/rx3lixir/lab_bc/internal/blockchain"
	"github.com/rx3lixir/lab_bc/internal/fork"
	"github.com/rx3lixir/lab_bc/internal/storage"
)

type Manager struct {
	forkMgr *fork.Manager
}

func NewManager(forkMgr *fork.Manager) *Manager {
	return &Manager{forkMgr: forkMgr}
}

func (m *Manager) Validate(chain1Name, chain2Name string) error {
	file1, err := m.forkMgr.GetChainFile(chain1Name)
	if err != nil {
		return err
	}

	file2, err := m.forkMgr.GetChainFile(chain2Name)
	if err != nil {
		return err
	}

	storage1 := storage.NewJSONStorage(file1)
	bc1, err := storage1.Load()
	if err != nil {
		return fmt.Errorf("failed to load chain '%s': %w", chain1Name, err)
	}

	storage2 := storage.NewJSONStorage(file2)
	bc2, err := storage2.Load()
	if err != nil {
		return fmt.Errorf("failed to load chain '%s': %w", chain2Name, err)
	}

	if err := bc1.Validate(); err != nil {
		return fmt.Errorf("chain '%s' validation failed: %w", chain1Name, err)
	}

	if err := bc2.Validate(); err != nil {
		return fmt.Errorf("chain '%s' validation failed: %w", chain2Name, err)
	}

	commonAncestor := m.forkMgr.FindCommonAncestor(bc1, bc2)
	if commonAncestor == -1 {
		return fmt.Errorf("chains have no common ancestor")
	}

	fmt.Printf("✓ Both chains are valid\n")
	fmt.Printf("✓ Common ancestor found at block #%d\n", commonAncestor)
	fmt.Printf("  Chain '%s': %d blocks\n", chain1Name, bc1.Length())
	fmt.Printf("  Chain '%s': %d blocks\n", chain2Name, bc2.Length())

	return nil
}

func (m *Manager) Resolve(chain1Name, chain2Name string) error {
	file1, err := m.forkMgr.GetChainFile(chain1Name)
	if err != nil {
		return err
	}

	file2, err := m.forkMgr.GetChainFile(chain2Name)
	if err != nil {
		return err
	}

	storage1 := storage.NewJSONStorage(file1)
	bc1, err := storage1.Load()
	if err != nil {
		return fmt.Errorf("failed to load chain '%s': %w", chain1Name, err)
	}

	storage2 := storage.NewJSONStorage(file2)
	bc2, err := storage2.Load()
	if err != nil {
		return fmt.Errorf("failed to load chain '%s': %w", chain2Name, err)
	}

	commonAncestor := m.forkMgr.FindCommonAncestor(bc1, bc2)
	if commonAncestor == -1 {
		return fmt.Errorf("chains have no common ancestor - cannot resolve")
	}

	var winner, loser *blockchain.Blockchain
	var winnerName, loserName string
	var winnerStorage, loserStorage *storage.JSONStorage

	info1, _ := m.forkMgr.Config.GetChain(chain1Name)
	info2, _ := m.forkMgr.Config.GetChain(chain2Name)

	if bc1.Length() > bc2.Length() {
		winner, loser = bc1, bc2
		winnerName, loserName = chain1Name, chain2Name
		winnerStorage, loserStorage = storage1, storage2
	} else if bc2.Length() > bc1.Length() {
		winner, loser = bc2, bc1
		winnerName, loserName = chain2Name, chain1Name
		winnerStorage, loserStorage = storage2, storage1
	} else {
		if info1.ForkFrom == nil {
			winner, loser = bc1, bc2
			winnerName, loserName = chain1Name, chain2Name
			winnerStorage, loserStorage = storage1, storage2
		} else if info2.ForkFrom == nil {
			winner, loser = bc2, bc1
			winnerName, loserName = chain2Name, chain1Name
			winnerStorage, loserStorage = storage2, storage1
		} else if *info1.ForkFrom == chain2Name {
			winner, loser = bc2, bc1
			winnerName, loserName = chain2Name, chain1Name
			winnerStorage, loserStorage = storage2, storage1
		} else {
			winner, loser = bc1, bc2
			winnerName, loserName = chain1Name, chain2Name
			winnerStorage, loserStorage = storage1, storage2
		}
	}

	fmt.Printf("Winner: '%s' (%d blocks)\n", winnerName, winner.Length())
	fmt.Printf("Loser: '%s' (%d blocks)\n", loserName, loser.Length())

	loserBlocks := loser.Blocks()
	winnerBlocks := winner.Blocks()

	existingIDs := make(map[string]bool)
	for _, block := range winnerBlocks {
		if block.Data.ID != "" {
			existingIDs[block.Data.ID] = true
		}
	}

	addedCount := 0
	for i := commonAncestor + 1; i < len(loserBlocks); i++ {
		record := loserBlocks[i].Data
		if record.ID == "" || !existingIDs[record.ID] {
			if _, err := winner.AddBlock(record); err != nil {
				return fmt.Errorf("failed to add block from loser chain: %w", err)
			}
			existingIDs[record.ID] = true
			addedCount++
		}
	}

	if err := winnerStorage.Save(winner); err != nil {
		return fmt.Errorf("failed to save winner chain: %w", err)
	}

	if err := loserStorage.Save(winner); err != nil {
		return fmt.Errorf("failed to save loser chain: %w", err)
	}

	fmt.Printf("✓ Resolve complete\n")
	fmt.Printf("  Added %d unique records from '%s' to '%s'\n", addedCount, loserName, winnerName)
	fmt.Printf("  Both chains now have %d blocks\n", winner.Length())

	return nil
}
