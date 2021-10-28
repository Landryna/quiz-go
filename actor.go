package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type actorGroup []actor
type runFunc func() error
type cancelFunc func() error

var actors actorGroup

type actor struct {
	Name   string
	Run    runFunc
	Cancel cancelFunc
}

func newActor(name string, rf runFunc, cf cancelFunc) *actor {
	return &actor{
		Name:   name,
		Run:    rf,
		Cancel: cf,
	}
}

func newActorGroup() *actorGroup {
	return &actorGroup{}
}

func (*actorGroup) add(a *actor) {
	actors = append(actors, *a)
}

func (*actorGroup) show() actorGroup {
	return actors
}

// run executes all provided actors and waits for first error - then cancels every actor
func (a *actorGroup) run() {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})

	errorCh := make(chan error)
	actorCh := make(chan string)
	for _, a := range actors {
		a := a
		go func() {
			err := a.Run()
			errorCh <- err
			actorCh <- a.Name
		}()
	}

	// wait for first error
	err := <-errorCh
	actor := <-actorCh
	logger.Errorf("actor %s: %v", actor, err)

	for _, a := range actors {
		a := a
		err := a.Cancel()
		if err != nil {
			logger.Errorf("actor cancel %s: %v", actor, err)
			os.Exit(1)
		}
	}
}
