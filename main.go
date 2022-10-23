package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"strings"
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
	// nworkers := 1
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
			log.Printf("Error %v in %s\n", err, dir)
			return
		}

		for _, path := range paths {

			if path.IsDir() {
				childpath := strings.TrimSuffix(dir+"/"+path.Name(), "/")
				childDirs = append(childDirs, childpath)
			}
			if path.Type().IsRegular() {
				w.out <- path.Name()
			}
		}
		w.addDir(childDirs)

		w.wg.Done()
	}
}

func (w *Walker) addDir(dirs []string) {
	if len(dirs) > 0 {

		w.wg.Add(len(dirs))
		go func() {
			for _, dir := range w.dirs {
				w.in <- dir
			}
		}()
	}
}

func main() {
	rootPath := flag.String("path", "/home/amine/Documents/programming/golang/", "path to crawl")
	flag.Parse()

	s := time.Now()

	w := InitWalker([]string{*rootPath})

	// harvest results
	results := []string{}
	go func() {
		for p := range w.out {
			results = append(results, p)
		}
	}()
	w.Run()

	e := time.Since(s)

	log.Default().Printf("%v files in %fs", len(results), e.Seconds())

}
