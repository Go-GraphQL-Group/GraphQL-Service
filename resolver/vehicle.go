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

func (r *queryResolver) Vehicle(ctx context.Context, id string) (*model.Vehicle, error) {
	err, vehicle := boltdb.GetVehicleByID(id, nil)
	checkErr(err)
	return vehicle, err
}
func (r *queryResolver) Vehicles(ctx context.Context, first *int, after *string) (model.VehicleConnection, error) {
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return model.VehicleConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return model.VehicleConnection{}, err
		}
		from = i
	}
	count := 0
	startID := ""
	hasPreviousPage := true
	hasNextPage := true
	// 获取edges
	edges := []model.VehicleEdge{}
	db, err := bolt.Open("data/data.db", 0600, nil)
	checkErr(err)
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
				_, vehicle := boltdb.GetVehicleByID(string(k), db)
				edges = append(edges, model.VehicleEdge{
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
					_, vehicle := boltdb.GetVehicleByID(string(k), db)
					edges = append(edges, model.VehicleEdge{
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
		return model.VehicleConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := model.PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     encodeCursor(startID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
	}

	return model.VehicleConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
func (r *queryResolver) VehiclesSearch(ctx context.Context, search string, first *int, after *string) (*model.VehicleConnection, error) {
	//搜索对象的类型
	searchType := 0
	if strings.HasPrefix(search, "Name:") {
		search = strings.TrimPrefix(search, "Name:")
	} else if strings.HasPrefix(search, "Model:") {
		search = strings.TrimPrefix(search, "Model:")
		searchType = 1
	} else {
		return &model.VehicleConnection{}, errors.New("Search content must be ' Name:<Vehicle's Name you want to get> ' OR ' Model:<Vehicle's Model yout want to get> ' ")
	}
	from := -1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return &model.VehicleConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return &model.VehicleConnection{}, err
		}
		from = i
	}
	count := 0
	hasPreviousPage := false
	hasNextPage := false
	// 获取edges
	edges := []model.VehicleEdge{}
	db, err := bolt.Open("data/data.db", 0600, nil)
	checkErr(err)
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(vehiclesBucket)).Cursor()
		k, _ := c.First()
		// 判断是否还有前向页
		if from != -1 {
			for k != nil {
				_, vehicle := boltdb.GetVehicleByID(string(k), db)
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
			_, vehicle := boltdb.GetVehicleByID(string(k), db)
			if searchType == 0 {
				if vehicle.Name == search {
					edges = append(edges, model.VehicleEdge{
						Node:   vehicle,
						Cursor: encodeCursor(string(k)),
					})
					count++
				}
			} else {
				if *(vehicle.Model) == search {
					edges = append(edges, model.VehicleEdge{
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
			_, vehicle := boltdb.GetVehicleByID(string(k), db)
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
		return &model.VehicleConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := model.PageInfo{
		StartCursor:     encodeCursor(edges[0].Node.ID),
		EndCursor:       encodeCursor(edges[count-1].Node.ID),
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
	}

	return &model.VehicleConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
