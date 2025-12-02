package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/rx3lixir/lab1-blockchain/models"
)

var (
	// Color palette - vibrant hex colors
	cyan   = lipgloss.Color("#00D9FF")
	green  = lipgloss.Color("#00FF9F")
	red    = lipgloss.Color("#FF6B6B")
	yellow = lipgloss.Color("#FFD93D")
	gray   = lipgloss.Color("#6C7A89")
	purple = lipgloss.Color("#A78BFA")

	// Basic styles
	bold      = lipgloss.NewStyle().Bold(true)
	success   = lipgloss.NewStyle().Foreground(green).Bold(true)
	errorSt   = lipgloss.NewStyle().Foreground(red).Bold(true)
	warning   = lipgloss.NewStyle().Foreground(yellow).Bold(true)
	info      = lipgloss.NewStyle().Foreground(cyan).Bold(true)
	muted     = lipgloss.NewStyle().Foreground(gray)
	highlight = lipgloss.NewStyle().Foreground(purple)
)

// Message helpers
func ShowSuccess(msg string) { fmt.Println(success.Render("✓ " + msg)) }
func ShowError(msg string)   { fmt.Println(errorSt.Render("✗ " + msg)) }
func ShowWarning(msg string) { fmt.Println(warning.Render("⚠ " + msg)) }
func ShowInfo(msg string)    { fmt.Println(info.Render("ℹ " + msg)) }

// ShowBlock displays a single block
func ShowBlock(block *models.Block) {
	if block == nil {
		ShowWarning("Block not found")
		return
	}

	fmt.Println()
	fmt.Println(bold.Render(fmt.Sprintf("BLOCK #%d", block.Index)))

	timeStr := time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05")
	fmt.Printf("%s %s\n", muted.Render("Time:"), timeStr)

	fmt.Println()
	fmt.Printf("%s %s\n", muted.Render("Name:"), block.Data.FullName)
	fmt.Printf("%s %s\n", muted.Render("Record:"), block.Data.Zachetka)
	fmt.Printf("%s %d / %s\n", muted.Render("Course/Group:"), block.Data.Course, block.Data.Group)
	fmt.Printf("%s %s\n", muted.Render("Subject:"), block.Data.Subject)
	fmt.Printf("%s %s\n", muted.Render("Grade:"), formatGrade(block.Data.Grade))

	fmt.Println()
	fmt.Printf("%s %s\n", muted.Render("Prev Hash:"), truncHash(block.PreviousHash))
	fmt.Printf("%s %s\n", muted.Render("Hash:"), truncHash(block.Hash))
	fmt.Printf("%s %d\n", muted.Render("Nonce:"), block.Nonce)
	fmt.Println(strings.Repeat("─", 60))
}

// ShowAllBlocks displays all blocks
func ShowAllBlocks(blocks []*models.Block) {
	fmt.Println(bold.Render(fmt.Sprintf("\nBLOCKCHAIN (%d blocks)", len(blocks))))
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

// Helpers
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
