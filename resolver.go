package GraphQL_Service

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
func (r *queryResolver) People(ctx context.Context, id string) (*People, error) {
	err, people := GetPeopleByID(id)
	checkErr(err)
	return people, err
}
func (r *queryResolver) Film(ctx context.Context, id string) (*Film, error) {
	err, film := GetFilmByID(id)
	checkErr(err)
	return film, err
}
func (r *queryResolver) Starship(ctx context.Context, id string) (*Starship, error) {
	err, starship := GetStarshipByID(id)
	checkErr(err)
	return starship, err
}
func (r *queryResolver) Vehicle(ctx context.Context, id string) (*Vehicle, error) {
	err, vehicle := GetVehicleByID(id)
	checkErr(err)
	return vehicle, err
}
func (r *queryResolver) Specie(ctx context.Context, id string) (*Specie, error) {
	err, specie := GetSpeciesByID(id)
	checkErr(err)
	return specie, err
}
func (r *queryResolver) Planet(ctx context.Context, id string) (*Planet, error) {
	err, planet := GetPlanetByID(id)
	checkErr(err)
	return planet, err
}

func encodeCursor(i uint) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor%d", i)))
}

func (r *queryResolver) Peoples(ctx context.Context, first *int, after *string) (PeopleConnection, error) {
	from := 1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return PeopleConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return PeopleConnection{}, err
		}
		from = i + 1
	}
	to := -1
	if first != nil {
		to = from + *first
	}

	// 获取edges
	edges := []PeopleEdge{}
	var i uint
	for i = uint(from); i < uint(to); i++ {
		_, people := GetPeopleByID(strconv.Itoa(int(i)))
		if people.ID == "" {
			break
		}
		edges = append(edges, PeopleEdge{
			Node:   people,
			Cursor: encodeCursor(i),
		})
	}

	// 获取pageInfo
	pageInfo := PageInfo{
		StartCursor: encodeCursor(uint(from)),
	}
	if from == 1 || len(edges) == 0 {
		pageInfo.HasPreviousPage = false
	} else {
		pageInfo.HasPreviousPage = true
	}
	if i < uint(to) || len(edges) == 0 {
		pageInfo.HasNextPage = false
	} else if i == uint(to) {
		_, people := GetPeopleByID(strconv.Itoa(int(i)))
		if people.ID == "" {
			pageInfo.HasNextPage = false
		} else {
			pageInfo.HasNextPage = true
		}
	}
	if len(edges) == 0 {
		pageInfo.EndCursor = encodeCursor(i)
	} else {
		pageInfo.EndCursor = encodeCursor(i - 1)
	}
	return PeopleConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: int(i) - from,
	}, nil
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

func (r *queryResolver) PeopleSearch(ctx context.Context, search string, first *int, after *string) (*PeopleConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) FilmsSearch(ctx context.Context, search string, first *int, after *string) (*FilmConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) StarshipsSearch(ctx context.Context, search string, first *int, after *string) (*StarshipConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) VehiclesSearch(ctx context.Context, search string, first *int, after *string) (*VehicleConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) SpeciesSearch(ctx context.Context, search string, first *int, after *string) (*SpecieConnection, error) {
	panic("not implemented")
}
func (r *queryResolver) PlanetsSearch(ctx context.Context, search string, first *int, after *string) (*PlanetConnection, error) {
	panic("not implemented")
}
