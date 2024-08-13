package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

const defaultSessionDuration time.Duration = time.Duration(25 * time.Minute)

type session struct {
	Start    time.Time `json:"start"`
	Duration int       `json:"duration"`
	Category string    `json:"category"`
}

func printElapsed(d time.Duration) {
	fmt.Printf("\033[1A\033[K")
	fmt.Println(d.Truncate(1 * time.Second))
}

func runSession(duration time.Duration, category string, timerEnabled bool) session {
	startTime := time.Now()
	elapsed := time.Since(startTime)

	if timerEnabled {
		fmt.Println("=============================")
		fmt.Println("=============================")
	}
	for elapsed < time.Duration(duration) {

		if timerEnabled {
			printElapsed(elapsed)
		}
		time.Sleep(100 * time.Millisecond)
		elapsed = time.Since(startTime)
	}
	return session{Start: startTime, Duration: int(duration.Minutes()), Category: category}
}

func loadSessions(filename string) ([]session, error) {

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sessions []session
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&sessions)
	if err != nil && err != io.EOF {
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

func sendNotification(msg string, silent bool) {
	if silent {
		fmt.Println(msg)
		fmt.Println()
	}
	err := exec.Command("notify-send", msg).Run()
	if err != nil {
		log.Fatal(err)
	}
}
