package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Configuration of kaput
// merged in order of
//	defaults
//	config file
//  env vars
//  cli flags
type Configuration struct {
	// server doe not run TLS
	path                    string
	listeningPort           string
	listeningHost           string
	callingPort             string
	callingDomain           string
	callingProtocol         string
	internalCallingPort     string
	internalCallingDomain   string
	internalCallingProtocol string
}

// Getter functions to make it external immutable
// Additional effort for this is ...

// Path configuration file
func Path() string {
	return configuration.path
}

// ListeningPort server port
func ListeningPort() string {
	return configuration.listeningPort
}

// ListeningHost server host
func ListeningHost() string {
	return configuration.listeningHost
}

// CallingPort calling port
func CallingPort() string {
	return configuration.callingPort
}

// CallingDomain calling domain
func CallingDomain() string {
	return configuration.callingDomain
}

// CallingProtocol calling protocol
func CallingProtocol() string {
	return configuration.callingProtocol
}

//InternalCallingPort itnernal calling port
func InternalCallingPort() string {
	return configuration.internalCallingPort
}

// InternalCallingDomain internal calling domain
func InternalCallingDomain() string {
	return configuration.internalCallingDomain
}

// InternalCallingProtocol internal calling protocol
func InternalCallingProtocol() string {
	return configuration.internalCallingProtocol
}

func (c *Configuration) merge(mc *Configuration) *Configuration {
	if mc.path != "" {
		c.path = mc.path
	}
	if mc.listeningPort != "" {
		c.listeningPort = mc.listeningPort
	}
	if mc.listeningHost != "" {
		c.listeningHost = mc.listeningHost
	}
	if mc.callingPort != "" {
		c.callingPort = mc.callingPort
	}
	if mc.callingDomain != "" {
		c.callingDomain = mc.callingDomain
	}
	if mc.callingProtocol != "" {
		c.callingProtocol = mc.callingProtocol
	}
	if mc.internalCallingPort != "" {
		c.internalCallingPort = mc.internalCallingPort
	}
	if mc.internalCallingDomain != "" {
		c.internalCallingDomain = mc.internalCallingDomain
	}
	if mc.internalCallingProtocol != "" {
		c.internalCallingProtocol = mc.internalCallingProtocol
	}
	return c
}

var configuration Configuration

var defaultConfiguration = &Configuration{
	path:                    "./config.json",
	listeningPort:           "8080",
	listeningHost:           "0.0.0.0",
	callingPort:             "8080",
	callingDomain:           "localhost",
	callingProtocol:         "http",
	internalCallingPort:     "8080",
	internalCallingDomain:   "localhost",
	internalCallingProtocol: "http",
}

type encodableConfiguration struct {
	ListeningPort           string `json:"listeningPort"`
	ListeningHost           string `json:"listeningHost"`
	CallingPort             string `json:"callingPort"`
	CallingDomain           string `json:"callingDomain"`
	CallingProtocol         string `json:"callingProtocol"`
	InternalCallingPort     string `json:"internalCallingPort"`
	InternalCallingDomain   string `json:"internalCallingDomain"`
	InternalCallingProtocol string `json:"internalCallingProtocol"`
}

func copy(from *encodableConfiguration, to *Configuration) {
	to.listeningPort = from.ListeningPort
	to.listeningHost = from.ListeningHost
	to.callingPort = from.CallingPort
	to.callingDomain = from.CallingDomain
	to.callingProtocol = from.CallingProtocol
	to.internalCallingPort = from.InternalCallingPort
	to.internalCallingDomain = from.InternalCallingDomain
	to.internalCallingProtocol = from.InternalCallingProtocol
}

func configfileConfiguration(path string) *Configuration {
	var configuration Configuration
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return &configuration
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err)
		return &configuration
	}
	// only json format supported for now
	var ec encodableConfiguration
	if err := json.Unmarshal(b, &ec); err != nil {
		log.Println(err)
		return &configuration
	}
	copy(&ec, &configuration)
	return &configuration
}

func environmentConfiguration() *Configuration {
	return &Configuration{
		listeningPort:           os.Getenv("K_LISTENING_PORT"),
		listeningHost:           os.Getenv("K_LISTENING_HOST"),
		callingPort:             os.Getenv("K_CALLING_PORT"),
		callingDomain:           os.Getenv("K_CALLING_DOMAIN"),
		callingProtocol:         os.Getenv("K_CALLING_PROTOCOL"),
		internalCallingPort:     os.Getenv("K_INTERNAL_CALLING_PORT"),
		internalCallingDomain:   os.Getenv("K_INTERNAL_CALLING_DOMAIN"),
		internalCallingProtocol: os.Getenv("K_INTERNAL_CALLING_PROTOCOL"),
	}
}

func flagsConfiguration() *Configuration {
	// TODO: support flags
	return &Configuration{}
}

func init() {
	fmt.Println("Reading configuration ...")
	// align configurations

	// K_LISTENING_PORT might overwrite PORT in environment
	if os.Getenv("PORT") != "" && os.Getenv("K_LISTENING_PORT") == "" {
		os.Setenv("K_LISTENING_PORT", os.Getenv("PORT"))
	}

	// configPath cannot be read from config file
	path := defaultConfiguration.path
	other := environmentConfiguration().
		merge(flagsConfiguration())
	if other.path != "" {
		path = other.path
	}

	configuration.merge(defaultConfiguration).
		merge(configfileConfiguration(path)).
		merge(other)
}
