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

func (r *queryResolver) Film(ctx context.Context, id string) (*model.Film, error) {
	err, film := boltdb.GetFilmByID(id, nil)
	checkErr(err)
	return film, err
}
func (r *queryResolver) Films(ctx context.Context, first *int, after *string) (model.FilmConnection, error) {
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return model.FilmConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return model.FilmConnection{}, err
		}
		from = i
	}
	count := 0
	startID := ""
	hasPreviousPage := true
	hasNextPage := true
	// 获取edges
	edges := []model.FilmEdge{}
	db, err := bolt.Open("data/data.db", 0600, nil)
	checkErr(err)
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
				_, film := boltdb.GetFilmByID(string(k), db)
				edges = append(edges, model.FilmEdge{
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
					_, film := boltdb.GetFilmByID(string(k), db)
					edges = append(edges, model.FilmEdge{
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
		return model.FilmConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := model.PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     encodeCursor(startID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
	}

	return model.FilmConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
func (r *queryResolver) FilmsSearch(ctx context.Context, search string, first *int, after *string) (*model.FilmConnection, error) {
	if strings.HasPrefix(search, "Title:") {
		search = strings.TrimPrefix(search, "Title:")
	} else {
		return &model.FilmConnection{}, errors.New("Search content must be ' Title:<Film's Title you want to get> ' ")
	}
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return &model.FilmConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return &model.FilmConnection{}, err
		}
		from = i
	}
	count := 0
	hasPreviousPage := false
	hasNextPage := false
	// 获取edges
	edges := []model.FilmEdge{}
	db, err := bolt.Open("data/data.db", 0600, nil)
	checkErr(err)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(filmsBucket)).Cursor()
		k, _ := c.First()
		// 判断是否还有前向页
		if from != -1 {
			for k != nil {
				_, film := boltdb.GetFilmByID(string(k), db)
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
			_, film := boltdb.GetFilmByID(string(k), db)
			if film.Title == search {
				edges = append(edges, model.FilmEdge{
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
			_, film := boltdb.GetFilmByID(string(k), db)
			if film.Title == search {
				hasNextPage = true
				break
			}
			k, _ = c.Next()
		}
		return nil
	})
	if count == 0 {
		return &model.FilmConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := model.PageInfo{
		StartCursor:     encodeCursor(edges[0].Node.ID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
	}

	return &model.FilmConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
