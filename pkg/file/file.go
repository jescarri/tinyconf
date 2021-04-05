package file

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{})
}

type File struct {
	Name           string      `yaml:"name"`
	GID            int         `yaml:"gid"`
	UID            int         `yaml:"uid"`
	Mode           os.FileMode `yaml:"mode"`
	Content        string      `yaml:"content"`
	Service        string      `yaml:"service"`
	ServiceRestart bool        `yaml:"serviceRestart"`
}

func (f *File) exists() (bool, error) {
	if _, err := os.Stat(f.Name); os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}

func (f *File) changeMode() error {
	exists, err := f.exists()
	if err != nil {
		return err
	}
	if exists {
		err := os.Chmod(f.Name, f.Mode)
		if err != nil {
			return err
		}
		err = os.Chown(f.Name, f.UID, f.GID)
		return err
	}
	return fmt.Errorf("File does not exists")
}

// Write writes the file to disk.
// ToDo: use a buffered writer.
func (f *File) Write() error {
	exists, err := f.exists()
	if !exists {
		err := ioutil.WriteFile(f.Name, []byte(f.Content), f.Mode)
		if err != nil {
			log.Errorf("%s", err)
		}
	}
	// Check if desired content matches the content on disk\
	diskContent, err := ioutil.ReadFile(f.Name)
	if err != nil {
		log.Errorf("%s", err)
		return err
	}
	fileHash := getHash(diskContent)
	configHash := getHash([]byte(f.Content))
	if fileHash != configHash {
		log.Infof("File: %s requires change", f.Name)
		err := os.Chmod(f.Name, 0777)
		if err != nil {
			return err
		}

		if f.Service != "" {
			log.Infof("Service: %s will require restart", f.Service)
			f.ServiceRestart = true
		}
		err = ioutil.WriteFile(f.Name, []byte(f.Content), f.Mode)
		if err != nil {
			return err
		}
	}
	// as cleanup always set the file permissions to whatever is in the config
	return f.changeMode()
}

func getHash(content []byte) string {
	md5HashInBytes := md5.Sum(content)
	return hex.EncodeToString(md5HashInBytes[:])
}
