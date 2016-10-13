package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
)

type Repository struct {
	RepositoryFolder string
	ManifestPath     string
	KeyPath          string
	Key              []byte
}

func (rep *Repository) GetManifest() *os.File {
	if !rep.isRepoInitialized() {
		panic(fmt.Errorf("Repository not initialized. Aborting."))
	}
	f, err := os.OpenFile(rep.ManifestPath, os.O_RDWR, 0666) // FIX: cambiar por readonly
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			fmt.Println("Manifest not found. Did you initialize the repository?")
			panic(err)
		}
	}
	return f
}

func (rep *Repository) checkIfAlreadyExists() bool {
	_, err := os.Stat(rep.RepositoryFolder)
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			return false
		}
	}
	return true
}

func (rep *Repository) InitializeRepository() {
	rep.RepositoryFolder = rep.generateRepositoryPath()
	if rep.checkIfAlreadyExists() {
		panic(fmt.Errorf("Found previously existing repository. Try rebuilding.\n"))
	}
	rep.ManifestPath = filepath.Join(rep.RepositoryFolder, "manifest")
	rep.KeyPath = filepath.Join(rep.RepositoryFolder, "key")

	ou := syscall.Umask(0)
	err := os.MkdirAll(rep.RepositoryFolder, 0755)
	if err != nil {
		fmt.Printf("Could not create registry folder: %s\n", rep.RepositoryFolder)
		panic(err)
	}
	syscall.Umask(ou)

	makeFile := func(fp string) {
		f, err := os.OpenFile(fp, os.O_RDWR, 0666)
		defer f.Close()
		if err != nil {
			if strings.Contains(err.Error(), "no such file") {
				fmt.Printf("File %s not found, creating.\n", fp)
				f, err = os.Create(fp)
				if err == nil {
					return
				}
			}
			fmt.Printf("Could not create file: %s! Aborting!\n", f)
			panic(err)
		} else {
			panic(fmt.Errorf("File %s already exists. Aborting.", fp))
		}
	}

	makeFile(rep.KeyPath)
	makeFile(rep.ManifestPath)

	rep.createKey()
}

func (rep *Repository) RebuildRepository() {
	rep.RepositoryFolder = rep.generateRepositoryPath()
	rep.ManifestPath = filepath.Join(rep.RepositoryFolder, "manifest")
	rep.KeyPath = filepath.Join(rep.RepositoryFolder, "key")
}

func (rep *Repository) generateRepositoryPath() string {
	u, err := user.Current()
	if err != nil {
		fmt.Println("Could not get current user!")
		panic(err)
	}
	return filepath.Join(u.HomeDir, "badgopher.mnf")
}

func (rep *Repository) createKey() []byte {
	k := GenerateAesKey()
	err := ioutil.WriteFile(rep.KeyPath, k, 0666)
	if err != nil {
		fmt.Printf("Could not save key: %s\n", rep.KeyPath)
		panic(err)
	} else {
		fmt.Printf("Saved key to %s\n", rep.KeyPath)
	}
	rep.Key = k
	return rep.Key
}

func (rep *Repository) GetDaKey() []byte {
	if !rep.isRepoInitialized() {
		panic(fmt.Errorf("Repository not initialized. Aborting."))
	}
	if len(rep.Key) != 0 {
		return rep.Key
	}
	k, err := ioutil.ReadFile(rep.KeyPath)
	if err != nil {
		panic(fmt.Errorf("Could not read key!"))
	}
	rep.Key = k
	return k
}

func (rep *Repository) isRepoInitialized() bool {
	if rep.KeyPath == "" || rep.ManifestPath == "" {
		return false
	}
	return true
}
