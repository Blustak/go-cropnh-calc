package main

import (
	"fmt"
	"os"

	"github.com/Blustak/go-cropnh-calc/internal/server"
	"github.com/joho/godotenv"
)

const (
	envDbString  string = "DB_STRING"
	envLogString string = "LOG_FILE"
	envHostname  string = "SERVER_HOSTNAME"
	envPort      string = "SERVER_PORT"
)

func main() {
	godotenv.Load()
	var setupOpts struct {
		DatabasePath string
		Hostname     string
		Port         string
		LogFilePath  *string
		LogIsJSON    bool
	} = struct {
		DatabasePath string
		Hostname     string
		Port         string
		LogFilePath  *string
		LogIsJSON    bool
	}{}
	setupOpts.DatabasePath = assertEnv(envDbString)
	setupOpts.Hostname = assertEnv(envHostname)
	setupOpts.Port = assertEnv(envPort)
    setupOpts.LogFilePath = envLookupOrFallback(envLogString,nil)
    setupOpts.LogIsJSON = false

	server.InitServer(&setupOpts)

	server.ListenAndServe()
}

func assertEnv(envKey string) string {
	s, ok := os.LookupEnv(envKey)
	if !ok {
		panic(fmt.Sprintf("env var %s not set", envKey))
	}
	return s
}

func envLookupOrFallback(envKey string, fallback *string) *string {
    sp := new(string)
    s, ok := os.LookupEnv(envKey)
    if !ok {
        return fallback
    }
    *sp = s
    return sp
}
