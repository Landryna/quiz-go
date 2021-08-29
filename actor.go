package main

import "fmt"

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

func (*actorGroup) run() {
	error := make(chan error)
	//TODO: zastanowic sie czy potrzebuje nazwy aktora
	actor := make(chan string)
	//TODO: moze jakies handlowanie panici?
	for _, a := range actors {
		a := a
		go func() {
			err := a.Run()
			error <- err
			actor <- a.Name
		}()
	}

	for {
		select {
		case err := <-error:
			actor := <-actor
			fmt.Printf("\nActor %s ended with msg: %s", actor, err.Error())
		}
		break
	}

	for _, a := range actors {
		a := a
		err := a.Cancel()
		if err != nil {
			fmt.Printf("unable to cancel actor: %s", a.Name)
		}
	}
}
