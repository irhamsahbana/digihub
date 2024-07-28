package main

import (
	"codebase-app/cmd"
	"codebase-app/internal/adapter"
	"codebase-app/internal/infrastructure/config"
	"os"
	"strings"

	"flag"

	"github.com/rs/zerolog/log"
)

func main() {
	configPath := flag.String("path", "./", "path to config file")
	configFilename := flag.String("filename", ".env", "config file name")
	flag.Parse()

	initialize(*configPath, *configFilename)
	newOsArgs := []string{}

	for _, arg := range os.Args {
		if strings.Contains(arg, "-path") || strings.Contains(arg, "-filename") {
			continue
		}

		newOsArgs = append(newOsArgs, arg)
	}
	os.Args = newOsArgs

	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	seedCmd := flag.NewFlagSet("seed", flag.ExitOnError)
	consumerCmd := flag.NewFlagSet("consumer", flag.ExitOnError)
	wsCmd := flag.NewFlagSet("ws", flag.ExitOnError)

	if len(os.Args) < 2 {
		log.Info().Msg("No command provided, defaulting to 'server'")
		cmd.RunServer(serverCmd, os.Args[1:])
		os.Exit(0)
	}

	switch os.Args[1] {
	case "seed":
		cmd.RunSeed(seedCmd, os.Args[2:])
	case "consumer":
		cmd.RunConsumer(consumerCmd, os.Args[2:])
	case "server":
		cmd.RunServer(serverCmd, os.Args[2:])
	case "ws":
		cmd.RunWebsocket(wsCmd, os.Args[2:])
	default:
		log.Info().Msg("Invalid command provided, defaulting to 'server' with provided flags")
		if os.Args[1][0] == '-' { // check if the first argument is a flag
			cmd.RunServer(serverCmd, os.Args[1:])
			os.Exit(0)
		}

		cmd.RunServer(serverCmd, os.Args[2:]) // default to server if invalid command and flags are provided
	}
}

func initialize(path, filename string) {
	log.Info().Msgf("Initializing configuration with path: %s and filename: %s", path, filename)

	config.Configuration(
		config.WithPath(path),
		config.WithFilename(filename),
	).Initialize()

	adapter.Adapters = &adapter.Adapter{}
}
