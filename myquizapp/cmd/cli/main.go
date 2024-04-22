//CLI

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type Question struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Choices []string `json:"choices"`
}

type ScoreResponse struct {
	Score      int     `json:"score"`
	Percentile float64 `json:"percentile"`
}

var rootCmd = &cobra.Command{
	Use:   "quiz",
	Short: "Quiz CLI",
	Run: func(cmd *cobra.Command, args []string) {
		startQuiz()
	},
}

func startQuiz() {
	questions := fetchQuestions()
	answers := make([]int, len(questions))
	for i, q := range questions {
		fmt.Printf("Question %d: %s\n", i+1, q.Text)
		for j, choice := range q.Choices {
			fmt.Printf("%d: %s\n", j+1, choice)
		}
		fmt.Printf("Your answer: ")
		fmt.Scan(&answers[i])
		answers[i]-- // Adjust for zero-based index
	}
	submitAnswers(answers)
}

func fetchQuestions() []Question {
	resp, err := http.Get("http://localhost:8080/quiz")
	if err != nil {
		fmt.Println("Error fetching questions:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var questions []Question
	json.Unmarshal(body, &questions)
	return questions
}

func submitAnswers(answers []int) {
	data, _ := json.Marshal(map[string][]int{"answers": answers})
	resp, err := http.Post("http://localhost:8080/submit", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error submitting answers:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var scoreResp ScoreResponse
	json.Unmarshal(body, &scoreResp)
	fmt.Printf("You scored %d/%d. You were better than %.2f%% of all quizzers.\n", scoreResp.Score, len(answers), scoreResp.Percentile)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
