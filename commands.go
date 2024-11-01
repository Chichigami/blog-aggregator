package main

import (
	"fmt"
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
	if len(cmd.args) == 0 {
		return fmt.Errorf("need a username. Usage: gator login USERNAME")
	}
	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("%s has been set \n", cmd.args[0])
	return nil
}

type commands struct {
	handlers map[string]func(*state, command) error
}

type command struct {
	name string
	args []string
}
