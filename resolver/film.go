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

func (r *queryResolver) Film(ctx context.Context, id string) (*model.Film, error) {
	film := db.MySQLGetFilmByID(id)
	return film, nil
}
func (r *queryResolver) Films(ctx context.Context, first *int, after *string) (model.FilmConnection, error) {
	from := 1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return model.FilmConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return model.FilmConnection{}, err
		}
		from = i + 1
	}

	hasPreviousPage := true
	hasNextPage := true

	if from == 1 {
		hasPreviousPage = false
	}
	// 获取edges
	edges := []model.FilmEdge{}

	for len(edges) < *first {
		film := db.MySQLGetFilmBy_id(strconv.Itoa(from))
		if film.ID == "" {
			break
		}
		edges = append(edges, model.FilmEdge{
			Node:   film,
			Cursor: encodeCursor(strconv.Itoa(from)),
		})
		from++
	}

	if len(edges) < *first {
		hasNextPage = false
	}

	if len(edges) == 0 {
		return model.FilmConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := model.PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     encodeCursor(strconv.Itoa(from - len(edges))),
		EndCursor:       encodeCursor(strconv.Itoa(from - 1)),
	}

	return model.FilmConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: len(edges),
	}, nil
}
func (r *queryResolver) FilmsSearch(ctx context.Context, search string, first *int, after *string) (*model.FilmConnection, error) {
	if strings.HasPrefix(search, "Title:") {
		search = strings.TrimPrefix(search, "Title:")
	} else {
		return &model.FilmConnection{}, errors.New("Search content must be ' Title:<Film's Title you want to get> ' ")
	}
	from := 1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return &model.FilmConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return &model.FilmConnection{}, err
		}
		from = i + 1
	}
	count := 0
	hasPreviousPage := false
	hasNextPage := false
	// 查询from之前是否存在满足search条件的条目
	for i := 1; i < from; i++ {
		film := db.MySQLGetFilmBy_id(strconv.Itoa(i))
		if film.ID == "" {
			break
		}
		if film.Title == search {
			hasPreviousPage = true
			break
		}
	}
	// 获取edges
	edges := []model.FilmEdge{}
	for len(edges) < *first {
		film := db.MySQLGetFilmBy_id(strconv.Itoa(from))
		if film.ID == "" {
			break
		}
		if film.Title == search {
			edges = append(edges, model.FilmEdge{
				Node:   film,
				Cursor: encodeCursor(strconv.Itoa(from)),
			})
		}
		from++
	}
	// 查询from之后是否存在满足search条件的条目
	for {
		film := db.MySQLGetFilmBy_id(strconv.Itoa(from))
		if film.ID == "" {
			break
		}
		if film.Title == search {
			hasNextPage = true
			break
		}
		from++
	}
	if len(edges) == 0 {
		return &model.FilmConnection{
			PageInfo: model.PageInfo{
				HasPreviousPage: hasPreviousPage,
				HasNextPage:     hasNextPage,
			},
		}, nil
	}
	// 获取pageInfo
	pageInfo := model.PageInfo{
		StartCursor:     edges[0].Cursor,
		EndCursor:       edges[len(edges)-1].Cursor,
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
	}

	return &model.FilmConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
