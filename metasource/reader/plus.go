package reader

/*
#cgo pkg-config: createrepo_c
#include "createrepo_c/xml_parser.h"
#include "createrepo_c/package.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"unsafe"
)

func MakeDatabase() (bool, error) {
	var gexp *C.GError
	var expt error
	var iter *C.cr_PkgIterator
	var wait sync.WaitGroup
	var prmypath, filepath, othrpath string
	var prmyconv, fileconv, othrconv *C.char
	var prmypack, filepack, othrpack chan *C.cr_Package
	var prmyrslt, filerslt, othrrslt chan bool
	var pack *C.cr_Package

	prmypath = "/home/fedohide-origin/projects/metasource/rawhide/10beaa5fb8bb9b8710f4608ea9bf84aff2fb68e5efc7e82bf12b421867ad3d8f-primary.xml"
	filepath = "/home/fedohide-origin/projects/metasource/rawhide/4182e96bacb8bb0ccdcc9d446977416a2a18b49a4aed13d6c550be45a1bf061e-filelists.xml"
	othrpath = "/home/fedohide-origin/projects/metasource/rawhide/e3ac902af73897fe77cbc4df42d1c87d72ff4d69c9e792224d0ff3857f630e92-other.xml"
	prmyconv = C.CString(prmypath)
	fileconv = C.CString(filepath)
	othrconv = C.CString(othrpath)
	defer C.free(unsafe.Pointer(prmyconv))
	defer C.free(unsafe.Pointer(fileconv))
	defer C.free(unsafe.Pointer(othrconv))

	prmypack = make(chan *C.cr_Package, 10)
	filepack = make(chan *C.cr_Package, 10)
	othrpack = make(chan *C.cr_Package, 10)
	prmyrslt = make(chan bool, 10)
	filerslt = make(chan bool, 10)
	othrrslt = make(chan bool, 10)

	wait.Add(3)
	go PopulatePrmy(&wait, prmypack, prmyrslt)
	go PopulateFile(&wait, filepack, filerslt)
	go PopulateOthr(&wait, othrpack, othrrslt)

	iter = C.cr_PkgIterator_new(prmyconv, fileconv, othrconv, nil, nil, nil, nil, &gexp)
	if iter == nil {
		expt = errors.New(fmt.Sprintf("%s", C.GoString(gexp.message)))
		slog.Log(nil, slog.LevelError, fmt.Sprintf("%s", expt.Error()))
		return false, expt
	}
	defer C.cr_PkgIterator_free(iter, &gexp)

	for C.cr_PkgIterator_is_finished(iter) == 0 {
		pack = C.cr_PkgIterator_parse_next(iter, &gexp)
		if pack == nil {
			break
		}
		prmypack <- pack
		filepack <- pack
		othrpack <- pack

		prmydone, _ := <-prmyrslt
		filedone, _ := <-filerslt
		othrdone, _ := <-othrrslt

		if prmydone && filedone && othrdone {
			head := fmt.Sprintf("%s %s:%s-%s.%s", C.GoString(pack.name), C.GoString(pack.epoch), C.GoString(pack.version), C.GoString(pack.release), C.GoString(pack.arch))
			slog.Log(nil, slog.LevelInfo, fmt.Sprintf("%s added.", head))
		}
	}
	close(prmypack)
	close(filepack)
	close(othrpack)
	close(prmyrslt)
	close(filerslt)
	close(othrrslt)

	wait.Wait()

	defer C.free(unsafe.Pointer(pack))
	defer C.free(unsafe.Pointer(gexp))
	return true, nil
}
