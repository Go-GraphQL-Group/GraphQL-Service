package resolver

import (
	"context"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

	boltdb "github.com/Go-GraphQL-Group/GraphQL-Service/db"
	"github.com/Go-GraphQL-Group/GraphQL-Service/model"
	"github.com/boltdb/bolt"
)

func (r *queryResolver) Planet(ctx context.Context, id string) (*model.Planet, error) {
	err, planet := boltdb.GetPlanetByID(id, nil)
	checkErr(err)
	return planet, err
}
func (r *queryResolver) Planets(ctx context.Context, first *int, after *string) (model.PlanetConnection, error) {
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return model.PlanetConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return model.PlanetConnection{}, err
		}
		from = i
	}
	count := 0
	startID := ""
	hasPreviousPage := true
	hasNextPage := true
	// 获取edges
	edges := []model.PlanetEdge{}
	db, err := bolt.Open("data/data.db", 0600, nil)
	checkErr(err)
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
				_, planet := boltdb.GetPlanetByID(string(k), db)
				edges = append(edges, model.PlanetEdge{
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
					_, planet := boltdb.GetPlanetByID(string(k), db)
					edges = append(edges, model.PlanetEdge{
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
		return model.PlanetConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := model.PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     encodeCursor(startID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
	}

	return model.PlanetConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
func (r *queryResolver) PlanetsSearch(ctx context.Context, search string, first *int, after *string) (*model.PlanetConnection, error) {
	if strings.HasPrefix(search, "Name:") {
		search = strings.TrimPrefix(search, "Name:")
	} else {
		return &model.PlanetConnection{}, errors.New("Search content must be ' Name:<Species's Name you want to get> ' ")
	}
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return &model.PlanetConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return &model.PlanetConnection{}, err
		}
		from = i
	}
	count := 0
	hasPreviousPage := false
	hasNextPage := false
	// 获取edges
	edges := []model.PlanetEdge{}
	db, err := bolt.Open("data/data.db", 0600, nil)
	checkErr(err)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(planetsBucket)).Cursor()
		k, _ := c.First()
		// 判断是否还有前向页
		if from != -1 {
			for k != nil {
				_, planet := boltdb.GetPlanetByID(string(k), db)
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
			_, planet := boltdb.GetPlanetByID(string(k), db)
			if planet.Name == search {
				edges = append(edges, model.PlanetEdge{
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
			_, planet := boltdb.GetPlanetByID(string(k), db)
			if planet.Name == search {
				hasNextPage = true
				break
			}
			k, _ = c.Next()
		}
		return nil
	})
	if count == 0 {
		return &model.PlanetConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := model.PageInfo{
		StartCursor:     encodeCursor(edges[0].Node.ID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
	}

	return &model.PlanetConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
