package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	AppStyle           = lipgloss.NewStyle().PaddingLeft(2).PaddingTop(1)
	TitleStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	focusedStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	focusButtonStyle   = lipgloss.NewStyle().Background(lipgloss.Color("205")).Foreground(lipgloss.Color("0"))
	blurredStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle        = focusedStyle
	noStyle            = lipgloss.NewStyle()
	HelpStyle          = blurredStyle
	listFocusedStyle   = lipgloss.NewStyle().Background(lipgloss.Color("#0066cc")).Foreground(lipgloss.Color("255"))
	listUnfocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

	focusedButton = focusButtonStyle.Render("  Submit  ")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))

	checkmarkChecked = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
)
