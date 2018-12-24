package db

import (
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
	// fmt.Println(string(confFile))
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

func mysqlGetFilm(film1 *model1.Film) *model2.Film {
	var has bool
	var err error
	film2 := convertFilm(film1)

	for _, specieID := range film1.Species {
		specie1 := &model1.Species{
			ID: specieID,
		}
		has, err = Engine.Get(specie1)
		CheckErr(err)
		if has {
			specie2 := convertSpecie(specie1)
			film2.Species = append(film2.Species, specie2)
		}
	}

	for _, starshipID := range film1.Starships {
		starship1 := &model1.Starship{
			ID: starshipID,
		}
		has, err = Engine.Get(starship1)
		CheckErr(err)
		if has {
			starship2 := convertStarship(starship1)
			film2.Starships = append(film2.Starships, starship2)
		}
	}

	for _, vehicleID := range film1.Vehicles {
		vehicle1 := &model1.Vehicle{
			ID: vehicleID,
		}
		has, err = Engine.Get(vehicle1)
		CheckErr(err)
		if has {
			vehicle2 := convertVehicle(vehicle1)
			film2.Vehicles = append(film2.Vehicles, vehicle2)
		}
	}

	for _, CharacterID := range film1.Character {
		character1 := &model1.People{
			ID: CharacterID,
		}
		has, err = Engine.Get(character1)
		CheckErr(err)
		if has {
			character2 := convertPeople(character1)
			film2.Characters = append(film2.Characters, character2)
		}
	}

	for _, PlanetID := range film1.Character {
		planet1 := &model1.Planet{
			ID: PlanetID,
		}
		has, err = Engine.Get(planet1)
		CheckErr(err)
		if has {
			planet2 := convertPlanet(planet1)
			film2.Planets = append(film2.Planets, planet2)
		}
	}

	return film2
}

func MySQLGetFilmBy_id(_id string) *model2.Film {
	film1 := &model1.Film{}
	has, err := Engine.Where("_id = ?", _id).Get(film1)
	CheckErr(err)
	film2 := &model2.Film{}
	if has {
		film2 = mysqlGetFilm(film1)
	}
	return film2
}

func MySQLGetFilmByID(ID string) *model2.Film {
	film1 := &model1.Film{
		ID: ID,
	}
	has, err := Engine.Get(film1)
	CheckErr(err)
	film2 := &model2.Film{}
	if has {
		film2 = mysqlGetFilm(film1)
	}
	return film2
}

func mysqlGetPlanet(planet1 *model1.Planet) *model2.Planet {
	var has bool
	var err error
	planet2 := convertPlanet(planet1)
	for _, peopleID := range planet1.Residents {
		people1 := &model1.People{
			ID: peopleID,
		}
		has, err = Engine.Get(people1)
		CheckErr(err)
		if has {
			people2 := convertPeople(people1)
			planet2.Residents = append(planet2.Residents, people2)
		}
	}

	for _, filmID := range planet1.Films {
		film1 := &model1.Film{
			ID: filmID,
		}
		has, err = Engine.Get(film1)
		CheckErr(err)
		if has {
			film2 := convertFilm(film1)
			planet2.Films = append(planet2.Films, film2)
		}
	}

	return planet2
}

func MySQLGetPlanetBy_id(_id string) *model2.Planet {
	planet1 := &model1.Planet{}
	has, err := Engine.Where("_id = ?", _id).Get(planet1)
	CheckErr(err)
	planet2 := &model2.Planet{}
	if has {
		planet2 = mysqlGetPlanet(planet1)
	}
	return planet2
}

func MySQLGetPlanetByID(ID string) *model2.Planet {
	planet1 := &model1.Planet{
		ID: ID,
	}
	has, err := Engine.Get(planet1)
	CheckErr(err)
	planet2 := &model2.Planet{}
	if has {
		planet2 = mysqlGetPlanet(planet1)
	}
	return planet2
}

func mysqlGetSpecie(specie1 *model1.Species) *model2.Specie {
	var has bool
	var err error
	specie2 := convertSpecie(specie1)

	for _, filmID := range specie1.Films {
		film1 := &model1.Film{
			ID: filmID,
		}
		has, err = Engine.Get(film1)
		CheckErr(err)
		if has {
			film2 := convertFilm(film1)
			specie2.Films = append(specie2.Films, film2)
		}
	}

	for _, PeopleID := range specie1.People {
		character1 := &model1.People{
			ID: PeopleID,
		}
		has, err = Engine.Get(character1)
		CheckErr(err)
		if has {
			character2 := convertPeople(character1)
			specie2.People = append(specie2.People, character2)
		}
	}

	return specie2
}

func MySQLGetSpecieBy_id(_id string) *model2.Specie {
	specie1 := &model1.Species{}
	has, err := Engine.Where("_id = ?", _id).Get(specie1)
	CheckErr(err)
	specie2 := &model2.Specie{}
	if has {
		specie2 = mysqlGetSpecie(specie1)
	}
	return specie2
}

func MySQLGetSpecieByID(ID string) *model2.Specie {
	specie1 := &model1.Species{
		ID: ID,
	}
	has, err := Engine.Get(specie1)
	CheckErr(err)
	specie2 := &model2.Specie{}
	if has {
		specie2 = mysqlGetSpecie(specie1)
	}
	return specie2
}

func mysqlGetStarship(starship1 *model1.Starship) *model2.Starship {
	var has bool
	var err error
	starship2 := convertStarship(starship1)

	for _, filmID := range starship1.Films {
		film1 := &model1.Film{
			ID: filmID,
		}
		has, err = Engine.Get(film1)
		CheckErr(err)
		if has {
			film2 := convertFilm(film1)
			starship2.Films = append(starship2.Films, film2)
		}
	}

	for _, PeopleID := range starship1.Pilots {
		character1 := &model1.People{
			ID: PeopleID,
		}
		has, err = Engine.Get(character1)
		CheckErr(err)
		if has {
			character2 := convertPeople(character1)
			starship2.Pilots = append(starship2.Pilots, character2)
		}
	}

	return starship2
}

func MySQLGetStarshipBy_id(_id string) *model2.Starship {
	starship1 := &model1.Starship{}
	has, err := Engine.Where("_id = ?", _id).Get(starship1)
	CheckErr(err)
	starship2 := &model2.Starship{}
	if has {
		starship2 = mysqlGetStarship(starship1)
	}
	return starship2
}

func MySQLGetStarshipByID(ID string) *model2.Starship {
	starship1 := &model1.Starship{
		ID: ID,
	}
	has, err := Engine.Get(starship1)
	CheckErr(err)
	starship2 := &model2.Starship{}
	if has {
		starship2 = mysqlGetStarship(starship1)
	}
	return starship2
}

func mysqlGetVehicle(vehicle1 *model1.Vehicle) *model2.Vehicle {
	var has bool
	var err error
	vehicle2 := convertVehicle(vehicle1)

	for _, filmID := range vehicle1.Films {
		film1 := &model1.Film{
			ID: filmID,
		}
		has, err = Engine.Get(film1)
		CheckErr(err)
		if has {
			film2 := convertFilm(film1)
			vehicle2.Films = append(vehicle2.Films, film2)
		}
	}

	for _, PeopleID := range vehicle1.Pilots {
		character1 := &model1.People{
			ID: PeopleID,
		}
		has, err = Engine.Get(character1)
		CheckErr(err)
		if has {
			character2 := convertPeople(character1)
			vehicle2.Pilots = append(vehicle2.Pilots, character2)
		}
	}

	return vehicle2
}

func MySQLGetVehicleBy_id(_id string) *model2.Vehicle {
	vehicle1 := &model1.Vehicle{}
	has, err := Engine.Where("_id = ?", _id).Get(vehicle1)
	CheckErr(err)
	vehicle2 := &model2.Vehicle{}
	if has {
		vehicle2 = mysqlGetVehicle(vehicle1)
	}
	return vehicle2
}

func MySQLGetVehicleByID(ID string) *model2.Vehicle {
	vehicle1 := &model1.Vehicle{
		ID: ID,
	}
	has, err := Engine.Get(vehicle1)
	CheckErr(err)
	vehicle2 := &model2.Vehicle{}
	if has {
		vehicle2 = mysqlGetVehicle(vehicle1)
	}
	return vehicle2
}
