package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/blackorder/configwatcher"
)

// ServerConfig represents server-specific configuration
type ServerConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	MaxClients   int           `json:"max_clients"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	DSN         string `json:"dsn"`
	MaxConns    int    `json:"max_conns"`
	MaxIdle     int    `json:"max_idle"`
	MaxLifetime string `json:"max_lifetime"`
}

// AppConfig represents the main application configuration
type AppConfig struct {
	Environment string         `json:"environment"`
	LogLevel    string         `json:"log_level"`
	Server      ServerConfig   `json:"server"`
	Database    DatabaseConfig `json:"database"`
}

func main() {
	fmt.Println("🔧 Multi-Config Watcher Example")
	fmt.Println("================================")

	// Default configurations
	defaultAppConfig := AppConfig{
		Environment: "development",
		LogLevel:    "info",
		Server: ServerConfig{
			Host:         "localhost",
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			MaxClients:   100,
		},
		Database: DatabaseConfig{
			DSN:         "postgres://localhost/myapp",
			MaxConns:    10,
			MaxIdle:     5,
			MaxLifetime: "1h",
		},
	}

	defaultServerConfig := ServerConfig{
		Host:         "0.0.0.0",
		Port:         3000,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		MaxClients:   50,
	}

	defaultDBConfig := DatabaseConfig{
		DSN:         "postgres://localhost/standalone",
		MaxConns:    5,
		MaxIdle:     2,
		MaxLifetime: "30m",
	}

	// Create error channels
	appErrChan := make(chan error, 10)
	serverErrChan := make(chan error, 10)
	dbErrChan := make(chan error, 10)

	// Create multiple watchers
	appWatcher := configwatcher.NewWatcher(
		defaultAppConfig,
		"app.json",
		configwatcher.WithErrorChan[AppConfig](appErrChan),
	)

	serverWatcher := configwatcher.NewWatcher(
		defaultServerConfig,
		"server.json",
		configwatcher.WithErrorChan[ServerConfig](serverErrChan),
	)

	dbWatcher := configwatcher.NewWatcher(
		defaultDBConfig,
		"database.json",
		configwatcher.WithErrorChan[DatabaseConfig](dbErrChan),
	)

	// Handle errors
	go handleErrors("App", appErrChan)
	go handleErrors("Server", serverErrChan)
	go handleErrors("Database", dbErrChan)

	// Subscribe to all configuration changes
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appUpdates := appWatcher.Subscribe(ctx)
	serverUpdates := serverWatcher.Subscribe(ctx)
	dbUpdates := dbWatcher.Subscribe(ctx)

	// Monitor all configuration changes
	go func() {
		for {
			select {
			case <-appUpdates:
				fmt.Println("📱 App config updated:")
				printAppConfig(appWatcher.Get())
			case <-serverUpdates:
				fmt.Println("🖥️  Server config updated:")
				printServerConfig(serverWatcher.Get())
			case <-dbUpdates:
				fmt.Println("🗄️  Database config updated:")
				printDBConfig(dbWatcher.Get())
			case <-ctx.Done():
				return
			}
		}
	}()

	// Print initial configurations
	fmt.Println("\n📋 Initial Configurations:")
	fmt.Println("App Config:")
	printAppConfig(appWatcher.Get())
	fmt.Println("\nServer Config:")
	printServerConfig(serverWatcher.Get())
	fmt.Println("\nDatabase Config:")
	printDBConfig(dbWatcher.Get())

	// Simulate configuration updates
	go simulateUpdates(appWatcher, serverWatcher, dbWatcher)

	fmt.Println("\n📝 Edit app.json, server.json, or database.json to see live updates!")

	// Keep running for demo
	time.Sleep(30 * time.Second)
	fmt.Println("👋 Demo finished.")
}

func handleErrors(name string, errChan <-chan error) {
	for err := range errChan {
		log.Printf("❌ %s config error: %v", name, err)
	}
}

func printAppConfig(config AppConfig) {
	fmt.Printf("  Environment: %s\n", config.Environment)
	fmt.Printf("  Log Level: %s\n", config.LogLevel)
	fmt.Printf("  Server Port: %d\n", config.Server.Port)
	fmt.Printf("  DB Max Conns: %d\n", config.Database.MaxConns)
}

func printServerConfig(config ServerConfig) {
	fmt.Printf("  Host: %s\n", config.Host)
	fmt.Printf("  Port: %d\n", config.Port)
	fmt.Printf("  Read Timeout: %v\n", config.ReadTimeout)
	fmt.Printf("  Max Clients: %d\n", config.MaxClients)
}

func printDBConfig(config DatabaseConfig) {
	fmt.Printf("  DSN: %s\n", config.DSN)
	fmt.Printf("  Max Connections: %d\n", config.MaxConns)
	fmt.Printf("  Max Idle: %d\n", config.MaxIdle)
	fmt.Printf("  Max Lifetime: %s\n", config.MaxLifetime)
}

func simulateUpdates(
	appWatcher *configwatcher.Watcher[AppConfig],
	serverWatcher *configwatcher.Watcher[ServerConfig],
	dbWatcher *configwatcher.Watcher[DatabaseConfig],
) {
	time.Sleep(3 * time.Second)

	// Update app config
	fmt.Println("\n💾 Updating app config...")
	appConfig := appWatcher.Get()
	appConfig.Environment = "production"
	appConfig.LogLevel = "warn"
	if err := appWatcher.Save(appConfig); err != nil {
		log.Printf("Failed to save app config: %v", err)
	}

	time.Sleep(2 * time.Second)

	// Update server config
	fmt.Println("\n💾 Updating server config...")
	serverConfig := serverWatcher.Get()
	serverConfig.Port = 8443
	serverConfig.MaxClients = 200
	if err := serverWatcher.Save(serverConfig); err != nil {
		log.Printf("Failed to save server config: %v", err)
	}

	time.Sleep(2 * time.Second)

	// Update database config
	fmt.Println("\n💾 Updating database config...")
	dbConfig := dbWatcher.Get()
	dbConfig.MaxConns = 20
	dbConfig.MaxLifetime = "2h"
	if err := dbWatcher.Save(dbConfig); err != nil {
		log.Printf("Failed to save database config: %v", err)
	}
}
