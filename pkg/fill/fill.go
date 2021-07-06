package fill

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Position struct {
	X int
	Y int
}

type Border struct {
	filePosition Position
	left         []int
	top          []int
	right        []int
	bottom       []int
}

func EqualRGB(a, b color.Color) bool {
	ar, ag, ab, _ := a.RGBA()
	br, bg, bb, _ := b.RGBA()
	return ar == br && ag == bg && ab == bb
}

func neighbors(p Position) []Position {
	var out []Position
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			np := Position{X: p.X + i, Y: p.Y + j}
			if p == np {
				continue
			}
			out = append(out, np)
		}
	}
	return out
}

func TileFill(rgba *image.RGBA, tilePositions []Position, fromColor, toColor color.Color) Border {
	bounds := rgba.Bounds()
	visited := map[Position]bool{}
	border := Border{}
	for len(tilePositions) > 0 {
		tp := tilePositions[0]
		tilePositions = tilePositions[1:]
		if visited[tp] {
			continue
		}
		visited[tp] = true
		if !EqualRGB(rgba.At(tp.X, tp.Y), fromColor) {
			continue
		}
		rgba.Set(tp.X, tp.Y, toColor)
		if tp.X == bounds.Min.X {
			border.left = append(border.left, tp.Y)
		}
		if tp.Y == bounds.Min.Y {
			border.top = append(border.top, tp.X)
		}
		if tp.X == (bounds.Max.X - 1) {
			border.right = append(border.right, tp.Y)
		}
		if tp.Y == (bounds.Max.Y - 1) {
			border.bottom = append(border.bottom, tp.X)
		}
		for _, v := range neighbors(tp) {
			if v.X < bounds.Min.X || v.Y < bounds.Min.Y || v.X >= bounds.Max.X || v.Y >= bounds.Max.Y {
				continue
			}
			tilePositions = append(tilePositions, v)
		}
	}
	return border
}

type Work struct {
	path         string
	filePosition Position
	positions    map[Position]struct{}
	edges        Border
	fromColor    color.Color
	toColor      color.Color
}

func imagePath(path string, p Position) string {
	return filepath.Join(path, fmt.Sprintf("tile_%d_%d.png", p.X, p.Y))
}

func WorkManager(workCh chan<- Work, resultsCh <-chan Border, initial Work, path string) {
	doWorkCh := workCh
	fromColor := initial.fromColor
	toColor := initial.toColor
	workSet := map[Position]Work{initial.filePosition: initial}
	todo := []Position{initial.filePosition}
	processing := map[Position]bool{}
	loopCount := 0
	for len(todo) > 0 || len(processing) > 0 {
		loopCount++
		var w Work
		i := 0
		for i = 0; i < len(todo); i++ {
			if !processing[todo[i]] {
				w = workSet[todo[i]]
				break
			}
		}
		if i < len(todo) {
			temp := todo[:i]
			todo = todo[i:]
			todo = append(todo, temp...)
			doWorkCh = workCh
		} else {
			doWorkCh = nil
		}
		select {
		case doWorkCh <- w:
			processing[w.filePosition] = true
			todo = todo[1:]
			delete(workSet, w.filePosition)
		case b := <-resultsCh:
			if len(b.left) > 0 {
				lp := Position{X: b.filePosition.X - 1, Y: b.filePosition.Y}
				if ww, ok := workSet[lp]; ok {
					ww.edges.right = append(ww.edges.right, b.left...)
					workSet[lp] = ww
				} else {
					workSet[lp] = Work{
						path:         imagePath(path, lp),
						filePosition: lp,
						edges:        Border{right: b.left},
						fromColor:    fromColor,
						toColor:      toColor,
					}
					todo = append(todo, lp)
				}
			}
			if len(b.top) > 0 {
				tp := Position{X: b.filePosition.X, Y: b.filePosition.Y - 1}
				if ww, ok := workSet[tp]; ok {
					ww.edges.bottom = append(ww.edges.bottom, b.top...)
					workSet[tp] = ww
				} else {
					workSet[tp] = Work{
						path:         imagePath(path, tp),
						filePosition: tp,
						edges:        Border{bottom: b.top},
						fromColor:    fromColor,
						toColor:      toColor,
					}
					todo = append(todo, tp)
				}
			}
			if len(b.right) > 0 {
				rp := Position{X: b.filePosition.X + 1, Y: b.filePosition.Y}
				if ww, ok := workSet[rp]; ok {
					ww.edges.left = append(ww.edges.left, b.right...)
					workSet[rp] = ww
				} else {
					workSet[rp] = Work{
						path:         imagePath(path, rp),
						filePosition: rp,
						edges:        Border{left: b.right},
						fromColor:    fromColor,
						toColor:      toColor,
					}
					todo = append(todo, rp)
				}
			}
			if len(b.bottom) > 0 {
				bp := Position{X: b.filePosition.X, Y: b.filePosition.Y + 1}
				if ww, ok := workSet[bp]; ok {
					ww.edges.top = append(ww.edges.top, b.bottom...)
					workSet[bp] = ww
				} else {
					workSet[bp] = Work{
						path:         imagePath(path, bp),
						filePosition: bp,
						edges:        Border{top: b.bottom},
						fromColor:    fromColor,
						toColor:      toColor,
					}
				}
				todo = append(todo, bp)
			}
			delete(processing, b.filePosition)
		}
		if (loopCount % 10) == 0 {
			log.Printf("loop %d, todo items %d, processing items %d", loopCount, len(todo), len(processing))
		}
	}
	log.Printf("loop %d", loopCount)
}

func Worker(workCh <-chan Work, resultsCh chan<- Border) {
	for work := range workCh {
		f, err := os.OpenFile(work.path, os.O_RDWR, 0755)
		if err != nil {
			resultsCh <- Border{filePosition: work.filePosition}
			continue
		}
		i, _, err := image.Decode(f)
		if err != nil {
			f.Close()
			log.Printf("error decoding %s: %v", work.path, err)
			resultsCh <- Border{filePosition: work.filePosition}
			continue
		}
		rgba := image.NewRGBA(image.Rect(0, 0, i.Bounds().Dx(), i.Bounds().Dy()))
		draw.Draw(rgba, rgba.Bounds(), i, i.Bounds().Min, draw.Src)
		err = f.Close()
		if err != nil {
			resultsCh <- Border{filePosition: work.filePosition}
			continue
		}
		tilePositions := work.positions
		if len(tilePositions) == 0 {
			tilePositions = make(map[Position]struct{})
			// use border
			b := rgba.Bounds()
			for _, y := range work.edges.left {
				tilePositions[Position{X: b.Min.X, Y: y}] = struct{}{}
			}
			for _, x := range work.edges.top {
				tilePositions[Position{X: x, Y: b.Min.Y}] = struct{}{}
			}
			for _, y := range work.edges.right {
				tilePositions[Position{X: b.Max.X - 1, Y: y}] = struct{}{}
			}
			for _, x := range work.edges.bottom {
				tilePositions[Position{X: x, Y: b.Max.Y - 1}] = struct{}{}
			}
		}
		positions := make([]Position, 0, len(tilePositions))
		for p := range tilePositions {
			positions = append(positions, p)
		}
		b := TileFill(rgba, positions, work.fromColor, work.toColor)
		b.filePosition = work.filePosition
		f, err = os.Create(work.path)
		if err != nil {
			log.Printf("error creating file %s: %v", work.path, err)
			resultsCh <- Border{filePosition: work.filePosition}
			continue
		}
		err = png.Encode(f, rgba)
		if err != nil {
			log.Printf("error encoding %s: %v", work.path, err)
			resultsCh <- Border{filePosition: work.filePosition}
			continue
		}
		err = f.Close()
		if err != nil {
			log.Printf("error closing %s: %v", work.path, err)
			resultsCh <- Border{filePosition: work.filePosition}
			continue
		}
		resultsCh <- b
	}
}

func WorldFill(path string, filePosition, tilePosition Position, toColor color.Color, numWorkers int) error {
	file := imagePath(path, filePosition)
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		return err
	}
	fromColor := i.At(tilePosition.X, tilePosition.Y)
	if EqualRGB(fromColor, toColor) {
		return nil
	}
	todo := Work{
		path:         file,
		filePosition: filePosition,
		positions: map[Position]struct{}{
			tilePosition: {},
		},
		fromColor: fromColor,
		toColor:   toColor,
	}

	resultsCh := make(chan Border)
	workCh := make(chan Work)

	wg := sync.WaitGroup{}
	wg.Add(1)
	log.Print("starting work manager")
	go func() {
		WorkManager(workCh, resultsCh, todo, path)
		close(workCh)
		wg.Done()
	}()
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		log.Printf("starting worker %d", i)
		go func() {
			Worker(workCh, resultsCh)
			wg.Done()
		}()
	}
	wg.Wait()
	close(resultsCh)
	return nil
}
