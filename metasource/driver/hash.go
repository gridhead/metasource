package driver

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"log/slog"
	"metasource/metasource/models/home"
	"os"
	"sync"
)

func VerifyChecksum(unit *home.FileUnit, vers *string, wait *sync.WaitGroup, cast *int) {
	defer wait.Done()

	var file *os.File
	var read hash.Hash
	var expt error
	var buff []byte
	var csum string

	file, expt = os.Open(unit.Path)
	if expt != nil {
		unit.Name = ""
		unit.Path = ""
		slog.Log(nil, slog.LevelDebug, fmt.Sprintf("[%s] Checksum mismatch for %s due to %s", *vers, unit.Name, expt.Error()))
		return
	}
	defer file.Close()

	if unit.Hash.Type != "sha256" {
		slog.Log(nil, slog.LevelDebug, fmt.Sprintf("[%s] Checksum mismatch for %s due to unknown checksum type", *vers, unit.Name))
		unit.Name = ""
		unit.Path = ""
		return
	}

	read = sha256.New()
	buff = make([]byte, 4*1024)
	for {
		var indx int
		indx, expt = file.Read(buff)
		if indx > 0 {
			read.Write(buff[:indx])
		}
		if expt == io.EOF {
			break
		}
		if expt != nil {
			unit.Name = ""
			unit.Path = ""
			slog.Log(nil, slog.LevelDebug, fmt.Sprintf("[%s] Checksum mismatch for %s due to %s", *vers, unit.Name, expt.Error()))
			return
		}
	}

	csum = fmt.Sprintf("%x", read.Sum(nil))
	if csum != unit.Hash.Data {
		slog.Log(nil, slog.LevelDebug, fmt.Sprintf("[%s] Checksum mismatch for %s", *vers, unit.Name))
		unit.Name = ""
		unit.Path = ""
		return
	}

	slog.Log(nil, slog.LevelDebug, fmt.Sprintf("[%s] Checksum verified for %s", *vers, unit.Name))
	*cast++
	return
}
