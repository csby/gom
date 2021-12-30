package main

import (
	"os"
	"path/filepath"
)

var (
	svcDir = ""
	server = &program{}
	log    = &LogWriter{}
)

// args['exec', 'svc-folder', 'log-folder']
func init() {
	args := os.Args
	path, _ := filepath.Abs(args[0])

	argc := len(args)
	if argc > 1 {
		svcDir = args[1]
	} else {
		svcDir = filepath.Dir(path)
	}
	if argc > 2 {
		log.folder = args[2]
	}

	log.Info("shell run at: ", path)
	log.Info("svc folder: ", svcDir)
	log.Info("log folder: ", log.folder)

	if len(svcDir) > 0 {
		os.Chdir(svcDir)
	}
	curDir, _ := os.Getwd()
	log.Info("current directory: ", curDir)

	server.shell.Directory = svcDir
}
