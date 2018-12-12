package GraphQL_Service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Go-GraphQL-Group/SW-Crawler/model"
	"github.com/boltdb/bolt"
)

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

const peopleBucket = "People"
const filmsBucket = "Film"
const planetsBucket = "Planet"
const speciesBucket = "Specie"
const starshipsBucket = "Starship"
const vehiclesBucket = "Vehicle"

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
func (r *queryResolver) People(ctx context.Context, id string) (*People, error) {
	fmt.Println(44444444)
	db, err := bolt.Open("server/data/data.db", 0600, nil)
	checkErr(err)
	defer db.Close()
	fmt.Println(db.Path())
	people := &model.People{}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(peopleBucket))
		json.Unmarshal(b.Get([]byte(id)), people)
		fmt.Println(people)
		return nil
	})

	people2 := &People{
		ID:        id,
		Name:      people.Name,
		BirthYear: &people.Birth_year,
		EyeColor:  &people.Eye_color,
		Gender:    &people.Gender,
		HairColor: &people.Hair_color,
		Height:    &people.Heigth,
		Mass:      &people.Mass,
		SkinColor: &people.Skin_color,
		Homeworld: nil,
		Films:     nil,
		Species:   nil,
		Starships: nil,
		Vehicles:  nil,
	}
	return people2, nil
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
	from := 1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return PeopleConnection{}, err
		}
		fmt.Println(string(b))
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return PeopleConnection{}, err
		}
		from = i
	}
	to := -1
	if first != nil {
		to = from + *first
	}

	db, err := bolt.Open("server/data/data.db", 0600, nil)
	checkErr(err)
	defer db.Close()
	fmt.Println(db.Path())
	nodes := []People{}
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(peopleBucket)).Cursor()
		for k, v := c.Seek([]byte(string(from))); k != nil && bytes.Compare(k, []byte(string(to))) < 0; k, v = c.Next() {
			people := &model.People{}
			json.Unmarshal(v, people)
			people2 := &People{
				ID:        string(k),
				Name:      people.Name,
				BirthYear: &people.Birth_year,
				EyeColor:  &people.Eye_color,
				Gender:    &people.Gender,
				HairColor: &people.Hair_color,
				Height:    &people.Heigth,
				Mass:      &people.Mass,
				SkinColor: &people.Skin_color,
				Homeworld: nil,
				Films:     nil,
				Species:   nil,
				Starships: nil,
				Vehicles:  nil,
			}
			fmt.Println(people2)
			nodes = append(nodes, *people2)
		}
		return nil
	})

	peopleConnection := PeopleConnection{
		PageInfo: PageInfo{
			HasNextPage:     false,
			HasPreviousPage: false,
			StartCursor:     base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor%d", from))),
			EndCursor:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor%d", to))),
		},
		Edges:      nil,
		Nodes:      nodes,
		TotalCount: to - from,
	}
	return peopleConnection, nil
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
