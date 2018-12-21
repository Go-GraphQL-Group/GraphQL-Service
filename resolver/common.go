package resolver

import (
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
)

const peopleBucket = "People"
const filmsBucket = "Film"
const planetsBucket = "Planet"
const speciesBucket = "Specie"
const starshipsBucket = "Starship"
const vehiclesBucket = "Vehicle"

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func encodeCursor(k string) string {
	i, _ := strconv.Atoi(k)
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor%d", i)))
}
