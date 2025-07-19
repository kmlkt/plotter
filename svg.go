package pltt

import (
	"fmt"
	"io"
	"iter"
	"strings"
)

func buildGraph(w io.Writer, iters []iter.Seq[record], cfg graphConfig) {
	if len(iters) == 0 {
		return
	}
	data := multiCollect(iters)
	desc := newRecordsDescriptor(data)
	desc.Times.Max = cfg.maxT
	if cfg.minT.Unix() != 0 {
		desc.Times.Min = cfg.minT
	}

	fmt.Fprint(w, `
			<svg version='1.1' preserveAspectRatio='none'
			style='overflow: visible' xmlns='http://www.w3.org/2000/svg'>
			`)

	for i, dataSource := range data {
		if len(dataSource) == 0 {
			continue
		}
		points := make([]point, len(dataSource))
		for i, r := range dataSource {
			points[i].X = desc.Times.position(r.Timestamp)
			points[i].Y = 1 - desc.Values.position(r.Value)
			points[i].Value = fmt.Sprintf("%v", r.Value)
		}
		color := dataSourceColor(i)
		drawLine(points, w, color)
		drawDots(points, w, color)
	}
	drawAxis(desc, cfg, w)
	drawLegend(w, cfg.keys)
	fmt.Fprint(w, "</svg>")
}

func drawLine(points []point, w io.Writer, color string) {
	fmt.Fprintf(w, `
		<svg viewBox='0 0 1 1' x='10%%' y='0' width='90%%' height='90%%'
		preserveAspectRatio='none' style='overflow: visible' xmlns='http://www.w3.org/2000/svg'>
		<path fill='transparent' stroke='%s' stroke-width='3' style='vector-effect: non-scaling-stroke' d='M
		`, color)
	points[0].DrawOnLine(w)
	for _, p := range points[1:] {
		fmt.Fprint(w, "L ")
		p.DrawOnLine(w)
	}
	fmt.Fprint(w, "'/></svg>")
}

func drawDots(points []point, w io.Writer, color string) {
	for _, p := range points {
		p.DrawSeparately(w, color)
	}
}

func drawAxis(data recordsDescriptor, cfg graphConfig, builder io.Writer) {
	times, values := labels(data, cfg)
	drawXAxis(times, builder)
	drawYAxis(values, builder)
}

func drawXAxis(times []label, w io.Writer) {
	for _, time := range times {
		parts := strings.Split(time.Value, " ")
		if len(parts) == 2 {
			fmt.Fprintf(w,
				"<text text-anchor='middle' dominant-baseline='middle' x='%f%%' y='%f%%'>%s</text>",
				10+time.Position*90, 93.0, parts[0])
			fmt.Fprintf(w,
				"<text text-anchor='middle' dominant-baseline='middle' x='%f%%' y='%f%%'>%s</text>",
				10+time.Position*90, 98.0, parts[1])
		} else {
			fmt.Fprintf(w,
				"<text text-anchor='middle' dominant-baseline='middle' x='%f%%' y='%f%%'>%s</text>",
				10+time.Position*90, 95.0, time.Value)
		}
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
	X     float64
	Y     float64
	Value string
}

func (p point) DrawOnLine(w io.Writer) {
	fmt.Fprintf(w, "%f %f", p.X, p.Y)
}

func (p point) DrawSeparately(w io.Writer, color string) {
	fmt.Fprintf(w, "<circle cx='%f%%' cy='%f%%' r='4' fill='%s'>%s</circle>", 10+p.X*90, p.Y*90, color, p.Value)
}

func drawLegend(w io.Writer, titles []string) {
	if len(titles) <= 1 {
		return
	}
	for i, title := range titles {
		fmt.Fprintf(w, `<text x='100%%' y='%d' text-anchor='end' dominant-baseline='hanging'
			fill='%s'>%s</text>`, i*16, dataSourceColor(i), title)
	}
}

var colors = []string{"red", "blue", "green", "yellow", "maroon", "aqua", "purple", "olive", "wheat"}

func dataSourceColor(i int) string {
	if i < len(colors) {
		return colors[i]
	}
	return "black"
}
