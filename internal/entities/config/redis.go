package config

import "time"

// Redis is a configuration for standalone redis.
type Redis struct {
	Address, Password         string
	DB                        int
	ReadTimeout, WriteTimeout time.Duration
}
