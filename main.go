package main

import (
	"fmt"
	"time"
	"ctc/backends/controller"
)

func main() {
	json := "{name: jim, skill: Not any}"
	c := ctc.CreateController()
	c.AddTeam("House")
	w, err := ctc.CreateWorker("John", "House")
	if err != nil {
		return
	}
	w.WorkerFunc = func(input []byte) {
		for _, one := range input {
			fmt.Println(string(one))
			time.Sleep(1 * time.Second)
		}
	}
	c.DistributeWorker(w)
	c.ListInsides()
	i := 0
	for {
		taskName := "#" + string(i)
		t, err := ctc.CreateTask(taskName, "House")
		if err != nil {
			return
		}
		t.AcceptParams(json)
		c.ReorganizeTeams(5)
		c.PostTask(t)
		c.ListInsides()
		time.Sleep(1 * time.Second)
		if i >= 1000000 {
			break
		}
		i++
	}
}