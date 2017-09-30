// A tiny go profiler based on wall time.
package tinyprof

import (
	"github.com/olekukonko/tablewriter"

	"io"
	"math"
	"os"
	"strconv"
	"time"
)

type ProfilerRegistry struct {
	underlying []*Profiler
}

type Profiler struct {
	buf []checkPoint
}

var defaultRegistry *ProfilerRegistry

func init() {
	defaultRegistry = &ProfilerRegistry{
		underlying: []*Profiler{},
	}
}

type checkPoint struct {
	timestamp time.Time
	id        string
}

type result struct {
	rows map[string]*resultRow
}

type resultRow struct {
	count int
	sum   time.Duration
	max   time.Duration
	min   time.Duration
	avg   time.Duration
}

func NewProfilerRegistry() *ProfilerRegistry {
	return &ProfilerRegistry{
		underlying: []*Profiler{},
	}
}

func NewProfiler(registry *ProfilerRegistry) *Profiler {
	p := Profiler{
		buf: []checkPoint{},
	}

	var r *ProfilerRegistry
	if registry == nil {
		r = defaultRegistry
	} else {
		r = registry
	}
	r.underlying = append(r.underlying, &p)
	p.buf = append(p.buf, checkPoint{timestamp: time.Now()})
	return &p
}

func (p *Profiler) Step(id string) {
	p.buf = append(p.buf, checkPoint{
		timestamp: time.Now(),
		id:        id,
	})
}

func Write(w io.Writer, registry *ProfilerRegistry) {
	var r *ProfilerRegistry
	if registry == nil {
		r = defaultRegistry
	} else {
		r = registry
	}

	result := r.aggregate()
	writer := tablewriter.NewWriter(w)
	writer.SetHeader([]string{
		"Name", "count", "sum(ms)", "max(ms)", "min(ms)", "avg(ms)",
	})
	writer.SetAutoFormatHeaders(false)
	writer.SetAlignment(tablewriter.ALIGN_RIGHT)

	for id, row := range result.rows {
		writer.Append([]string{
			id,
			strconv.Itoa(row.count),
			formatDuration(row.sum),
			formatDuration(row.max),
			formatDuration(row.min),
			formatDuration(row.avg),
		})
	}

	writer.Render()
}

func Print(registry *ProfilerRegistry) {
	Write(os.Stdout, registry)
}

func (r *ProfilerRegistry) aggregate() result {
	rows := make(map[string]*resultRow)

	for _, p := range r.underlying {

		for i, c := range p.buf {
			if i > 0 {
				p := p.buf[i-1].timestamp
				d := c.timestamp.Sub(p)

				var row *resultRow
				if r := rows[c.id]; r != nil {
					row = r
				} else {
					row = &resultRow{
						min: math.MaxInt64,
					}
					rows[c.id] = row
				}

				row.count += 1
				row.sum += d
				if d <= row.min {
					row.min = d
				}
				if d >= row.max {
					row.max = d
				}
				row.avg = row.sum / time.Duration(row.count)
			}
		}
	}

	return result{
		rows: rows,
	}
}

func formatDuration(d time.Duration) string {
	return strconv.FormatFloat(float64(d)/(1000*1000), 'f', -1, 64)
}
