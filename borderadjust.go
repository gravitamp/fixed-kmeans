package main

import (
	"math/rand"
	"sort"
)

type diffsort struct {
	differ float64
	data   Observation
}

func (c Clusters) borderadjust(A int) {
	// rand.Seed(time.Now().UnixNano())
	rand.Seed(20)
	// For each point p in area A
	var diff []diffsort
	var diffB []diffsort
	var diffA []diffsort
	var obsA Observations
	var obsB Observations
	var r int
	var clust []int

	for _, p := range c[A].Observations {
		r, _ = c.Neighbour(p, A)
		distA := p.Distance(c[A].Center)
		distB := p.Distance(c[r].Center)
		//  Calculate diff(p, B) based on (2);
		diff = append(diff, diffsort{distB - distA, p})
		clust = append(clust, r)
	}
	//  Sort all the diff(p, B) ascending;
	sort.SliceStable(diff, func(i, j int) bool {
		return diff[i].differ < diff[j].differ
	})
	//  Move the first m point in area A based on sorted
	// diff(p, B) to area B;
	// n := rand.Intn(103-101) + 101
	chunkSize := (len(d) + 20 - 1) / 20
	if len(c[A].Observations) > chunkSize { // && len(c[A].Observations) > len(c[r].Observations) {
		m := len(c[A].Observations) - chunkSize
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
	// if len(obsA) != 0 && len(diffB) != 0 {
	// 	c[A].Observations = obsA
	// 	for j := 0; j < len(diffB); j++ {
	// 		d := clust[j]		//neighbour cluster
	// 		c[d].Observations = append(c[d].Observations, diffB[j].data)
	// 	}
	// }
	if len(diffA) != 0 && len(diffB) != 0 {
		c[A].Observations = obsA
		if A < (len(c) - 1) {
			for j := 0; j < len(diffB); j++ {
				c[A+1].Observations = append(c[A+1].Observations, diffB[j].data)
			}
		} else {
			for j := 0; j < len(diffB); j++ {
				c[0].Observations = append(c[0].Observations, diffB[j].data)
			}
		}
	}
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
