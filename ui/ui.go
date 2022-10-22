package ui

import (
	"fmt"
	"jobber/data"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-redis/redis"
)

type jobMeta struct {
	job    data.Job
	output string
}

type Model struct {
	cursor int
	jobs   []jobMeta
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

	var jobs []jobMeta
	for _, key := range keys {
		job, err := rdb.Get(key).Result()
		if err != nil {
			continue
		}
		jobData, err := data.UnmarshalBinary([]byte(job))
		if err != nil {
			continue
		}

		jobMeta := jobMeta{
			job: jobData,
		}
		jobs = append(jobs, jobMeta)
	}

	pubsub := rdb.Subscribe("jobber-new-job")

	return Model{pubsub: pubsub, jobs: jobs}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.subscribe(),
		m.readJobOuput(),
	)
}

type jobOuput struct {
	id     string
	output string
}

func (m *Model) readJobOuput() tea.Cmd {
	return func() tea.Msg {
		if len(m.jobs) == 0 {
			return nil
		}

		jobMeta := m.jobs[m.cursor]

		b, err := os.ReadFile(jobMeta.job.OutPath)
		var output string
		if err != nil {
			output = err.Error()
		} else {
			output = string(b)
		}

		return jobOuput{
			id:     jobMeta.job.Id,
			output: output,
		}
	}
}

func (m *Model) subscribe() tea.Cmd {
	return func() tea.Msg {
		msg, err := m.pubsub.ReceiveMessage()
		if err != nil {
			return nil
		}

		return msg
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case *redis.Message:
		job, err := data.UnmarshalBinary([]byte(msg.Payload))
		if err != nil {
			return m, nil
		}
		m.jobs = append(m.jobs, jobMeta{
			job: job,
		})
		return m, m.subscribe()

	case jobOuput:
		jId := -1
		var meta jobMeta
		for i, j := range m.jobs {
			if j.job.Id == msg.id {
				jId = i
				meta = j
				break
			}
		}

		if jId == -1 {
			return m, nil
		}

		meta.output = msg.output
		m.jobs[jId] = meta

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "j":
			m.cursor = (m.cursor + 1) % len(m.jobs)
			return m, m.readJobOuput()

		case "k":
			c := m.cursor - 1
			if c == -1 {
				c = len(m.jobs) - 1
			}
			m.cursor = c
			return m, m.readJobOuput()

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
	for i, job := range m.jobs {
		s := lipgloss.NewStyle()
		var jMeta string

		if i == m.cursor {
			s = s.Bold(true)
			jMeta = lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Render(
				fmt.Sprintf("OutPath: %v\nPid: %v\nPwd: %v\n", job.job.OutPath, job.job.Pid, job.job.Pwd),
			)
		}

		jTitle := s.Render(fmt.Sprintf("[%v] %v", job.job.Id, job.job.Cmdline))

		rJobs = append(rJobs, lipgloss.JoinVertical(lipgloss.Left, jTitle, jMeta))

	}

	jobsList := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Align(lipgloss.Left).
		Render(
			lipgloss.JoinVertical(lipgloss.Left, rJobs...),
		)

	jobOutput := lipgloss.NewStyle().Width(20).Border(lipgloss.NormalBorder()).Render(m.jobs[m.cursor].output)

	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, jobsList, jobOutput))

	return s.String()
}
