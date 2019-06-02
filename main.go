package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	rundir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	var (
		dirpath       = flag.String("path", rundir, "dir path")
		whitelistpath = flag.String("white", "default", "whitelist dir path")
		createnowwhitelist = flag.Bool("newwhitelist",false,"create new now whitelist")
		newwhitelistname = flag.String("new","","Name and path of new whitelist")
		newwhitelistpath = flag.String("newpath","","Path of new whitelist path")
	)

	flag.Parse()

	if *createnowwhitelist {
			if len(*newwhitelistname) != 0{
					if len(*newwhitelistpath) != 0 {
							CreateWhitelist(*newwhitelistname,*newwhitelistpath)
							os.Exit(1)
					}else{
							fmt.Println("usage : -newwhitelist -new PATH/WHITELISTNAME -newpath PATH/DIRCTORY")
							os.Exit(1)
					}
			}else{
					fmt.Println("usage : -newwhitelist -new PATH/WHITELISTNAME -newpath PATH/DIRCTORY")
					os.Exit(1)
			}
	}

	whitelists := OpenWhitelist(*whitelistpath)
	//fmt.Println(whitelists)

	/*for _, i := range DirExplore(dirpath) {
		fmt.Println(i)
	}*/
	detectedfiles := DirExplore(*dirpath)
	//fmt.Println(detectedfiles)
	var result []string

	for _, detectedfile := range detectedfiles {
			if CheckFileExsist(detectedfile,whitelists) {
					result = append(result,detectedfile)
			}
	}
	if len(result) >= 1 {
			for num, notpermitedfile := range result {
					fullpathfile, err := filepath.Abs(notpermitedfile)
					if err != nil{
							log.Fatal(err)
					}
					fmt.Printf("%d : [ %s ] --  NOT PERMITED putting Web ROOT directory\n",num,fullpathfile)
			}
	}
}


func DirExplore(path string) []string {
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var detectpaths []string

	for _, dir := range dirs {
		/*currentdir, err := filepath.Abs(dir.Name())
		if err != nil{
				log.Fatal(err)
		}
		fmt.Printf("%d : %s\n", n,currentdir)*/
		if dir.IsDir() {
			detectpaths = append(detectpaths, DirExplore(filepath.Join(path, dir.Name()))...)
			//DirExplore(filepath.Join(path,dir.Name()))
		}
		detectpaths = append(detectpaths, dir.Name())
	}
	return detectpaths
}

func OpenWhitelist(listpath string)[]string {

	whitelist, err := os.Open(listpath)
	if err != nil {
		fmt.Println("[USAGE : -path TARGETDIRECTORY -white PATH/WHITELIST")
		log.Fatal(err)
	}
	defer whitelist.Close()

	var whitelists []string

	sc := bufio.NewScanner(whitelist)

	for i := 1; sc.Scan(); i++ {
		if err := sc.Err(); err != nil {
			log.Fatal(err)
		}
		whitelists = append(whitelists,sc.Text())
	}
	return whitelists
}

func CheckFileExsist(unknownfile string, whitelists []string)bool{
		for _, whitelist := range whitelists {
				//fmt.Printf(" %d : whitelist : %s -- unknownfile : %s\n",i,whitelist,unknownfile)
				if whitelist == unknownfile {
						//fmt.Println("MATCHED!!")
						return false
				}
		}
		return true
}

func CreateWhitelist(newfilename string,rundir string){
		newfile, err := os.Create(newfilename)
		if err != nil {
				log.Fatal(err)
		}
		defer newfile.Close()

		NewComponents := DirExplore(rundir)
		for _, NewComponent := range NewComponents {
				newfile.WriteString(NewComponent+"\n")
		}


}
