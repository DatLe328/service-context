package core

import "log"

func Recover() {
	if err := recover(); err != nil {
		log.Println("Recovery error:", err)
	}
}
