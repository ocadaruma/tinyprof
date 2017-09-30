# tinyprof
A tiny go profiler based on wall time.

## Usage

```bash
$ go get github.com/ocadaruma/tinyprof
```

app.go:

```go
package main

import (
	"github.com/ocadaruma/tinyprof"
	
	"time"
)

func main() {
	profiler := tinyprof.NewProfiler(nil)
	
	time.Sleep(100 * time.Millisecond)
	profiler.Step("point A")
	
	time.Sleep(500 * time.Millisecond)
	profiler.Step("point B")
	
	time.Sleep(1000 * time.Millisecond)
	profiler.Step("point C")
	
	tinyprof.Print(nil)
}
```

```bash
$ go run app.go
+---------+-------+-------------+-------------+-------------+-------------+
|  Name   | count |  total(ms)  |   max(ms)   |   min(ms)   |   avg(ms)   |
+---------+-------+-------------+-------------+-------------+-------------+
| point A |     1 |  102.870693 |  102.870693 |  102.870693 |  102.870693 |
| point B |     1 |  505.045537 |  505.045537 |  505.045537 |  505.045537 |
| point C |     1 | 1004.036785 | 1004.036785 | 1004.036785 | 1004.036785 |
+---------+-------+-------------+-------------+-------------+-------------+
```
