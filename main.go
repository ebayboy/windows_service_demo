package main

import (
	"log"
	"time"

	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() {
	// Do work here
	for {
		log.Println("Service is running...")
		time.Sleep(time.Minute)
	}
}

func (p *program) Stop(s service.Service) error {
	// Here you can decide what happens when the service is requested to stop.
	// To ignore the stop request, you can simply return nil and not actually stop the service.
	log.Println("Service stop request ignored.")
	return nil // Ignore stop request
}

func main() {
	config := &service.Config{
		Name:        "GoServiceExample",
		DisplayName: "Go Service Example",
		Description: "This is an example Go service that ignores stop requests.",
	}

	prg := &program{}
	s, err := service.New(prg, config)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
