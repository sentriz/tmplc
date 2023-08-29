package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(files []string) error {
	var leaves, partials []string
	for _, name := range files {
		switch filepath.Ext(name) {
		case ".tmpl":
			leaves = append(leaves, name)
		default:
			partials = append(partials, name)
		}
	}

	var pool = sync.Pool{
		New: func() any { return new(bytes.Buffer) },
	}

	var leafErrs []error
	for _, name := range leaves {
		if err := processLeaf(&pool, name, partials); err != nil {
			leafErrs = append(leafErrs, fmt.Errorf("file %q: %w", name, err))
			continue
		}
	}

	return errors.Join(leafErrs...)
}

func processLeaf(pool *sync.Pool, name string, partials []string) error {
	var toParse []string
	toParse = append(toParse, name)        // the leaf should be the first defined, main template
	toParse = append(toParse, partials...) // partial(s) if any come after, in a new separate tree

	t, err := template.ParseFiles(toParse...)
	if err != nil {
		return fmt.Errorf("parsing leaf: %v", err)
	}

	buff := pool.Get().(*bytes.Buffer)
	defer pool.Put(buff)
	buff.Reset()

	if err := t.Execute(buff, nil); err != nil {
		return fmt.Errorf("executing template: %v", err)
	}

	destName := strings.TrimSuffix(name, filepath.Ext(name))
	dest, err := os.Create(destName)
	if err != nil {
		return fmt.Errorf("creating destination file: %v", err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, buff); err != nil {
		return fmt.Errorf("copy to destination file: %v", err)
	}

	return nil
}
