package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	selectedWindow int
	methodBox      textarea.Model
	urlInput       textarea.Model
	body           textarea.Model
	response       textarea.Model
	output         string
	url            string
	input          string
	method         string
	senderStyle    lipgloss.Style
	err            error
}

func initialModel() model {
	//url box
	ta := textarea.New()
	ta.Placeholder = "Enter URL here...."
	// ta.Focus()

	// ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(50)
	ta.SetHeight(1)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false
	// input box
	inta := textarea.New()
	inta.Placeholder = "Enter body here...."
	// inta.Focus()
	// inta.Prompt = "┃ "
	inta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	inta.ShowLineNumbers = false
	// output box
	outta := textarea.New()
	outta.Placeholder = "Response will be here...."
	// outta.Focus()
	// outta.Prompt = "┃ "
	outta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	outta.ShowLineNumbers = false
	outta.SetHeight(10)
	// method box
	mta := textarea.New()
	mta.Placeholder = "GET"
	// mta.Focus()
	// mta.Prompt = "┃ "
	mta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	mta.ShowLineNumbers = false
	mta.SetWidth(10)
	mta.SetHeight(1)
	mta.Focus()

	return model{
		urlInput:    ta,
		body:        inta,
		response:    outta,
		output:      "",
		url:         "",
		input:       "",
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		methodBox:   mta,
		method:      "GET",
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		urlCmd    tea.Cmd
		inputCmd  tea.Cmd
		outputCmd tea.Cmd
		methodCmd tea.Cmd
	)
	// m.textarea, tiCmd = m.textarea.Update(msg)
	// m.viewport, vpCmd = m.viewport.Update(msg)
	switch m.selectedWindow {
	case 1:
		m.urlInput, urlCmd = m.urlInput.Update(msg)
	case 2:
		m.body, inputCmd = m.body.Update(msg)
	case 0:
		m.methodBox, methodCmd = m.methodBox.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			fmt.Println("Exiting...")
			fmt.Println(urlCmd)
			fmt.Println(inputCmd)
			fmt.Println(outputCmd)
			fmt.Println(methodCmd)
			return m, tea.Quit
		case tea.KeyTab:
			m.selectedWindow = (m.selectedWindow + 1) % 4
			switch m.selectedWindow {
			case 1:
				m.urlInput.Focus()
				// remove cursor from other windows
				m.body.Blur()
				m.response.Blur()
				m.methodBox.Blur()
			case 2:
				m.body.Focus()
				m.urlInput.Blur()
				m.response.Blur()
				m.methodBox.Blur()
			case 3:
				m.response.Focus()
				m.urlInput.Blur()
				m.body.Blur()
				m.methodBox.Blur()
			case 0:
				m.methodBox.Focus()
				m.urlInput.Blur()
				m.body.Blur()
				m.response.Blur()
			}
		case tea.KeyCtrlR:
			client := &http.Client{}
			var data = strings.NewReader(m.body.Value())
			req, err := http.NewRequest(m.methodBox.Value(), m.body.Value(), data)
			if err != nil {
				log.Fatal(err)
			}
			req.Header.Set("accept", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			bodyText, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			m.response.SetValue(string(bodyText))
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(urlCmd, inputCmd, outputCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"URL:\t%s\t%s\nBody:\n%s\nResponse:\n%s",
		m.methodBox.View(),
		m.urlInput.View(),
		m.body.View(),
		m.response.View(),
	) + "\n\n"
}
