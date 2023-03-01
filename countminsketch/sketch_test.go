package countminsketch

import (
	"io/ioutil"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var (
	normalSamples   []string
	overflowSamples []string
	fileSamples     []string
)

func init() {
	sample := "a b b c c c e e e e e d d d d g g g g g g g f f f f f f"
	normalSamples = strings.Split(sample, " ")

	overflowSamples = strings.Split(sample, " ")
	for i := 0; i < 200; i++ {
		overflowSamples = append(overflowSamples, "a")
	}

	// tr -c '[:alnum:]' '[\n*]' < testdata/lorem_ipsum.txt | sort | uniq -c | sort -nr | head  -6
	bytes, _ := ioutil.ReadFile("../testdata/lorem_ipsum.txt")
	nonAlnum := regexp.MustCompile("[^[:alnum:]]+")
	content := nonAlnum.ReplaceAllString(string(bytes), " ")
	fileSamples = strings.Split(content, " ")
}

func TestCountMinSketch(t *testing.T) {
	type args struct {
		epsilon float64
		k       int
		data    []string
	}
	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		{
			"normal",
			args{
				0.01,
				3,
				normalSamples,
			},
			map[string]int{"g": 7, "f": 6, "e": 5},
		},
		{
			"overflow",
			args{
				0.01,
				3,
				overflowSamples,
			},
			map[string]int{"a": 201, "g": 7, "f": 6},
		},
		{
			"file",
			args{
				0.002,
				5,
				fileSamples,
			},
			map[string]int{"et": 43},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(WithEstimates(tt.args.epsilon, 0.01), tt.args.k)
			for _, v := range tt.args.data {
				s.Incr(v)
			}
			if got := s.TopK(); !IsSubset(got, tt.want) {
				t.Errorf("CountMinSketch.TopK() = %v, want %v", got, tt.want)
			}
			t.Logf("CountMinSketch.TopK() = %v, want %v", s.TopK(), tt.want)
		})
	}
}

func IsSubset[K, V comparable](m, sub map[K]V) bool {
	if len(sub) > len(m) {
		return false
	}
	for k := range sub {
		if _, found := m[k]; !found {
			return false
		}
	}
	return true
}

func Benchmark_CMS_Incr_ε0_001_δ0_1(b *testing.B) {
	s := WithEstimates(0.001, 0.1)
	summary := New(s, 50)
	for i := 0; i < b.N; i++ {
		summary.Incr(strconv.Itoa(int(rand.Int31())))
	}
}
