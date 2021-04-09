//go:generate go run main.go

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime err")
	}

	rootDir := path.Join(path.Dir(filename), "..")

	dstDir := filepath.Join(rootDir, "internal/libwebp")
	srcDir := filepath.Join(rootDir, "libwebp_src", "src")

	// The Go and the Webp source must live side-by-side in the same
	// directory.
	//
	// The custom bindings are named with a "a__" prefix. Keep those.
	fis, err := ioutil.ReadDir(dstDir)
	if err != nil {
		log.Fatal(err)
	}

	keepRe := regexp.MustCompile(`^(a__|\.)`)

	for _, fi := range fis {
		if keepRe.MatchString(fi.Name()) {
			continue
		}
		os.Remove(filepath.Join(dstDir, fi.Name()))
	}

	csourceRe := regexp.MustCompile(`\.[ch]$`)

	err = filepath.Walk(srcDir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() || !csourceRe.MatchString(fi.Name()) {
			return nil
		}

		filename := filepath.ToSlash(strings.TrimPrefix(path, srcDir))
		filename = strings.TrimPrefix(filename, "/")
		target := filepath.Join(dstDir, fi.Name())

		if err := ioutil.WriteFile(target, []byte(fmt.Sprintf(`#ifndef LIBWEBP_NO_SRC
#include "../../libwebp_src/src/%s"
#endif
`, filename)), 0o644); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
