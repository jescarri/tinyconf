package main

import (
	"github.com/jescarri/tinyconf/pkg/config"
	"github.com/jescarri/tinyconf/pkg/systemd"
)

func main() {
	c, err := config.LoadConfig("./a.yaml")
	if err != nil {
		panic(err)
	}
	//install pkgs
	for _, p := range c.InstallPkgs {
		err := p.Install()
		if err != nil {
			panic(err)
		}
	}
	//uinstall pkgs
	for _, d := range c.RemovePkgs {
		err := d.Remove()
		if err != nil {
			panic(err)
		}

	}
	// Write Files.
	for _, f := range c.Files {
		err := f.Write()
		if err != nil {
			panic(err)
		}
		if f.Service != "" && f.ServiceRestart {
			err := systemd.RestartUnit(f.Service)
			if err != nil {
				panic(err)
			}
		}
	}
}
