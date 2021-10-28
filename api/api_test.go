package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"net/http"
	"quizGO/api"
	"testing"
	"time"
)

const URI = "http://localhost:8000/api/"

func initializeTestBackground(t *testing.T) {
	logger := logrus.New()

	quizAPI := api.NewAPI(logger)
	adminService := api.NewAdminService(logger)
	questionService := api.NewQuestionService(logger)

	// set routes for admin related endpoints
	quizAPI.Router = adminService.Route(quizAPI.Router)
	// set routes for question related endpoints
	quizAPI.Router = questionService.Route(quizAPI.Router)

	errCh := make(chan error)
	go func() {
		if err := quizAPI.Run(); err != nil {
			logger.Error(err)
			errCh <- err
		}
		defer quizAPI.Close()
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	var counter int
	for {
		// tick every 1 second
		select {
		case <-ticker.C:
			counter++
		}
		// see if errCh received error, if so fail test now
		select {
		case <-errCh:
			t.FailNow()
		default:
		}
		// wait for errCh for 3 seconds
		if counter == 3 {
			break
		}
	}
}

func TestQuizAPI(t *testing.T) {
	initializeTestBackground(t)

	type want struct {
		statusCode int
		question   api.Question
	}
	type params struct {
		question api.Question
	}
	tests := []struct {
		name   string
		want   want
		params params
		mock   func(params params, t *testing.T) *http.Response
	}{
		{
			name: "successful create question",
			want: want{
				statusCode: 200,
				question:   api.Question{},
			},
			params: params{
				question: api.Question{
					Title: "simple question",
				},
			},
			mock: func(params params, t *testing.T) *http.Response {
				url := fmt.Sprintf("%s/%s", URI, "v1/question")
				body, err := json.Marshal(params.question)
				require.Nil(t, err)

				resp, err := http.Post(url, "aplication/json", bytes.NewReader(body))
				require.Nil(t, err)
				defer resp.Body.Close()
				return resp
			},
		},
		{
			name: "successful create and get question",
			want: want{
				statusCode: 200,
				question: api.Question{
					Title: "simple question",
				},
			},
			params: params{
				question: api.Question{
					Title: "simple question",
				},
			},
			mock: func(params params, t *testing.T) *http.Response {
				// create question
				url := fmt.Sprintf("%s/%s", URI, "v1/question")
				body, err := json.Marshal(params.question)
				require.Nil(t, err)

				resp, err := http.Post(url, "application/json", bytes.NewReader(body))
				require.Nil(t, err)
				defer resp.Body.Close()

				// get question
				url = fmt.Sprintf("%s/%s", URI, "v1/question/0")
				resp, err = http.Get(url)
				require.Nil(t, err)

				return resp
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mock(tt.params, t)
			require.Equal(t, tt.want.statusCode, got.StatusCode)
			// we have declared the question as expected output
			if len(tt.want.question.Title) > 0 {
				var question api.Question
				err := json.NewDecoder(got.Body).Decode(&question)
				require.Nil(t, err)
				require.Equal(t, tt.want.question.Id, question.Id)
				require.Equal(t, tt.want.question.Title, question.Title)
			}
		})
	}
}
