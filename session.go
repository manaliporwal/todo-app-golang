package main

import (
	"fmt"
	"crypto/rand"
)
type Session struct {
	userName string
	expiry	string
}

var sessionCache map[string]Session

func initSessionCache() {
	sessionCache = make(map[string]Session)
}

func generateUniqueSessionId() (string, error){
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
	        b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	fmt.Println(uuid)
	return uuid, nil
}
func addSession(userName, expiry string) (string, error){
	//session := Session{sessionID, expiry}
	//Generate an unique sessionID
	sessionID, err := generateUniqueSessionId()
	if err == nil {
		sessionCache[sessionID] =  Session{userName, expiry}
		return sessionID, nil
	}
	return "", err
}

func validateSession(sessionID string) bool {
	if _, ok := sessionCache[sessionID]; ok {
		return true
	}
	return false
}

func getUserNameFromSession(sessionID string) (string, bool) {
	if session, ok := sessionCache[sessionID]; ok {
		return session.userName, true
	}
	return "", false
}

