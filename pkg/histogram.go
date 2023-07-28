package tq

import (
	"fmt"
	"github.com/ekimeel/sabal-pb/pb"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	histogramLimit = 50
	binSize        = 25
)

type Histogram struct {
	Id          uint32         `json:"id"`
	PointId     uint32         `json:"point-id"`
	LastUpdated time.Time      `json:"last-updated"`
	KeyCount    int            `json:"key-count"`
	ValueCount  int            `json:"value-count"`
	Histogram   map[string]int `json:"histogram"`
}

type bin struct {
	start float64
	end   float64
	count int
}

func (hist *Histogram) update(metrics []*pb.Metric) {
	if hist.Histogram == nil {
		hist.Histogram = make(map[string]int, 0)
	}

	for i := range metrics {
		hist.Histogram[fmt.Sprintf("%g", metrics[i].Value)] += 1
	}

	if len(hist.Histogram) >= histogramLimit {
		hist.bin(binSize)
	}

	hist.KeyCount, hist.ValueCount = hist.binStats()
}

func (hist *Histogram) binStats() (int, int) {
	var totalCount int
	binCount := len(hist.Histogram)

	for _, count := range hist.Histogram {
		totalCount += count
	}

	return binCount, totalCount
}

func roundToSigFigs(num float64, sigFigs int) float64 {
	if num == 0 {
		return 0
	}

	magnitude := math.Floor(math.Log10(math.Abs(num)))
	scale := math.Pow(10, float64(sigFigs)-1-magnitude)

	return math.Round(num*scale) / scale
}

func (hist *Histogram) bin(maxBins int) {
	originalBins := make([]bin, 0, len(hist.Histogram))

	for k, v := range hist.Histogram {
		if strings.Contains(k, "-") {
			split := strings.Split(k, "-")
			start, _ := strconv.ParseFloat(split[0], 64)
			end, _ := strconv.ParseFloat(split[1], 64)
			originalBins = append(originalBins, bin{start: start, end: end, count: v})
		} else {
			val, _ := strconv.ParseFloat(k, 64)
			originalBins = append(originalBins, bin{start: val, end: val, count: v})
		}
	}

	sort.Slice(originalBins, func(i, j int) bool {
		return originalBins[i].start < originalBins[j].start
	})

	newHist := make(map[string]int)

	if len(originalBins) > 0 {
		minVal := originalBins[0].start
		maxVal := originalBins[len(originalBins)-1].end
		binSize := (maxVal - minVal) / float64(maxBins)

		for _, bin := range originalBins {
			newBinStart := math.Floor((bin.start-minVal)/binSize)*binSize + minVal
			newBinEnd := newBinStart + binSize
			newBinStart = roundToSigFigs(newBinStart, 3) // Round to 3 significant figures
			newBinEnd = roundToSigFigs(newBinEnd, 3)     // Round to 3 significant figures
			newBinKey := fmt.Sprintf("%g-%g", newBinStart, newBinEnd)
			newHist[newBinKey] += bin.count
		}
	}

	hist.Histogram = newHist
}
