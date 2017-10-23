package main

import (
	"gonum.org/v1/gonum/stat"
	"testing"
)

var testData = Data{
	X: []float64{-4.32679099, -5.15888702, 1.6152995, 12.27539741,
		10.97570561, 7.31776081, -8.59303456, 8.97395768,
		1.90158568, -1.09838831},
	Y: []float64{15.01726764, 12.97752127, -5.4206861, -1.73396994,
		15.52848773, -9.81706237, 1.64974166, -5.76083789,
		0.20345582, 10.57435605}}

func TestMinMax(t *testing.T) {
	min, max := MinMax(testData.X)
	if min != -8.59303456 {
		t.Error(`MinMax(testData.X) failed: min != -8.59303456`)
	}
	if max != 12.27539741 {
		t.Error(`MinMax(testData.X) failed: max != 12.27539741`)
	}
}

func TestRegression(t *testing.T) {
	a, b := stat.LinearRegression(testData.X, testData.Y, nil, false)
	if a <= 4.42 || a >= 4.43 {
		t.Error(`stat.LinearRegression(testData.X, testData.Y, nil, false) failed: a != 4.427`)
	}
	if b <= -0.47 || b >= -0.46 {
		t.Error(`stat.LinearRegression(testData.X, testData.Y, nil, false) failed: b != -0.463`)
	}
}

func TestResidSumOfSquares(t *testing.T) {
	rss := ResidSumOfSquares(testData.X, testData.Y, 4.4274524, -0.46294153)
	if rss <= 699.49 || rss >= 699.50 {
		t.Error(`ResidSumOfSquares failed: rss !=699.498`)
	}
}
