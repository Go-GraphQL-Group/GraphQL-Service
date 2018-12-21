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

func (r *queryResolver) Starship(ctx context.Context, id string) (*model.Starship, error) {
	err, starship := boltdb.GetStarshipByID(id, nil)
	checkErr(err)
	return starship, err
}
func (r *queryResolver) Starships(ctx context.Context, first *int, after *string) (model.StarshipConnection, error) {
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return model.StarshipConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return model.StarshipConnection{}, err
		}
		from = i
	}
	count := 0
	startID := ""
	hasPreviousPage := true
	hasNextPage := true
	// 获取edges
	edges := []model.StarshipEdge{}
	db, err := bolt.Open("data/data.db", 0600, nil)
	checkErr(err)
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
				_, starship := boltdb.GetStarshipByID(string(k), db)
				edges = append(edges, model.StarshipEdge{
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
					_, starship := boltdb.GetStarshipByID(string(k), db)
					edges = append(edges, model.StarshipEdge{
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
		return model.StarshipConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := model.PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     encodeCursor(startID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
	}

	return model.StarshipConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
func (r *queryResolver) StarshipsSearch(ctx context.Context, search string, first *int, after *string) (*model.StarshipConnection, error) {
	//搜索对象的类型
	searchType := 0
	if strings.HasPrefix(search, "Name:") {
		search = strings.TrimPrefix(search, "Name:")
	} else if strings.HasPrefix(search, "Model:") {
		search = strings.TrimPrefix(search, "Model:")
		searchType = 1
	} else {
		return &model.StarshipConnection{}, errors.New("Search content must be ' Name:<StarShip's Name you want to get> ' OR ' Model:<StarShip's Model yout want to get> ' ")
	}
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return &model.StarshipConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return &model.StarshipConnection{}, err
		}
		from = i
	}
	count := 0
	hasPreviousPage := false
	hasNextPage := false
	// 获取edges
	edges := []model.StarshipEdge{}
	db, err := bolt.Open("data/data.db", 0600, nil)
	checkErr(err)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(starshipsBucket)).Cursor()
		k, _ := c.First()
		// 判断是否还有前向页
		if from != -1 {
			for k != nil {
				_, starship := boltdb.GetStarshipByID(string(k), db)
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
			_, starship := boltdb.GetStarshipByID(string(k), db)
			if searchType == 0 {
				if starship.Name == search {
					edges = append(edges, model.StarshipEdge{
						Node:   starship,
						Cursor: encodeCursor(string(k)),
					})
					count++
				}
			} else {
				if *(starship.Model) == search {
					edges = append(edges, model.StarshipEdge{
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
			_, starship := boltdb.GetStarshipByID(string(k), db)
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
		return &model.StarshipConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := model.PageInfo{
		StartCursor:     encodeCursor(edges[0].Node.ID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
	}

	return &model.StarshipConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
