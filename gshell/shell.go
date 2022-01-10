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
	Prepares  []*InfoPrepare

	kill func() error
}

func (s *shell) Run() {
	s.runPrepares()

	name := s.cmdExePath(s.Exec)
	args := s.cmdArgs(s.Args)
	cmd := exec.Command(name, args...)
	cmd.Dir = s.Directory
	cmd.Stdout = log
	cmd.Stderr = log
	err := cmd.Start()
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

func (s *shell) runPrepares() {
	c := len(s.Prepares)
	for i := 0; i < c; i++ {
		prepare := s.Prepares[i]
		if prepare == nil {
			continue
		}

		s.runPrepare(prepare)
	}
}

func (s *shell) runPrepare(prepare *InfoPrepare) {
	if prepare == nil {
		return
	}
	name := s.cmdExePath(prepare.Exec)
	if len(name) < 1 {
		return
	}

	args := s.cmdArgs(prepare.Args)
	cmd := exec.Command(name, args...)
	cmd.Dir = s.Directory
	cmd.Stdout = log
	cmd.Stderr = log
	err := cmd.Start()
	if err != nil {
		log.Error(fmt.Sprintf("start prepare (%s) fail: ", prepare.Exec), err)
		log.Output("")
	} else {
		log.Info(fmt.Sprintf("start prepare (%s) [%d] success: %s", prepare.Exec, cmd.Process.Pid, cmd.String()))
		defer cmd.Process.Kill()
		defer log.Output("")
		defer log.Info(fmt.Sprintf("prepare (%s) [%d] exited", prepare.Exec, cmd.Process.Pid))
	}

	cmd.Wait()
}

func (s *shell) cmdExePath(name string) string {
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

	return name
}

func (s *shell) cmdArgs(v string) []string {
	args := make([]string, 0)
	values := strings.Split(v, " ")
	c := len(values)
	for i := 0; i < c; i++ {
		value := strings.TrimSpace(values[i])
		if len(value) > 0 {
			args = append(args, value)
		}
	}

	return args
}
