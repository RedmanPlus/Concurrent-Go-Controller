package shell

import (
	"os"
	"fmt"
	"bufio"
	"log"
	"ctc/backends/controller"
)

func shell(c controller.Controller) {
	running := true
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
			worker, err := controller.createWorker(name, spec)
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

func RunShell(c controller.Controller) {
	go shell(c)
}