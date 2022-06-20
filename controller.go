package main

import (
	"fmt"
	"os"
	"bufio"
	"strconv"
	"errors"
	"log"
)

// Basic functions

func min(slice []int) int {
	m := 0
	for i, v := range slice {
		if i == 0 || v < m {
			m = v
		} 
	}
	return m
}

func isInInt(slice []int, obj int) bool {
	for _, one := range slice {
		if one == obj {
			return true
		}
	}
	return false
}

func isInStr(slice []string, obj string) bool {
	for _, one := range slice {
		if one == obj {
			return true
		}
	}
	return false
}

func isInSpec(slice []Spec, obj Spec) (bool, int) {
	for i, one := range slice {
		if one.name == obj.name {
			return true, i
		}
	}
	return false, -1
}

func strToIntSlice(str string) []int {
	var returnSlice []int
	var constr string
	for i := 0; i < len(str); i++ {
		if str[i] == ' ' {
			intArg, err := strconv.Atoi(string(str[i]))
			if err != nil {
				log.Fatal(err)
			}
			returnSlice = append(returnSlice, intArg)
			constr = ""
		} else {
			constr += string(str[i])
		}
	}
	return returnSlice
}

func strToFloatSlice(str string) []float64 {
	var returnSlice []float64
	var constr string
	for i := 0; i < len(str); i++ {
		if str[i] == ' ' {
			constr = ""
		} else {
			constr += string(str[i])
			floatArg, err := strconv.ParseFloat(constr, 64)
			if err != nil {
				continue
			}
			returnSlice = append(returnSlice, floatArg)
		}
	}
	return returnSlice
}

func findSpecByName(name string) (*Spec, error) {
	for i := 0; i < len(programSpecs); i++ {
		if name == programSpecs[i].name {
			return &programSpecs[i], nil
		}
	}
	err := errors.New("No Spec with that name")
	return &Spec{}, err
}

func deleteTask(tasks []Task, index int) []Task {
	return append(tasks[:index], tasks[index+1:]...)
}

func deleteTeam(teams []Team, index int) []Team {
	return append(teams[:index], teams[index+1:]...)
}

// Basic types

type Spec struct {
	name string
	desc string
	id   int
}

type Task struct {
	task         string
	spec         Spec
	body         TaskBody
	id           int
}

type TaskBody struct {
	str string
	num []int
	dec []float64
}

// Worker Type

type Worker struct {
	name  string
	spec  Spec
	tasks []Task
	id    int
}

func (w *Worker) assignTask(task Task) {
	w.tasks = append(w.tasks, task)
}

// Team Type

type Team struct {
	spec           Spec
	workers        []Worker
	unasignedTasks []Task
	id             int
}


func (t *Team) assignWorker(worker Worker) bool {
	if t.spec.name == worker.spec.name {
		t.workers = append(t.workers, worker)
		return true
	}
	return false
}

func (t *Team) distributeTasks() {
	for i := 0; i < len(t.unasignedTasks); i++ {
		busyMap := make(map[int]Worker)
		for _, worker := range t.workers {
			busyMap[len(worker.tasks)] = worker
		}
		var keySlice []int
		for key := range busyMap {
			keySlice = append(keySlice, key)
		}
		least := min(keySlice)
		for j := range t.workers {
			if len(t.workers[j].tasks) == least {
				t.workers[j].tasks = append(t.workers[j].tasks, t.unasignedTasks[i])
				t.unasignedTasks = deleteTask(t.unasignedTasks, i)
				break
			}
		}
	}
}

func (t *Team) assignTaskToTeam(task Task) {
	t.unasignedTasks = append(t.unasignedTasks, task)
	t.distributeTasks()
}

type Controller struct {
	teams             []Team
	allowedTasksSpecs []Spec
}

// Controller Type

func (c *Controller) addTeam(spec string) {
	isAlready := false
	for _, team := range c.teams {
		if team.spec.name == spec {
			isAlready = true
		}
	}
	if !isAlready {
		t := createTeam(spec)
		c.teams = append(c.teams, t)
		isIn := false
		for _, sp := range c.allowedTasksSpecs {
			if sp.name == spec {
				isIn = true
			}
		}
		if !isIn {
			newSpec, err := createSpec(spec)
			if err != nil {
				log.Fatal(err)
			}
			c.allowedTasksSpecs = append(c.allowedTasksSpecs, newSpec)
		}
	} else {
		fmt.Printf("Team with spec %v already exists\n", spec)
	}
}

func (c *Controller) distributeWorker(worker Worker) {
	isOk := false
	for i := range c.teams {
		isOk = c.teams[i].assignWorker(worker)
		if isOk {
			break
		}
	}
	if !isOk {
		fmt.Printf("Worker %v cannot be assigned - no matching teams for spec %v. Create a new team with sufficient Spec requirements\n", worker.name, worker.spec.name)
	}
}

func (c *Controller) addWorker(attrs ...string) {
	w, err := createWorker(attrs...)
	if err != nil {
		log.Fatal(err)
	}
	c.distributeWorker(w)
} 

func (c *Controller) postTask(task Task) {
	isOk := false
	for i := range c.teams {
		if c.teams[i].spec == task.spec {
			c.teams[i].assignTaskToTeam(task)
			isOk = true
		}
	}
	if !isOk {
		fmt.Printf("Couldn't post a task %v to any team. Create a team with spec %v", task.task, task.spec)
	}
}

func (c *Controller) wrapTask(attrs ...string) {
	var (
		taskArgsInt string
		taskArgsStr string
		taskArgsFlo string
		taskArgsIntSlice   []int
		taskArgsFloatSlice []float64
	)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Enter a list of textual arguments for a task: \n")
	scanner.Scan()
	taskArgsStr = scanner.Text()
	fmt.Printf("Enter a list of integer arguments for a task: \n")
	scanner.Scan()
	taskArgsInt = scanner.Text()
	
	if len(taskArgsInt) > 0 {
		taskArgsIntSlice = strToIntSlice(taskArgsInt)
	}
	fmt.Printf("Enter a list of float arguments for a task:\n")
	scanner.Scan()
	taskArgsFlo = scanner.Text()

	if len(taskArgsFlo) > 0 {
		taskArgsFloatSlice = strToFloatSlice(taskArgsFlo)
	}

	body := TaskBody{
		taskArgsStr,
		taskArgsIntSlice,
		taskArgsFloatSlice,
	}
	switch len(attrs) {
	case 2:	
		spec, err := findSpecByName(attrs[0])
		if err != nil {
			log.Fatal(err)
		}
		task := Task{
			task: attrs[1],
			spec: *spec,
			body: body,
			id:   len(programTasks) + 1,
		}
		programTasks = append(programTasks, task)
		c.postTask(task)
	case 1:
		spec, err := findSpecByName(attrs[0])
		if err != nil {
			log.Fatal(err)
		}
		task := Task{
			spec: *spec,
			body: body,
			id:   len(programTasks) + 1,
		}
		programTasks = append(programTasks, task)
		c.postTask(task)
	default:
		fmt.Printf("ERROR: Passed %v arguments posting a Task to the system, whereas only 2 are needed\n", len(attrs))
	}
	fmt.Printf("Task created\n")
}

func (c *Controller) listInsides() {
	fmt.Printf("\nController stats:\n")
	if len(c.teams) > 0 {
		for _, team := range c.teams {
			fmt.Printf("\tTeam: %v\n", team.spec.name)
			if len(team.workers) > 0 {
				for _, worker := range team.workers {
					fmt.Printf("\t\tWorker: %v\n", worker.name)
					if len(worker.tasks) > 0 {
						for _, task := range worker.tasks {
							fmt.Printf("\t\t\tTask: %v\n", task.task)
							fmt.Printf("\t\t\t%v\n", task.body.str)
							fmt.Printf("\t\t\t%v\n", task.body.num)
							fmt.Printf("\t\t\t%v\n", task.body.dec)
						}
					} else {
						fmt.Printf("\t\t\tWorker %v doesn't have assigned tasks yet\n", worker.name)
					}
				}
				fmt.Printf("\t\tTeam %v currently has %v unasigned tasks\n", team.spec.name, len(team.unasignedTasks))
			} else {
				fmt.Printf("\t\tTeam %v doesn't have assigned workers yet\n", team.spec.name)
				fmt.Printf("\t\tTeam %v currently has %v unasigned tasks\n", team.spec.name, len(team.unasignedTasks))
			}
		}
	} else {
		fmt.Printf("\tController is empty\n")
	}
	fmt.Printf("\n")
}

// Operational functions

var programSpecs   []Spec
var programTasks   []Task
var programWorkers []Worker
var programTeams   []Team

func createSpec(attrs ...string) (Spec, error) {
	switch len(attrs) {
	case 2:
		spec := Spec{
			attrs[0],
			attrs[1],
			len(programSpecs) + 1,
		}
		programSpecs = append(programSpecs, spec)

		return spec, nil

	case 1:
		spec := Spec{
			name: attrs[0],
			id:   len(programSpecs) + 1,
		}
		programSpecs = append(programSpecs, spec)

		return spec, nil

	default:
		fmt.Printf("ERROR: Passed %v arguments defining a Spec, whereas only 2 are needed\n", len(attrs))

		return Spec{}, errors.New("Too many args")
	}
}

func createTask(attrs ...string) (Task, error) {
	switch len(attrs) {
	case 2:
		isSpec := false
		var reqSpec Spec
		for _, spec := range programSpecs {
			if spec.name == attrs[1] {
				reqSpec = spec
				isSpec = true
				break
			}
		}
		if isSpec {
			task := Task{
				task: attrs[0],
				spec: reqSpec,
				id:   len(programTasks) + 1,
			}
			programTasks = append(programTasks, task)
			return task, nil
		} else {
			spec, err := createSpec(attrs[1])
			if err != nil {
				log.Fatal(err)
			}
			task := Task{
				task: attrs[0],
				spec: spec,
				id:   len(programTasks) + 1,
			}
			programTasks = append(programTasks, task)
			return task, nil
		}
	case 1:
		task := Task{
			task: attrs[0],
			id:   len(programTasks) + 1,
		}
		programTasks = append(programTasks, task)
		return task, nil
	default:
		fmt.Printf("ERROR: Passed %v arguments defining a Task, whereas only 2 are needed\n", len(attrs))
		return Task{}, errors.New("Too many args")
	}
}

func createWorker(attrs ...string) (Worker, error) {
	switch len(attrs) {
	case 2:
		isSpec := false
		var reqSpec Spec
		for _, spec := range programSpecs {
			if spec.name == attrs[1] {
				reqSpec = spec
				isSpec = true
				break
			}
		}
		if isSpec {
			worker := Worker{
				name: attrs[0],
				spec: reqSpec,
				id:   len(programWorkers) + 1,
			}
			programWorkers = append(programWorkers, worker)
			return worker, nil
		} else {
			spec, err := createSpec(attrs[1])
			if err != nil {
				log.Fatal(err)
			}
			worker := Worker{
				name: attrs[0],
				spec: spec,
				id:   len(programWorkers) + 1,
			}
			programWorkers = append(programWorkers, worker)
			return worker, nil
		}
	case 1:
		worker := Worker{
			name: attrs[0],
			id:   len(programWorkers) + 1,
		}
		programWorkers = append(programWorkers, worker)
		return worker, nil
	default:
		fmt.Printf("ERROR: Passed %v arguments defining a Worker, whereas only 2 are needed\n", len(attrs))
		return Worker{}, errors.New("Too many args")
	}
}

// Unlike other /create/ functions, createTeam is more strict, cause it makes just an empty
// team with only a Spec associated with it.

func createTeam(spec string) Team {
	for _, one := range programSpecs {
		if one.name == spec {
			team := Team{
				spec: one,
				id: len(programTeams) + 1,
			}
			programTeams = append(programTeams, team)
			return team
		}
	}
	one, err := createSpec(spec)
	if err != nil {
		log.Fatal(err)
	}
	team := Team{
		spec: one,
		id: len(programTeams) + 1,
	}
	programTeams = append(programTeams, team)
	return team
}

func createController(attrs ...string) Controller {
	var teams []Team
	var specs []Spec
	for _, team := range attrs {
		oneSpec, err := createSpec(team)
		if err != nil {
			log.Fatal(err)
		}
		oneTeam := createTeam(team)
		teams = append(teams, oneTeam)
		specs = append(specs, oneSpec)
	}
	controller := Controller{teams, specs}
	return controller
}

func generateWorkerSet(idents ...string) []Worker {
	var workers []Worker
	for i := 0; i < len(idents); i += 2 {
		worker, err := createWorker(idents[i], idents[i+1])
		if err != nil {
			log.Fatal(err)
		}
		workers = append(workers, worker)
	}
	return workers
}

func main() {
	running := true
	c := createController()
	scanner := bufio.NewScanner(os.Stdin)
	for running {
		fmt.Printf(">")
		var input string
		fmt.Scanln(&input)
		switch input{
		case "add_team":
			fmt.Printf("Enter a new team's spec: ")
			var spec string
			fmt.Scanln(&spec)
			c.addTeam(spec)
		case "add_worker":
			fmt.Printf("Enter a worker's name: ")
			var name string
			fmt.Scanln(&name)
			fmt.Printf("Enter a worker's spec: ")
			var spec string
			fmt.Scanln(&spec)
			worker, err := createWorker(name, spec)
			if err != nil {
				log.Fatal(err)
			}
			c.distributeWorker(worker)
		case "add_task":
			fmt.Printf("Enter a task's name: ")
			scanner.Scan()
			name := scanner.Text()
			fmt.Printf("Enter a task's spec: ")
			scanner.Scan()
			spec := scanner.Text()
			c.wrapTask(spec, name)
		case "list":
			c.listInsides()
		case "quit":
			running = false
		}
	}
}