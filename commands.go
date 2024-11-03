package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chichigami/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if f, ok := c.handlers[cmd.name]; ok {
		err := f(s, cmd)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("command not registered")
}

func handlerLogin(s *state, cmd command) error {
	givenName := cmd.args[0]
	if len(givenName) == 0 {
		return fmt.Errorf("need a username. Usage: gator login USERNAME")
	}
	if _, err := s.db.GetUser(context.Background(), givenName); err != nil {
		fmt.Println("cant login; username DNE")
		os.Exit(1)
	}
	err := s.cfg.SetUser(givenName)
	if err != nil {
		return err
	}
	fmt.Printf("%s has been set \n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("need a username. Usage: go run . register USERNAME")
	}
	givenName := cmd.args[0]
	if _, err := s.db.GetUser(context.Background(), givenName); err == nil {
		fmt.Println("user exists")
		os.Exit(1)
	}
	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      givenName,
	}
	_, err := s.db.CreateUser(context.Background(), newUser)
	if err != nil {
		return err
	}
	s.cfg.SetUser(givenName)
	fmt.Printf("%s has been created\n%#v\n", givenName, newUser)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		fmt.Println("Reset failed")
		os.Exit(1)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not print users")
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	result, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("fetchFeed fail")
	}
	fmt.Println(result)
	return nil
}

type commands struct {
	handlers map[string]func(*state, command) error
}

type command struct {
	name string
	args []string
}
