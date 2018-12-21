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

func (r *queryResolver) Specie(ctx context.Context, id string) (*model.Specie, error) {
	specie := db.MySQLGetSpecieByID(id)
	return specie, nil
}
func (r *queryResolver) Species(ctx context.Context, first *int, after *string) (model.SpecieConnection, error) {
	from := 1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return model.SpecieConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return model.SpecieConnection{}, err
		}
		from = i + 1
	}
	hasPreviousPage := true
	hasNextPage := true

	if from == 1 {
		hasPreviousPage = false
	}
	// 获取edges
	edges := []model.SpecieEdge{}

	for len(edges) < *first {
		specie := db.MySQLGetSpecieBy_id(strconv.Itoa(from))
		if specie.ID == "" {
			break
		}
		edges = append(edges, model.SpecieEdge{
			Node:   specie,
			Cursor: encodeCursor(strconv.Itoa(from)),
		})
		from++
	}

	if len(edges) < *first {
		hasNextPage = false
	}

	if len(edges) == 0 {
		return model.SpecieConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := model.PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     encodeCursor(strconv.Itoa(from - len(edges))),
		EndCursor:       encodeCursor(strconv.Itoa(from - 1)),
	}

	return model.SpecieConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: len(edges),
	}, nil
}
func (r *queryResolver) SpeciesSearch(ctx context.Context, search string, first *int, after *string) (*model.SpecieConnection, error) {
	if strings.HasPrefix(search, "Name:") {
		search = strings.TrimSpace(strings.TrimPrefix(search, "Name:"))
	} else {
		return &model.SpecieConnection{}, errors.New("Search content must be ' Name:<Specie's Name you want to get> ' ")
	}
	from := 1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return &model.SpecieConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return &model.SpecieConnection{}, err
		}
		from = i + 1
	}
	count := 0
	hasNextPage := false
	hasPreviousPage := false
	// 查询from之前是否存在满足search条件的条目
	for i := 1; i < from; i++ {
		specie := db.MySQLGetSpecieBy_id(strconv.Itoa(i))
		if specie.ID == "" {
			break
		}
		if specie.Name == search {
			hasPreviousPage = true
			break
		}
	}
	// 获取edges
	edges := []model.SpecieEdge{}
	for len(edges) < *first {
		specie := db.MySQLGetSpecieBy_id(strconv.Itoa(from))
		if specie.ID == "" {
			break
		}
		if specie.Name == search {
			edges = append(edges, model.SpecieEdge{
				Node:   specie,
				Cursor: encodeCursor(strconv.Itoa(from)),
			})
		}
		from++
	}

	// 查询from之后是否存在满足search条件的条目
	for {
		specie := db.MySQLGetSpecieBy_id(strconv.Itoa(from))
		if specie.ID == "" {
			break
		}
		if specie.Name == search {
			hasNextPage = true
			break
		}
		from++
	}
	if len(edges) == 0 {
		return &model.SpecieConnection{
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

	return &model.SpecieConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
