package controller

import (
	"testing"
)

func TestService_getStatus(t *testing.T) {
	s := &Service{}

	name := "sshd"
	status, err := s.getStatus(name)
	if err != nil {
		t.Error(name, " => error: ", err)
	} else {
		t.Log(name, " => status: ", status)
	}

	name = "openvpn"
	status, err = s.getStatus(name)
	if err != nil {
		t.Error(name, " => error: ", err)
	} else {
		t.Log(name, " => status: ", status)
	}
}
