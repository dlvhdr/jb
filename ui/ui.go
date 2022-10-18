package ui

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-redis/redis"
)

type Job struct {
	key     string
	cmdline string
	command *exec.Cmd
}

type Model struct {
	jobs   []Job
	pubsub *redis.PubSub
}

func NewModel() Model {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	keys, err := rdb.Keys("job-*").Result()
	if err != nil {
		panic(err)
	}

	var jobs []Job
	for _, key := range keys {
		job, err := rdb.Get(key).Result()
		if err != nil {
			continue
		}
		jobs = append(jobs, Job{
			key:     key,
			cmdline: job,
		})
	}

	pubsub := rdb.Subscribe("jobber-new-job")
	defer pubsub.Close()

	return Model{pubsub: pubsub, jobs: jobs}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		// case newJobMsg:

		// case jobOutput:

	}
	return m, nil
}

func (m Model) View() string {
	s := strings.Builder{}
	s.WriteString("JOBBER\n\n")

	s.WriteString("JOBS\n")

	var rJobs []string
	for _, job := range m.jobs {
		rJobs = append(rJobs, fmt.Sprintf("[%v] %v", job.key, job.cmdline))
	}
	s.WriteString(lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Align(lipgloss.Left).
		Render(
			lipgloss.JoinVertical(lipgloss.Left, rJobs...),
		),
	)

	return s.String()
}
