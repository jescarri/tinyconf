package main

import (
	"fmt"
	"github.com/jescarri/tinyconf/pkg/config"
	"github.com/jescarri/tinyconf/pkg/systemd"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"os"
)

var cfgFile string

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{})
	flag.StringVar(&cfgFile, "config-file", "tinyconf.yaml", "config file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "tinyconf usage:.\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()
}

func main() {
	c, err := config.LoadConfig(cfgFile)
	if err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}
	//install pkgs
	for _, p := range c.InstallPkgs {
		err := p.Install()
		if err != nil {
			log.Errorf("%s", err)
			os.Exit(1)
		}
	}
	// Write Files.
	for _, f := range c.Files {
		err := f.Write()
		if err != nil {
			log.Errorf("%s", err)
			os.Exit(1)
		}
		if f.Service != "" && f.ServiceRestart {
			err := systemd.RestartUnit(f.Service)
			if err != nil {
				log.Errorf("%s", err)
				os.Exit(1)
			}
		}
	}
	//uinstall pkgs
	for _, d := range c.RemovePkgs {
		err := d.Remove()
		if err != nil {
			log.Errorf("%s", err)
			os.Exit(1)
		}
	}
	for _, s := range c.OnBootSvcs {
		log.Infof("Enabling systemd unit: %s", s)
		err := systemd.EnableUnit(s)
		if err != nil {
			log.Errorf("%s", err)
		}
		log.Infof("Making shure that unit: %s is started", s)
		err = systemd.StartUnit(s)
		if err != nil {
			log.Errorf("%s", err)
		}
	}
}
