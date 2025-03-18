package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var boardStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder())

var positionStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF0000"))

type errMsg error

type cord struct {
	x int
	y int
}

type model struct {
	quitting bool
	err      error
	board    [][]int
	size     int
	pos      cord
	picking  int
}

var quitKeys = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)

type MoveKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Write key.Binding
}

var KeyMap = MoveKeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("l", "move right"),
	),
	Write: key.NewBinding(
		key.WithKeys("z"),
		key.WithHelp("z", "set board"),
	),
}

func initialModel() model {
	size := 10
	board := make([][]int, size)
	for i := range size {
		board[i] = make([]int, size)
	}

	return model{
		board: board,
		size:  size,
		pos: cord{
			x: 0,
			y: 0,
		},
		picking: 1,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.quitting = true
			return m, tea.Quit

		}
		if key.Matches(msg, KeyMap.Up) {
			if m.pos.x > 0 {
				m.pos.x -= 1
			}
		}
		if key.Matches(msg, KeyMap.Down) {
			if m.pos.x < m.size-1 {
				m.pos.x += 1
			}
		}
		if key.Matches(msg, KeyMap.Left) {
			if m.pos.y > 0 {
				m.pos.y -= 1
			}
		}
		if key.Matches(msg, KeyMap.Right) {
			if m.pos.y < m.size-1 {
				m.pos.y += 1
			}
		}
		if key.Matches(msg, KeyMap.Write) {
			if m.board[m.pos.x][m.pos.y] == 0 {
				m.board[m.pos.x][m.pos.y] = m.picking
				if m.gameset() {
					m.quitting = true
					return m, tea.Quit
				}
				if m.picking == 1 {
					m.picking = 2
				} else {
					m.picking = 1
				}
			}
		}
		return m, nil
	case errMsg:
		m.err = msg
		return m, nil

	default:
		return m, nil
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	cpicking := "O"
	if m.picking == 2 {
		cpicking = "X"
	}
	if m.quitting == true {
		return "GameSet, winner is " + cpicking + "\n"
	}
	str := ""
	for x := range m.size {
		row := ""
		for y := range m.size {
			next := ""
			if m.board[x][y] == 0 {
				next = ". "
			}
			if m.board[x][y] == 1 {
				next = "O "
			}
			if m.board[x][y] == 2 {
				next = "X "
			}
			if m.pos.x == x && m.pos.y == y {
				next = positionStyle.Render(next)
			}
			row += next
		}
		str += row + "\n"
	}
	return boardStyle.Render(str) + "\n" + fmt.Sprintf("curr: %s at (%d, %d)", cpicking, m.pos.x, m.pos.y)
}

func (m model) gameset() bool {
	for x := range m.size {
		for y := range m.size {
			if m.board[x][y] == 0 {
				continue
			}

			if x >= 2 && x <= m.size-3 && m.board[x-2][y] == m.board[x-1][y] && m.board[x-1][y] == m.board[x][y] && m.board[x+1][y] == m.board[x][y] && m.board[x+1][y] == m.board[x+2][y] {
				return true
			}
			if y >= 2 && y <= m.size-3 && m.board[x][y-2] == m.board[x][y-1] && m.board[x][y-1] == m.board[x][y] && m.board[x][y+1] == m.board[x][y] && m.board[x][y+1] == m.board[x][y+2] {
				return true
			}
			if x >= 2 && x <= m.size-3 && y >= 2 && y <= m.size-3 && m.board[x-2][y-2] == m.board[x-1][y-1] && m.board[x-1][y-1] == m.board[x][y] && m.board[x+1][y+1] == m.board[x][y] && m.board[x+1][y+1] == m.board[x+2][y+2] {
				return true
			}
			if x >= 2 && x <= m.size-3 && y >= 2 && y <= m.size-3 && m.board[x-2][y+2] == m.board[x-1][y+1] && m.board[x-1][y+1] == m.board[x][y] && m.board[x+1][y-1] == m.board[x][y] && m.board[x+1][y-1] == m.board[x+2][y-2] {
				return true
			}
		}
	}
	return false
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
