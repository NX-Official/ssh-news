package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"os/exec"
	"runtime"
	"ssh-news/provider"
)

var currWindowsHeight int
var currWindowsWidth int

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	defer func() {
		switch m.selected {
		case selectedSources:
			// Set style of news-selected item to normal
			delegate := list.NewDefaultDelegate()
			delegate.SetSpacing(0)
			delegate.Styles.SelectedTitle = delegate.Styles.NormalTitle
			delegate.Styles.SelectedDesc = delegate.Styles.NormalDesc
			selectedNews := m.news[m.sources.SelectedItem().(sourcesItem).Title()]
			selectedNews.SetDelegate(delegate)
			selectedNews.ResetSelected()
			m.news[m.sources.SelectedItem().(sourcesItem).Title()] = selectedNews

		case selectedNews:
			// Set style of news-selected item to default
			selectedNews := m.news[m.sources.SelectedItem().(sourcesItem).Title()]
			delegate := list.NewDefaultDelegate()
			delegate.SetSpacing(0)
			selectedNews.SetDelegate(delegate)
			m.news[m.sources.SelectedItem().(sourcesItem).Title()] = selectedNews
		}

		// Set height of news and sources to current window height
		selectedNews := m.news[m.sources.SelectedItem().(sourcesItem).Title()]
		selectedNews.SetHeight(currWindowsHeight)
		selectedNews.SetWidth(currWindowsWidth - m.sources.Width())
		m.news[m.sources.SelectedItem().(sourcesItem).Title()] = selectedNews
		m.sources.SetHeight(currWindowsHeight)
	}()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "left":
			if m.news[m.sources.SelectedItem().(sourcesItem).Title()].Paginator.Page == 0 {
				m.selected = 0
			}
		case "right":
			if m.selected == 0 {
				m.selected = 1
				return m, cmd
			}
		case "enter":
			if m.selected == 1 {
				openBrowser(m.news[m.sources.SelectedItem().(sourcesItem).Title()].SelectedItem().(newsItem).url)
			}
		}
	case tea.WindowSizeMsg:
		currWindowsHeight = msg.Height
		currWindowsWidth = msg.Width
	}

	// pass the message to the selected model
	switch m.selected {
	case selectedSources:
		m.sources, cmd = m.sources.Update(msg)
	case selectedNews:
		m.news[m.sources.SelectedItem().(sourcesItem).Title()], cmd = m.news[m.sources.SelectedItem().(sourcesItem).Title()].Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Left, m.sources.View(), m.news[m.sources.SelectedItem().(sourcesItem).Title()].View())
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	provider.Fetch()
	UpdateModel(provider.Get())

	m.sources.SetFilteringEnabled(false)
	m.sources.SetShowTitle(false)
	m.sources.SetShowStatusBar(false)
	m.sources.SetShowHelp(false)

	for k := range m.news {
		v := m.news[k]
		v.SetShowTitle(false)
		v.SetShowHelp(false)
		v.SetFilteringEnabled(false)
		v.SetShowStatusBar(false)

		m.news[k] = v
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal("Error running program:", err)
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
		log.Println(err)
	}
}
