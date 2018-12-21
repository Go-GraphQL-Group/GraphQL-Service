package resolver

import (
	"github.com/Go-GraphQL-Group/GraphQL-Service/graphql"
)

type Resolver struct{}

func (r *Resolver) Query() graphql.QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }
