package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/rx3lixir/lab1-blockchain/models"
)

var (
	// Color palette
	cyan   = lipgloss.Color("#00D9FF")
	green  = lipgloss.Color("#00FF9F")
	red    = lipgloss.Color("#FF6B6B")
	yellow = lipgloss.Color("#FFD93D")
	gray   = lipgloss.Color("#6C7A89")
	purple = lipgloss.Color("#A78BFA")
	orange = lipgloss.Color("#FF9F43")

	// Styles
	bold      = lipgloss.NewStyle().Bold(true)
	success   = lipgloss.NewStyle().Foreground(green).Bold(true)
	errorSt   = lipgloss.NewStyle().Foreground(red).Bold(true)
	warning   = lipgloss.NewStyle().Foreground(yellow).Bold(true)
	info      = lipgloss.NewStyle().Foreground(cyan).Bold(true)
	muted     = lipgloss.NewStyle().Foreground(gray)
	highlight = lipgloss.NewStyle().Foreground(purple)
	forkBadge = lipgloss.NewStyle().Foreground(orange).Bold(true)
)

// Message helpers
func ShowSuccess(msg string) { fmt.Println(success.Render("✓ " + msg)) }
func ShowError(msg string)   { fmt.Println(errorSt.Render("✗ " + msg)) }
func ShowWarning(msg string) { fmt.Println(warning.Render("⚠ " + msg)) }
func ShowInfo(msg string)    { fmt.Println(info.Render("ℹ " + msg)) }

// ShowBlock displays a single block with fork information
func ShowBlock(block *models.Block) {
	if block == nil {
		ShowWarning("Block not found")
		return
	}

	fmt.Println()

	// Заголовок с версией форка
	title := fmt.Sprintf("BLOCK #%d", block.Index)
	forkLabel := getForkLabel(block.ForkVersion)
	if forkLabel != "" {
		title += " " + forkBadge.Render(forkLabel)
	}

	fmt.Println(bold.Render(title))

	timeStr := time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05")
	fmt.Printf("%s %s\n", muted.Render("Time:"), timeStr)

	fmt.Println()
	fmt.Printf("%s %s\n", muted.Render("Name:"), block.Data.FullName)
	fmt.Printf("%s %s\n", muted.Render("Record:"), block.Data.Zachetka)
	fmt.Printf("%s %d / %s\n", muted.Render("Course/Group:"), block.Data.Course, block.Data.Group)
	fmt.Printf("%s %s\n", muted.Render("Subject:"), block.Data.Subject)

	// Показываем Teacher если есть (Soft Fork+)
	if block.Data.Teacher != "" {
		fmt.Printf("%s %s %s\n",
			muted.Render("Teacher:"),
			block.Data.Teacher,
			forkBadge.Render("[Soft Fork field]"))
	} else if block.ForkVersion >= models.ForkVersionSoft {
		fmt.Printf("%s %s\n", muted.Render("Teacher:"), muted.Render("(not specified)"))
	}

	fmt.Printf("%s %s\n", muted.Render("Grade:"), formatGrade(block.Data.Grade))

	fmt.Println()

	// Показываем информацию о хэше с учётом форка
	hashInfo := getHashInfo(block)
	fmt.Printf("%s %s\n", muted.Render("Prev Hash:"), truncHash(block.PreviousHash))
	fmt.Printf("%s %s %s\n", muted.Render("Hash:"), truncHash(block.Hash), hashInfo)
	fmt.Printf("%s %d\n", muted.Render("Nonce:"), block.Nonce)

	// Дополнительная информация о сложности
	difficulty := getDifficulty(block)
	if difficulty != "" {
		fmt.Printf("%s %s\n", muted.Render("Difficulty:"), difficulty)
	}

	fmt.Println(strings.Repeat("─", 70))
}

// ShowAllBlocks displays all blocks with fork indicators
func ShowAllBlocks(blocks []*models.Block) {
	// Подсчитываем статистику форков
	forkStats := make(map[int]int)
	for _, block := range blocks {
		forkStats[block.ForkVersion]++
	}

	fmt.Println(bold.Render(fmt.Sprintf("\nBLOCKCHAIN (%d blocks)", len(blocks))))

	// Показываем статистику форков
	if len(forkStats) > 1 {
		fmt.Println()
		if count, ok := forkStats[models.ForkVersionOriginal]; ok && count > 0 {
			fmt.Printf("  %s %d blocks\n", muted.Render("Original:"), count)
		}
		if count, ok := forkStats[models.ForkVersionSoft]; ok && count > 0 {
			fmt.Printf("  %s %d blocks\n", forkBadge.Render("Soft Fork:"), count)
		}
		if count, ok := forkStats[models.ForkVersionHard]; ok && count > 0 {
			fmt.Printf("  %s %d blocks\n", errorSt.Render("Hard Fork:"), count)
		}
	}

	fmt.Println()

	for _, block := range blocks {
		ShowBlock(block)
	}
}

// ShowSearchResults displays search results
func ShowSearchResults(blocks []*models.Block, query string) {
	fmt.Println(bold.Render(fmt.Sprintf("\nSearch: \"%s\" (%d found)", query, len(blocks))))
	for _, block := range blocks {
		ShowBlock(block)
	}
}

// Helper functions

func getForkLabel(forkVersion int) string {
	switch forkVersion {
	case models.ForkVersionSoft:
		return "[SOFT FORK]"
	case models.ForkVersionHard:
		return "[HARD FORK]"
	default:
		return ""
	}
}

func getHashInfo(block *models.Block) string {
	hashLen := len(block.Hash)

	switch block.ForkVersion {
	case models.ForkVersionHard:
		return forkBadge.Render(fmt.Sprintf("[SHA3-512, %d chars]", hashLen))
	case models.ForkVersionSoft:
		return forkBadge.Render(fmt.Sprintf("[SHA3-384, %d chars]", hashLen))
	default:
		return muted.Render(fmt.Sprintf("[SHA3-384, %d chars]", hashLen))
	}
}

func getDifficulty(block *models.Block) string {
	// Определяем сложность по префиксу хэша
	if strings.HasPrefix(block.Hash, "000") {
		return forkBadge.Render("000 (Soft Fork+)")
	} else if strings.HasPrefix(block.Hash, "00") {
		return muted.Render("00 (Original)")
	}
	return ""
}

func truncHash(hash string) string {
	if len(hash) <= 20 {
		return highlight.Render(hash)
	}
	return highlight.Render(hash[:12] + "..." + hash[len(hash)-4:])
}

func formatGrade(grade int) string {
	var style lipgloss.Style
	switch grade {
	case 5:
		style = success
	case 4:
		style = info
	case 3:
		style = warning
	case 2:
		style = errorSt
	default:
		style = muted
	}
	return style.Render(fmt.Sprintf("%d", grade))
}

// ShowForkComparison показывает сравнение двух веток после Hard Fork
func ShowForkComparison(jsonBlocks, csvBlocks []*models.Block) {
	fmt.Println(bold.Render("\n=== FORK COMPARISON ==="))
	fmt.Println()

	fmt.Println(info.Render("📁 JSON Branch (Original + Soft Fork):"))
	fmt.Printf("   Blocks: %d\n", len(jsonBlocks))
	fmt.Printf("   Hash: SHA3-384\n")
	fmt.Printf("   Format: JSON\n")
	fmt.Println()

	fmt.Println(errorSt.Render("📁 CSV Branch (Hard Fork):"))
	fmt.Printf("   Blocks: %d\n", len(csvBlocks))
	fmt.Printf("   Hash: SHA3-512\n")
	fmt.Printf("   Format: CSV\n")
	fmt.Println()

	fmt.Println(warning.Render("⚠️  These branches are INCOMPATIBLE!"))
	fmt.Println("   New blocks on one branch cannot be validated on the other.")
	fmt.Println()
}
