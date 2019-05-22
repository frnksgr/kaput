package config

import (
	"fmt"
	"os"
)

// Listening data for server
type Listening struct {
	Port string
	Host string
}

// Calling data for recursive calls
type Calling struct {
	Protocol string
	Domain   string
	Port     string
}

// Config provides process wide configuration
type Config struct {
	Listening Listening
	Calling   Calling
}

// Data the configuration read upon start up
var Data *Config

func defaults() *Config {
	return &Config{
		Listening: Listening{
			Port: "8080",
			Host: "0.0.0.0",
		},

		Calling: Calling{
			Protocol: "http",
			Domain:   "localhost",
			Port:     "8080",
		},
	}
}

func fromEnv() *Config {
	return &Config{
		Listening: Listening{
			Port: os.Getenv("PORT"),
			Host: os.Getenv("HOST"),
		},

		Calling: Calling{
			Protocol: os.Getenv("CALLING_PROTOCOL"),
			Domain:   os.Getenv("CALLING_DOMAIN"),
			Port:     os.Getenv("CALLING_PORT"),
		},
	}
}

func (c *Config) merge(mc *Config) *Config {
	if mc.Listening.Port != "" {
		c.Listening.Port = mc.Listening.Port
	}
	if mc.Listening.Host != "" {
		c.Listening.Host = mc.Listening.Host
	}
	if mc.Calling.Protocol != "" {
		c.Calling.Protocol = mc.Calling.Protocol
	}
	if mc.Calling.Domain != "" {
		c.Calling.Domain = mc.Calling.Domain
	}
	if mc.Calling.Port != "" {
		c.Calling.Port = mc.Calling.Port
	}
	return c
}

func init() {
	fmt.Println("Reading configuration ...")
	Data = defaults().merge(fromEnv())
}
