package main

import (
	"os"
	"path/filepath"
)

const (
	moduleVersion = "1.0.1.0"
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

	if len(svcDir) > 0 {
		os.Chdir(svcDir)
	}
	curDir, _ := os.Getwd()

	log.Info("shell run at: ", path)
	log.Info("shell version: ", moduleVersion)
	log.Info("log folder: ", log.folder)
	log.Info("cur folder: ", curDir)
	log.Info("svc folder: ", svcDir)

	server.shell.Directory = svcDir
}
