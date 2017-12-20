package alfred

import (
	"regexp"
	"strings"
	"sync"
)

// TaskGroup contains a task name and it's arguments
type TaskGroup struct {
	Name string
	Args []string
}

// ParseTaskGroup takes in a string, and parses it into a TaskGroup
func (t *Task) ParseTaskGroup(group string) []TaskGroup {
	tg := make([]TaskGroup, 0)
	group = strings.TrimSpace(group)

	if group == "" {
		return tg
	}

	if strings.Index(group, "\n") == -1 {
		// This means we have a regular space delimited list
		tasks := strings.Split(group, " ")
		for _, task := range tasks {
			tg = append(tg, TaskGroup{Name: task, Args: []string{}})
		}
	} else {
		// mix and match here
		tasks := strings.Split(group, "\n")
		for _, task := range tasks {
			re := regexp.MustCompile(`(.*?)\((.*?)\)`)
			results := re.FindStringSubmatch(task)
			if len(results) == 0 {
				tg = append(tg, TaskGroup{Name: strings.TrimSpace(task), Args: []string{}})
			} else {
				args := strings.Split(results[2], ",")
				for idx, a := range args {
					// trim the extra whitespace
					args[idx] = strings.TrimSpace(a)
				}
				tg = append(tg, TaskGroup{Name: strings.TrimSpace(results[1]), Args: args})
			}
		}
	}

	return tg
}

func execTaskGroup(taskGroups []TaskGroup, task Task, context *Context, tasks map[string]Task) {
	for _, tg := range taskGroups {
		NewTask(tg.Name, InitialContext(tg.Args), tasks)
	}
}

func goExecTaskGroup(taskGroups []TaskGroup, task Task, context *Context, tasks map[string]Task) {
	var wg sync.WaitGroup
	for _, tg := range taskGroups {
		wg.Add(1)
		go func() {
			NewTask(tg.Name, InitialContext(tg.Args), tasks)
			wg.Done()
		}()
		wg.Wait()
	}
}
