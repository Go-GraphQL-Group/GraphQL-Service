package GraphQL_Service

import (
	"context"
)

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) People(ctx context.Context, id string) (*People, error) {
	panic("not implemented")
}
func (r *queryResolver) Film(ctx context.Context, id string) (*Film, error) {
	panic("not implemented")
}
func (r *queryResolver) Starship(ctx context.Context, id string) (*Starship, error) {
	panic("not implemented")
}
func (r *queryResolver) Vehicle(ctx context.Context, id string) (*Vehicle, error) {
	panic("not implemented")
}
func (r *queryResolver) Specie(ctx context.Context, id string) (*Specie, error) {
	panic("not implemented")
}
func (r *queryResolver) Planet(ctx context.Context, id string) (*Planet, error) {
	panic("not implemented")
}
func (r *queryResolver) Peoples(ctx context.Context, first *int, after *string) (PeopleConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) Films(ctx context.Context, first *int, after *string) (FilmConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) Starships(ctx context.Context, first *int, after *string) (StarshipConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) Vehicles(ctx context.Context, first *int, after *string) (VehicleConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) Species(ctx context.Context, first *int, after *string) (SpecieConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) Planets(ctx context.Context, first *int, after *string) (PlanetConnection, error) {
	panic("not implemented")
}
