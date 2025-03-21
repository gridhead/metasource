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

func MakeDatabase() (int64, error) {
	var gexp *C.GError
	var expt error
	var iter *C.cr_PkgIterator
	var wait sync.WaitGroup
	var prmypath, filepath, othrpath string
	var prmyconv, fileconv, othrconv *C.char
	var prmypack, filepack, othrpack chan *C.cr_Package
	var prmydone, filedone, othrdone chan bool
	var prmyover, fileover, othrover chan bool
	var prmyover_main, fileover_main, othrover_main bool
	var pack *C.cr_Package
	var numb int64

	prmypath = "/home/fedohide-origin/projects/metasource/rawhide/10beaa5fb8bb9b8710f4608ea9bf84aff2fb68e5efc7e82bf12b421867ad3d8f-primary.xml"
	filepath = "/home/fedohide-origin/projects/metasource/rawhide/4182e96bacb8bb0ccdcc9d446977416a2a18b49a4aed13d6c550be45a1bf061e-filelists.xml"
	othrpath = "/home/fedohide-origin/projects/metasource/rawhide/e3ac902af73897fe77cbc4df42d1c87d72ff4d69c9e792224d0ff3857f630e92-other.xml"
	prmyconv = C.CString(prmypath)
	fileconv = C.CString(filepath)
	othrconv = C.CString(othrpath)
	defer C.free(unsafe.Pointer(prmyconv))
	defer C.free(unsafe.Pointer(fileconv))
	defer C.free(unsafe.Pointer(othrconv))

	prmypack, filepack, othrpack = make(chan *C.cr_Package, 10), make(chan *C.cr_Package, 10), make(chan *C.cr_Package, 10)
	prmydone, filedone, othrdone = make(chan bool, 10), make(chan bool, 10), make(chan bool, 10)
	prmyover, fileover, othrover = make(chan bool), make(chan bool), make(chan bool)

	wait.Add(3)
	go PopulatePrmy(&wait, prmypack, prmydone, prmyover)
	go PopulateFile(&wait, filepack, filedone, fileover)
	go PopulateOthr(&wait, othrpack, othrdone, othrover)

	prmyover_main, _ = <-prmyover
	fileover_main, _ = <-fileover
	othrover_main, _ = <-othrover

	close(prmyover)
	close(fileover)
	close(othrover)

	if prmyover_main || fileover_main || othrover_main {
		expt = errors.New("Metadata databases already exist or opening failed")
		slog.Log(nil, slog.LevelError, fmt.Sprintf("%s", expt.Error()))
		return numb, nil
	}

	iter = C.cr_PkgIterator_new(prmyconv, fileconv, othrconv, nil, nil, nil, nil, &gexp)
	if iter == nil {
		expt = errors.New(fmt.Sprintf("%s", C.GoString(gexp.message)))
		slog.Log(nil, slog.LevelError, fmt.Sprintf("%s", expt.Error()))
		return numb, expt
	}
	defer C.cr_PkgIterator_free(iter, &gexp)

	for C.cr_PkgIterator_is_finished(iter) == 0 {
		var prmydone_main, filedone_main, othrdone_main bool

		pack = C.cr_PkgIterator_parse_next(iter, &gexp)
		if pack == nil {
			break
		}

		prmypack <- pack
		filepack <- pack
		othrpack <- pack

		prmydone_main, _ = <-prmydone
		filedone_main, _ = <-filedone
		othrdone_main, _ = <-othrdone

		if prmydone_main && filedone_main && othrdone_main {
			numb += 1
			head := fmt.Sprintf("%s %s:%s-%s.%s", C.GoString(pack.name), C.GoString(pack.epoch), C.GoString(pack.version), C.GoString(pack.release), C.GoString(pack.arch))
			slog.Log(nil, slog.LevelInfo, fmt.Sprintf("%s added.", head))
		}
	}
	slog.Log(nil, slog.LevelWarn, fmt.Sprintf("%d package(s) added.", numb))

	close(prmypack)
	close(filepack)
	close(othrpack)
	close(prmydone)
	close(filedone)
	close(othrdone)

	wait.Wait()

	defer C.free(unsafe.Pointer(pack))
	defer C.free(unsafe.Pointer(gexp))
	return numb, nil
}
