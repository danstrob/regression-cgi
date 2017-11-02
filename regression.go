package main

import (
	"image/color"
	"log"
	"math"
	"path/filepath"
	"strconv"
	"time"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Regression takes in a Data struct, runs a bivariate OLS regression and
// returns the residual sum of squares for the OLS model as well as for the predicted values
// based on the user's guesses for the intercept and slope.
func Regression(d *Data) (float64, float64) {
	a, b := stat.LinearRegression(d.X, d.Y, nil, false)
	rssOLS := ResidSumOfSquares(d.X, d.Y, a, b)
	rssGuess := ResidSumOfSquares(d.X, d.Y, d.Intercept, d.Slope)
	return rssOLS, rssGuess
}

// DrawPlot takes in a Data struct and produces an image file
// of a bivariate scatter plot with a line graph based on the user's guesses
// for the intercept and slope. It returns the file path of the image.
func DrawPlot(d *Data) string {
	scatterData := plotterData(d.X, d.Y)
	lineData := lineData(d.Intercept, d.Slope, d.X)

	// Create a new plot, set its title, axis labels and grid.
	p, err := plot.New()
	if err != nil {
		log.Println(err)
	}
	font, err := vg.MakeFont("Helvetica", 18)
	if err != nil {
		log.Println(err)
	}
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"
	p.X.Label.TextStyle.Font = font
	p.Y.Label.TextStyle.Font = font
	p.X.Tick.Label.Font = font
	p.Y.Tick.Label.Font = font
	p.Add(plotter.NewGrid())

	// Make a scatter and line plot and set their style.
	s, err := plotter.NewScatter(scatterData)
	if err != nil {
		log.Println(err)
	}
	s.GlyphStyle.Shape = draw.CircleGlyph{}
	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
	s.GlyphStyle.Radius = vg.Length(4)

	l, err := plotter.NewLine(lineData)
	if err != nil {
		log.Println(err)
	}
	l.LineStyle.Width = vg.Points(2)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	// Add the plotters to the plot and save to a PNG file.
	p.Add(s, l)
	imagePath := makeFilePath()

	if err := p.Save(6*vg.Inch, 6*vg.Inch, imagePath); err != nil {
		log.Println(err)
	}
	return imagePath
}

// lineData returns the plotter data for a line plot based on user guesses
// for the intercept and slope, using the range of the x values of the scatter plot data.
func lineData(intercept, slope float64, xdata []float64) plotter.XYs {
	min := floats.Min(xdata)
	max := floats.Max(xdata)
	x := []float64{min}
	y := []float64{intercept + slope*min}
	for xval := min; xval <= max; xval++ {
		x = append(x, xval)
		y = append(y, intercept+slope*xval)
	}
	return plotterData(x, y)
}

func makeFilePath() string {
	return filepath.Join("images", strconv.FormatInt(time.Now().UnixNano(), 10)+".png")
}

func plotterData(x, y []float64) plotter.XYs {
	pts := make(plotter.XYs, len(x))
	for i, xvals := range x {
		pts[i].X = xvals
		pts[i].Y = y[i]
	}
	return pts
}

func ResidSumOfSquares(x, y []float64, alpha, beta float64) (rss float64) {
	if len(x) != len(y) {
		log.Fatalln("Length mismatch of slices X and Y.")
	}
	for i, xval := range x {
		resid := (alpha + beta*xval) - y[i]
		rss += resid * resid
	}
	return round(rss)
}

func round(x float64) float64 {
	return math.Floor((x*100)+0.5) / 100
}
