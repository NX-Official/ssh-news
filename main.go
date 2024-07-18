package main

import (
	"context"
	"errors"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
	"net"
	"os"
	"os/signal"
	"ssh-news/provider"
	"syscall"
	"time"
)

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
		selectedNews.SetHeight(m.currWindowsHeight)
		selectedNews.SetWidth(m.currWindowsWidth)
		m.news[m.sources.SelectedItem().(sourcesItem).Title()] = selectedNews
		m.sources.SetHeight(m.currWindowsHeight)
	}()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "left":
			if m.news[m.sources.SelectedItem().(sourcesItem).Title()].Paginator.Page == 0 {
				m.selected = selectedSources
				selectedSourcesIdx := m.sources.Index()
				newM := newModel(provider.Get(), m.session)
				m.news = newM.news
				m.sources = newM.sources
				m.sources.Select(selectedSourcesIdx)
			}
		case "right":
			if m.selected == selectedSources {
				m.selected = selectedNews
				return m, cmd
			}
		case "enter", " ":
			if m.selected == selectedNews {

				news := m.news[m.sources.SelectedItem().(sourcesItem).Title()]
				selected := news.SelectedItem().(newsItem)
				selectedIdx := news.Index()
				copyURL(selected.url, m.session)
				selected.desc = "已复制链接到剪贴板"
				news.SetItem(selectedIdx, selected)
				m.news[m.sources.SelectedItem().(sourcesItem).Title()] = news

			}
		}
	case tea.WindowSizeMsg:
		m.currWindowsHeight = msg.Height
		m.currWindowsWidth = msg.Width
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

func init() {
	//f, err := tea.LogToFile("debug.log", "")
	//if err != nil {
	//	// log.Fatal(err)
	//}
	//defer f.Close()

}

const (
	host = "localhost"
	port = "23234"
)

func main() {
	cronFetch()

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			myCustomBubbleteaMiddleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}

}

func copyURL(url string, s ssh.Session) {
	bubbletea.MakeRenderer(s).Output().Copy(url)
}

func myCustomBubbleteaMiddleware() wish.Middleware {
	newProg := func(m tea.Model, opts ...tea.ProgramOption) *tea.Program {
		p := tea.NewProgram(m, opts...)
		return p
	}
	teaHandler := func(s ssh.Session) *tea.Program {
		m := newModel(provider.Get(), s)
		return newProg(m, append(bubbletea.MakeOptions(s), tea.WithAltScreen())...)
	}
	return bubbletea.MiddlewareWithProgramHandler(teaHandler, termenv.ANSI256)
}
