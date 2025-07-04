package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func openAndClose(id int, path string) error {
	fmt.Println(path)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var buf [1]byte
	f.Read(buf[:])

	return nil
}

func newPool(concurrency int) (chan string, func()) {
	jobs := make(chan string, concurrency)
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for path := range jobs {
				if err := openAndClose(i, path); err != nil {
					fmt.Fprintf(os.Stderr, "[worker %d] error: %v\n", id, err)
				}
			}
		}(i)
	}
	return jobs, func() {
		close(jobs)
		wg.Wait()
	}
}

func tryOneDir(root string, jobs chan string, count *int) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "walk error: %v\n", err)
			return nil
		}
		if !d.IsDir() {
			jobs <- path
			(*count)++
		}
		return nil
	})
}

var concurrency = flag.Uint("c", 8, "Number of files to open concurrently (recommended: 4–16)")

func mains(args []string) error {
	if len(args) <= 0 {
		flag.Usage()
		return nil
	}

	c := int(*concurrency)
	if c <= 0 || c > 16 {
		return errors.New("invalid value for -c: must be between 1 and 16")
	}

	count := 0
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Println("\nDone.")
		fmt.Printf("Elapsed time: %s\n", elapsed)
		fmt.Printf("Found files: %d\n", count)
	}()

	jobs, closer := newPool(c)
	defer closer()

	var errs []error
	for _, arg1 := range args {
		err := tryOneDir(arg1, jobs, &count)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func main() {
	flag.Parse()
	if err := mains(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
