package GraphQL_Service

import (
	// "encoding/json"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	// "github.com/Go-GraphQL-Group/GraphQL-Service/server/service"
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
	err, people := GetPeopleByID(id, nil)
	checkErr(err)
	return people, err
}
func (r *queryResolver) Film(ctx context.Context, id string) (*Film, error) {
	err, film := GetFilmByID(id, nil)
	checkErr(err)
	return film, err
}
func (r *queryResolver) Starship(ctx context.Context, id string) (*Starship, error) {
	err, starship := GetStarshipByID(id, nil)
	checkErr(err)
	return starship, err
}
func (r *queryResolver) Vehicle(ctx context.Context, id string) (*Vehicle, error) {
	err, vehicle := GetVehicleByID(id, nil)
	checkErr(err)
	return vehicle, err
}
func (r *queryResolver) Specie(ctx context.Context, id string) (*Specie, error) {
	err, specie := GetSpeciesByID(id, nil)
	checkErr(err)
	return specie, err
}
func (r *queryResolver) Planet(ctx context.Context, id string) (*Planet, error) {
	err, planet := GetPlanetByID(id, nil)
	checkErr(err)
	return planet, err
}

func encodeCursor(k string) string {
	i, _ := strconv.Atoi(k)
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor%d", i)))
}

func (r *queryResolver) Peoples(ctx context.Context, first *int, after *string) (PeopleConnection, error) {
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return PeopleConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return PeopleConnection{}, err
		}
		from = i
	}
	count := 0
	startId := ""
	endId := ""
	hasPreviousPage := true
	hasNextPage := true
	// 获取edges
	edges := []PeopleEdge{}
	db, err := bolt.Open("./data/data.db", 0600, nil)
	CheckErr(err)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(peopleBucket)).Cursor()

		// 判断是否还有前向页
		k, v := c.First()
		if from == -1 || strconv.Itoa(from) == string(k) {
			startId = string(k)
			hasPreviousPage = false
		}

		if from == -1 {
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				_, people := GetPeopleByID(string(k), db)
				edges = append(edges, PeopleEdge{
					Node:   people,
					Cursor: encodeCursor(string(k)),
				})
				count++
				if count == *first {
					endId = string(k)
					break
				}
			}
		} else {
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				if strconv.Itoa(from) == string(k) {
					k, _ = c.Next()
					startId = string(k)
				}
				if startId != "" {
					_, people := GetPeopleByID(string(k), db)
					edges = append(edges, PeopleEdge{
						Node:   people,
						Cursor: encodeCursor(string(k)),
					})
					count++
					if count == *first {
						endId = string(k)
						break
					}
				}
			}
		}

		k, v = c.Next()
		if k == nil && v == nil {
			hasNextPage = false
		}
		return nil
	})
	if count == 0 {
		return PeopleConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     encodeCursor(startId),
		EndCursor:       encodeCursor(endId),
	}

	return PeopleConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
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
	// fmt.Println(ctx)
	// token := &service.Token{}
	// tokenJson, _ := ctx.Value(service.Issuer).(string)
	// json.Unmarshal([]byte(tokenJson), token)
	// service.ParseToken(token.SW_TOKEN, []byte(service.SecretKey))

	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return &PeopleConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return &PeopleConnection{}, err
		}
		from = i
	}
	count := 0
	hasPreviousPage := false
	hasNextPage := false
	// 获取edges
	edges := []PeopleEdge{}
	db, err := bolt.Open("./data/data.db", 0600, nil)
	CheckErr(err)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(peopleBucket)).Cursor()
		k, _ := c.First()
		// 判断是否还有前向页
		if from != -1 {
			for k != nil {
				_, people := GetPeopleByID(string(k), db)
				if people.Name == search {
					hasPreviousPage = true
				}
				if strconv.Itoa(from) == string(k) {
					k, _ = c.Next()
					break
				}
				k, _ = c.Next()
			}
		}

		// 添加edge
		for k != nil {
			_, people := GetPeopleByID(string(k), db)
			if people.Name == search {
				edges = append(edges, PeopleEdge{
					Node:   people,
					Cursor: encodeCursor(string(k)),
				})
				count++
			}
			k, _ = c.Next()
			if count == *first {
				break
			}
		}

		// 判断是否还有后向页
		for k != nil {
			_, people := GetPeopleByID(string(k), db)
			if people.Name == search {
				hasNextPage = true
				break
			}
			k, _ = c.Next()
		}
		return nil
	})
	if count == 0 {
		return &PeopleConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := PageInfo{
		StartCursor:     encodeCursor(edges[0].Node.ID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
	}

	return &PeopleConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil

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
