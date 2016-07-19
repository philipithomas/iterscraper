package main

import "fmt"

func emitTasks(taskChan chan site) {
	for i := idLow; i < idHigh; i++ {
		url := fmt.Sprintf(urlBase, i)
		taskChan <- site{url: url, id: i}
	}
}
