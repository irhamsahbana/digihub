package main

import (
	"codebase-app/cmd"
	"codebase-app/internal/adapter"
	"codebase-app/internal/infrastructure/config"
	"os"

	"flag"

	"github.com/rs/zerolog/log"
)

func main() {
	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	seedCmd := flag.NewFlagSet("seed", flag.ExitOnError)
	consumerCmd := flag.NewFlagSet("consumer", flag.ExitOnError)
	ulidCmd := flag.NewFlagSet("ulid", flag.ExitOnError)

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
		cmd.RunWebsocket(ulidCmd, os.Args[2:])
	default:
		log.Info().Msg("Invalid command provided, defaulting to 'server' with provided flags")
		if os.Args[1][0] == '-' { // check if the first argument is a flag
			cmd.RunServer(serverCmd, os.Args[1:])
			os.Exit(0)
		}

		cmd.RunServer(serverCmd, os.Args[2:]) // default to server if invalid command and flags are provided
	}
}

func init() {
	config.Configuration(
		config.WithPath("./"),
		config.WithFilename(".env"),
	).Initialize()

	adapter.Adapters = &adapter.Adapter{}
}
