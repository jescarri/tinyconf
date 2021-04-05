package apt

import (
	"fmt"
	"github.com/jescarri/tinyconf/pkg/exec"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{})
}

const pkgNotInstalled = "not-installed"
const pkgInstalled = "installed"
const pkgUpdated = "updated"
const pkgRemoved = "removed"
const pkgNotFound = "not-found"
const pkgUpdateNeeded = "update-needed"

var e = &exec.Exec{Shell: "/bin/sh"}

type AptPkg struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	State   string `yaml:"state"`
}

// getPkgName returns the full name of the package in debian
// based distros is name=version throw an error if pkg version
// is not defined
func (p *AptPkg) getPkgName() (string, error) {
	if p.Version == "" {
		log.Errorf("Package: %s requires a version")
		return "", fmt.Errorf("Package: %s requires a version")
	}
	return fmt.Sprintf("%s=%s", p.Name, p.Version), nil
}

// pkgSearch pkgList runs apt and searches for a pkg name installed, returns array of
// strings and error if there's an error executing the commnand
func (p *AptPkg) pkgSearch() ([]string, error) {
	stdout, stderr, err := e.Run(fmt.Sprintf("apt list --installed %s", p.Name))
	if err != nil {
		log.Errorf("Can't determine if pkg is installed: %s", stdout)
		log.Debugf("%s", stderr)
		return []string{}, err
	}
	return strings.Split(stdout, "\n"), nil
}

// isInstalled checks if a pkg is already installed
// Returns string(pkgNotInstalled,pkgInstalled,pkgNotFound,pkgUpdateNeeded) and Error if the pkg is not a valid
// pkg name
func (p *AptPkg) isInstalled() (bool, error) {
	lines, err := p.pkgSearch()
	if err != nil {
		return false, err
	}
	for _, l := range lines {
		if strings.Contains(l, p.Name) {
			if strings.Contains(l, p.Version) {
				p.State = pkgInstalled
			} else {
				p.State = pkgUpdateNeeded
			}
			return true, nil
		}
	}
	return false, nil
}

func (p *AptPkg) runAptInstall() error {
	fullPkgName, err := p.getPkgName()
	if err != nil {
		return err
	}
	stdout, stderr, err := e.Run(fmt.Sprintf("apt install -y %s", fullPkgName))

	if err != nil {
		log.Errorf("Error occurred while installing package '%s': %s", fullPkgName, stderr)
	}
	log.Debug(stdout)
	return err
}

func (p *AptPkg) Install() error {
	installed, err := p.isInstalled()
	if err != nil {
		return err
	}
	if installed {
		if p.State == pkgUpdateNeeded {
			log.Infof("Pkg: %s is not on required version, proceeding to upgrade it to: %s", p.Name, p.Version)
			err := p.runAptInstall()
			p.State = pkgUpdated
			return err
		}
		log.Infof("Pkg: %s is already installed and on required version: %s", p.Name, p.Version)
		return nil
	}
	log.Infof("Installing pkg: %s version: %s", p.Name, p.Version)
	err = p.runAptInstall()
	if err != nil {
		p.State = pkgInstalled
	}
	return err
}

func (p *AptPkg) Remove() error {
	installed, err := p.isInstalled()
	if err != nil {
		return err
	}
	if installed {
		fullPkgName, err := p.getPkgName()
		if err != nil {
			return err
		}
		log.Infof("Removing pkg: %s version :%s", p.Name, p.Version)
		stdout, stderr, err := e.Run(fmt.Sprintf("apt remove -y %s", fullPkgName))
		if err != nil {
			log.Errorf("Error occurred while removing package '%s': %s", fullPkgName, stderr)
			log.Debug("%s", err)
			return err
		}
		log.Debug(stdout)
		p.State = pkgRemoved
	}
	return nil
}
