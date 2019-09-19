package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"
)

const (
	ready = iota
	firing
	refactory

	defaultRows     = 30
	defaultCols     = 60
	defaultInterval = 100
)

var (
	sb strings.Builder // string builder for drawing

	clearScreen   = "\u001b[2J"
	cellReady     = []byte("\u001b[48;5;252m  \u001b[0m") // white
	cellFiring    = []byte("\u001b[48;5;28m  \u001b[0m")  // green
	cellRefactory = []byte("\u001b[48;5;220m  \u001b[0m") // yellow
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type game struct {
	rows     int
	cols     int
	interval time.Duration

	prev    [][]uint8
	current [][]uint8

	frozen bool // if frozen, game will stop
}

func newGame(rows, cols, interval int) *game {
	prev := make([][]uint8, rows)
	current := make([][]uint8, rows)

	for x := 0; x < rows; x++ {
		prev[x] = make([]uint8, cols)
		current[x] = make([]uint8, cols)
		for y := 0; y < cols; y++ {
			current[x][y] = uint8(rand.Intn(3))
		}
	}

	return &game{
		rows:     rows,
		cols:     cols,
		interval: time.Millisecond * time.Duration(interval),
		prev:     prev,
		current:  current,
		frozen:   false,
	}
}

func (g *game) run() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fmt.Print(clearScreen) // clean the whole screen
	for {
		select {
		case <-c:
			log.Println("Bye.")
			return
		default:
			g.draw()
			g.next()
			if g.frozen {
				log.Println("lives are frozen, game stopped!")
				return
			}
			time.Sleep(g.interval)
		}
	}
}

func (g *game) draw() {
	sb.Reset()
	sb.WriteString(fmt.Sprintf("\u001b[%dD", g.rows*(g.cols+1)))
	sb.WriteString(fmt.Sprintf("\u001b[%dA", g.rows))
	for _, col := range g.current {
		for _, cell := range col {
			switch cell {
			case ready:
				sb.Write(cellReady)
			case firing:
				sb.Write(cellFiring)
			case refactory:
				sb.Write(cellRefactory)
			default:
				log.Fatalf("invalid cell value: %d", cell)
			}
		}
		sb.WriteByte('\n')
	}
	fmt.Print(sb.String())
}

func (g *game) next() {
	g.prev, g.current = g.current, g.prev // switch prev and current

	hasChange := false

	for x := 0; x < g.rows; x++ {
		for y := 0; y < g.cols; y++ {
			switch g.prev[x][y] {
			case ready:
				// get alive count around cell
				aliveCount := 0
				if x > 0 && y > 0 && g.prev[x-1][y-1] == firing { // x-1, y-1
					aliveCount++
				}
				if x > 0 && g.prev[x-1][y] == firing { // x-1, y
					aliveCount++
				}
				if x > 0 && y < g.cols-1 && g.prev[x-1][y+1] == firing { // x-1, y+1
					aliveCount++
				}
				if y > 0 && g.prev[x][y-1] == firing { // x, y-1
					aliveCount++
				}
				if y < g.cols-1 && g.prev[x][y+1] == firing { // x, y+1
					aliveCount++
				}
				if x < g.rows-1 && y > 0 && g.prev[x+1][y-1] == firing { // x+1, y-1
					aliveCount++
				}
				if x < g.rows-1 && g.prev[x+1][y] == firing { // x+1, y
					aliveCount++
				}
				if x < g.rows-1 && y < g.cols-1 && g.prev[x+1][y+1] == firing { // x+1, y+1
					aliveCount++
				}

				if aliveCount == 2 {
					g.current[x][y] = firing
					if !hasChange {
						hasChange = true
					}
				} else {
					g.current[x][y] = ready
				}
			case firing:
				g.current[x][y] = refactory
				if !hasChange {
					hasChange = true
				}
			case refactory:
				g.current[x][y] = ready
				if !hasChange {
					hasChange = true
				}
			}
		}
	}
	if !hasChange { // no change in this round, game freeze
		g.frozen = true
	}
}

func main() {
	var rows, cols, interval int

	flag.IntVar(&rows, "r", defaultRows, "rows count")
	flag.IntVar(&cols, "c", defaultCols, "columns count")
	flag.IntVar(&interval, "i", defaultInterval, "sleep interval between iterations (ms)")
	flag.Parse()

	if rows <= 0 {
		log.Printf("Invalid rows %d, use default value", rows)
		rows = defaultRows
	}
	if cols <= 0 {
		log.Printf("Invalid columns %d, use default value", cols)
		cols = defaultCols
	}
	if interval <= 0 {
		log.Printf("Invalid interval %d, use default value", interval)
		interval = defaultInterval
	}

	g := newGame(rows, cols, interval)
	g.run()
}
