package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type grf func(apath string) interface{}

var Dgrf grf = func(apath string) interface{} {
	return nil
}

func IterFiles(root string, goDot bool, getRootElem grf, w func(level int32, path string, apath string, rootElem interface{}) error) error {
	f, err := os.Stat(root)
	if err != nil {
		return err
	}
	if !f.IsDir() {
		return fmt.Errorf("%s is not a dir", root)
	}

	l := int32(0)

	pathSep := string(os.PathSeparator)

	var rf func(root string, level int32, apath string) error
	rf = func(root string, level int32, apath string) error {
		dir, err := ioutil.ReadDir(root)
		if err != nil {
			return err
		}
		rootElem := getRootElem(apath)
		for _, fi := range dir {
			if strings.Index(fi.Name(), ".") == 0 && goDot == false {
				continue
			}
			tPath := apath
			if tPath != "" {
				tPath = tPath + "/"
			}
			if fi.IsDir() {
				if fi.Name() == "." || fi.Name() == ".." {
					continue
				}

				err = rf(root+pathSep+fi.Name(), level+1, tPath+fi.Name())
				if err != nil {
					return err
				}
			} else {
				err := w(level, root+pathSep+fi.Name(), tPath+fi.Name(), rootElem)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	err = rf(root, l, "")
	return err
}

func IterDir(root string, goDot bool, rootElem interface{}, w func(level int32, path string, apath string, rootElem interface{}) (interface{}, error)) error {
	f, err := os.Stat(root)
	if err != nil {
		return err
	}
	if !f.IsDir() {
		return fmt.Errorf("%s is not a dir", root)
	}

	l := int32(0)

	pathSep := string(os.PathSeparator)

	var rf func(root string, level int32, apath string, rootElem interface{}) error
	rf = func(root string, level int32, apath string, rootElem interface{}) error {
		dir, err := ioutil.ReadDir(root)
		if err != nil {
			return err
		}
		for _, fi := range dir {
			if strings.Index(fi.Name(), ".") == 0 && goDot == false {
				continue
			}
			tPath := apath
			if tPath != "" {
				tPath = tPath + "/"
			}
			if fi.IsDir() {
				if fi.Name() == "." || fi.Name() == ".." {
					continue
				}
				nre, err := w(level, root+pathSep+fi.Name(), tPath+fi.Name(), rootElem)
				if err != nil {
					return err
				}
				err = rf(root+pathSep+fi.Name(), level+1, tPath+fi.Name(), nre)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	err = rf(root, l, "", rootElem)
	return err
}
