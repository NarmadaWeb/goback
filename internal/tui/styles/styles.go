package styles

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	PrimaryColor    = lipgloss.Color("#7C3AED") // Purple
	SecondaryColor  = lipgloss.Color("#06B6D4") // Cyan
	AccentColor     = lipgloss.Color("#10B981") // Green
	ErrorColor      = lipgloss.Color("#EF4444") // Red
	WarningColor    = lipgloss.Color("#F59E0B") // Yellow
	SuccessColor    = lipgloss.Color("#10B981") // Green
	InfoColor       = lipgloss.Color("#3B82F6") // Blue
	MutedColor      = lipgloss.Color("#6B7280") // Gray
	BorderColor     = lipgloss.Color("#374151") // Dark Gray
	BackgroundColor = lipgloss.Color("#1F2937") // Dark Background
	TextColor       = lipgloss.Color("#F9FAFB") // Light Text
	HelpColor       = lipgloss.Color("#9CA3AF") // Help Text
)

// Base styles
var (
	// Title styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Align(lipgloss.Center).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true).
			Align(lipgloss.Center).
			MarginBottom(1)

	// Text styles
	DescriptionStyle = lipgloss.NewStyle().
				Foreground(MutedColor).
				Align(lipgloss.Center).
				MarginBottom(1)

	InfoStyle = lipgloss.NewStyle().
			Foreground(InfoColor)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor)

	ErrorTitleStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true).
			Underline(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(WarningColor).
			Bold(true)

	MutedStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Faint(true)

	AccentStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)

	// Menu and navigation styles
	ItemStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			PaddingLeft(2)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(AccentColor).
				Bold(true).
				PaddingLeft(2)

	CursorStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)

	// Help and instruction styles
	HelpStyle = lipgloss.NewStyle().
			Foreground(HelpColor).
			Italic(true).
			Align(lipgloss.Center).
			MarginTop(2)

	// Container styles
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// Form styles
	InputStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Background(BackgroundColor).
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(BorderColor)

	FocusedInputStyle = lipgloss.NewStyle().
				Foreground(TextColor).
				Background(BackgroundColor).
				Padding(0, 1).
				Border(lipgloss.NormalBorder()).
				BorderForeground(AccentColor)

	LabelStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Bold(true).
			MarginRight(1)

	// Progress styles
	ProgressBarStyle = lipgloss.NewStyle().
				Foreground(AccentColor)

	ProgressBackgroundStyle = lipgloss.NewStyle().
				Foreground(MutedColor).
				Faint(true)
)

const (
	StyleSuccess = "success"
	StyleError   = "error"
	StyleWarning = "warning"
	StyleInfo    = "info"
)

// RenderLogo for rendering
func RenderLogo() string {
	logo := `
   ▄██████▄   ▄██████▄  ▀█████████▄     ▄████████  ▄████████    ▄█   ▄█▄
  ███    ███ ███    ███   ███    ███   ███    ███ ███    ███   ███ ▄███▀
  ███    █▀  ███    ███   ███    ███   ███    ███ ███    █▀    ███▐██▀
 ▄███        ███    ███  ▄███▄▄▄██▀    ███    ███ ███         ▄█████▀
▀▀███ ████▄  ███    ███ ▀▀███▀▀▀██▄  ▀███████████ ███        ▀▀█████▄
  ███    ███ ███    ███   ███    ██▄   ███    ███ ███    █▄    ███▐██▄
  ███    ███ ███    ███   ███    ███   ███    ███ ███    ███   ███ ▀███▄
  ████████▀   ▀██████▀  ▄█████████▀    ███    █▀  ████████▀    ███   ▀█▀
                                                               ▀
`

	return lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Align(lipgloss.Center).
		Render(logo)
}

// Progress bar renderer
func RenderProgress(current, total int) string {
	if total <= 0 {
		return ""
	}

	width := 50
	progress := float64(current) / float64(total)
	filled := int(progress * float64(width))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	percentage := int(progress * 100)

	progressText := fmt.Sprintf("[%s] %d%%", bar, percentage)

	return ProgressBarStyle.Render(progressText)
}

// Key help renderer
func RenderKeyHelp(shortcuts map[string]string) string {
	if len(shortcuts) == 0 {
		return ""
	}

	var help []string
	for key, desc := range shortcuts {
		help = append(help,
			lipgloss.NewStyle().Foreground(AccentColor).Render(key)+
				": "+
				lipgloss.NewStyle().Foreground(HelpColor).Render(desc))
	}

	return HelpStyle.Render(strings.Join(help, " • "))
}

// Status message styles
func RenderSuccess(message string) string {
	return SuccessStyle.Render("✅ " + message)
}

func RenderError(message string) string {
	return ErrorStyle.Render("❌ " + message)
}

func RenderWarning(message string) string {
	return WarningStyle.Render("⚠️  " + message)
}

func RenderInfo(message string) string {
	return InfoStyle.Render("ℹ️  " + message)
}

// Box styles for different contexts
func RenderBox(title, content string, boxType string) string {
	var style lipgloss.Style

	switch boxType {
	case StyleSuccess:
		style = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SuccessColor).
			Padding(1, 2)
	case StyleError:
		style = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ErrorColor).
			Padding(1, 2)
	case StyleWarning:
		style = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(WarningColor).
			Padding(1, 2)
	case StyleInfo:
		style = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(InfoColor).
			Padding(1, 2)
	default:
		style = BorderStyle
	}

	var titleStyle lipgloss.Style
	switch boxType {
	case StyleSuccess:
		titleStyle = SuccessStyle
	case StyleError:
		titleStyle = ErrorStyle
	case StyleWarning:
		titleStyle = WarningStyle
	case StyleInfo:
		titleStyle = InfoStyle
	default:
		titleStyle = TitleStyle
	}

	if title != "" {
		return style.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				titleStyle.Render(title),
				"",
				content,
			),
		)
	}

	return style.Render(content)
}

// Loading spinner frames
var SpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func RenderSpinner(frame int, message string) string {
	spinner := SpinnerFrames[frame%len(SpinnerFrames)]
	return AccentStyle.Render(spinner) + " " + InfoStyle.Render(message)
}

// Menu container style
func RenderMenu(title, description string, items []string, cursor int) string {
	var menuItems []string

	for i, item := range items {
		if i == cursor {
			menuItems = append(menuItems, CursorStyle.Render("❯")+" "+SelectedItemStyle.Render(item))
		} else {
			menuItems = append(menuItems, "  "+ItemStyle.Render(item))
		}
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		TitleStyle.Render(title),
		DescriptionStyle.Render(description),
		"",
		strings.Join(menuItems, "\n"),
	)

	return BorderStyle.Render(content)
}

// Form field renderer
func RenderFormField(label, value, placeholder string, focused bool, required bool) string {
	labelText := label
	if required {
		labelText += " *"
	}

	style := InputStyle
	if focused {
		style = FocusedInputStyle
	}

	displayValue := value
	if value == "" && placeholder != "" {
		displayValue = MutedStyle.Render(placeholder)
	}

	if focused && value != "" {
		displayValue += "▎" // Cursor
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		LabelStyle.Render(labelText),
		style.Render(displayValue),
	)
}

// Configuration review styles
func RenderConfigReview(config map[string]string) string {
	var lines []string

	for key, value := range config {
		lines = append(lines,
			LabelStyle.Render(key+":")+" "+
				InfoStyle.Render(value))
	}

	return strings.Join(lines, "\n")
}

// Center content utility
func Center(content string, width, height int) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// Responsive width calculation
func GetContentWidth(terminalWidth int) int {
	maxWidth := 120
	minWidth := 60

	if terminalWidth < minWidth {
		return terminalWidth - 4 // Leave some margin
	}
	if terminalWidth > maxWidth {
		return maxWidth
	}
	return terminalWidth - 10 // Leave margin on sides
}

// Adaptive styles based on terminal capabilities
func GetAdaptiveStyle(hasColor bool) lipgloss.Style {
	if hasColor {
		return lipgloss.NewStyle().Foreground(PrimaryColor)
	}
	return lipgloss.NewStyle().Bold(true)
}

// Animation helper for transitions
func FadeIn(content string, opacity float64) string {
	if opacity >= 1.0 {
		return content
	}
	if opacity <= 0.0 {
		return ""
	}

	// Simple fade by adjusting faint property
	style := lipgloss.NewStyle()
	if opacity < 0.5 {
		style = style.Faint(true)
	}

	return style.Render(content)
}

// Multi-column layout helper
func RenderColumns(left, right string, totalWidth int) string {
	leftWidth := totalWidth/2 - 2
	rightWidth := totalWidth/2 - 2

	leftStyle := lipgloss.NewStyle().Width(leftWidth)
	rightStyle := lipgloss.NewStyle().Width(rightWidth)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(left),
		rightStyle.Render(right),
	)
}

// Configuration-specific styles
var (
	OptionStyle = lipgloss.NewStyle().
			Foreground(TextColor)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)

	DisabledStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Faint(true)

	ContinueStyle = lipgloss.NewStyle().
			Foreground(InfoColor).
			Bold(true)

	FieldLabelStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true)

	FieldValueStyle = lipgloss.NewStyle().
			Foreground(TextColor)

	PlaceholderStyle = lipgloss.NewStyle().
				Foreground(MutedColor).
				Italic(true)

	ConfirmStyle = lipgloss.NewStyle().
			Foreground(WarningColor).
			Bold(true).
			MarginTop(1).
			MarginBottom(1)
)

// Status badge renderer
func RenderBadge(text string, badgeType string) string {
	var style lipgloss.Style

	switch badgeType {
	case StyleSuccess:
		style = lipgloss.NewStyle().
			Background(SuccessColor).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 1).
			Bold(true)
	case StyleError:
		style = lipgloss.NewStyle().
			Background(ErrorColor).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).
			Bold(true)
	case StyleWarning:
		style = lipgloss.NewStyle().
			Background(WarningColor).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 1).
			Bold(true)
	case StyleInfo:
		style = lipgloss.NewStyle().
			Background(InfoColor).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).
			Bold(true)
	default:
		style = lipgloss.NewStyle().
			Background(MutedColor).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)
	}

	return style.Render(text)
}
