package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"quizGO/api"
	"quizGO/config"
	"strings"
	"text/template"

	"github.com/valyala/fasthttp"
)

func includeTemplates() (*template.Template, error) {
	// for some reason templates\\* does not work on Windows :)
	templates, err := template.ParseFiles("templates\\title", "templates\\welcome", "templates\\question", "templates\\end")
	if err != nil {
		return nil, err
	}

	if err = templates.ExecuteTemplate(os.Stdout, "title", nil); err != nil {
		return nil, err
	}

	if err = templates.ExecuteTemplate(os.Stdout, "welcome", nil); err != nil {
		return nil, err
	}
	return templates, nil
}

func getQuestions() ([]api.Question, error) {
	res, err := http.Get(fmt.Sprintf("%s://%s:%d/api/%s", config.Protocol, config.Host, config.Port, "v1"))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != fasthttp.StatusOK {
		return nil, errors.New("bad request")
	}
	defer res.Body.Close()

	var questions []api.Question
	if err = json.NewDecoder(res.Body).Decode(&questions); err != nil {
		return nil, err
	}
	return questions, nil
}

func showResult(templates *template.Template, username string, score int) error {
	result := struct {
		Username string
		Score    int
	}{
		Username: username,
		Score:    score,
	}

	return templates.ExecuteTemplate(os.Stdout, "end", result)
}

func startQuiz() error {
	templates, err := includeTemplates()
	if err != nil {
		return err
	}

	var username string
	_, err = fmt.Scan(&username)
	if err != nil {
		return err
	}

	questions, err := getQuestions()
	if err != nil {
		return err
	}

	var score int
	for _, q := range questions {
		if err = templates.ExecuteTemplate(os.Stdout, "question", q); err != nil {
			return err
		}
		var answer string
		if _, err = fmt.Scan(&answer); err != nil {
			return err
		}
		if q.CorrectAnswer != strings.ToLower(answer) {
			continue
		}
		score++
	}

	return showResult(templates, username, score)
}
