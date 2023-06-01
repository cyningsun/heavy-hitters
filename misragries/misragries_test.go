package misragries

import (
	"os"
	"regexp"
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
	bytes, _ := os.ReadFile("../testdata/lorem_ipsum.txt")
	nonAlnum := regexp.MustCompile("[^[:alnum:]]+")
	content := nonAlnum.ReplaceAllString(string(bytes), " ")
	fileSamples = strings.Split(content, " ")
}

func TestMisraGries(t *testing.T) {
	type args struct {
		k    int
		data []string
	}

	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		{
			"normal",
			args{
				5,
				normalSamples,
			},
			map[string]int{"g": 7, "f": 6, "e": 5},
		},
		{
			"overflow",
			args{
				5,
				overflowSamples,
			},
			map[string]int{"a": 201, "g": 7, "f": 6},
		},
		{
			"file",
			args{
				60,
				fileSamples,
			},
			map[string]int{"et": 43},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mg := NewMisraGries(tt.args.k)
			for _, element := range tt.args.data {
				mg.ProcessElement(element)
			}
			if got := mg.TopK(); !IsSubset(got, tt.want) {
				t.Errorf("MisraGries.TopK() = %v, want %v", got, tt.want)
			}
			t.Logf("MisraGries.TopK() = %v, want %v", mg.TopK(), tt.want)
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
