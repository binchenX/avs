// Package utils provides helper functions
package utils

// FileExists check if file or dir exsit
import (
	"fmt"
	"os"
)

// FileExists is file exsit
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// CreateGenDirIfNotExsit create gen dir
func CreateGenDirIfNotExsit(genDirAbs string) bool {
	if r, _ := FileExists(genDirAbs); r == true {
		return true
	}

	var err = os.MkdirAll(genDirAbs, 0777)
	if err == nil {
		return true
	}

	fmt.Println("err: faild to mkdiir", genDirAbs, err)
	return false
}

// IncludedIn check if all items in a is also in b
func IncludedIn(a, b []string) bool {
	for _, t := range a {
		find := false
		for _, p := range b {
			if t == p {
				find = true
				break
			}
		}
		if find == false {
			return false
		}
	}
	return true
}
