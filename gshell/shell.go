package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type shell struct {
	Directory string
	Exec      string
	Args      string

	kill func() error
}

func (s *shell) Run() {
	name := s.Exec
	path, err := filepath.Abs(name)
	if err == nil {
		fi, fe := os.Stat(path)
		if !os.IsNotExist(fe) {
			if !fi.IsDir() {
				name = path

				if runtime.GOOS == "linux" {
					err = os.Chmod(name, 0700)
					if err != nil {
						log.Error("赋予启动文件可执行权限失败: ", err)
					}
				}
			}
		}
	}

	args := strings.Split(s.Args, " ")
	cmd := exec.Command(name, args...)
	cmd.Dir = s.Directory
	cmd.Stdout = log
	cmd.Stderr = log
	err = cmd.Start()
	if err != nil {
		log.Error("start exec fail: ", err)
	} else {
		log.Info(fmt.Sprintf("start exec success [%d]: %s", cmd.Process.Pid, cmd.String()))
		s.kill = cmd.Process.Kill
	}

	cmd.Wait()
}

func (s *shell) Shut() {
	if s.kill != nil {
		s.kill()
	}
}
