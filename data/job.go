package data

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Job struct {
	Id      string   `json:"id"`
	Pid     int      `json:"pid"`
	Cmdline string   `json:"cmdline"`
	Name    string   `json:"name"`
	Args    []string `json:"args"`
}

func (j Job) MarshalBinary() ([]byte, error) {
	return json.Marshal(j)
}

func (j *Job) Parse() error {
	args := strings.Split(j.Cmdline, " ")
	if len(args) < 1 {
		return fmt.Errorf("not enough arguments in %v", j.Cmdline)
	}

	j.Name = args[0]
	if len(args) > 1 {
		j.Args = args[1:]
	} else {
		j.Args = []string{}
	}

	return nil
}
