package main

import (
	"log"
	"time"
)

func (a *App) apiPoller() {
	for {
		if a.Alpha.checkSession() {
			log.Println("Session not valid, creating new session")
			a.Alpha.getAuthToken()
		}

		s, err := a.Alpha.getRTStats()
		if err != nil {
			log.Println(err)
		}
		if s.Status {
			a.createMetrics(s)
		} else {
			log.Println("Error acquiring stats, will refresh token")
      a.Alpha.getAuthToken()
		}

		time.Sleep(3 * time.Second)

	}
}
