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

var docStyle = lipgloss.NewStyle()

type sourcesItem struct {
	name string
}

func (i sourcesItem) Title() string       { return i.name }
func (i sourcesItem) Description() string { return "" }
func (i sourcesItem) FilterValue() string { return "" }

type newsItem struct {
	title string
	desc  string
	url   string
}

func (i newsItem) Title() string       { return i.title }
func (i newsItem) Description() string { return i.desc }
func (i newsItem) FilterValue() string { return "" }

type model struct {
	sources  list.Model
	news     map[string]list.Model
	selected int
}

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
			openBrowser(m.news[m.sources.SelectedItem().(sourcesItem).Title()].SelectedItem().(newsItem).url)
		}

	case tea.WindowSizeMsg:
		selectedNews := m.news[m.sources.SelectedItem().(sourcesItem).Title()]
		selectedNews.SetHeight(msg.Height)
		m.news[m.sources.SelectedItem().(sourcesItem).Title()] = selectedNews
	}

	if m.selected == 0 {
		m.sources, cmd = m.sources.Update(msg)
	} else {
		m.news[m.sources.SelectedItem().(sourcesItem).Title()], cmd = m.news[m.sources.SelectedItem().(sourcesItem).Title()].Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Left, m.sources.View(), m.news[m.sources.SelectedItem().(sourcesItem).Title()].View())
}

func main() {
	sources := []list.Item{
		sourcesItem{name: "微博"},
		sourcesItem{name: "今日头条"},
	}

	news := map[string][]list.Item{
		"微博": {
			newsItem{title: "平台否认陆川被盗号", desc: "168.7万", url: "https://s.weibo.com/weibo?q=%23%E5%B9%B3%E5%8F%B0%E5%90%A6%E8%AE%A4%E9%99%86%E5%B7%9D%E8%A2%AB%E7%9B%97%E5%8F%B7%23&Refer=index"},
			newsItem{title: "天生臭脸和天生笑脸的区别", desc: "25.0万", url: "https://s.weibo.com/weibo?q=%23%E5%A4%A9%E7%94%9F%E8%87%AD%E8%84%B8%E5%92%8C%E5%A4%A9%E7%94%9F%E7%AC%91%E8%84%B8%E7%9A%84%E5%8C%BA%E5%88%AB%23&Refer=index"},
			newsItem{title: "为什么中国制造能行", desc: "4.7万", url: "https://s.weibo.com/weibo?q=%23%E4%B8%BA%E4%BB%80%E4%B9%88%E4%B8%AD%E5%9B%BD%E5%88%B6%E9%80%A0%E8%83%BD%E8%A1%8C%23&Refer=index"},
			newsItem{title: "饿了么17吃货节全月最优惠", desc: "NaN万", url: "https://s.weibo.com/weibo?q=%23%E9%A5%BF%E4%BA%86%E4%B9%8817%E5%90%83%E8%B4%A7%E8%8A%82%E5%85%A8%E6%9C%88%E6%9C%80%E4%BC%98%E6%83%A0%23&Refer=index"},
		},
		"今日头条": {
			newsItem{title: "今日头条 1", desc: "1111111111"},
			newsItem{title: "今日头条 2", desc: "1111111111"},
			newsItem{title: "今日头条 3", desc: "1111111111"},
		},
	}

	m := model{
		sources: list.New(sources, list.DefaultDelegate{
			ShowDescription: false,
			Styles:          list.NewDefaultItemStyles(),
		}, 20, 15),
		news: map[string]list.Model{
			"微博":   list.New(news["微博"], list.NewDefaultDelegate(), 30, 15),
			"今日头条": list.New(news["今日头条"], list.NewDefaultDelegate(), 30, 15),
		},
	}
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
