package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

const (
	_ParallelismFactor int = 2
)

type Walker struct {
	nWorkers int
	dirs     []string
	in       chan string
	out      chan string
	wg       *sync.WaitGroup
	done     chan bool
}

func InitWalker(dirs []string) *Walker {
	// Init the Walker
	nworkers := runtime.NumCPU() * _ParallelismFactor

	log.Default().Printf("Walking %v\n", dirs)
	return &Walker{
		nWorkers: nworkers,
		dirs:     dirs,
		wg:       new(sync.WaitGroup),
		in:       make(chan string),
		out:      make(chan string),
		done:     make(chan bool),
	}
}

func (w *Walker) Run() {
	// start work pool
	for i := 0; i < w.nWorkers; i++ {
		go w.worker()
	}

	w.wg.Add(len(w.dirs))
	w.addDir(w.dirs)

	go func() {
		w.wg.Wait()
		close(w.in)
		close(w.out)
		w.done <- true
	}()

	<-w.done
}

func (w *Walker) worker() {
	for dir := range w.in {
		childDirs := []string{}

		paths, err := os.ReadDir(dir)
		if err != nil {
			w.wg.Done()
			continue
		}

		for _, path := range paths {
			if path.IsDir() {
				childpath := filepath.Join(dir, path.Name())
				childDirs = append(childDirs, childpath)
			}
			if path.Type().IsRegular() {
				w.out <- path.Name()
			}
		}

		if len(childDirs) > 0 {
			w.wg.Add(len(childDirs))
			go w.addDir(childDirs)
		}

		w.wg.Done()
	}
}

func (w *Walker) addDir(dirs []string) {
	for _, dir := range dirs {
		w.in <- dir
	}
}

func main() {
	rootPath := flag.String("path", "/home/amine", "path to crawl")
	flag.Parse()

	s := time.Now()

	w := InitWalker([]string{*rootPath})

	// harvest results
	rDone := make(chan bool)
	results := []string{}
	go func() {
		for p := range w.out {
			results = append(results, p)
		}
		rDone <- true
	}()
	w.Run()

	e := time.Since(s)
	<-rDone
	log.Default().Printf("%v files in %fs", len(results), e.Seconds())

}
