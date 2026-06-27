package main

import (
	"Spot-Sync/internal/config"
	"Spot-Sync/internal/server"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	// Connect to the database
	db := config.ConnectDatabase(cfg)
	// Start the server
	server.Start(cfg, db)

}
