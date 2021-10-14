// Package kmeans implements the k-means clustering algorithm
// See: https://en.wikipedia.org/wiki/K-means_clustering
package main

import (
	"fmt"
	"math/rand"
	"sort"
)

// Kmeans configuration/option struct
type Kmeans struct {
	// when a plotter is set, Plot gets called after each iteration
	plotter Plotter
	// deltaThreshold (in percent between 0.0 and 0.1) aborts processing if
	// less than n% of data points shifted clusters in the last iteration
	deltaThreshold float64
	// iterationThreshold aborts processing when the specified amount of
	// algorithm iterations was reached
	iterationThreshold int
}

// The Plotter interface lets you implement your own plotters
type Plotter interface {
	Plot(cc Clusters, iteration int) error
	Plot2(cc Clusters, iteration int) error
}

// NewWithOptions returns a Kmeans configuration struct with custom settings
func NewWithOptions(deltaThreshold float64, plotter Plotter) (Kmeans, error) {
	if deltaThreshold <= 0.0 || deltaThreshold >= 1.0 {
		return Kmeans{}, fmt.Errorf("threshold is out of bounds (must be >0.0 and <1.0, in percent)")
	}

	return Kmeans{
		plotter:            plotter,
		deltaThreshold:     deltaThreshold,
		iterationThreshold: 96,
	}, nil
}

// New returns a Kmeans configuration struct with default settings
func NewK() Kmeans {
	m, _ := NewWithOptions(0.01, nil)
	return m
}

// Partition executes the k-means algorithm on the given dataset and
// partitions it into k clusters
func (m Kmeans) Partition(dataset Observations, k int) (Clusters, error) {
	if k > len(dataset) {
		return Clusters{}, fmt.Errorf("the size of the data set must at least equal k")
	}

	cc, err := New(k, dataset)
	if err != nil {
		return cc, err
	}
	// fmt.Println(cc)
	points := make([]int, len(dataset))
	changes := 1

	for i := 0; changes > 0; i++ {
		changes = 0
		cc.Reset()

		// HERE

		for p, point := range dataset {
			ci := cc.Nearest(point)
			//HERE!!!

			cc[ci].Append(point)
			if points[p] != ci {
				points[p] = ci
				changes++
			}
		}

		// fmt.Println(len(cc))
		for ci := 0; ci < len(cc); ci++ {
			if len(cc[ci].Observations) == 0 {
				// During the iterations, if any of the cluster centers has no
				// data points associated with it, assign a random data point (HERE, why random?)
				// to it.
				// Also see: http://user.ceng.metu.edu.tr/~tcan/ceng465_f1314/Schedule/KMeansEmpty.html
				var ri int
				for {
					// find a cluster with at least two data points, otherwise
					// we're just emptying one cluster to fill another
					ri = rand.Intn(len(dataset))
					if len(cc[points[ri]].Observations) > 1 {
						break
					}
				}
				cc[ci].Append(dataset[ri])
				points[ri] = ci

				// Ensure that we always see at least one more iteration after
				// randomly assigning a data point to a cluster
				changes = len(dataset)
			}
		}

		//HERE
		if changes > 0 { //&& (isi cc<min && >max)
			cc.Recenter()
		}
		if m.plotter != nil {
			err := m.plotter.Plot(cc, i)
			if err != nil {
				return nil, fmt.Errorf("failed to plot chart: %s", err)
			}
		}
		if i == m.iterationThreshold ||
			changes < int(float64(len(dataset))*m.deltaThreshold) {
			// fmt.Println("Aborting:", changes, int(float64(len(dataset))*m.TerminationThreshold))
			break
		}
	}

	return cc, nil
}

func (m *Kmeans) NewPartition(dataset Observations, k int, seed int64) (Clusters, error) {
	if k > len(dataset) {
		return Clusters{}, fmt.Errorf("the size of the data set must at least equal to k (%d)", k)
	}

	cc, err := NewClusters(seed, k, dataset)
	if err != nil {
		return cc, err
	}

	points := make([]int, len(dataset))
	changes := 1

	for i := 0; changes > 0; i++ {
		changes = 0
		cc.Reset()

		for p, point := range dataset {
			ci := cc.Nearest(point)
			cc[ci].Append(point)
			if points[p] != ci {
				points[p] = ci
				changes++
			}
		}

		for ci := 0; ci < len(cc); ci++ {
			if len(cc[ci].Observations) == 0 {
				// During the iterations, if any of the cluster centers has no
				// data points associated with it, assign a random data point
				// to it.
				// Also see: http://user.ceng.metu.edu.tr/~tcan/ceng465_f1314/Schedule/KMeansEmpty.html
				var ri int
				for {
					// find a cluster with at least two data points, otherwise
					// we're just emptying one cluster to fill another
					ri = rand.Intn(len(dataset))
					if len(cc[points[ri]].Observations) > 1 {
						break
					}
				}
				cc[ci].Append(dataset[ri])
				points[ri] = ci

				// Ensure that we always see at least one more iteration after
				// randomly assigning a data point to a cluster
				changes = len(dataset)
			}
		}

		if changes > 0 {
			cc.Recenter()
		}
		if m.plotter != nil {
			err := m.plotter.Plot(cc, i)
			if err != nil {
				return nil, fmt.Errorf("failed to plot chart: %s", err)
			}
		}
		if i == m.iterationThreshold ||
			changes < int(float64(len(dataset))*m.deltaThreshold) {
			// return Clusters{}, fmt.Errorf("iteration threshold '%d' reached", m.iterationThreshold)
			break
		}
	}

	return cc, nil
}

func (c Clusters) newborderadjust(A int, B int) (Observations, []diffsort, []int) {
	// For each point p in area A
	var diff []diffsort
	var clust []int
	var diffB []diffsort
	var diffA []diffsort
	var obsA Observations
	var obsB Observations
	var r int

	for _, p := range c[A].Observations {
		r, _ = c.Neighbour(p, A)
		distA := p.Distance(c[A].Center)
		distB := p.Distance(c[B].Center)
		//  Calculate diff(p, B) based on (2);
		diff = append(diff, diffsort{distB - distA, p})
		clust = append(clust, r)
	}

	//  Move the first m point in area A based on sorted
	// diff(p, B) to area B;
	// n := rand.Intn(103-101) + 101
	chunkSize := (len(d) + 20 - 1) / 20
	if len(c[A].Observations) > chunkSize && len(c[A].Observations) > len(c[B].Observations) {
		m := len(c[A].Observations) - chunkSize
		//  Sort all the diff(p, B) ascending;
		sort.SliceStable(diff, func(i, j int) bool {
			return diff[i].differ < diff[j].differ
		})
		for i := 0; i < m; i++ {
			// move to B
			diffB = append(diffB, diff[i])
		}
		diffA = diff[m:]
		for i := 0; i < len(diffA); i++ {
			obsA = append(obsA, diffA[i].data)
		}
		for i := 0; i < len(diffB); i++ {
			obsB = append(obsB, diffB[i].data)
		}
	}
	return obsA, diffB, clust
}
