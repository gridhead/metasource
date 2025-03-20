package reader

/*
   #cgo pkg-config: createrepo_c
   #include <stdlib.h>
   #include "createrepo_c/xml_parser.h"
   #include <createrepo_c/sqlite.h>
   #include <createrepo_c/package.h>
   #include <createrepo_c/createrepo_c.h>
*/
import "C"

import (
	"sync"
	"unsafe"
)

func MakeDatabase() (bool, error) {
	var wait sync.WaitGroup
	var prmypath, filepath, othrpath string
	var prmyconv, fileconv, othrconv *C.char

	prmypath = "/home/fedohide-origin/projects/metasource/rawhide/10beaa5fb8bb9b8710f4608ea9bf84aff2fb68e5efc7e82bf12b421867ad3d8f-primary.xml"
	filepath = "/home/fedohide-origin/projects/metasource/rawhide/4182e96bacb8bb0ccdcc9d446977416a2a18b49a4aed13d6c550be45a1bf061e-filelists.xml"
	othrpath = "/home/fedohide-origin/projects/metasource/rawhide/e3ac902af73897fe77cbc4df42d1c87d72ff4d69c9e792224d0ff3857f630e92-other.xml"
	prmyconv = C.CString(prmypath)
	fileconv = C.CString(filepath)
	othrconv = C.CString(othrpath)
	defer C.free(unsafe.Pointer(prmyconv))
	defer C.free(unsafe.Pointer(fileconv))
	defer C.free(unsafe.Pointer(othrconv))

	wait.Add(3)
	go PopulatePrmy(&wait, &prmyconv, &fileconv, &othrconv)
	go PopulateFile(&wait, &prmyconv, &fileconv, &othrconv)
	go PopulateOthr(&wait, &prmyconv, &fileconv, &othrconv)
	wait.Wait()

	return true, nil
}
