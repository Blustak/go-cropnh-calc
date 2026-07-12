package server

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Blustak/go-cropnh-calc/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

var Server struct {
	Hostname, Port string
	DB             *database.Queries
	Log            *slog.Logger
	logFile        *os.File
	dbFile         *sql.DB
}

var serverInit = false

var defaultOpts = func() struct {
	DatabasePath, Hostname, Port string
	LogFilePath                  *string
	LogIsJSON                    bool
} {
	wdPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return struct {
		DatabasePath string
		Hostname     string
		Port         string
		LogFilePath  *string
		LogIsJSON    bool
	}{
		DatabasePath: wdPath + "data/data.db",
		LogFilePath:  nil,
		LogIsJSON:    false,
		Hostname:     "localhost",
		Port:         "8080",
	}
}()

func InitServer(opts *struct {
	DatabasePath, Hostname, Port string
	LogFilePath                  *string
	LogIsJSON                    bool
}) {
	if serverInit {
		panic("double initialisation!")
	}
	defer func() { serverInit = true }()
	if opts == nil {
		opts = &defaultOpts
	}

	var logFile *os.File
	var logHandler slog.Handler
	var err error

	if opts.LogFilePath == nil {
		logFile = os.Stdout
        logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelDebug,
        })
	} else {
		logFile, err = os.OpenFile(*opts.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("error opening log file: %v\n", err))
		}
        logHandler = slog.NewMultiHandler(
            slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}),
            func() slog.Handler {
                if opts.LogIsJSON {
                    return slog.NewJSONHandler(logFile, &slog.HandlerOptions{
                        Level:slog.LevelDebug,
                    })
                } else {
                    return slog.NewTextHandler(logFile, &slog.HandlerOptions{Level:slog.LevelDebug})
                }
            }(),
        )
	}
	serverLog := slog.New(logHandler)

	dbCon, err := sql.Open("sqlite3", opts.DatabasePath)
	if err != nil {
		panic(err)
	}
	Server = struct {
		Hostname string
		Port     string
		DB       *database.Queries
		Log      *slog.Logger
		logFile  *os.File
		dbFile   *sql.DB
	}{
		Hostname: opts.Hostname,
		Port:     opts.Port,
		DB:       database.New(dbCon),
		dbFile:   dbCon,
		Log:      serverLog,
		logFile:  logFile,
	}
}

func DeinitServer() {
	Server.dbFile.Close()
	Server.logFile.Close()
}

func ListenAndServe() error {
	err := listenAndServe()
	if err != nil {
		Server.Log.Error("error", slog.Any("err", err), slog.Any("server", Server))
	}
	return err
}

func listenAndServe() error {
	serveMux := http.NewServeMux()
	if serveMux == nil {
		return fmt.Errorf("serveMux nil")
	}
    if err := registerPaths(serveMux); err != nil {
        panic(err)
    }
    serve := http.Server{
        Addr: Server.Hostname + ":" + Server.Port,
        Handler: serveMux,
    }
    fmt.Printf("server started on %s\n", serve.Addr)
    Server.Log.Info("started listening", slog.Any("serverInfo", Server), slog.String("address", serve.Addr))
    return serve.ListenAndServe()
}
