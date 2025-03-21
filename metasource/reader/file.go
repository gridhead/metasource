package reader

/*
#cgo pkg-config: createrepo_c
#include "createrepo_c/sqlite.h"
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

func PopulateFile(wait *sync.WaitGroup, dbpk <-chan *C.cr_Package, done chan<- bool, over chan<- bool) {
	defer wait.Done()

	var path string
	var conv *C.char
	var base *C.cr_SqliteDb
	var pack *C.cr_Package
	var rslt int
	var gexp *C.GError
	var expt error

	path = "akashdeeep-filelist.sqlite"
	conv = C.CString(path)
	defer C.free(unsafe.Pointer(conv))

	base = C.cr_db_open(conv, C.CR_DB_FILELISTS, &gexp)
	defer C.cr_db_close(base, &gexp)
	if gexp != nil {
		expt = errors.New(fmt.Sprintf("%s", C.GoString(gexp.message)))
		slog.Log(nil, slog.LevelError, fmt.Sprintf("%s", expt.Error()))
		over <- true
	} else {
		over <- false
	}

	for pack = range dbpk {
		rslt = int(C.cr_db_add_pkg(base, pack, &gexp))
		if rslt != 0 {
			expt = errors.New(fmt.Sprintf("%s", C.GoString(gexp.message)))
			slog.Log(nil, slog.LevelWarn, fmt.Sprintf("%s", expt.Error()))
			done <- false
		} else {
			done <- true
		}
	}
}
