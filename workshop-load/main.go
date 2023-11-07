package main

import (
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"workshop-server/workshop"
)

const (
	ApplicationName = "workshop-load"
	EnvVarEndpoint  = "SERVER_ENDPOINTS"
	EnvVarLoadUsers = "LOAD_USERS"
)

var (
	MinSleepMs         = 100
	MaxSleepMs         = 1900
	KeepRunning        = true
	GRand              *rand.Rand
	MaxConcurrentUsers = 2
)

func main() {
	slog.Info(fmt.Sprintf("Running %s ...", ApplicationName))

	endpoints := []string{"localhost:8080"}
	if str, exists := os.LookupEnv(EnvVarEndpoint); exists {
		slog.Info("Using endpoints from env", "endpoint", str)
		endpoints = strings.Split(str, ",")
		slog.Info(fmt.Sprintf("Found %d endpoints to target", len(endpoints)))
	}

	if str, exists := os.LookupEnv(EnvVarLoadUsers); exists {
		slog.Info("Using load user number from env", "users", str)
		var err error
		if MaxConcurrentUsers, err = strconv.Atoi(str); err != nil {
			log.Fatal(err)
		}
	}

	GRand = rand.New(rand.NewSource(1337))

	for i := 0; i < MaxConcurrentUsers; i++ {
		go doUserRequests(endpoints[i%len(endpoints)], i+1)
	}

	// Add static nginx endpoint to see nginx logs as well
	go doUserRequests("workshop-example-nginx-1", 99)

	for {
		time.Sleep(250 * time.Millisecond)
	}
}

func doUserRequests(endpoint string, userId int) {
	client := workshop.NewClient(endpoint)
	for KeepRunning {
		switch GRand.Intn(4) {
		case 1:
			slog.Info(fmt.Sprintf("User %d calling login on %s", userId, endpoint))
			if err := client.Login(); err != nil {
				log.Fatal(err)
			}
		case 2:
			slog.Info(fmt.Sprintf("User %d calling logout on %s", userId, endpoint))
			if err := client.Logout(); err != nil {
				log.Fatal(err)
			}
		case 3:
			slog.Info(fmt.Sprintf("User %d calling getProject on %s", userId, endpoint))
			if err := client.GetProject(); err != nil {
				log.Fatal(err)
			}
		case 4:
			slog.Info(fmt.Sprintf("User %d calling updateProject on %s", userId, endpoint))
			if err := client.UpdateProject(); err != nil {
				log.Fatal(err)
			}

		}

		time.Sleep(time.Duration(MinSleepMs+GRand.Intn(MaxSleepMs)) * time.Millisecond)
	}
}
