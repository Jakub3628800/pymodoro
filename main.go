package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

const defaultSessionDuration time.Duration = time.Duration(5 * time.Second)

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

func loadSessions(filename string) ([]session, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sessions []session
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&sessions)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func saveSessions(filename string, sessions []session) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(sessions)
	if err != nil {
		return err
	}

	return nil
}

func sendNotification(msg string) {
	err := exec.Command("notify-send", msg).Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	allSessions, err := loadSessions("sessions.json")
	if err != nil {
		log.Fatal(err)
	}

	duration := defaultSessionDuration
	sendNotification(fmt.Sprintf("Session started\n %s", duration.String()))
	session := runSession(duration, "unknown")
	sendNotification(fmt.Sprintf("Session ended\n %s", duration.String()))

	allSessions = append(allSessions, session)
	fmt.Println(allSessions)

	//sessionsJson, _ := json.Marshal(allSessions)
	//fmt.Println(string(sessionsJson))
	saveSessions("sessions.json", allSessions)

}
