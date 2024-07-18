package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"os/exec"
	"runtime"
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "left":
			m.selected = 0
		case "right":
			m.selected = 1
		case "enter":
			if m.selected == 1 {
				openBrowser(m.news[m.sources.SelectedItem().(sourcesItem).Title()].SelectedItem().(newsItem).url)
			}
		}

	case tea.WindowSizeMsg:
		selectedNews := m.news[m.sources.SelectedItem().(sourcesItem).Title()]
		selectedNews.SetHeight(msg.Height)
		m.news[m.sources.SelectedItem().(sourcesItem).Title()] = selectedNews
		m.sources.SetHeight(msg.Height)
	}

	styles := list.NewDefaultItemStyles()
	styles.SelectedTitle = selectedTitleStyle
	styles.SelectedDesc = selectedDescStyle

	if m.selected == 0 {
		m.sources.SetDelegate(list.DefaultDelegate{
			ShowDescription: false,
			Styles:          styles,
		})
		selectedNews := m.news[m.sources.SelectedItem().(sourcesItem).Title()]
		selectedNews.SetDelegate(list.NewDefaultDelegate())
		m.news[m.sources.SelectedItem().(sourcesItem).Title()] = selectedNews
		m.sources, cmd = m.sources.Update(msg)
	} else {
		m.news[m.sources.SelectedItem().(sourcesItem).Title()], cmd = m.news[m.sources.SelectedItem().(sourcesItem).Title()].Update(msg)
		m.sources.SetDelegate(list.DefaultDelegate{
			ShowDescription: false,
			Styles:          list.NewDefaultItemStyles(),
		})
		selectedNews := m.news[m.sources.SelectedItem().(sourcesItem).Title()]

		delegate := list.NewDefaultDelegate()
		delegate.Styles = styles

		selectedNews.SetDelegate(delegate)
		m.news[m.sources.SelectedItem().(sourcesItem).Title()] = selectedNews

	}
	return m, cmd
}

func (m model) View() string {

	return lipgloss.JoinHorizontal(lipgloss.Left, m.sources.View(), m.news[m.sources.SelectedItem().(sourcesItem).Title()].View())
}

func main() {
	UpdateNews([]byte(rawNews))
	UpdateModel()

	m.sources.SetFilteringEnabled(false)
	m.sources.SetShowTitle(false)
	m.sources.SetShowStatusBar(false)
	m.sources.SetShowHelp(false)

	for k := range m.news {
		v := m.news[k]
		v.SetShowTitle(false)
		v.SetShowHelp(false)
		v.SetFilteringEnabled(false)
		m.news[k] = v
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Println(err)
	}
}
