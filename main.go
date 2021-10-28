package main

import (
	log "github.com/sirupsen/logrus"
	"quizGO/api"
)

func main() {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	actors := newActorGroup()
	{
		quizAPI := api.NewAPI(logger)
		adminService := api.NewAdminService(logger)
		questionService := api.NewQuestionService(logger)

		// set routes for admin related endpoints
		quizAPI.Router = adminService.Route(quizAPI.Router)
		// set routes for question related endpoints
		quizAPI.Router = questionService.Route(quizAPI.Router)

		runFunc := func() error { return quizAPI.Run() }
		cancelFunc := func() error { return quizAPI.Close() }

		actor := newActor("api", runFunc, cancelFunc)
		actors.add(actor)
	}
	go actors.run()
	err := startQuiz()
	if err != nil {
		logger.Errorf("startQuiz: %v", err)
	}
}
