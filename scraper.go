package main

func scrape(taskChan, dataChan chan site) {
	for {
		site := <-taskChan
		site.fetch()
		dataChan <- site
	}
}
