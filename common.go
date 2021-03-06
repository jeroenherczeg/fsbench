package fsbench

import (
	"math/rand"
	"time"

	"github.com/dripolles/histogram"
	"github.com/mxk/go-flowrate/flowrate"
)

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

type Status struct {
	Duration time.Duration // Time period covered by the statistics
	Bytes    int64         // Total number of bytes transferred
	AvgRate  int64         // Average transfer rate (Bytes / Duration)
}

func NewStatus(fs flowrate.Status) Status {
	return Status{
		Duration: fs.Duration,
		Bytes:    fs.Bytes,
		AvgRate:  fs.AvgRate,
	}
}

type AggregatedStatus struct {
	Status
	Files            int                  // Number of files transferred
	Errors           int                  // Number of errors
	HistogramAvgRate *histogram.Histogram `json:"-"`
}

func NewAggregatedStatus() *AggregatedStatus {
	return &AggregatedStatus{
		HistogramAvgRate: histogram.NewHistogram(),
	}
}

func (s *AggregatedStatus) Add(a Status) {
	s.Files++
	s.Bytes += a.Bytes
	s.Duration += a.Duration
	s.AvgRate = int64(float64(s.Bytes) / s.Duration.Seconds())
	s.HistogramAvgRate.Add((int(float64(s.Bytes) / s.Duration.Seconds())))
}

func (s *AggregatedStatus) Sum(a *AggregatedStatus) {
	s.Files += a.Files
	s.Errors += a.Errors
	s.Bytes += a.Bytes
	s.Duration += a.Duration
	s.AvgRate = int64(float64(s.Bytes) / s.Duration.Seconds())
	s.HistogramAvgRate.Update(a.HistogramAvgRate)
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// radomString generates a random string of any length, code extracted from:
// http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func randomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
