package config

import (
	"github.com/jescarri/tinyconf/pkg/apt"
	"github.com/jescarri/tinyconf/pkg/file"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	InstallPkgs []apt.AptPkg `yaml:"install_pkgs"`
	RemovePkgs  []apt.AptPkg `yaml:"remove_pkgs"`
	Files       []file.File  `yaml:"files"`
	OnBootSvcs  []string     `yaml:"onboot_svcs"`
}

func LoadConfig(file string) (*Config, error) {
	c := &Config{}
	conf, err := ioutil.ReadFile(file)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(conf, c)
	if err != nil {
		return c, err
	}
	return c, nil
}
