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

func (r *queryResolver) Vehicle(ctx context.Context, id string) (*model.Vehicle, error) {
	vehicle := db.MySQLGetVehicleByID(id)
	return vehicle, nil
}
func (r *queryResolver) Vehicles(ctx context.Context, first *int, after *string) (model.VehicleConnection, error) {
	from := 1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return model.VehicleConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return model.VehicleConnection{}, err
		}
		from = i + 1
	}
	hasPreviousPage := true
	hasNextPage := true

	if from == 1 {
		hasPreviousPage = false
	}
	// 获取edges
	edges := []model.VehicleEdge{}

	for len(edges) < *first {
		vehicle := db.MySQLGetVehicleBy_id(strconv.Itoa(from))
		if vehicle.ID == "" {
			break
		}
		edges = append(edges, model.VehicleEdge{
			Node:   vehicle,
			Cursor: encodeCursor(strconv.Itoa(from)),
		})
		from++
	}

	if len(edges) < *first {
		hasNextPage = false
	}

	if len(edges) == 0 {
		return model.VehicleConnection{}, nil
	}
	// 获取pageInfo
	pageInfo := model.PageInfo{
		HasPreviousPage: hasPreviousPage,
		HasNextPage:     hasNextPage,
		StartCursor:     encodeCursor(strconv.Itoa(from - len(edges))),
		EndCursor:       encodeCursor(strconv.Itoa(from - 1)),
	}

	return model.VehicleConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: len(edges),
	}, nil
}
func (r *queryResolver) VehiclesSearch(ctx context.Context, search string, first *int, after *string) (*model.VehicleConnection, error) {
	if strings.HasPrefix(search, "Name:") {
		search = strings.TrimSpace(strings.TrimPrefix(search, "Name:"))
	} else {
		return &model.VehicleConnection{}, errors.New("Search content must be ' Name:<Vehicle's Name you want to get> ' ")
	}
	from := 1
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return &model.VehicleConnection{}, err
		}
		i, err := strconv.Atoi(strings.TrimPrefix(string(b), "cursor"))
		if err != nil {
			return &model.VehicleConnection{}, err
		}
		from = i + 1
	}
	count := 0
	hasNextPage := false
	hasPreviousPage := false
	// 查询from之前是否存在满足search条件的条目
	for i := 1; i < from; i++ {
		vehicle := db.MySQLGetVehicleBy_id(strconv.Itoa(i))
		if vehicle.ID == "" {
			break
		}
		if vehicle.Name == search {
			hasPreviousPage = true
			break
		}
	}
	// 获取edges
	edges := []model.VehicleEdge{}
	for len(edges) < *first {
		vehicle := db.MySQLGetVehicleBy_id(strconv.Itoa(from))
		if vehicle.ID == "" {
			break
		}
		if vehicle.Name == search {
			edges = append(edges, model.VehicleEdge{
				Node:   vehicle,
				Cursor: encodeCursor(strconv.Itoa(from)),
			})
		}
		from++
	}

	// 查询from之后是否存在满足search条件的条目
	for {
		vehicle := db.MySQLGetVehicleBy_id(strconv.Itoa(from))
		if vehicle.ID == "" {
			break
		}
		if vehicle.Name == search {
			hasNextPage = true
			break
		}
		from++
	}
	if len(edges) == 0 {
		return &model.VehicleConnection{
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

	return &model.VehicleConnection{
		PageInfo:   pageInfo,
		Edges:      edges,
		TotalCount: count,
	}, nil
}
