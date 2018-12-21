package db

import (
	"encoding/json"
	"fmt"
	"regexp"

	model2 "github.com/Go-GraphQL-Group/GraphQL-Service/model"
	model1 "github.com/Go-GraphQL-Group/SW-Crawler/model"
	"github.com/boltdb/bolt"
)

const peopleBucket = "People"
const filmsBucket = "Film"
const planetsBucket = "Planet"
const speciesBucket = "Specie"
const starshipsBucket = "Starship"
const vehiclesBucket = "Vehicle"

// 正则替换
const origin = "http://localhost:8080/query/+[a-zA-Z_]+/"

const replace = ""

const preURL = "http://localhost:8080/query/"

// convert database model to graphql model

func convertPeople(people1 *model1.People) *model2.People {
	people2 := &model2.People{}
	people2.ID = people1.ID
	people2.Name = people1.Name
	people2.BirthYear = &people1.Birth_year
	people2.EyeColor = &people1.Eye_color
	people2.Gender = &people1.Gender
	people2.HairColor = &people1.Hair_color
	people2.Height = &people1.Heigth
	people2.Mass = &people1.Mass
	people2.SkinColor = &people1.Skin_color
	return people2
}

func convertPlanet(planet1 *model1.Planet) *model2.Planet {
	planet2 := &model2.Planet{}
	planet2.ID = planet1.ID
	planet2.Name = planet1.Name
	planet2.Diameter = &planet1.Diameter
	planet2.RotationPeriod = &planet1.Rotation_period
	planet2.OrbitalPeriod = &planet1.Orbital_period
	planet2.Gravity = &planet1.Gravity
	planet2.Population = &planet1.Population
	planet2.Climate = &planet1.Climate
	planet2.Terrain = &planet1.Terrain
	planet2.SurfaceWater = &planet1.Surface_water
	return planet2
}

func convertFilm(film1 *model1.Film) *model2.Film {
	film2 := &model2.Film{}
	film2.ID = film1.ID
	film2.Title = film1.Title
	film2.EpisodeID = &film1.Episode_id
	film2.OpeningCrawl = &film1.Opening_crawl
	film2.Director = &film1.Director
	film2.Producer = &film1.Producer
	film2.ReleaseDate = &film1.Release_data
	return film2
}

func convertSpecie(specie1 *model1.Species) *model2.Specie {
	specie2 := &model2.Specie{}
	specie2.ID = specie1.ID
	specie2.Name = specie1.Name
	specie2.Classification = &specie1.Classification
	specie2.AverageHeight = &specie1.Average_height
	specie2.AverageLifespan = &specie1.Average_lifespan
	specie2.EyeColors = &specie1.Eye_colors
	specie2.HairColors = &specie1.Hair_colors
	specie2.SkinColors = &specie1.Skin_colors
	specie2.Language = &specie1.Language
	specie2.Designation = &specie1.Designation
	return specie2
}

func convertStarship(starship1 *model1.Starship) *model2.Starship {
	starship2 := &model2.Starship{}
	starship2.ID = starship1.ID
	starship2.Name = starship1.Name
	starship2.Model = &starship1.Model
	starship2.StarshipClass = &starship1.Starship_class
	starship2.Manufacturer = &starship1.Manufacturer
	starship2.CostInCredits = &starship1.Cost_in_credits
	starship2.Length = &starship1.Length
	starship2.Crew = &starship1.Crew
	starship2.Passengers = &starship1.Passenger
	starship2.MaxAtmospheringSpeed = &starship1.Max_atmosphering_speed
	starship2.HyperdriveRating = &starship1.Hyperdrive_rating
	starship2.MGLT = &starship1.MGLT
	starship2.CargoCapacity = &starship1.Cargo_capacity
	starship2.Consumables = &starship1.Consumables
	return starship2
}

func convertVehicle(vehicle1 *model1.Vehicle) *model2.Vehicle {
	vehicle2 := &model2.Vehicle{}
	vehicle2.ID = vehicle1.ID
	vehicle2.Name = vehicle1.Name
	vehicle2.Model = &vehicle1.Model
	vehicle2.Manufacturer = &vehicle1.Manufacturer
	vehicle2.Length = &vehicle1.Length
	vehicle2.CostInCredits = &vehicle1.Cost_in_credits
	vehicle2.Crew = &vehicle1.Crew
	vehicle2.Passengers = &vehicle1.Passenger
	vehicle2.MaxAtmospheringSpeed = &vehicle1.Max_atmosphering_speed
	vehicle2.CargoCapacity = &vehicle1.Cargo_capacity
	vehicle2.Consumables = &vehicle1.Consumables
	return vehicle2
}

func GetPeopleByID(ID string, db *bolt.DB) (error, *model2.People) {
	var err error
	if db == nil {
		db, err = bolt.Open("data/data.db", 0600, nil)
		CheckErr(err)
		defer db.Close()
	}
	people1 := &model2.People{}

	people := &model1.People{}
	err = db.View(func(tx *bolt.Tx) error {
		peoBuck := tx.Bucket([]byte(peopleBucket))
		v := peoBuck.Get([]byte(ID))

		if v == nil {
			return err
		}

		// 正则替换
		re, _ := regexp.Compile(origin)
		rep := re.ReplaceAllString(string(v), replace)
		err = json.Unmarshal([]byte(rep), people)
		CheckErr(err)
		people1 = convertPeople(people)

		// Home
		homeId := people.Homeworld
		homeId = homeId[0 : len(homeId)-1]
		homeBuck := tx.Bucket([]byte(planetsBucket))
		planet := homeBuck.Get([]byte(homeId))
		planet1 := &model1.Planet{}
		err = json.Unmarshal([]byte(planet), planet1)
		people1.Homeworld = convertPlanet(planet1)

		// films
		for _, it := range people.Films {
			it = it[0 : len(it)-1]
			filmBuck := tx.Bucket([]byte(filmsBucket))
			film := filmBuck.Get([]byte(it))
			film1 := &model1.Film{}
			err = json.Unmarshal([]byte(film), film1)
			CheckErr(err)
			people1.Films = append(people1.Films, convertFilm(film1))
		}

		// species
		for _, it := range people.Species {
			it = it[0 : len(it)-1]
			specBuck := tx.Bucket([]byte(speciesBucket))
			specie := specBuck.Get([]byte(it))
			specie1 := &model1.Species{}
			err = json.Unmarshal([]byte(specie), specie1)
			CheckErr(err)
			people1.Species = append(people1.Species, convertSpecie(specie1))
		}

		// vehicles
		for _, it := range people.Vehicles {
			it = it[0 : len(it)-1]
			vehBuck := tx.Bucket([]byte(vehiclesBucket))
			vehicle := vehBuck.Get([]byte(it))
			vehicle1 := &model1.Vehicle{}
			err = json.Unmarshal([]byte(vehicle), vehicle1)
			CheckErr(err)
			people1.Vehicles = append(people1.Vehicles, convertVehicle(vehicle1))
		}

		// starships
		for _, it := range people.Starships {
			it = it[0 : len(it)-1]
			starBuck := tx.Bucket([]byte(starshipsBucket))
			starship := starBuck.Get([]byte(it))
			starship1 := &model1.Starship{}
			err = json.Unmarshal([]byte(starship), starship1)
			CheckErr(err)
			people1.Starships = append(people1.Starships, convertStarship(starship1))
		}

		people.Url = preURL + "people/" + people.Url
		return nil
	})
	return err, people1
}

func GetFilmByID(ID string, db *bolt.DB) (error, *model2.Film) {
	var err error
	if db == nil {
		db, err = bolt.Open("data/data.db", 0600, nil)
		CheckErr(err)
		defer db.Close()
	}

	film1 := &model2.Film{}
	film := &model1.Film{}
	err = db.View(func(tx *bolt.Tx) error {
		filmBuck := tx.Bucket([]byte(filmsBucket))
		v := filmBuck.Get([]byte(ID))

		if v == nil {
			return err
		}

		// 正则替换
		re, _ := regexp.Compile(origin)
		rep := re.ReplaceAllString(string(v), replace)
		err = json.Unmarshal([]byte(rep), film)
		CheckErr(err)

		film1 = convertFilm(film)

		// character
		for _, it := range film.Character {
			it = it[0 : len(it)-1]
			peoBuck := tx.Bucket([]byte(peopleBucket))
			people := peoBuck.Get([]byte(it))
			people1 := &model1.People{}
			err = json.Unmarshal([]byte(people), people1)
			CheckErr(err)
			film1.Characters = append(film1.Characters, convertPeople(people1))
		}

		// planet
		for _, it := range film.Planets {
			it = it[0 : len(it)-1]
			plaBuck := tx.Bucket([]byte(planetsBucket))
			planet := plaBuck.Get([]byte(it))
			planet1 := &model1.Planet{}
			err = json.Unmarshal([]byte(planet), planet1)
			CheckErr(err)
			film1.Planets = append(film1.Planets, convertPlanet(planet1))
		}

		// starship
		for _, it := range film.Starships {
			it = it[0 : len(it)-1]
			starBuck := tx.Bucket([]byte(starshipsBucket))
			starship := starBuck.Get([]byte(it))
			starship1 := &model1.Starship{}
			err = json.Unmarshal([]byte(starship), starship1)
			CheckErr(err)
			film1.Starships = append(film1.Starships, convertStarship(starship1))
		}

		// vehicle
		for _, it := range film.Vehicles {
			it = it[0 : len(it)-1]
			vehBuck := tx.Bucket([]byte(vehiclesBucket))
			vehicle := vehBuck.Get([]byte(it))
			vehicle1 := &model1.Vehicle{}
			err = json.Unmarshal([]byte(vehicle), vehicle1)
			CheckErr(err)
			film1.Vehicles = append(film1.Vehicles, convertVehicle(vehicle1))
		}

		// specie
		for _, it := range film.Species {
			it = it[0 : len(it)-1]
			specBuck := tx.Bucket([]byte(speciesBucket))
			specie := specBuck.Get([]byte(it))
			specie1 := &model1.Species{}
			err = json.Unmarshal([]byte(specie), specie1)
			CheckErr(err)
			film1.Species = append(film1.Species, convertSpecie(specie1))
		}

		film.Url = preURL + "films/" + film.Url
		return nil
	})
	return err, film1
}

func GetPlanetByID(ID string, db *bolt.DB) (error, *model2.Planet) {
	var err error
	if db == nil {
		db, err = bolt.Open("data/data.db", 0600, nil)
		CheckErr(err)
		defer db.Close()
	}

	planet1 := &model2.Planet{}
	planet := &model1.Planet{}
	err = db.View(func(tx *bolt.Tx) error {
		plaBuck := tx.Bucket([]byte(planetsBucket))
		v := plaBuck.Get([]byte(ID))

		if v == nil {
			return err
		}

		// 正则替换
		re, _ := regexp.Compile(origin)
		rep := re.ReplaceAllString(string(v), replace)
		err = json.Unmarshal([]byte(rep), planet)
		CheckErr(err)
		planet1 = convertPlanet(planet)

		// resident
		for _, it := range planet.Residents {
			it = it[0 : len(it)-1]
			peoBuck := tx.Bucket([]byte(peopleBucket))
			people := peoBuck.Get([]byte(it))
			people1 := &model1.People{}
			err = json.Unmarshal([]byte(people), people1)
			CheckErr(err)
			planet1.Residents = append(planet1.Residents, convertPeople(people1))
		}

		// film
		for _, it := range planet.Films {
			it = it[0 : len(it)-1]
			filmBuck := tx.Bucket([]byte(filmsBucket))
			film := filmBuck.Get([]byte(it))
			film1 := &model1.Film{}
			err = json.Unmarshal([]byte(film), film1)
			CheckErr(err)
			planet1.Films = append(planet1.Films, convertFilm(film1))
		}
		return nil
	})
	return err, planet1
}

func GetSpeciesByID(ID string, db *bolt.DB) (error, *model2.Specie) {
	var err error
	if db == nil {
		db, err = bolt.Open("data/data.db", 0600, nil)
		CheckErr(err)
		defer db.Close()
	}

	specie1 := &model2.Specie{}
	specie := &model1.Species{}
	err = db.View(func(tx *bolt.Tx) error {
		specBuck := tx.Bucket([]byte(speciesBucket))
		v := specBuck.Get([]byte(ID))

		if v == nil {
			return err
		}

		// 正则替换
		re, _ := regexp.Compile(origin)
		rep := re.ReplaceAllString(string(v), replace)
		// 替换后结果
		// fmt.Println("After regexp: " + rep)

		err = json.Unmarshal([]byte(rep), specie)
		CheckErr(err)

		specie1 = convertSpecie(specie)

		// Home
		homeId := specie.Homeworld
		homeId = homeId[0 : len(homeId)-1]
		homeBuck := tx.Bucket([]byte(planetsBucket))
		planet := homeBuck.Get([]byte(homeId))
		planet1 := &model1.Planet{}
		err = json.Unmarshal([]byte(planet), planet1)
		CheckErr(err)
		specie1.Homeworld = convertPlanet(planet1)

		// people
		for _, it := range specie.People {
			it = it[0 : len(it)-1]
			peoBuck := tx.Bucket([]byte(peopleBucket))
			people := peoBuck.Get([]byte(it))
			people1 := &model1.People{}
			err = json.Unmarshal([]byte(people), people1)
			CheckErr(err)
			specie1.People = append(specie1.People, convertPeople(people1))
		}

		// film
		for _, it := range specie.Films {
			it = it[0 : len(it)-1]
			filmBuck := tx.Bucket([]byte(filmsBucket))
			film := filmBuck.Get([]byte(it))
			film1 := &model1.Film{}
			err = json.Unmarshal([]byte(film), film1)
			CheckErr(err)
			specie1.Films = append(specie1.Films, convertFilm(film1))
		}
		return nil
	})
	return err, specie1
}

func GetStarshipByID(ID string, db *bolt.DB) (error, *model2.Starship) {
	var err error
	if db == nil {
		db, err = bolt.Open("data/data.db", 0600, nil)
		CheckErr(err)
		defer db.Close()
	}

	starship1 := &model2.Starship{}
	starship := &model1.Starship{}
	err = db.View(func(tx *bolt.Tx) error {
		starBuck := tx.Bucket([]byte(starshipsBucket))
		v := starBuck.Get([]byte(ID))

		if v == nil {
			return err
		}

		// 正则替换
		re, _ := regexp.Compile(origin)
		rep := re.ReplaceAllString(string(v), replace)
		// 替换后结果
		// fmt.Println("After regexp: " + rep)

		err = json.Unmarshal([]byte(rep), starship)
		CheckErr(err)
		starship1 = convertStarship(starship)

		// people
		for _, it := range starship.Pilots {
			it = it[0 : len(it)-1]
			peoBuck := tx.Bucket([]byte(peopleBucket))
			people := peoBuck.Get([]byte(it))
			people1 := &model1.People{}
			err = json.Unmarshal([]byte(people), people1)
			CheckErr(err)
			starship1.Pilots = append(starship1.Pilots, convertPeople(people1))
		}

		// film
		for _, it := range starship.Films {
			it = it[0 : len(it)-1]
			filmBuck := tx.Bucket([]byte(filmsBucket))
			film := filmBuck.Get([]byte(it))
			film1 := &model1.Film{}
			err = json.Unmarshal([]byte(film), film1)
			CheckErr(err)
			starship1.Films = append(starship1.Films, convertFilm(film1))
		}
		return nil
	})
	return err, starship1
}

func GetVehicleByID(ID string, db *bolt.DB) (error, *model2.Vehicle) {
	var err error
	if db == nil {
		db, err = bolt.Open("data/data.db", 0600, nil)
		CheckErr(err)
		defer db.Close()
	}

	vehicle1 := &model2.Vehicle{}
	vehicle := &model1.Vehicle{}
	err = db.View(func(tx *bolt.Tx) error {
		vehicBuck := tx.Bucket([]byte(vehiclesBucket))
		v := vehicBuck.Get([]byte(ID))

		if v == nil {
			return err
		}

		// 正则替换
		re, _ := regexp.Compile(origin)
		rep := re.ReplaceAllString(string(v), replace)
		// 替换后结果
		// fmt.Println("After regexp: " + rep)

		err = json.Unmarshal([]byte(rep), vehicle)
		CheckErr(err)

		vehicle1 = convertVehicle(vehicle)

		// people
		for _, it := range vehicle.Pilots {
			it = it[0 : len(it)-1]
			peoBuck := tx.Bucket([]byte(peopleBucket))
			people := peoBuck.Get([]byte(it))
			people1 := &model1.People{}
			err = json.Unmarshal([]byte(people), people1)
			CheckErr(err)
			vehicle1.Pilots = append(vehicle1.Pilots, convertPeople(people1))
		}

		// film
		for _, it := range vehicle.Films {
			it = it[0 : len(it)-1]
			filmBuck := tx.Bucket([]byte(filmsBucket))
			film := filmBuck.Get([]byte(it))
			film1 := &model1.Film{}
			err = json.Unmarshal([]byte(film), film1)
			CheckErr(err)
			vehicle1.Films = append(vehicle1.Films, convertFilm(film1))

		}
		return nil
	})
	return err, vehicle1
}

// err
func CheckErr(err error) {
	if err != nil {
		fmt.Println("Error occur: ", err)
		// os.Exit(1)
	}
}
