package commits

import (
	"fmt"
	"github.com/LanceLRQ/commits-go/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

// ** Status View **

const viewPortWidth = 60

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
var manualListTitleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFDF5")).
	Background(lipgloss.Color("#25A065")).
	Padding(0, 1)
var manualListTipMessageStyle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).Render
var manualListStyle = lipgloss.NewStyle().Padding(1, 2)

type GitStatusViewModel struct {
	stepMode     int // 0 - files view , 1 - manually add
	view         viewport.Model
	manualList   list.Model
	fileStatus   []GitFileStatus
	changeStatus GitChangeStatus
}

type delegateKeyMap struct {
	choose key.Binding
}

func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
	}
}
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
		},
	}
}

func newListItemDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
	}
}

func newListItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(GitFileStatus); ok {
			title = i.Title()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				return m.NewStatusMessage(manualListTipMessageStyle("You chose " + title))
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

func renderViewPort(fileStatus []GitFileStatus) (*viewport.Model, error) {
	view := viewport.New(viewPortWidth, 10)
	view.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	const glamourGutter = 2
	glamourRenderWidth := viewPortWidth - view.Style.GetHorizontalFrameSize() - glamourGutter

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(glamourRenderWidth),
	)
	if err != nil {
		return nil, err
	}

	tmpl := utils.TranslateF("view_git_files_status", map[string]any{
		"FileStatus": renderFileStatusList(fileStatus),
	})

	str, err := renderer.Render(tmpl)
	if err != nil {
		return nil, err
	}
	view.SetContent(str)
	return &view, nil
}

func renderManualAddFileView(fileStatus []GitFileStatus) (*list.Model, error) {
	// Make initial list of items
	numItems := len(fileStatus)
	items := make([]list.Item, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = fileStatus[i]
	}

	// Setup list
	manualList := list.New(items, newListItemDelegate(newListItemDelegateKeyMap()), 60, 30)
	manualList.Title = "手动添加文件"
	manualList.Styles.Title = manualListTitleStyle
	manualList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{}
	}
	return &manualList, nil
}

func renderFileStatusList(fileStatus []GitFileStatus) string {
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
	for _, file := range fileStatus {
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
		stepMode:     0,
		fileStatus:   fileStatus,
		changeStatus: gitChangeStatus,
	}

	vp, err := renderViewPort(fileStatus)
	if err != nil {
		return nil, err
	}
	instance.view = *vp

	ml, err := renderManualAddFileView(fileStatus)
	if err != nil {
		return nil, err
	}
	instance.manualList = *ml
	return instance, nil
}

func (m GitStatusViewModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m GitStatusViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.stepMode == 1 {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			h, v := manualListStyle.GetFrameSize()
			m.manualList.SetSize(msg.Width-h, msg.Height-v)

		case tea.KeyMsg:
			// Don't match any of the keys below if we're actively filtering.
			if m.manualList.FilterState() == list.Filtering {
				break
			}
		}

		// This will also call our delegate's update function.
		newListModel, cmd := m.manualList.Update(msg)
		m.manualList = newListModel
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)

	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c", "esc":
				return m, tea.Quit
			case "e":
				m.stepMode = 1
				m.manualList, cmd = m.manualList.Update(msg)
				return m, cmd
			default:
				return m, nil
			}
		default:
			m.view, cmd = m.view.Update(msg)
			return m, cmd
		}
	}
}

func (m GitStatusViewModel) View() string {
	if m.stepMode == 1 {
		return manualListStyle.Render(m.manualList.View())
	}
	return m.view.View() + m.helpView()
}

func (m GitStatusViewModel) helpView() string {
	return helpStyle("\n↑/↓: Navigate • a: Add all files to staged • e: Add files manually • q: Quit\n")
}

// ** End Status View **
