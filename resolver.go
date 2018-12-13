package GraphQL_Service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
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
	startID := ""
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
			startID = string(k)
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
					break
				}
			}
		} else {
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				if strconv.Itoa(from) == string(k) {
					k, _ = c.Next()
					startID = string(k)
				}
				if startID != "" {
					_, people := GetPeopleByID(string(k), db)
					edges = append(edges, PeopleEdge{
						Node:   people,
						Cursor: encodeCursor(string(k)),
					})
					count++
					if count == *first {
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
		StartCursor:     encodeCursor(startID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
	}

	return PeopleConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
func (r *queryResolver) Films(ctx context.Context, first *int, after *string) (FilmConnection, error) {
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return FilmConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return FilmConnection{}, err
		}
		from = i
	}
	count := 0
	startID := ""
	hasPreviousPage := true
	hasNextPage := true
	// 获取edges
	edges := []FilmEdge{}
	db, err := bolt.Open("./data/data.db", 0600, nil)
	CheckErr(err)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(filmsBucket)).Cursor()

		// 判断是否还有前向页
		k, v := c.First()
		if from == -1 || strconv.Itoa(from) == string(k) {
			startID = string(k)
			hasPreviousPage = false
		}

		if from == -1 {
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				_, film := GetFilmByID(string(k), db)
				edges = append(edges, FilmEdge{
					Node:   film,
					Cursor: encodeCursor(string(k)),
				})
				count++
				if count == *first {
					break
				}
			}
		} else {
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				if strconv.Itoa(from) == string(k) {
					k, _ = c.Next()
					startID = string(k)
				}
				if startID != "" {
					_, film := GetFilmByID(string(k), db)
					edges = append(edges, FilmEdge{
						Node:   film,
						Cursor: encodeCursor(string(k)),
					})
					count++
					if count == *first {
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
		return FilmConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     encodeCursor(startID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
	}

	return FilmConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
func (r *queryResolver) Starships(ctx context.Context, first *int, after *string) (StarshipConnection, error) {
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return StarshipConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return StarshipConnection{}, err
		}
		from = i
	}
	count := 0
	startID := ""
	hasPreviousPage := true
	hasNextPage := true
	// 获取edges
	edges := []StarshipEdge{}
	db, err := bolt.Open("./data/data.db", 0600, nil)
	CheckErr(err)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(starshipsBucket)).Cursor()

		// 判断是否还有前向页
		k, v := c.First()
		if from == -1 || strconv.Itoa(from) == string(k) {
			startID = string(k)
			hasPreviousPage = false
		}

		if from == -1 {
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				_, starship := GetStarshipByID(string(k), db)
				edges = append(edges, StarshipEdge{
					Node:   starship,
					Cursor: encodeCursor(string(k)),
				})
				count++
				if count == *first {
					break
				}
			}
		} else {
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				if strconv.Itoa(from) == string(k) {
					k, _ = c.Next()
					startID = string(k)
				}
				if startID != "" {
					_, starship := GetStarshipByID(string(k), db)
					edges = append(edges, StarshipEdge{
						Node:   starship,
						Cursor: encodeCursor(string(k)),
					})
					count++
					if count == *first {
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
		return StarshipConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     encodeCursor(startID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
	}

	return StarshipConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
func (r *queryResolver) Vehicles(ctx context.Context, first *int, after *string) (VehicleConnection, error) {
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return VehicleConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return VehicleConnection{}, err
		}
		from = i
	}
	count := 0
	startID := ""
	hasPreviousPage := true
	hasNextPage := true
	// 获取edges
	edges := []VehicleEdge{}
	db, err := bolt.Open("./data/data.db", 0600, nil)
	CheckErr(err)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(vehiclesBucket)).Cursor()

		// 判断是否还有前向页
		k, v := c.First()
		if from == -1 || strconv.Itoa(from) == string(k) {
			startID = string(k)
			hasPreviousPage = false
		}

		if from == -1 {
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				_, vehicle := GetVehicleByID(string(k), db)
				edges = append(edges, VehicleEdge{
					Node:   vehicle,
					Cursor: encodeCursor(string(k)),
				})
				count++
				if count == *first {
					break
				}
			}
		} else {
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				if strconv.Itoa(from) == string(k) {
					k, _ = c.Next()
					startID = string(k)
				}
				if startID != "" {
					_, vehicle := GetVehicleByID(string(k), db)
					edges = append(edges, VehicleEdge{
						Node:   vehicle,
						Cursor: encodeCursor(string(k)),
					})
					count++
					if count == *first {
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
		return VehicleConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     encodeCursor(startID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
	}

	return VehicleConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
func (r *queryResolver) Species(ctx context.Context, first *int, after *string) (SpecieConnection, error) {
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return SpecieConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return SpecieConnection{}, err
		}
		from = i
	}
	count := 0
	startID := ""
	hasPreviousPage := true
	hasNextPage := true
	// 获取edges
	edges := []SpecieEdge{}
	db, err := bolt.Open("./data/data.db", 0600, nil)
	CheckErr(err)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(speciesBucket)).Cursor()

		// 判断是否还有前向页
		k, v := c.First()
		if from == -1 || strconv.Itoa(from) == string(k) {
			startID = string(k)
			hasPreviousPage = false
		}

		if from == -1 {
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				_, specie := GetSpeciesByID(string(k), db)
				edges = append(edges, SpecieEdge{
					Node:   specie,
					Cursor: encodeCursor(string(k)),
				})
				count++
				if count == *first {
					break
				}
			}
		} else {
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				if strconv.Itoa(from) == string(k) {
					k, _ = c.Next()
					startID = string(k)
				}
				if startID != "" {
					_, specie := GetSpeciesByID(string(k), db)
					edges = append(edges, SpecieEdge{
						Node:   specie,
						Cursor: encodeCursor(string(k)),
					})
					count++
					if count == *first {
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
		return SpecieConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     encodeCursor(startID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
	}

	return SpecieConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
func (r *queryResolver) Planets(ctx context.Context, first *int, after *string) (PlanetConnection, error) {
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return PlanetConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return PlanetConnection{}, err
		}
		from = i
	}
	count := 0
	startID := ""
	hasPreviousPage := true
	hasNextPage := true
	// 获取edges
	edges := []PlanetEdge{}
	db, err := bolt.Open("./data/data.db", 0600, nil)
	CheckErr(err)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(planetsBucket)).Cursor()

		// 判断是否还有前向页
		k, v := c.First()
		if from == -1 || strconv.Itoa(from) == string(k) {
			startID = string(k)
			hasPreviousPage = false
		}

		if from == -1 {
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				_, planet := GetPlanetByID(string(k), db)
				edges = append(edges, PlanetEdge{
					Node:   planet,
					Cursor: encodeCursor(string(k)),
				})
				count++
				if count == *first {
					break
				}
			}
		} else {
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				if strconv.Itoa(from) == string(k) {
					k, _ = c.Next()
					startID = string(k)
				}
				if startID != "" {
					_, planet := GetPlanetByID(string(k), db)
					edges = append(edges, PlanetEdge{
						Node:   planet,
						Cursor: encodeCursor(string(k)),
					})
					count++
					if count == *first {
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
		return PlanetConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     encodeCursor(startID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
	}

	return PlanetConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}

func (r *queryResolver) PeopleSearch(ctx context.Context, search string, first *int, after *string) (*PeopleConnection, error) {
	// fmt.Println(ctx)
	// token := &service.Token{}
	// tokenJson, _ := ctx.Value(service.Issuer).(string)
	// json.Unmarshal([]byte(tokenJson), token)
	// service.ParseToken(token.SW_TOKEN, []byte(service.SecretKey))
	if strings.HasPrefix(search, "Name:") {
		search = strings.TrimPrefix(search, "Name:")
	} else {
		return &PeopleConnection{}, errors.New("Search content must be ' Name:<People's Name you want to get> ' ")
	}
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
			if first != nil && count == *first {
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
	if strings.HasPrefix(search, "Title:") {
		search = strings.TrimPrefix(search, "Title:")
	} else {
		return &FilmConnection{}, errors.New("Search content must be ' Title:<Film's Title you want to get> ' ")
	}
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return &FilmConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return &FilmConnection{}, err
		}
		from = i
	}
	count := 0
	hasPreviousPage := false
	hasNextPage := false
	// 获取edges
	edges := []FilmEdge{}
	db, err := bolt.Open("./data/data.db", 0600, nil)
	CheckErr(err)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(filmsBucket)).Cursor()
		k, _ := c.First()
		// 判断是否还有前向页
		if from != -1 {
			for k != nil {
				_, film := GetFilmByID(string(k), db)
				if film.Title == search {
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
			_, film := GetFilmByID(string(k), db)
			if film.Title == search {
				edges = append(edges, FilmEdge{
					Node:   film,
					Cursor: encodeCursor(string(k)),
				})
				count++
			}
			k, _ = c.Next()
			if first != nil && count == *first {
				break
			}
		}

		// 判断是否还有后向页
		for k != nil {
			_, film := GetFilmByID(string(k), db)
			if film.Title == search {
				hasNextPage = true
				break
			}
			k, _ = c.Next()
		}
		return nil
	})
	if count == 0 {
		return &FilmConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := PageInfo{
		StartCursor:     encodeCursor(edges[0].Node.ID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
	}

	return &FilmConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
func (r *queryResolver) StarshipsSearch(ctx context.Context, search string, first *int, after *string) (*StarshipConnection, error) {
	//搜索对象的类型
	searchType := 0
	if strings.HasPrefix(search, "Name:") {
		search = strings.TrimPrefix(search, "Name:")
	} else if strings.HasPrefix(search, "Model:") {
		search = strings.TrimPrefix(search, "Model:")
		searchType = 1
	} else {
		return &StarshipConnection{}, errors.New("Search content must be ' Name:<StarShip's Name you want to get> ' OR ' Model:<StarShip's Model yout want to get> ' ")
	}
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return &StarshipConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return &StarshipConnection{}, err
		}
		from = i
	}
	count := 0
	hasPreviousPage := false
	hasNextPage := false
	// 获取edges
	edges := []StarshipEdge{}
	db, err := bolt.Open("./data/data.db", 0600, nil)
	CheckErr(err)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(starshipsBucket)).Cursor()
		k, _ := c.First()
		// 判断是否还有前向页
		if from != -1 {
			for k != nil {
				_, starship := GetStarshipByID(string(k), db)
				if searchType == 0 {
					if starship.Name == search {
						hasPreviousPage = true
					}
				} else {
					if *(starship.Model) == search {
						hasPreviousPage = true
					}
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
			_, starship := GetStarshipByID(string(k), db)
			if searchType == 0 {
				if starship.Name == search {
					edges = append(edges, StarshipEdge{
						Node:   starship,
						Cursor: encodeCursor(string(k)),
					})
					count++
				}
			} else {
				if *(starship.Model) == search {
					edges = append(edges, StarshipEdge{
						Node:   starship,
						Cursor: encodeCursor(string(k)),
					})
					count++
				}
			}
			k, _ = c.Next()
			if first != nil && count == *first {
				break
			}
		}

		// 判断是否还有后向页
		for k != nil {
			_, starship := GetStarshipByID(string(k), db)
			if searchType == 0 {
				if starship.Name == search {
					hasNextPage = true
					break
				}
			} else {
				if *(starship.Model) == search {
					hasNextPage = true
					break
				}
			}
			k, _ = c.Next()
		}
		return nil
	})
	if count == 0 {
		return &StarshipConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := PageInfo{
		StartCursor:     encodeCursor(edges[0].Node.ID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
	}

	return &StarshipConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
func (r *queryResolver) VehiclesSearch(ctx context.Context, search string, first *int, after *string) (*VehicleConnection, error) {
	//搜索对象的类型
	searchType := 0
	if strings.HasPrefix(search, "Name:") {
		search = strings.TrimPrefix(search, "Name:")
	} else if strings.HasPrefix(search, "Model:") {
		search = strings.TrimPrefix(search, "Model:")
		searchType = 1
	} else {
		return &VehicleConnection{}, errors.New("Search content must be ' Name:<Vehicle's Name you want to get> ' OR ' Model:<Vehicle's Model yout want to get> ' ")
	}
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return &VehicleConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return &VehicleConnection{}, err
		}
		from = i
	}
	count := 0
	hasPreviousPage := false
	hasNextPage := false
	// 获取edges
	edges := []VehicleEdge{}
	db, err := bolt.Open("./data/data.db", 0600, nil)
	CheckErr(err)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(vehiclesBucket)).Cursor()
		k, _ := c.First()
		// 判断是否还有前向页
		if from != -1 {
			for k != nil {
				_, vehicle := GetVehicleByID(string(k), db)
				if searchType == 0 {
					if vehicle.Name == search {
						hasPreviousPage = true
					}
				} else {
					if *(vehicle.Model) == search {
						hasPreviousPage = true
					}
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
			_, vehicle := GetVehicleByID(string(k), db)
			if searchType == 0 {
				if vehicle.Name == search {
					edges = append(edges, VehicleEdge{
						Node:   vehicle,
						Cursor: encodeCursor(string(k)),
					})
					count++
				}
			} else {
				if *(vehicle.Model) == search {
					edges = append(edges, VehicleEdge{
						Node:   vehicle,
						Cursor: encodeCursor(string(k)),
					})
					count++
				}
			}
			k, _ = c.Next()
			if first != nil && count == *first {
				break
			}
		}

		// 判断是否还有后向页
		for k != nil {
			_, vehicle := GetVehicleByID(string(k), db)
			if searchType == 0 {
				if vehicle.Name == search {
					hasNextPage = true
					break
				}
			} else {
				if *(vehicle.Model) == search {
					hasNextPage = true
					break
				}
			}
			k, _ = c.Next()
		}
		return nil
	})
	if count == 0 {
		return &VehicleConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := PageInfo{
		StartCursor:     encodeCursor(edges[0].Node.ID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
	}

	return &VehicleConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
func (r *queryResolver) SpeciesSearch(ctx context.Context, search string, first *int, after *string) (*SpecieConnection, error) {
	if strings.HasPrefix(search, "Name:") {
		search = strings.TrimPrefix(search, "Name:")
	} else {
		return &SpecieConnection{}, errors.New("Search content must be ' Name:<Species's Name you want to get> ' ")
	}
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return &SpecieConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return &SpecieConnection{}, err
		}
		from = i
	}
	count := 0
	hasPreviousPage := false
	hasNextPage := false
	// 获取edges
	edges := []SpecieEdge{}
	db, err := bolt.Open("./data/data.db", 0600, nil)
	CheckErr(err)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(speciesBucket)).Cursor()
		k, _ := c.First()
		// 判断是否还有前向页
		if from != -1 {
			for k != nil {
				_, specie := GetSpeciesByID(string(k), db)
				if specie.Name == search {
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
			_, specie := GetSpeciesByID(string(k), db)
			if specie.Name == search {
				edges = append(edges, SpecieEdge{
					Node:   specie,
					Cursor: encodeCursor(string(k)),
				})
				count++
			}
			k, _ = c.Next()
			if first != nil && count == *first {
				break
			}
		}

		// 判断是否还有后向页
		for k != nil {
			_, specie := GetSpeciesByID(string(k), db)
			if specie.Name == search {
				hasNextPage = true
				break
			}
			k, _ = c.Next()
		}
		return nil
	})
	if count == 0 {
		return &SpecieConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := PageInfo{
		StartCursor:     encodeCursor(edges[0].Node.ID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
	}

	return &SpecieConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
func (r *queryResolver) PlanetsSearch(ctx context.Context, search string, first *int, after *string) (*PlanetConnection, error) {
	if strings.HasPrefix(search, "Name:") {
		search = strings.TrimPrefix(search, "Name:")
	} else {
		return &PlanetConnection{}, errors.New("Search content must be ' Name:<Species's Name you want to get> ' ")
	}
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return &PlanetConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return &PlanetConnection{}, err
		}
		from = i
	}
	count := 0
	hasPreviousPage := false
	hasNextPage := false
	// 获取edges
	edges := []PlanetEdge{}
	db, err := bolt.Open("./data/data.db", 0600, nil)
	CheckErr(err)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(planetsBucket)).Cursor()
		k, _ := c.First()
		// 判断是否还有前向页
		if from != -1 {
			for k != nil {
				_, planet := GetPlanetByID(string(k), db)
				if planet.Name == search {
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
			_, planet := GetPlanetByID(string(k), db)
			if planet.Name == search {
				edges = append(edges, PlanetEdge{
					Node:   planet,
					Cursor: encodeCursor(string(k)),
				})
				count++
			}
			k, _ = c.Next()
			if first != nil && count == *first {
				break
			}
		}

		// 判断是否还有后向页
		for k != nil {
			_, planet := GetPlanetByID(string(k), db)
			if planet.Name == search {
				hasNextPage = true
				break
			}
			k, _ = c.Next()
		}
		return nil
	})
	if count == 0 {
		return &PlanetConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := PageInfo{
		StartCursor:     encodeCursor(edges[0].Node.ID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
	}

	return &PlanetConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
