package commits

import (
	"fmt"
	"github.com/LanceLRQ/commits-go/utils"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

const viewPortWidth = 60

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

type GitStatusViewModel struct {
	view         viewport.Model
	fileStatus   []GitFileStatus
	changeStatus GitChangeStatus
}

func (m GitStatusViewModel) renderViewPort() (*viewport.Model, error) {

	vp := viewport.New(viewPortWidth, 10)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	const glamourGutter = 2
	glamourRenderWidth := viewPortWidth - vp.Style.GetHorizontalFrameSize() - glamourGutter

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(glamourRenderWidth),
	)
	if err != nil {
		return nil, err
	}

	tmpl := utils.TranslateF("view_git_files_status", map[string]any{
		"FileStatus": m.renderFileStatusList(),
	})

	str, err := renderer.Render(tmpl)
	if err != nil {
		return nil, err
	}

	vp.SetContent(str)
	return &vp, nil
}

func (m GitStatusViewModel) renderFileStatusList() string {
	gitStatusNameMap := map[rune]string{
		0:   "",
		'M': utils.Translate("git_file_status_m"),
		'A': utils.Translate("git_file_status_a"),
		'D': utils.Translate("git_file_status_d"),
		'R': utils.Translate("git_file_status_r"),
		'C': utils.Translate("git_file_status_c"),
		'U': utils.Translate("git_file_status_u"),
		'?': utils.Translate("git_file_status_?"),
		'!': utils.Translate("git_file_status_!"),
	}
	sb := strings.Builder{}
	for _, file := range m.fileStatus {
		s := gitStatusNameMap[file.StagedStatus]
		us := gitStatusNameMap[file.UnstagedStatus]
		if file.TargetPath != "" {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s -> %s |\n", s, us, file.Path, file.TargetPath))
		}
		sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", s, us, file.Path))
	}
	return sb.String()
}

func NewGitStatusView() (*GitStatusViewModel, error) {
	// get git status
	fileStatus, gitChangeStatus, err := getGitStatus()
	if err != nil {
		return nil, err
	}

	instance := &GitStatusViewModel{
		fileStatus:   fileStatus,
		changeStatus: gitChangeStatus,
	}

	vp, err := instance.renderViewPort()
	if err != nil {
		return nil, err
	}
	instance.view = *vp

	return instance, nil
}

func (m GitStatusViewModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m GitStatusViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		default:
			var cmd tea.Cmd
			m.view, cmd = m.view.Update(msg)
			return m, cmd
		}
	default:
		return m, nil
	}
}

func (m GitStatusViewModel) View() string {
	return m.view.View() + m.helpView()
}

func (m GitStatusViewModel) helpView() string {
	return helpStyle("\n  ↑/↓: Navigate • q: Quit\n")
}
