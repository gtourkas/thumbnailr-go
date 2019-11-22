// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/nozzle/throttler"
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

const (
	DefaultWorkingDir = "../src"
	DefaultHandlersSrcPath = "../src/lambdas"
	DefaultHandlersDestPath = "../build"
	DefaultMaxConcurBuilds = 2
)

var conf struct {
	WorkingDir string
	HandlersSrcPath string
	HandlersDestPath string
	MaxConcurBuilds int
}

func readConf() error {
	ini, err := ini.Load("conf.ini")
	if err != nil {
		return errors.Wrap(err,"cannot read conf.ini")
	}

	conf.WorkingDir = ini.Section("").Key("working_dir").MustString(DefaultWorkingDir)
	conf.HandlersSrcPath = ini.Section("").Key("handlers_src_path").MustString(DefaultHandlersSrcPath)
	conf.HandlersDestPath = ini.Section("").Key("handlers_dest_path").MustString(DefaultHandlersDestPath)
	conf.MaxConcurBuilds = ini.Section("").Key("max_concur_builds").MustInt(DefaultMaxConcurBuilds)

	return nil
}

func Build() error {
	mg.Deps(readConf)

	workingDir := conf.WorkingDir
	srcPath := conf.HandlersSrcPath
	destPath := conf.HandlersDestPath
	maxConcurBuilds := conf.MaxConcurBuilds
	debug := false
	if strings.ToLower(os.Getenv("TN_DEBUG")) == "true" {
		debug = true
	}

	if err := os.Chdir(workingDir); err != nil {
		return errors.Wrapf(err, "cannot change to the working dir %s", workingDir)
	}

	files, err := ioutil.ReadDir(srcPath)
	if err != nil {
		return errors.Wrapf(err, "cannot read handlers src path %s", srcPath)
	}

	// count handler dirs first
	dirsCnt := 0
	singleHandler := os.Getenv("TNHANDLER")
	if singleHandler != "" {
		if _, err := os.Stat(path.Join(srcPath,singleHandler)); os.IsNotExist(err) {
			return fmt.Errorf("single handler %s does not exist", singleHandler)
		}

		err := buildSingle(srcPath, destPath, singleHandler, debug)
		if err != nil {
			log.Printf("handler %s build failed: %v\n", singleHandler, err)
		}
		return nil
	}

	for _, f := range files {
		if f.IsDir() {
			dirsCnt ++
		}
	}

	// concurrently run the handler builds
	t := throttler.New(maxConcurBuilds, dirsCnt)
	for _, f := range files {
		if f.IsDir() {
			go func(handlerDir string) {
				err := buildSingle(srcPath, destPath, handlerDir, debug)
				if err != nil {
					log.Printf("handler %s build failed: %v\n", handlerDir, err)
				}
				t.Done(err)
			}(f.Name())
			errCnt := t.Throttle()
			if errCnt > 0 {
				return errors.New("one or more handler builds have failed")
			}
		}
	}

	return nil
}

func buildSingle(srcPath string, destPath string, handlerDir string, debug bool) error {
	env := map[string]string{
		"GOOS": "linux",
		"GOARCH": "amd64",
	}
	output := path.Join(destPath, handlerDir, "main")
	input := path.Join(srcPath, handlerDir)

	log.Printf("building handler %s ...\n", handlerDir)

	debugArg := "-gcflags="
	if debug {
		debugArg = "-gcflags=-N -l"
	}
	if err := sh.RunWith(env, "go", "build", debugArg, "-o", output, input); err != nil {
		return errors.Wrapf(err,"cannot go build %s", input)
	}
	return nil
}


func Clean() error {
	mg.Deps(readConf)

	destPath := conf.HandlersDestPath

	files, err := ioutil.ReadDir(destPath)
	if err != nil {
		return errors.Wrapf(err, "cannot read handlers dest path %s", destPath)
	}

	for _, f := range files {
		if f.IsDir() {
			handlerDir := f.Name()
			handlerPath := path.Join(destPath, handlerDir)
			if err := os.RemoveAll(handlerPath); err != nil {
				return errors.Wrapf(err, "cannot delete handler dest path %s", handlerPath)
			}
			log.Printf("deleted built handler %s", handlerDir)
		}
	}
	return nil
}
