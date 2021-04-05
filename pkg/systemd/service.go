package systemd

import (
	"fmt"
	"github.com/jescarri/tinyconf/pkg/exec"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{})
}

var e = &exec.Exec{Shell: "/bin/sh"}

func EnableUnit(unit string) error {
	log.Infof("Enabling systemd unit: %s", unit)
	stdout, stderr, err := e.Run(fmt.Sprintf("systemctl enable %s", unit))
	if err != nil {
		log.Errorf("err: %s, stderr: %s", err, stderr)
	}
	log.Debugf("%s", stdout)
	return err
}

func StartUnit(unit string) error {
	log.Infof("starting systemd unit: %s", unit)
	stdout, stderr, err := e.Run(fmt.Sprintf("systemctl start %s", unit))
	if err != nil {
		log.Errorf("err: %s, stderr: %s", err, stderr)
	}
	log.Debugf("%s", stdout)
	return err
}

func RestartUnit(unit string) error {
	log.Infof("Restarting systemd unit: %s", unit)
	stdout, stderr, err := e.Run(fmt.Sprintf("systemctl restart %s", unit))
	if err != nil {
		log.Errorf("err: %s, stderr: %s", err, stderr)
	}
	log.Debugf("%s", stdout)
	return err
}
