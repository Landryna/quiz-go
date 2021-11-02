# quiz-go

Simple program that allows you to modify questions via REST API
and play the game via CMD

### Using the quiz-go

#### After staring the server, you need to create a question/questions
```bigquery
curl --location --request POST 'localhost:8000/api/v1/question' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "Best programming language?",
    "answers": {
        "a":"Go",
        "b":"Java",
        "c":"C#",
        "d":"Python"
    },
    "correct_answer": "a"
}
```
### Start the CMD running using ```go run``` and follow the instructions provided in the main loop :stuck_out_tongue_winking_eye:
## HF recruiter :smile:

