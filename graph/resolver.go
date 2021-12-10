package graph

import (
	"market/auth"
	"market/repositories"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ProductsRepo repositories.Products
	Auth         auth.Client
}
