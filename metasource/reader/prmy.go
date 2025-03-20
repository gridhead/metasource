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
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"unsafe"
)

func PopulatePrmy(wait *sync.WaitGroup, prmyconv **C.char, fileconv **C.char, othrconv **C.char) {
	defer wait.Done()

	var path string
	var conv *C.char
	var base *C.cr_SqliteDb
	var pack *C.cr_Package
	var rslt int
	var gexp *C.GError
	var expt error
	var iter *C.cr_PkgIterator
	// var head string

	iter = C.cr_PkgIterator_new(*prmyconv, *fileconv, *othrconv, nil, nil, nil, nil, &gexp)
	if iter == nil {
		expt = errors.New(fmt.Sprintf("%s", C.GoString(gexp.message)))
		slog.Log(nil, slog.LevelError, fmt.Sprintf("Failed to create package iterator. %s", expt.Error()))
		return
	}
	defer C.cr_PkgIterator_free(iter, &gexp)

	path = "akashdeeep-primary.sqlite"
	conv = C.CString(path)
	defer C.free(unsafe.Pointer(conv))
	base = C.cr_db_open(conv, C.CR_DB_PRIMARY, &gexp)
	defer C.cr_db_close(base, &gexp)
	if gexp != nil {
		expt = errors.New(fmt.Sprintf("%s", C.GoString(gexp.message)))
		slog.Log(nil, slog.LevelError, fmt.Sprintf("%s", expt.Error()))
		return
	}

	for C.cr_PkgIterator_is_finished(iter) == 0 {
		pack = C.cr_PkgIterator_parse_next(iter, &gexp)
		if pack == nil {
			break
		}
		// head = fmt.Sprintf("%s %s-%s.%s", C.GoString(pack.name), C.GoString(pack.version), C.GoString(pack.epoch), C.GoString(pack.release))
		// slog.Log(nil, slog.LevelInfo, fmt.Sprintf("%s added.", head))
		rslt = int(C.cr_db_add_pkg(base, pack, &gexp))
		if rslt != 0 {
			expt = errors.New(fmt.Sprintf("%s", C.GoString(gexp.message)))
			slog.Log(nil, slog.LevelWarn, fmt.Sprintf("Failed to add package. %s", expt.Error()))
		}
	}

	defer C.free(unsafe.Pointer(pack))
	defer C.free(unsafe.Pointer(gexp))
}
