package pltt

import (
	"fmt"
	"io"
)

func buildGraph(data recordsDescriptor, w io.Writer) {
	if len(data.Data) == 0 {
		return
	}
	szV := data.Values.Max - data.Values.Min
	szT := data.Times.Max - data.Times.Min
	points := make([]point, len(data.Data))
	for i, r := range data.Data {
		points[i].X = (r.TimeFloat64() - data.Times.Min) / szT
		points[i].Y = 1 - (r.Value-data.Values.Min)/szV
	}

	fmt.Fprint(w, `
		<svg version='1.1' preserveAspectRatio='none'
		style='overflow: visible' xmlns='http://www.w3.org/2000/svg'>
		`)
	drawLine(points, w)
	drawDots(points, w)
	drawAxis(data, w)
	fmt.Fprint(w, "</svg>")
}

func drawLine(points []point, w io.Writer) {
	fmt.Fprint(w, `
		<svg viewBox='0 0 1 1' x='10%' y='0' width='90%' height='90%'
		preserveAspectRatio='none' style='overflow: visible' xmlns='http://www.w3.org/2000/svg'>
		<path fill='transparent' stroke='black' stroke-width='3' style='vector-effect: non-scaling-stroke' d='M
		`)
	points[0].DrawOnLine(w)
	for _, p := range points[1:] {
		fmt.Fprint(w, "L ")
		p.DrawOnLine(w)
	}
	fmt.Fprint(w, "'/></svg>")
}

func drawDots(points []point, w io.Writer) {
	for _, p := range points {
		p.DrawSeparately(w)
	}
}

func drawAxis(data recordsDescriptor, builder io.Writer) {
	times, values := labels(data)
	drawXAxis(times, builder)
	drawYAxis(values, builder)
}

func drawXAxis(times []label, w io.Writer) {
	for _, time := range times {
		fmt.Fprintf(w,
			"<text text-anchor='middle' dominant-baseline='middle' x='%f%%' y='%f%%'>%s</text>",
			10+time.Position*90, 95.0, time.Value)
	}
}

func drawYAxis(values []label, w io.Writer) {
	for _, value := range values {
		fmt.Fprintf(w,
			"<text text-anchor='middle' dominant-baseline='middle' x='%f%%' y='%f%%'>%s</text>",
			5.0, 90-value.Position*90, value.Value)
	}
}

// coordiantes from 0 to 1
type point struct {
	X float64
	Y float64
}

func (p point) DrawOnLine(w io.Writer) {
	fmt.Fprintf(w, "%f %f", p.X, p.Y)
}

func (p point) DrawSeparately(w io.Writer) {
	fmt.Fprintf(w, "<circle cx='%f%%' cy='%f%%' r='4' fill='black'/>", 10+p.X*90, p.Y*90)
}
