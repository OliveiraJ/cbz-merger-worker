package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var rootFolderPath = os.Args[1]
var destinyFolder = os.Args[2]
var directorys []string
var pages []string
var pagenumber = 0

func main() {
	fmt.Println("Merging cbz files inside of " + os.Args[1])
	unzipCbzFiles(rootFolderPath)
	//Cimnhando pelos arquivos e pegando os nomes de cada página e pasta
	err := filepath.WalkDir(rootFolderPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			directorys = append(directorys, d.Name())

		} else {
			if filepath.Ext(d.Name()) == ".jpg" {
				pages = append(pages, d.Name())
			}
		}
		return nil
	})

	directorys = directorys[1:]
	if err != nil {
		fmt.Println(err)
	}
	//return the ammount of pages in the new .cbz file
	fmt.Println("The final cbz file will have: " + strconv.Itoa(len(pages)) + " pages")

	//create a destiny folder
	os.Mkdir(rootFolderPath+"/"+destinyFolder, 0755)

	//Copy files to the detiny folder and rename them to keep the right order
	for _, comicFolder := range directorys {
		err := filepath.WalkDir(rootFolderPath+"/"+comicFolder, func(path string, d os.DirEntry, err error) error {

			if err != nil {
				return err
			}
			if !d.IsDir() && filepath.Ext(d.Name()) == ".jpg" {
				pagenumber++

				originalPage, err := os.Open(rootFolderPath + "/" + comicFolder + "/" + d.Name())
				if err != nil {
					return err
				}
				defer originalPage.Close()

				if pagenumber < 10 {
					copyPage, err := os.Create(rootFolderPath + "/" + destinyFolder + "/" + "00" + strconv.Itoa(pagenumber) + ".jpg")
					if err != nil {
						return err
					}
					defer copyPage.Close()

					_, err = io.Copy(copyPage, originalPage)
					if err != nil {
						return err
					}
				} else if pagenumber < 100 {
					copyPage, err := os.Create(rootFolderPath + "/" + destinyFolder + "/" + "0" + strconv.Itoa(pagenumber) + ".jpg")
					if err != nil {
						return err
					}
					defer copyPage.Close()

					_, err = io.Copy(copyPage, originalPage)
					if err != nil {
						return err
					}
				} else {
					copyPage, err := os.Create(rootFolderPath + "/" + destinyFolder + "/" + strconv.Itoa(pagenumber) + ".jpg")
					if err != nil {
						return err
					}
					defer copyPage.Close()

					_, err = io.Copy(copyPage, originalPage)
					if err != nil {
						return err
					}
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	//Compress folder into a .cbz
	finalFile, err := os.Create(rootFolderPath + "/" + destinyFolder + ".zip")
	if err != nil {
		panic(err)
	}
	defer finalFile.Close()

	renamedFiles, err := ioutil.ReadDir(rootFolderPath + "/" + destinyFolder)
	if err != nil {
		panic(err)
	}

	zipWriter := zip.NewWriter(finalFile)

	for _, file := range renamedFiles {

		f, err := os.Open(rootFolderPath + "/" + destinyFolder + "/" + file.Name())
		if err != nil {
			panic(err)
		}
		defer f.Close()

		w, err := zipWriter.Create(destinyFolder + "/" + file.Name())
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(w, f); err != nil {
			panic(err)
		}
	}
	zipWriter.Close()

	err = os.Rename(rootFolderPath+"/"+destinyFolder+".zip", rootFolderPath+"/"+destinyFolder+".cbz")
	if err != nil {
		panic(err)
	}

	cleanFiles()

}

func unzipCbzFiles(rootFolderPath string) {
	files := []string{}
	err := filepath.WalkDir(rootFolderPath, func(path string, d os.DirEntry, err error) error {
		if !d.IsDir() && filepath.Ext(d.Name()) == ".cbz" {
			files = append(files, d.Name())
		}
		return err
	})
	if err != nil {
		panic(err)
	}

	for _, name := range files {
		pathInZip := strings.Replace(name, ".cbz", ".zip", 1)
		renameFiles(rootFolderPath+"/"+name, rootFolderPath+"/"+pathInZip)
	}

	for _, f := range files {
		f = strings.Replace(f, ".cbz", ".zip", 1)
		unzipSource(f, f)
	}

	for _, name := range files {
		name = strings.Replace(name, ".cbz", ".zip", 1)
		pathInZip := strings.Replace(name, ".zip", ".cbz", 1)
		renameFiles(rootFolderPath+"/"+name, rootFolderPath+"/"+pathInZip)
	}
}

func unzipSource(cbzName string, destination string) {
	destination = strings.Replace(destination, ".zip", "", 1)
	cbzOpened, err := zip.OpenReader(rootFolderPath + "/" + cbzName)
	if err != nil {
		panic(err)
	}
	defer cbzOpened.Close()

	if err := os.MkdirAll(rootFolderPath+"/"+destination, os.ModePerm); err != nil {
		panic(err)
	}

	for _, f := range cbzOpened.File {
		dstFile, err := os.OpenFile(rootFolderPath+"/"+destination+"/"+f.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInCbz, err := f.Open()
		if err != nil {
			panic(err)
		}
		fileInCbz.Close()

		if _, err := io.Copy(dstFile, fileInCbz); err != nil {
			panic(err)
		}

		dstFile.Close()
	}
}

func renameFiles(oldName, newName string) {
	err := os.Rename(oldName, newName)
	if err != nil {
		panic(err)
	}
}

func cleanFiles() {
	createdDirs := []string{}
	err := filepath.WalkDir(rootFolderPath, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			createdDirs = append(createdDirs, d.Name())
		}
		return err
	})
	if err != nil {
		panic(err)
	}

	for _, dir := range createdDirs {
		os.RemoveAll(rootFolderPath + "/" + dir)
	}

}
