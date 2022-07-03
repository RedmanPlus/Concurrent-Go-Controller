package ctc

import (
	"fmt"
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

func strToByteSlice(str string) []byte {
	var returnSlice []byte
	for i := 0; i < len(str); i++ {
		newByte := byte(str[i])
		returnSlice = append(returnSlice, newByte)
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

// Spec and Task Types

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
	data []byte
}

// Task operating logic

func (b *TaskBody) AddData(data byte) {
	b.data = append(b.data, data)
}

func (b TaskBody) DumpData() []byte {
	return b.data
}

func (b TaskBody) IterateData(iterator func(b byte) byte) []byte {
	var newData []byte
	for i := 0; i < len(b.data); i++ {
		newByte := iterator(b.data[i])
		newData = append(newData, newByte)
	}
	return newData
}

func (t *Task) AcceptParams(params ...string) {
	for i := 0; i < len(params); i++ {
		bytes := strToByteSlice(params[i])
		for _, one := range bytes {
			t.body.AddData(one)
		}
	}
}

// Worker Type

type Worker struct {
	name       string
	spec       Spec
	tasks      []Task
	id         int
	WorkerFunc func([]byte)
}

func (w *Worker) AssignTask(task Task) {
	w.tasks = append(w.tasks, task)
	go w.SolveTask(task)
}

func (w *Worker) SolveTask(task Task) {
	w.WorkerFunc(task.body.data)
	for i := 0; i < len(w.tasks); i++ {
		if w.tasks[i].id == task.id {
			w.tasks = deleteTask(w.tasks, i)
		}
	}
}

// Team Type

type Team struct {
	spec           Spec
	workers        []Worker
	unasignedTasks []Task
	id             int
}


func (t *Team) AssignWorker(worker Worker) bool {
	if t.spec.name == worker.spec.name {
		t.workers = append(t.workers, worker)
		return true
	}
	return false
}

func (t *Team) DuplicateWorker(worker Worker, threshold int) {
	duplicate := worker
	duplicate.tasks = []Task {}
	t.workers = append(t.workers, duplicate)
}

func (t *Team) DistributeTasks() {
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

func (t *Team) AssignTaskToTeam(task Task) {
	t.unasignedTasks = append(t.unasignedTasks, task)
	t.DistributeTasks()
}

func (t *Team) ManageTeamBisyness(threshold int) {
	for i := 0; i < len(t.workers); i++ {
		if len(t.workers[i].tasks) >= threshold {
			t.DuplicateWorker(t.workers[i], threshold)
		}
	}
}

type Controller struct {
	teams             []Team
	allowedTasksSpecs []Spec
}

// Controller Type

func (c *Controller) AddTeam(spec string) {
	isAlready := false
	for _, team := range c.teams {
		if team.spec.name == spec {
			isAlready = true
		}
	}
	if !isAlready {
		t := CreateTeam(spec)
		c.teams = append(c.teams, t)
		isIn := false
		for _, sp := range c.allowedTasksSpecs {
			if sp.name == spec {
				isIn = true
			}
		}
		if !isIn {
			newSpec, err := CreateSpec(spec)
			if err != nil {
				log.Fatal(err)
			}
			c.allowedTasksSpecs = append(c.allowedTasksSpecs, newSpec)
		}
	} else {
		fmt.Printf("Team with spec %v already exists\n", spec)
	}
}

func (c *Controller) ReorganizeTeams(threshold int) {
	for i := 0; i < len(c.teams); i++ {
		c.teams[i].ManageTeamBisyness(threshold)
	}
} 

func (c *Controller) DistributeWorker(worker Worker) {
	isOk := false
	for i := range c.teams {
		isOk = c.teams[i].AssignWorker(worker)
		if isOk {
			break
		}
	}
	if !isOk {
		fmt.Printf("Worker %v cannot be Assigned - no matching teams for spec %v. Create a new team with sufficient Spec requirements\n", worker.name, worker.spec.name)
	}
}

func (c *Controller) AddWorker(attrs ...string) {
	w, err := CreateWorker(attrs...)
	if err != nil {
		log.Fatal(err)
	}
	c.DistributeWorker(w)
} 

func (c *Controller) PostTask(task Task) {
	isOk := false
	for i := range c.teams {
		if c.teams[i].spec.name == task.spec.name {
			c.teams[i].AssignTaskToTeam(task)
			isOk = true
		}
	}
	if !isOk {
		fmt.Printf("Couldn't Post a task %v to any team. Create a team with spec %v", task.task, task.spec)
	}
}

func (c *Controller) ListInsides() {
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
							fmt.Printf("\t\t\t%v\n", task.body)
						}
					} else {
						fmt.Printf("\t\t\tWorker %v doesn't have Assigned tasks yet\n", worker.name)
					}
				}
				fmt.Printf("\t\tTeam %v currently has %v unasigned tasks\n", team.spec.name, len(team.unasignedTasks))
			} else {
				fmt.Printf("\t\tTeam %v doesn't have Assigned workers yet\n", team.spec.name)
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

func CreateSpec(attrs ...string) (Spec, error) {
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

func CreateTask(attrs ...string) (Task, error) {
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
			spec, err := CreateSpec(attrs[1])
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

func CreateWorker(attrs ...string) (Worker, error) {
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
			spec, err := CreateSpec(attrs[1])
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

// Unlike other /Create/ functions, CreateTeam is more strict, cause it makes just an empty
// team with only a Spec associated with it.

func CreateTeam(spec string) Team {
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
	one, err := CreateSpec(spec)
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

func CreateController(attrs ...string) Controller {
	var teams []Team
	var specs []Spec
	for _, team := range attrs {
		oneSpec, err := CreateSpec(team)
		if err != nil {
			log.Fatal(err)
		}
		oneTeam := CreateTeam(team)
		teams = append(teams, oneTeam)
		specs = append(specs, oneSpec)
	}
	controller := Controller{teams, specs}
	return controller
}