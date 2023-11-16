package main

import (
	"fmt"
	"time"
)

const defaultSessionDuration time.Duration = time.Duration(1 * time.Minute)

type session struct {
	start    time.Time
	duration int
	category string
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
	return session{start: startTime, duration: int(duration), category: category}
}

func main() {
	session := runSession(5*time.Second, "unknown")
	fmt.Print(session)
}
