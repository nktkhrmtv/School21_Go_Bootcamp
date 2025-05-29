package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	slPtr := flag.Bool("sl", false, "Find symbolic links")
	dPtr := flag.Bool("d", false, "Find directories")
	fPtr := flag.Bool("f", false, "Find files")
	extPtr := flag.String("ext", "", "Filter files by extension (requires -f)")
	flag.Parse()

	if !*slPtr && !*dPtr && !*fPtr {
		*slPtr, *dPtr, *fPtr = true, true, true
	}
	
	err := filepath.Walk(flag.Arg(0), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if *slPtr && info.Mode()&os.ModeSymlink != 0 {
			dest, err := os.Readlink(path)
			if err != nil {
				fmt.Printf("%s -> [broken]\n", path)
			} else {
				_, err = os.Stat(dest)
				if os.IsNotExist(err) {
					fmt.Printf("%s -> [broken]\n", path)
				} else {
					fmt.Printf("%s -> %s\n", path, dest)
				}
			}
		} else if *dPtr && info.IsDir() {
			fmt.Println(path)
		} else if *fPtr && info.Mode().IsRegular() {
			if *extPtr == "" || strings.HasSuffix(info.Name(), "."+*extPtr) {
				fmt.Println(path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}
}