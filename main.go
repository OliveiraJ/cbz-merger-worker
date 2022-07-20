package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	fmt.Println("Merging cbz files inside of " + os.Args[1])

	var rootFolderPath = os.Args[1]
	var destinyFolder = os.Args[2]
	var directorys []string
	var pages []string
	var pagenumber = 0

	//Cimnhando pelos arquivos e pegando os nomes de cada página e pasta
	err := filepath.WalkDir(rootFolderPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			directorys = append(directorys, d.Name())

		} else {
			pages = append(pages, d.Name())
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
			if !d.IsDir() {
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

}
