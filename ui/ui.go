package ui

import (
	"bufio"
	"fmt"
	"jobber/utils"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis"
)

type Job struct {
	cmdline string
	command *exec.Cmd
}

type Model struct {
	jobs    []Job
	watcher *fsnotify.Watcher
	pubsub  *redis.PubSub
}

func NewModel() Model {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pubsub := rdb.Subscribe("jobber-new-job")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	xdgDataDir := os.Getenv("XDG_DATA_HOME")
	jobberDataDir := path.Join(xdgDataDir, "jobber")
	jobberDataFile := path.Join(jobberDataDir, "data")

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					lastLine := utils.GetLastLineWithSeek(jobberDataFile)
					log.Println("Line added: ", lastLine)
					cmdArgs := strings.Split(lastLine, " ")
					if len(cmdArgs) < 1 {
						log.Println("Empty command")
						return
					}

					job := exec.Command("bash", "-c", lastLine)
					job.Dir = "/Users/dolevh/code/personal/github.com/dlvhdr/jobber"
					stdout, _ := job.StdoutPipe()

					err = job.Start()
					if err != nil {
						panic(err)
					}
					log.Printf(
						"Started cmd. pid: %v | path: %v, dir: %v, args: %v",
						job.Process.Pid,
						job.Path,
						job.Dir,
						strings.Join(job.Args, " "),
					)

					scanner := bufio.NewScanner(stdout)
					scanner.Split(bufio.ScanLines)
					for scanner.Scan() {
						m := scanner.Text()
						fmt.Println(m)
					}
					job.Wait()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(jobberDataFile)
	if err != nil {
		log.Fatal(err)
	}

	dataContent, err := os.ReadFile(jobberDataFile)
	var jobs []Job
	if err != nil {
		jobs = []Job{}
	} else {
		cmdlines := strings.Split(string(dataContent), "\n")
		for _, cmdline := range cmdlines {
			jobs = append(jobs, Job{
				cmdline: cmdline,
			})
		}
	}

	return Model{watcher: watcher, pubsub: pubsub}
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
	// s.WriteString(lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Align(lipgloss.Left).Render(
	// 	lipgloss.JoinVertical(lipgloss.Left, m.jobs...),
	// ))

	return s.String()
}
