package main

import (
	"github.com/chichigami/blog-aggregator/internal/config"
	"github.com/chichigami/blog-aggregator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}
