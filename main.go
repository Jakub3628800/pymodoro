package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const defaultSessionDuration time.Duration = time.Duration(1 * time.Minute)

type session struct {
	Start    time.Time `json:"start"`
	Duration int       `json:"duration"`
	Category string    `json:"category"`
}

func printElapsed(d time.Duration) {
	fmt.Printf("\033[1A\033[K")
	fmt.Println(d.Truncate(1 * time.Second))
}

func runSession(duration time.Duration, category string) session {
	startTime := time.Now()
	elapsed := time.Since(startTime)

	fmt.Println("=============================")
	for elapsed < time.Duration(duration) {

		printElapsed(elapsed)
		time.Sleep(100 * time.Millisecond)
		elapsed = time.Since(startTime)
	}
	return session{Start: startTime, Duration: int(duration.Minutes()), Category: category}
}

func loadSessions(filepath string) (string, error) {
	content, err := ioutil.ReadFile(filepath)
	return string(content), err
}

func saveSessions(filepath string, content string) error {
	f, err := os.Open(filepath)
	defer f.Close()
	if err == nil {
		return err
	}
	f.WriteString(content)
	return nil
}

func main() {
	s, err := loadSessions("sessions.json")
	if err != nil {
		fmt.Println("error loading sessions")
	}
	fmt.Println(s)
	allSessions := []session{}
	json.Unmarshal([]byte(s), &allSessions)
	fmt.Println(allSessions)

	session := runSession(5*time.Second, "unknown")
	allSessions = append(allSessions, session)
	fmt.Println(allSessions)

	sessionsJson, _ := json.Marshal(allSessions)
	saveSessions("sessions.json", string(sessionsJson))

}
