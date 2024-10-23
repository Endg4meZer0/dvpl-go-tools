package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Endg4meZer0/dvpl-go"
	"github.com/TwiN/go-color"
)

// Colored outputs
var (
	errorStr   = color.InBold(color.Ize(color.Red, "ERROR! "))
	warningStr = color.InBold(color.Ize(color.Yellow, "WARNING! "))
	successStr = color.InBold(color.Ize(color.Green, "SUCCESS! "))
)

// Flags
var (
	compressMode   = flag.Bool("c", false, "Sets the mode to 'compression'.")
	decompressMode = flag.Bool("d", false, "Sets the mode to 'decompression'.")
	recursive      = flag.Bool("r", false, "Recursively convert all files and the contents of all folders inside the set path.")
	force          = flag.Bool("f", false, "Force the compression algorithm to always use compression instead of detecting .tex files and applying no compression on them.")
	deleteOld      = flag.Bool("n", false, "Delete the old file after converting.")
	noColorOutput  = flag.Bool("p", false, "Use plain output (disable colored output).")
)

// Counters
var (
	total     = 0
	completed = 0
	failed    = 0
)

func main() {
	flag.CommandLine.Usage = printUsage
	flag.Parse()

	if *noColorOutput {
		errorStr = "ERROR! "
		warningStr = "WARNING! "
		successStr = "SUCCESS! "
	}

	if !*compressMode && !*decompressMode {
		fmt.Fprintln(os.Stderr, errorStr+"No mode set. Use --help for more information.")
		os.Exit(1)
	}

	paths := make([]string, 0, len(os.Args[2:]))
	for _, v := range os.Args[2:] {
		if !strings.HasPrefix(v, "-") {
			path, err := filepath.Abs(v)
			if err != nil {
				fmt.Fprintf(os.Stderr, errorStr+"An unknown error occurred when trying to read the set path (%s). Skipping...", v)
				continue
			}
			paths = append(paths, path)
		}
	}

	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, errorStr+"No paths specified. Use --help for more information.")
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
			fmt.Fprintf(os.Stderr, errorStr+"An unknown error occurred when trying to read %s. Does it exist, or maybe there's an issue with permissions?", path)
			total++
			failed++
			continue
		}

		if fileInfo.IsDir() && !*recursive {
			failed++
			fmt.Fprintf(os.Stderr, warningStr+"Ignoring %s since it's a directory and the recursion is disabled...", path)
			continue
		} else if fileInfo.IsDir() {
			dir, err := os.ReadDir(path)
			if err != nil {
				total++
				failed++
				fmt.Fprintf(os.Stderr, errorStr+"An unknown error occurred when trying to read %s as a directory. Does it exist, or maybe there's an issue with permissions?", path)
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
				fmt.Fprintf(os.Stderr, errorStr+"An unknown error occurred when trying to read %s. Does it exist, or maybe there's an issue with permissions?", path)
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
				fmt.Fprintf(os.Stderr, errorStr+"An error occurred during the conversion of the file %s: %s", path, err.Error())
				failed++
				continue
			}
			os.WriteFile(newPath, convertedFile, 0777)
			completed++
			fmt.Fprintf(os.Stdout, successStr+"Successfully converted %s", path)

			if *deleteOld {
				err := os.Remove(path)
				if err != nil {
					fmt.Fprintf(os.Stderr, warningStr+"The deletion flag is set, but I could not delete %s after the conversion!", path)
				}
			}
		}
	}
}

const helpText = "DVPL Converter by Endg4me_\n\nUsage: dvpl-go-tools <-c|-d> [-f] [-n] [-r] [-p] PATH [PATH...]\n\nOptions:\n"

func printUsage() {
	fmt.Print(helpText)
	flag.PrintDefaults()
}
