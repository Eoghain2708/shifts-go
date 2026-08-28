package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Basic colours
	Black   = lipgloss.NewStyle().Foreground(lipgloss.Color("0"))
	Red     = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	Green   = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	Yellow  = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	Blue    = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	Magenta = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	Cyan    = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	White   = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

	// Bright colours
	LightBlack   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	LightRed     = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	LightGreen   = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	LightYellow  = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	LightBlue    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	LightMagenta = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	LightCyan    = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	LightWhite   = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))

	// Bold basic colours
	BoldBlack   = Black.Bold(true)
	BoldRed     = Red.Bold(true)
	BoldGreen   = Green.Bold(true)
	BoldYellow  = Yellow.Bold(true)
	BoldBlue    = Blue.Bold(true)
	BoldMagenta = Magenta.Bold(true)
	BoldCyan    = Cyan.Bold(true)
	BoldWhite   = White.Bold(true)

	// Bold bright colours
	BoldLightBlack   = LightBlack.Bold(true)
	BoldLightRed     = LightRed.Bold(true)
	BoldLightGreen   = LightGreen.Bold(true)
	BoldLightYellow  = LightYellow.Bold(true)
	BoldLightBlue    = LightBlue.Bold(true)
	BoldLightMagenta = LightMagenta.Bold(true)
	BoldLightCyan    = LightCyan.Bold(true)
	BoldLightWhite   = LightWhite.Bold(true)
)
