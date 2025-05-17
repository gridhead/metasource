package test

import (
	"fmt"
	"metasource/metasource/config"
	"metasource/metasource/driver"
	"os"
)

func WipeGeneration(loca string) {
	_, expt := os.Stat(loca)
	if expt != nil {
		return
	} else {
		_ = os.RemoveAll(loca)
	}
}

func CopyGeneration(srce string, dest string) error {
	srcedata, expt := os.ReadFile(srce)
	if expt != nil {
		return expt
	}
	expt = os.WriteFile(dest, srcedata, os.ModePerm)
	if expt != nil {
		return expt
	}
	return nil
}

func Path_UnInit(temp string) (revert func()) {
	origpath := config.DBFOLDER
	config.DBFOLDER = temp
	return func() {
		config.DBFOLDER = origpath
	}
}

func Path_Init() (revert func()) {
	origpath := config.DBFOLDER
	basepath := "./assets"
	destpath := fmt.Sprintf("%s/test-%s", basepath, driver.GenerateIdentity(&config.RANDOM_LENGTH))
	config.DBFOLDER = destpath
	_ = os.MkdirAll(destpath, 0755)
	for _, item := range []string{"primary", "filelists", "other"} {
		_ = CopyGeneration(
			fmt.Sprintf("%s/%s", basepath, fmt.Sprintf("testbase_%s.sqlite", item)),
			fmt.Sprintf("%s/%s", destpath, fmt.Sprintf("metasource-rawhide-%s.sqlite", item)),
		)
	}
	return func() {
		config.DBFOLDER = origpath
		WipeGeneration(destpath)
	}
}
