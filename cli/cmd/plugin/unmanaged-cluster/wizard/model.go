// Copyright 2020-2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package wizard utilizes the bubbletea library to display
// and generate a interactive terminal user interface
package wizard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Bold(true)
	blurredStyle = lipgloss.NewStyle().Faint(true)
	noStyle      = lipgloss.NewStyle()

	focusedButton = focusedStyle.Copy().Render("[ -- Submit -- ]")
	blurredButton = fmt.Sprintf("[ -- %s -- ]", blurredStyle.Render("Submit"))
)

type Model struct {
	// Components to be accessed by unmanaged-cluster
	Clustername           string
	WorkerNodeCount       string
	ControlPlaneNodeCount string

	// Err is the error generated during interactive terminal
	Err error

	// Internal components of the bubbletea control loop
	inputs     []textinput.Model
	focusIndex int
}

//nolint:gocritic
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

//nolint:gocritic
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "shift+enter":
			m.populateWizardModel()
			return m, tea.Quit

		case "ctrl-c", "esc":
			m.setError(fmt.Errorf("exiting interactive terminal"))
			return m, tea.Quit

		// Set focus to next input
		//nolint:goconst
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				m.populateWizardModel()
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *Model) setError(e error) {
	m.Err = e
}

func (m *Model) populateWizardModel() {
	for i := range m.inputs {
		switch i {
		case 0:
			m.Clustername = m.inputs[i].Value()
		case 1:
			m.WorkerNodeCount = m.inputs[i].Value()
		case 2:
			m.ControlPlaneNodeCount = m.inputs[i].Value()
		}
	}
}

var welcomeMessage = `
🧙 Welcome! I am your friendly Unmanaged Cluster Wizard!

Use <tab> or the arrow keys to move between options.
Hit enter when over the Submit button or <shift + enter> to start bootstrapping!

You can read more about all these options in our documentation:
https://tanzucommunityedition.io/docs/

--------------------------------------
`

//nolint:gocritic
func (m Model) View() string {
	var b strings.Builder

	b.WriteString(welcomeMessage)

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}

func InitalModel() Model {
	m := Model{
		inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = focusedStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Prompt = "\nCluster Name: "
			t.Placeholder = "my-cluster-name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Prompt = "\n(optional) Worker node count: "
			t.Placeholder = "1"
			t.CharLimit = 3
		case 2:
			t.Prompt = "\n(optional) Control plane node count: "
			t.Placeholder = "1"
			t.CharLimit = 3
		}

		m.inputs[i] = t
	}

	return m
}
