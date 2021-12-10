package repositories

import (
	"github.com/go-pg/pg/v10"
)

// New - new db connection client
func New(opts *pg.Options) *pg.DB {
	return pg.Connect(opts)
}
