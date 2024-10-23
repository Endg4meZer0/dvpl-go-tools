# dvpl-go-tools | A binary for simple operations with DVPL compression format.

## Overview
World of Tanks Blitz and Tanks Blitz use a custom compression format named DVPL that is actually a LZ4 compression, usually of a level 2 with the exception of level 0 for .tex files, with the addition of a special footer. This binary makes use of my [specific Go package](https://github.com/Endg4meZer0/dvpl-go) to make any kind of simple (de-)compression work.

## Install
Assuming you have the go toolchain installed

```
go install github.com/Endg4meZer0/dvpl-go-tools
```

## Usage
```
dvpl-go-tools <-c|-d> [-f] [-n] [-r] [-p] PATH [PATH...]

Options:
  -c    Sets the mode to 'compression'.
  -d    Sets the mode to 'decompression'.
  -f    Force the compression algorithm to always use compression instead of detecting .tex files and applying no compression on them.
  -n    Delete the old file after converting.
  -p    Use plain output / disable colored output.
  -r    Recursively convert all files and the contents of all folders inside the set path.
```
