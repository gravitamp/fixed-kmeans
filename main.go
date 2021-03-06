// begin
//    specify the number k of clustering to assign.
//    randomly initialize k centroids.
//    repeat
//       expectation: Assign each point to its closest centroid.
//       maximization: Compute the new centroid (mean) of each cluster.
//    until The centroid position do not change.
// end
// Clustering with Constrained Problem for cluster result to have an equal number of member cluster.
// must learn weighted clustering

package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

var d Observations
var count []int

func main() {
	//setup data
	setupData("Traffic4.csv")
	// Partition the data points into 20 clusters
	km, _ := NewWithOptions(0.01, SimplePlotter{})
	// clusters, _ := km.Partition(d, 20)
	clusters, _ := km.NewPartition(d, 20, 20)

	for _, c := range clusters {
		count = append(count, len(c.Observations))
	}
	for i, c := range clusters {
		fmt.Printf("Centered at x: %.2f y: %.2f\n", c.Center[0], c.Center[1])
		// fmt.Printf("Matching data points: %+v\n", c.Observations)
		fmt.Printf("total %d: %d\n", i, len(c.Observations))
	}

	fmt.Println(sum(count))
	min, max := MinMax(count)
	fmt.Println(min, max)
	iter := 0
	for max-min > 10 {
		var count3 []int
		//  Plan the steps of adjustment among clusters;
		for i := 0; i < len(clusters); i++ {
			if len(clusters[i].Observations) > 102 {
				var diffA Observations
				var diffB []diffsort
				var clust []int
				// call borderadjust, get new cluster A & B
				// clusters.borderadjust(i)

				//call borderadjust, get new cluster A & B
				r, _ := clusters.Neighbour(clusters[i].Center, i)
				diffA, diffB, clust = clusters.newborderadjust(i, r)

				//FIX THIS (what?)
				if len(diffA) == 0 && len(diffB) == 0 {
					continue
				} else if len(diffA) != 0 && len(diffB) != 0 {
					clusters[i].Observations = diffA
					for j := 0; j < len(diffB); j++ {
						c := clust[j]
						clusters[c].Observations = append(clusters[c].Observations, diffB[j].data)
					}
				}
				clusters.Recenter()
			}
			// 	if len(diffA) == 0 && len(diffB) == 0 {
			// 		continue
			// 	} else if len(diffA) != 0 && len(diffB) != 0 {
			// 		clusters[i].Observations = diffA
			// 		if i < (len(clusters) - 1) {
			// 			for j := 0; j < len(diffB); j++ {
			// 				clusters[i+1].Observations = append(clusters[i+1].Observations, diffB[j].data)
			// 			}
			// 		} else {
			// 			for j := 0; j < len(diffB); j++ {
			// 				clusters[0].Observations = append(clusters[0].Observations, diffB[j].data)
			// 			}
			// 		}
			// 	}
			// 	clusters.Recenter()
			// }
		}
		// recenter

		iter++
		fmt.Println("iterasi ke-", iter)

		for i := 0; i < len(clusters); i++ {
			count3 = append(count3, len(clusters[i].Observations))
		}
		min, max = MinMax(count3)
		fmt.Println(min, max)
		fmt.Println("jarak", max-min)

		//plot
		if km.plotter != nil {
			err := km.plotter.Plot2(clusters, iter)
			if err != nil {
				return //nil, fmt.Errorf("failed to plot chart: %s", err)
			}
		}
		// if max-min > 11 {
		// 	break
		// }
		// max iter
		if iter == 10 {
			break
		}
		// if km.plotter != nil {
		// 	err := km.plotter.Plot2(clusters, 1)
		// 	if err != nil {
		// 		return //nil, fmt.Errorf("failed to plot chart: %s", err)
		// 	}
	}

	var count2 = 0
	//get balanced cluster
	for i, c := range clusters {
		fmt.Printf("Centered at x: %.2f y: %.2f\n", c.Center[0], c.Center[1])
		// fmt.Printf("Matching data points: %+v\n", c.Observations)
		fmt.Printf("total %d: %d\n", i, len(c.Observations))
		count2 += len(c.Observations)
	}
	fmt.Println(count2)
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

func sum(arr []int) int {
	var res int
	res = 0
	for i := 0; i < len(arr); i++ {
		res += arr[i]
	}
	return res
}

func MinMax(array []int) (int, int) {
	var max int = array[0]
	var min int = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}
