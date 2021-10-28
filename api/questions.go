package api

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

var _ Service = &questionService{}

type questionService struct {
	logger *logrus.Logger
}

type indexType struct {
	sync.Mutex
	counter int
}

var questionIndex = &indexType{counter: 0}

type Answers struct {
	A string `json:"a"`
	B string `json:"b"`
	C string `json:"c"`
	D string `json:"d"`
}

type Question struct {
	Id            int     `json:"id"`
	CreatedAt     int64   `json:"created_at"`
	UpdatedAt     int64   `json:"updated_at"`
	CreatedBy     string  `json:"created_by"`
	Title         string  `json:"title"`
	Answers       Answers `json:"answers"`
	CorrectAnswer string  `json:"correct_answer"`
}

type questions struct {
	sync.Mutex
	questionsList []Question
}

var questionsContainer = &questions{}

func NewQuestionService(log *logrus.Logger) *questionService {
	return &questionService{
		logger: log,
	}
}

func (i *indexType) increment() {
	i.Lock()
	defer i.Unlock()

	i.counter++
}

func (i *indexType) get() int {
	i.Lock()
	defer i.Unlock()

	return i.counter
}

func (qs *questions) addQuestion(q Question) {
	qs.Lock()
	defer qs.Unlock()

	qs.questionsList = append(qs.questionsList, q)
}

func (qs *questions) getQuestion(id int) Question {
	qs.Lock()
	defer qs.Unlock()

	for _, q := range qs.questionsList {
		if q.Id == id {
			return q
		}
	}
	return Question{}
}

func (q *questionService) Route(r *router.Router) *router.Router {
	r.GET("/api/v1/questions", q.listQuestions)
	r.GET("/api/v1/question/{id}", q.getQuestion)
	r.POST("/api/v1/question", q.addQuestion)
	return r
}

func (q *questionService) listQuestions(ctx *fasthttp.RequestCtx) {
	questions, err := json.Marshal(questionsContainer.questionsList)
	if err != nil {
		q.logger.WithError(err).Error("listQuestions: marshal")
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.SetBody(questions)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (q *questionService) getQuestion(ctx *fasthttp.RequestCtx) {
	id, ok := ctx.UserValue("id").(string)
	if !ok {
		q.logger.Error("getQuestion: id")
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}
	questionID, err := strconv.Atoi(id)
	if err != nil {
		q.logger.Errorf("getQuestion: %v", err)
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}
	question := questionsContainer.getQuestion(questionID)
	if question.Title == "" {
		q.logger.Errorf("getQuestion: question not found")
		ctx.Error("not found", fasthttp.StatusNotFound)
		return
	}
	body, err := json.Marshal(question)
	if err != nil {
		q.logger.WithError(err).Error("getQuestion: marshal")
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetBody(body)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (q *questionService) addQuestion(ctx *fasthttp.RequestCtx) {
	var question Question
	if err := decode(ctx.PostBody(), &question); err != nil {
		q.logger.WithError(err).Error("addQuestion: invalid payload")
		ctx.Error("invalid payload", fasthttp.StatusBadRequest)
		return
	}

	question.CreatedAt, question.UpdatedAt = time.Now().Unix(), time.Now().Unix()
	question.Id = questionIndex.get()
	questionIndex.increment()
	questionsContainer.addQuestion(question)
}
