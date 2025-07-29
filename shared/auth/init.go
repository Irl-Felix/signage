package auth

import "log"

func Init() {
	if err := InitDB(); err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
}
