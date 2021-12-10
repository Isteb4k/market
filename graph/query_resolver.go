package graph

import (
	"context"
	"market/auth"
	"market/graph/generated"
	"market/models"
)

type queryResolver struct{ *Resolver }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

func (r *queryResolver) Products(ctx context.Context) ([]*models.Product, error) {
	return r.ProductsRepo.GetProducts()
}

func (r *queryResolver) Viewer(ctx context.Context) (*models.Viewer, error) {
	user, err := auth.GetCurrentUserFromCTX(ctx)
	if err != nil {
		return nil, err
	}

	return &models.Viewer{
		User: user,
	}, nil
}
