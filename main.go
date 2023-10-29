package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

const default_timeout = time.Minute * 4

type model struct {
	timer        timer.Model
	keymap       keymap
	help         help.Model
	quitting     bool
	soundPlaying bool
	timeout      time.Duration
}

type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
}

//go:embed timer_done.wav
var soundFile []byte

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case timer.TickMsg:
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		m.timer, cmd = m.timer.Update(msg)
		m.keymap.stop.SetEnabled(m.timer.Running())
		m.keymap.start.SetEnabled(!m.timer.Running())
		return m, cmd

	case timer.TimeoutMsg:
		m.quitting = true
		m.soundPlaying = true
		return m, play_sound()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.reset):
			m.timer.Timeout = default_timeout
			if m.timeout != 0 {
				m.timer.Timeout = m.timeout
			}
			var cmds []tea.Cmd
			if m.soundPlaying {
				cmds = append(cmds, stopSound(), m.timer.Start())
				m.soundPlaying = false
				m.quitting = false
			}
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			return m, m.timer.Toggle()
		}
	}

	return m, cmd
}

func (m model) RunningHelpView() string {
	return "\n     " + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
	})
}

func (m model) DoneHelpView() string {
	return "\n     " + m.help.ShortHelpView([]key.Binding{
		m.keymap.reset,
		m.keymap.quit,
	})
}

func (m model) View() string {
	s := "\n     " + m.timer.View()

	if m.timer.Timedout() {
		s = "\n\n     All done!\n\n" +
			m.DoneHelpView()
	}
	s += "\n"
	if !m.quitting {
		s = "\n\n     Time remaining: " + s
		s += m.RunningHelpView()
	}

	return s
}

type processFinishedMsg bool

func play_sound() tea.Cmd {
	return func() tea.Msg {
		f := bytes.NewReader(soundFile)

		streamer, format, err := wav.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		defer streamer.Close()

		sr := format.SampleRate * 2
		speaker.Init(sr, sr.N(time.Second/10))

		resampled := beep.Resample(4, format.SampleRate, sr, streamer)

		done := make(chan bool)
		speaker.Play(beep.Seq(resampled, beep.Callback(func() {
			done <- true
		})))
		<-done
		return processFinishedMsg(true)
	}
}

func stopSound() tea.Cmd {
	return func() tea.Msg {
		speaker.Clear()
		speaker.Close()
		return nil
	}
}

func main() {
	args := os.Args[1:]
	timeout := default_timeout

	if len(args) == 1 {
		if customTimeout, err := time.ParseDuration(args[0]); err == nil {
			timeout = customTimeout
		}
	}

	m := model{
		timer:   timer.NewWithInterval(timeout, time.Millisecond),
		timeout: timeout,
		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "reset"),
			),
			quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}
	m.keymap.start.SetEnabled(false)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Uh oh, we encountered an error:", err)
		os.Exit(1)
	}
}
