//go:build !Test && !Shipping

package config

// Default (Development) configuration for the application.
const (
	ApiUrl        = "https://api.example.com/"
	Configuration = "Development"
)
