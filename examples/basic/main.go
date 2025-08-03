package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/blackorder/configwatcher"
)

// AppConfig represents the application configuration
type AppConfig struct {
	AppName  string            `json:"app_name"`
	Port     int               `json:"port"`
	Debug    bool              `json:"debug"`
	Database DatabaseConfig    `json:"database"`
	Features map[string]string `json:"features"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Database string `json:"database"`
}

func main() {
	// Define default configuration
	defaultConfig := AppConfig{
		AppName: "example-app",
		Port:    8080,
		Debug:   false,
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "admin",
			Database: "myapp",
		},
		Features: map[string]string{
			"feature1": "enabled",
			"feature2": "disabled",
		},
	}

	// Create error channel for handling configuration errors
	errChan := make(chan error, 10)

	// Create configuration watcher
	watcher := configwatcher.NewWatcher(
		defaultConfig,
		"app-config.json",
		configwatcher.WithErrorChan[AppConfig](errChan),
	)

	fmt.Println("🚀 Starting configuration watcher example...")

	// Handle errors in a separate goroutine
	go func() {
		for err := range errChan {
			log.Printf("❌ Config error: %v", err)
		}
	}()

	// Get and display initial configuration
	config := watcher.Get()
	fmt.Printf("📋 Initial configuration:\n")
	printConfig(config)

	// Subscribe to configuration changes
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updateChan := watcher.Subscribe(ctx)
	go func() {
		for range updateChan {
			newConfig := watcher.Get()
			fmt.Printf("\n🔄 Configuration updated:\n")
			printConfig(newConfig)
		}
	}()

	// Demonstrate programmatic configuration update
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("\n💾 Updating configuration programmatically...")

		config := watcher.Get()
		config.Port = 9090
		config.Debug = true
		config.Features["feature2"] = "enabled"

		if err := watcher.Save(config); err != nil {
			log.Printf("Failed to save config: %v", err)
		}
	}()

	// Another update after some time
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("\n💾 Another configuration update...")

		config := watcher.Get()
		config.AppName = "updated-example-app"
		config.Database.Host = "remote-db.example.com"
		config.Features["feature3"] = "experimental"

		if err := watcher.Save(config); err != nil {
			log.Printf("Failed to save config: %v", err)
		}
	}()

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("\n📝 Try editing 'app-config.json' manually to see live updates!")
	fmt.Println("⏹️  Press Ctrl+C to stop...")

	<-sigChan
	fmt.Println("\n👋 Shutting down gracefully...")
}

func printConfig(config AppConfig) {
	fmt.Printf("  App Name: %s\n", config.AppName)
	fmt.Printf("  Port: %d\n", config.Port)
	fmt.Printf("  Debug: %t\n", config.Debug)
	fmt.Printf("  Database:\n")
	fmt.Printf("    Host: %s\n", config.Database.Host)
	fmt.Printf("    Port: %d\n", config.Database.Port)
	fmt.Printf("    Username: %s\n", config.Database.Username)
	fmt.Printf("    Database: %s\n", config.Database.Database)
	fmt.Printf("  Features:\n")
	for key, value := range config.Features {
		fmt.Printf("    %s: %s\n", key, value)
	}
}
