package pltt

import (
	"math"
	"slices"
	"strconv"
	"strings"
)

func buildGraph(data []record) string {
	minV := math.MaxFloat64
	maxV := -math.MaxFloat64
	minT := int64(math.MaxInt64)
	maxT := int64(math.MinInt64)
	for _, r := range data {
		minV = min(minV, r.Value)
		maxV = max(maxV, r.Value)
		minT = min(minT, r.Timestamp.UnixNano())
		maxT = max(maxT, r.Timestamp.UnixNano())
	}
	szV := maxV - minV
	szT := float64(maxT - minT)
	points := make([]point, len(data))
	for i, r := range data {
		points[i].X = float64(r.Timestamp.UnixNano()-minT) / szT
		points[i].Y = (r.Value - minV) / szV
	}

	builder := strings.Builder{}
	builder.WriteString("<svg preserveAspectRatio='none' xmlns='http://www.w3.org/2000/svg'>")
	buildLine(points, &builder)
	builder.WriteString("</svg>")
	return builder.String()
}

func buildLine(points []point, builder *strings.Builder) {
	slices.SortFunc(points, func(a point, b point) int {
		if a.X < b.X {
			return -1
		}
		return 1
	})
	if len(points) == 0 {
		return
	}
	builder.WriteString("<svg viewBox='0 0 1 1' x='10%' y='0' width='90%' height='90%' preserveAspectRatio='none' xmlns='http://www.w3.org/2000/svg'>")
	builder.WriteString("<path fill='transparent' stroke='black' style='vector-effect: non-scaling-stroke' d='M ")
	points[0].StringB(builder)
	for _, p := range points[1:] {
		builder.WriteString("L ")
		p.StringB(builder)
	}
	builder.WriteString("'/>")
	builder.WriteString("</svg>")
}

// coordiantes from 0 to 1
type point struct {
	X float64
	Y float64
}

func (p point) StringB(builder *strings.Builder) {
	builder.WriteString(formatFloat(p.X))
	builder.WriteRune(' ')
	builder.WriteString(formatFloat(p.Y))
	builder.WriteRune(' ')
}

func formatFloat(x float64) string {
	return strconv.FormatFloat(x, 'f', -1, 64)
}
