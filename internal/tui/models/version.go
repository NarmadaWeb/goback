// internal/tui/models/version.go

package models

import (
	"github.com/NarmadaWeb/goback/internal/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var version string = "0.1.7"

type VersionModel struct {
	selected bool
}

func NewVersionModel() *VersionModel {
	return &VersionModel{}
}

func (m *VersionModel) Init() tea.Cmd {
	return nil
}

func (m *VersionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case keyEnter, keyEsc, keyQ:
			m.selected = true
		case keyCtrlC:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *VersionModel) View() string {
	title := styles.TitleStyle.Render("🚀 GoBack TUI Generator")

	version := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginTop(1).
		MarginBottom(1).
		Render(version)

	description := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250")).
		MarginBottom(2).
		Render("GoBack is a TUI (Terminal User Interface) built with Bubble Tea\n" +
			"to make it easier for backend developers to create backend projects with\n" +
			"various choices of frameworks, databases, ORMs, architectures, and DevOps tools.")

	features := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true).
		MarginBottom(1).
		Render("✨ Main Features:")

	featuresList := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250")).
		MarginLeft(2).
		MarginBottom(2).
		Render(`• Framework: Go Fiber, Gin, Chi, Echo
• Database: postgresQL, MySQL, SQLite
• Tools: SQLX, SQLC
• Arsitektur: Simple, DDD, Clean, Hexagonal
• DevOps: Kubernetes, Helm, Terraform, Ansible
• Auto-generate: CRUD operations, config, docs
• Template: Production-ready boilerplate`)

	tech := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true).
		MarginBottom(1).
		Render("⚡ Technology:")

	techList := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250")).
		MarginLeft(2).
		MarginBottom(2).
		Render(`• Built with: Go + Bubble Tea TUI
• Templates: Go templates with validation
• Generator: Automatic scaffolding
• Config: YAML/JSON support
• Validation: Input and business rules`)

	author := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(2).
		MarginBottom(1).
		Render("👨‍💻 Developed by: GoBack Team")

	repo := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginBottom(2).
		Render("🔗 Repository: https://github.com/NarmadaWeb/goback")

	help := styles.HelpStyle.Render("enter/esc: back to menu • ctrl+c: exit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		version,
		description,
		features,
		featuresList,
		tech,
		techList,
		author,
		repo,
		help,
	)
}

func (m *VersionModel) ShouldClose() bool {
	return m.selected
}

func (m *VersionModel) Reset() {
	m.selected = false
}
