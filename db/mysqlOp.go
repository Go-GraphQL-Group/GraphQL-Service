package db

import (
	"fmt"
	"io/ioutil"

	model2 "github.com/Go-GraphQL-Group/GraphQL-Service/model"
	model1 "github.com/Go-GraphQL-Group/SW-Crawler/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	yaml "gopkg.in/yaml.v2"
)

type config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Schema   string `yaml:"schema"`
}

var Engine *xorm.Engine

func init() {
	conf := &config{}
	confFile, err := ioutil.ReadFile("db/conf.yml")
	CheckErr(err)
	fmt.Println(string(confFile))
	err = yaml.Unmarshal(confFile, conf)
	CheckErr(err)
	dataSourceName := conf.Username + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.Schema + "?charset=utf8"
	// fmt.Println(dataSourceName)
	Engine, err = xorm.NewEngine("mysql", dataSourceName)
	CheckErr(err)
}

func mysqlGetPeople(people1 *model1.People) *model2.People {
	var has bool
	var err error
	people2 := convertPeople(people1)
	for _, filmID := range people1.Films {
		film1 := &model1.Film{
			ID: filmID,
		}
		has, err = Engine.Get(film1)
		CheckErr(err)
		if has {
			film2 := convertFilm(film1)
			people2.Films = append(people2.Films, film2)
		}
	}

	for _, specieID := range people1.Species {
		specie1 := &model1.Species{
			ID: specieID,
		}
		has, err = Engine.Get(specie1)
		CheckErr(err)
		if has {
			specie2 := convertSpecie(specie1)
			people2.Species = append(people2.Species, specie2)
		}
	}

	for _, vehicleID := range people1.Vehicles {
		vehicle1 := &model1.Vehicle{
			ID: vehicleID,
		}
		has, err = Engine.Get(vehicle1)
		CheckErr(err)
		if has {
			vehicle2 := convertVehicle(vehicle1)
			people2.Vehicles = append(people2.Vehicles, vehicle2)
		}
	}

	for _, starshipID := range people1.Starships {
		starship1 := &model1.Starship{
			ID: starshipID,
		}
		has, err = Engine.Get(starship1)
		CheckErr(err)
		if has {
			starship2 := convertStarship(starship1)
			people2.Starships = append(people2.Starships, starship2)
		}
	}
	return people2
}

func MySQLGetPeopleBy_id(_id string) *model2.People {
	people1 := &model1.People{}
	has, err := Engine.Where("_id = ?", _id).Get(people1)
	CheckErr(err)
	people2 := &model2.People{}
	if has {
		people2 = mysqlGetPeople(people1)
	}
	return people2
}

func MySQLGetPeopleByID(ID string) *model2.People {
	people1 := &model1.People{
		ID: ID,
	}
	has, err := Engine.Get(people1)
	CheckErr(err)
	people2 := &model2.People{}
	if has {
		people2 = mysqlGetPeople(people1)
	}
	return people2
}
