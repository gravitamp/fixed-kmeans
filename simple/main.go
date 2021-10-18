package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

var d Observations

func main() {
	rand.Seed(20)

	setupData("Traffic4.csv")

	// set up a random two-dimensional data set (float64 values between 0.0 and 1.0)

	fmt.Printf("%d data points\n", len(d))

	// Partition the data points into 7 clusters
	km, _ := NewWithOptions(0.01, SimplePlotter{})
	clusters, _ := km.Partition(d, 20)

	for i, c := range clusters {
		fmt.Printf("Cluster: %d\n", i)
		fmt.Printf("Centered at x: %.2f y: %.2f\n", c.Center[0], c.Center[1])
	}
}

func setupData(file string) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	csvReader := csv.NewReader(f)
	csvData, _ := csvReader.ReadAll()

	//read without header
	for i := 1; i < len(csvData); i++ {
		val, _ := strconv.Atoi(csvData[i][3])
		for j := 0; j < val; j++ {
			lat, _ := strconv.ParseFloat(csvData[i][1], 64)
			lng, _ := strconv.ParseFloat(csvData[i][2], 64)
			d = append(d, Coordinates{
				lng,
				lat,
			})
		}

	}
}
