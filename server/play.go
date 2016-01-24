package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func playHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseGameID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Println("taking action id", id)
	action, err := parseAction(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	g, ok := gameMap[id]
	if !ok {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("id %d is not registered", id))
		return
	}
	switch action {
	case "state":
		g.stateHandler(w, r)
	case "move":
		log.Println("move action")
		g.moveHandler(w, r, id)
	case "wait":
		g.waitHandler(w, r, id)
	default:
		writeError(w, http.StatusBadRequest, fmt.Sprintf("%s is not a valid play action", action))
	}
}

func parseGameID(r *http.Request) (GameID, error) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.SplitN(path, "/", 2)
	id, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("id %s error: %s", parts[0], err.Error())
	}
	return GameID(id), nil
}

func parseAction(r *http.Request) (string, error) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("Missing play action")
	}
	return parts[1], nil
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	w.Write([]byte(message))
	log.Println(message)
}
