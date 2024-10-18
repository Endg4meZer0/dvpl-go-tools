package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Endg4meZer0/dvpl-go"
)

var (
	compressMode   = flag.Bool("c", false, "Sets the mode to 'compression'.")
	decompressMode = flag.Bool("d", false, "Sets the mode to 'decompression'.")
	recursive      = flag.Bool("r", false, "Recursively convert all files and the contents of all folders inside the set path.")
	force          = flag.Bool("f", false, "Force the compression algorithm to always use compression instead of detecting .tex files and applying no compression on them.")
	deleteOld      = flag.Bool("n", false, "Delete the old file after converting.")
)

var (
	total     = 0
	completed = 0
	failed    = 0
)

func main() {
	flag.CommandLine.Usage = printUsage
	flag.Parse()

	if !*compressMode && !*decompressMode {
		fmt.Fprintln(os.Stderr, "No mode set. Use --help for more information.")
		os.Exit(1)
	}

	paths := make([]string, 0, len(os.Args[2:]))
	for _, v := range os.Args[2:] {
		if !strings.HasPrefix(v, "-") {
			path, err := filepath.Abs(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, "An unknown error occured when trying to read the set path (%s). Skipping...", v)
				continue
			}
			paths = append(paths, path)
		}
	}

	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "No paths specified. Use --help for more information.")
		os.Exit(2)
	}

	convertPaths(paths)

	fmt.Fprintf(os.Stdout, "\n\nFinished converting files.\nTotal: %d\nFailed: %d\nCompleted: %d", total, failed, completed)

	os.Exit(0)
}

func convertPaths(paths []string) {
	for _, path := range paths {
		fileInfo, err := os.Stat(path)

		if err != nil {
			fmt.Fprintf(os.Stderr, "An unknown error occured when trying to read %s. An issue with permissions?", path)
			total++
			failed++
			continue
		}

		if fileInfo.IsDir() && !*recursive {
			failed++
			fmt.Fprintf(os.Stderr, "Ignoring %s since it's a directory and the recursive flag (-r) is NOT set...", path)
			continue
		} else if fileInfo.IsDir() {
			dir, err := os.ReadDir(path)
			if err != nil {
				total++
				failed++
				fmt.Fprintf(os.Stderr, "An unknown error occured when trying to read %s as a directory. An issue with permissions?", path)
				continue
			}
			dirPaths := make([]string, 0, len(dir))
			for _, dirEntry := range dir {
				dirPaths = append(dirPaths, path+"/"+dirEntry.Name())
			}
			convertPaths(dirPaths)
		} else {
			total++
			file, err := os.ReadFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "An unknown error occured when trying to read %s. An issue with permissions?", path)
				failed++
				continue
			}

			var convertedFile []byte
			var newPath string

			if *decompressMode {
				convertedFile, err = dvpl.DecompressDVPL(file)
				var hasSuffix bool
				newPath, hasSuffix = strings.CutSuffix(path, ".dvpl")
				if !hasSuffix {
					newPath = path + ".nodvpl"
				}
			} else {
				convertedFile, err = dvpl.CompressDVPL(file, strings.HasSuffix(path, ".tex") && !*force)
				newPath = path + ".dvpl"
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "An error occured during the conversion of the file %s:\n%s", path, err.Error())
				failed++
				continue
			}
			os.WriteFile(newPath, convertedFile, 0777)
			completed++

			if *deleteOld {
				err := os.Remove(path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "The deletion flag (-n) is set, but I could not delete %s after the conversion!", path)
				}
			}
		}
	}
}

const usageText = "Usage: dvpl-go-tools <-c|-d> [-f] [-n] [-r] PATH [PATH...]\n\nOptions:\n"

func printUsage() {
	fmt.Fprint(os.Stderr, usageText)
	flag.PrintDefaults()
}
