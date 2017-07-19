package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fritzing/fzp/src/go"
	"github.com/urfave/cli"
)

var commandZipFlags = []cli.Flag{}

func commandZipAction(c *cli.Context) error {
	tmpArgs := c.Args()

	if len(tmpArgs) == 0 {
		log.Fatalf("Missing part filepath")
		os.Exit(1)
	}

	// Load fzp
	fzpPath := tmpArgs[0]
	fzpPath, err := filepath.Abs(fzpPath)
	if err != nil {
		log.Fatalf("Failed to make absolute path for %s: %#v", tmpArgs[0], err)
	}
	fzpDir := filepath.Dir(fzpPath)
	part, _, err := fzp.ReadFzp(fzpPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(127)
	}

	buf := &bytes.Buffer{}
	zipFile := zip.NewWriter(buf)

	// Load files to include in zip
	files := make(map[string]struct{})
	mustAddFile := func(src string, prefix string, content []byte) string {
		if src == "" {
			return ""
		}
		zipName := strings.Replace(src, "/", ".", -1)
		zipName = strings.Replace(zipName, "\\", ".", -1)
		if prefix != "" && !strings.HasPrefix(zipName, prefix) {
			zipName = prefix + zipName
		}
		if _, found := files[zipName]; found {
			return zipName
		}
		if content == nil {
			contentPrefix := strings.Replace(prefix, ".", "", -1)
			fullPath := filepath.Join(fzpDir, contentPrefix, src)
			content, err = ioutil.ReadFile(fullPath)
			if err != nil {
				log.Fatalf("Failed to read %s: %#v", fullPath, err)
			}
		}
		hdr := &zip.FileHeader{
			Name:   zipName,
			Method: zip.Deflate,
		}
		hdr.SetModTime(time.Now())
		f, err := zipFile.CreateHeader(hdr)
		if err != nil {
			log.Fatalf("Failed to create zip entry %s: %#v", src, err)
		}
		if _, err := f.Write(content); err != nil {
			log.Fatalf("Failed to add %s to zip: %#v", src, err)
		}
		zipFile.Flush()
		files[zipName] = struct{}{}
		return zipName
	}
	mustAddFile(part.Views.Breadboard.Image, "svg.", nil)
	mustAddFile(part.Views.Icon.Image, "svg.", nil)
	mustAddFile(part.Views.Pcb.Image, "svg.", nil)
	mustAddFile(part.Views.Schematic.Image, "svg.", nil)
	partContent, err := part.Marshal(fzp.FormatFzp)
	if err != nil {
		log.Fatalf("Failed to marshal part: %#v", err)
	}
	mustAddFile(filepath.Base(fzpPath), "part.", partContent)
	zipFile.Close()

	if errs := part.Check(); len(errs) > 0 {
		log.Println("Part contains errors:")
		for _, err := range errs {
			log.Println(err.Error())
		}
	}

	ext := filepath.Ext(fzpPath)
	outFile := filepath.Base(fzpPath[:len(fzpPath)-len(ext)]) + ".fzpz"
	if err := ioutil.WriteFile(outFile, buf.Bytes(), 0644); err != nil {
		log.Fatalf("Failed to write zip file %s: %#v", outFile, err)
	}

	return nil
}
