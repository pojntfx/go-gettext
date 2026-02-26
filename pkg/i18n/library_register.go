// Inspired by https://github.com/diamondburned/gotk4/blob/d44ab4b5b24e200c90c6a9ffb6632ada7e166a79/pkg/core/glib/glib.go#L50-L72

package i18n

import (
	"errors"

	"codeberg.org/puregotk/purego"
)

var (
	libc uintptr
)

func registerLibrary() error {
	if libc != 0 {
		return nil
	}

	gettextLibNames, err := getGettextLibraryNames()
	if err != nil {
		return errors.Join(errors.New("could get gettext library names"), err)
	}

	for _, gettextLibName := range gettextLibNames {
		libc, err = openLibrary(gettextLibName)
		if err != nil {
			return errors.Join(errors.New("could not open gettext library"), err)
		} else {
			break
		}
	}

	purego.RegisterLibFunc(&setlocale, libc, "setlocale")
	purego.RegisterLibFunc(&bindtextdomain, libc, "bindtextdomain")
	purego.RegisterLibFunc(&bindTextdomainCodeset, libc, "bind_textdomain_codeset")
	purego.RegisterLibFunc(&textdomain, libc, "textdomain")
	purego.RegisterLibFunc(&gettext, libc, "gettext")
	purego.RegisterLibFunc(&dgettext, libc, "dgettext")

	return nil
}
