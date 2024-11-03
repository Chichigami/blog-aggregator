package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/chichigami/blog-aggregator/internal/config"
	"github.com/chichigami/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	configFile, err := config.GetConfigFilePath()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	cfg := config.ReadAndParse(configFile)

	dbURL := cfg.DbURL
	db, db_err := sql.Open("postgres", dbURL)
	if db_err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	userCfg := state{
		cfg: &cfg,
		db:  dbQueries,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)

	input := os.Args
	if len(input) < 2 {
		log.Fatal("not enough arguments")
	}
	commandInput := command{
		name: input[1],
		args: input[2:],
	}
	err = cmds.run(&userCfg, commandInput)
	if err != nil {
		log.Fatal(err)
	}
}
