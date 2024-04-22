//Server

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Question struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Choices []string `json:"choices"`
	Answer  int      `json:"-"` // Index of the correct answer, hidden from JSON
}

type Submission struct {
	Answers []int `json:"answers"` // Indices of user's answers
}

type ScoreResponse struct {
	Score      int     `json:"score"`
	Percentile float64 `json:"percentile"`
}

var (
	questions = []Question{
		{ID: 1, Text: "What is the capital of Japan?", Choices: []string{"Seoul", "Beijing", "Tokyo", "Bangkok"}, Answer: 2},
		{ID: 2, Text: "Which planet is known as the Red Planet?", Choices: []string{"Venus", "Mars", "Jupiter", "Saturn"}, Answer: 1},
		{ID: 3, Text: "What is the largest ocean on Earth?", Choices: []string{"Atlantic Ocean", "Indian Ocean", "Arctic Ocean", "Pacific Ocean"}, Answer: 3},
		{ID: 4, Text: "What gas do plants absorb from the atmosphere to perform photosynthesis?", Choices: []string{"Carbon Dioxide", "Oxygen", "Nitrogen", "Hydrogen"}, Answer: 0},
		{ID: 5, Text: "Who wrote the Harry Potter series?", Choices: []string{"J.R.R. Tolkien", "J.K. Rowling", "Stephen King", "Suzanne Collins"}, Answer: 1},
	}
	scoreBoard []int // Store scores for percentile calculation
	mutex      sync.Mutex
)

func getQuiz(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(questions)
}

func submitQuiz(w http.ResponseWriter, r *http.Request) {
	var sub Submission
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	score, err := calculateScore(sub.Answers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mutex.Lock()
	scoreBoard = append(scoreBoard, score)
	mutex.Unlock()
	percentile := calculatePercentile(score)
	resp := ScoreResponse{Score: score, Percentile: percentile}
	json.NewEncoder(w).Encode(resp)
}

func calculateScore(userAnswers []int) (int, error) {
	if len(userAnswers) != len(questions) {
		return 0, fmt.Errorf("invalid number of answers")
	}
	score := 0
	for i, answer := range userAnswers {
		if answer == questions[i].Answer {
			score++
		}
	}
	return score, nil
}

func calculatePercentile(currentScore int) float64 {
	mutex.Lock()
	defer mutex.Unlock()
	var countBetter int
	for _, score := range scoreBoard {
		if currentScore > score {
			countBetter++
		}
	}
	return float64(countBetter) / float64(len(scoreBoard)) * 100
}

func main() {
	http.HandleFunc("/quiz", getQuiz)
	http.HandleFunc("/submit", submitQuiz)
	http.ListenAndServe(":8080", nil)
}
