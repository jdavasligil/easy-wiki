/*
Ewik generates static personalized wikis.

Without an explicit path, it assumes the current working directory as a target.
Given a directory, it generates necessary static files in that directory from
the options specified in a _config.toml file and markdown pages in the pages
directory. Basic directory structure and config files are created with init.

Usage:

	ewik [flags] [path]

Flags:

	-watch
	    Watch directory for live changes. Any changes to _config.toml or pages
	    in the pages directory will trigger a render pass.

	-init
	    Initialize the directory structure and config files.
*/
package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

func initialize() {
	fmt.Println("Generate _config.toml")
	fmt.Println("Create directories")
	fmt.Println("Generate static files")
}

func main() {
	var watchFlag bool
	var watchFlagHelp = `Watch directory for live changes to _config.toml and markdown pages in
the pages directory.`

	var initFlag bool
	var initFlagHelp = `Generate _config.toml, directory structure, and static web files for wiki.`

	flag.BoolVar(&watchFlag, "watch", false, watchFlagHelp)
	flag.BoolVar(&initFlag, "init", false, initFlagHelp)
	flag.Parse()

	filePathString := flag.Arg(0)
	if filePathString == "" {
		filePathString = "."
	}
	filePath := filepath.Join(filePathString)
	fmt.Println("Wiki Path:", filePath)

	ssg := NewStaticSiteGenerator(filePath)

	if initFlag {
		ssg.Initialize()
	}
}
