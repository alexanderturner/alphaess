package main

import (
	"log"
	"time"
)

func (a *App) ApiPoller() {
	for {
		if a.Alpha.CheckSession() {
			log.Println("Session not valid, creating new session")
			a.Alpha.GetAuthToken()
		}

		s, err := a.Alpha.GetRTStats()
		if err != nil {
			log.Println(err)
		}
		if s.Status {
			a.CreateMetrics(s)
		} else {
			log.Println("Error acquiring stats, will refresh token")
			a.Alpha.GetAuthToken()
		}

		time.Sleep(3 * time.Second)

	}
}
