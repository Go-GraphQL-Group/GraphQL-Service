package resolver

import (
	"context"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

	"github.com/Go-GraphQL-Group/GraphQL-Service/db"
	"github.com/Go-GraphQL-Group/GraphQL-Service/model"
)

func (r *queryResolver) Planet(ctx context.Context, id string) (*model.Planet, error) {
	planet := db.MySQLGetPlanetByID(id)
	return planet, nil
}
func (r *queryResolver) Planets(ctx context.Context, first *int, after *string) (model.PlanetConnection, error) {
	from := 1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return model.PlanetConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return model.PlanetConnection{}, err
		}
		from = i + 1
	}
	hasPreviousPage := true
	hasNextPage := true

	if from == 1 {
		hasPreviousPage = false
	}
	// 获取edges
	edges := []model.PlanetEdge{}

	for len(edges) < *first {
		planet := db.MySQLGetPlanetBy_id(strconv.Itoa(from))
		if planet.ID == "" {
			break
		}
		edges = append(edges, model.PlanetEdge{
			Node:   planet,
			Cursor: encodeCursor(strconv.Itoa(from)),
		})
		from++
	}

	if len(edges) < *first {
		hasNextPage = false
	}

	if len(edges) == 0 {
		return model.PlanetConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := model.PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     encodeCursor(strconv.Itoa(from - len(edges))),
		EndCursor:       encodeCursor(strconv.Itoa(from - 1)),
	}

	return model.PlanetConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: len(edges),
	}, nil
}
func (r *queryResolver) PlanetsSearch(ctx context.Context, search string, first *int, after *string) (*model.PlanetConnection, error) {
	if strings.HasPrefix(search, "Name:") {
		search = strings.TrimPrefix(search, "Name:")
	} else {
		return &model.PlanetConnection{}, errors.New("Search content must be ' Name:<Species's Name you want to get> ' ")
	}
	from := 1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return &model.PlanetConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return &model.PlanetConnection{}, err
		}
		from = i + 1
	}
	count := 0
	hasPreviousPage := false
	hasNextPage := false
	// 查询from之前是否存在满足search条件的条目
	for i := 1; i < from; i++ {
		planet := db.MySQLGetPlanetBy_id(strconv.Itoa(i))
		if planet.ID == "" {
			break
		}
		if planet.Name == search {
			hasPreviousPage = true
			break
		}
	}
	// 获取edges
	edges := []model.PlanetEdge{}
	for len(edges) < *first {
		planet := db.MySQLGetPlanetBy_id(strconv.Itoa(from))
		if planet.ID == "" {
			break
		}
		if planet.Name == search {
			edges = append(edges, model.PlanetEdge{
				Node:   planet,
				Cursor: encodeCursor(strconv.Itoa(from)),
			})
		}
		from++
	}
	// 查询from之后是否存在满足search条件的条目
	for {
		planet := db.MySQLGetPlanetBy_id(strconv.Itoa(from))
		if planet.ID == "" {
			break
		}
		if planet.Name == search {
			hasNextPage = true
			break
		}
		from++
	}
	if len(edges) == 0 {
		return &model.PlanetConnection{
			PageInfo: model.PageInfo{
				HasPreviousPage: hasPreviousPage,
				HasNextPage:     hasNextPage,
			},
		}, nil
	}
	pageInfo := model.PageInfo{
		StartCursor:     edges[0].Cursor,
		EndCursor:       edges[len(edges)-1].Cursor,
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
	}

	return &model.PlanetConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
