package config

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
)

func SetLogger(lglvtext *string) {
	var lglvoptn slog.Level

	switch *lglvtext {
	case "info":
		lglvoptn = slog.LevelInfo
	case "warn":
		lglvoptn = slog.LevelWarn
	case "debug":
		lglvoptn = slog.LevelDebug
	default:
		lglvoptn = slog.LevelInfo
	}

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level: lglvoptn,
	}))
	slog.SetDefault(logger)
}

var BODHIURL string = "https://bodhi.fedoraproject.org"
var KOJIREPO string = "https://kojipkgs.fedoraproject.org/repos"
var DLSERVER string = "https://dl.fedoraproject.org"

var LINKDICT = map[string][]string{
	"Fedora Linux": {
		"%s/pub/fedora/linux/releases/%s/Everything/x86_64/os/repodata/",
		"%s/pub/fedora/linux/updates/%s/Everything/x86_64/repodata/",
		"%s/pub/fedora/linux/updates/testing/%s/Everything/x86_64/repodata/",
	},
	"Fedora EPEL": {
		"%s/pub/epel/%s/Everything/x86_64/repodata/",
		"%s/pub/epel/testing/%s/Everything/x86_64/repodata/",
	},
	"Fedora EPEL Next": {
		"%s/pub/epel/next/%s/Everything/x86_64/repodata/",
		"%s/pub/epel/next/testing/%s/Everything/x86_64/repodata/",
	},
}

var REPODICT = map[string][]string{
	"fedora":    {"%s", "%s-updates", "%s-updates-testing"},
	"epel":      {"%s", "%s-testing"},
	"epel-next": {"%s", "%s-testing"},
}

var DBFOLDER string = "/var/tmp/metasource"

var WAITTIME int = 10

var FILENAME string = "metasource-%s-%s"

var ATTEMPTS int64 = 4
