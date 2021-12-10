package graph

import (
	"context"
	"market/graph/generated"
	"market/models"
)

type mutationResolver struct{ *Resolver }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

func (r *mutationResolver) RequestSignInCode(ctx context.Context, input models.RequestSignInCodeInput) (*models.ErrorPayload, error) {
	_, err := r.Auth.SendCode(input.Phone)
	if err != nil {
		return &models.ErrorPayload{
			Message: err.Error(),
		}, nil
	}

	return nil, nil
}

func (r *mutationResolver) SignInByCode(ctx context.Context, input models.SignInByCodeInput) (models.SignInOrErrorPayload, error) {
	res, err := r.Auth.SignIn(input.Phone, input.Code)
	if err != nil {
		return &models.ErrorPayload{
			Message: err.Error(),
		}, nil
	}

	return models.SignInPayload{
		Token: res.AuthToken,
		Viewer: &models.Viewer{
			User: &res.User,
		},
	}, nil
}
