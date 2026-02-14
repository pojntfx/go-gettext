// Inspired by https://github.com/diamondburned/gotk4/blob/d44ab4b5b24e200c90c6a9ffb6632ada7e166a79/pkg/core/glib/glib.go#L50-L72

package i18n

import (
	"errors"
	"log/slog"
)

var (
	setlocale             func(category int, locale string) string
	bindtextdomain        func(domainname string, dirname string) string
	bindTextdomainCodeset func(domainname string, codeset string) string
	textdomain            func(domainname string) string
	gettext               func(msgid string) string
	dgettext              func(domainname string, msgid string) string
)

// InitI18n initializes the i18n subsystem and sets the text domain. It runs the following C code:
//
//	setlocale(LC_ALL, "");
//	bindtextdomain(domain, dir);
//	bind_textdomain_codeset(domain, "UTF-8");
//	textdomain(domain);
//
// Use this for your main application. For libraries that need their own translation
// domain without changing the global text domain, use BindI18n instead.
func InitI18n(domain, dir string, logger *slog.Logger) error {
	if err := registerLibrary(); err != nil {
		return errors.Join(errors.New("could register gettext library"), err)
	}

	lcAll, err := getLCALL()
	if err != nil {
		return errors.Join(errors.New("could get LC_ALL value"), err)
	}

	if setlocale(lcAll, "") == "" {
		logger.Debug("failed to set locale, verify that system locale for this LC_ALL value is installed; ignoring and continuing", "LC_ALL", lcAll)
	}

	if bindtextdomain(domain, dir) == "" {
		return errors.New("failed to bind text domain")
	}

	if bindTextdomainCodeset(domain, "UTF-8") == "" {
		return errors.New("failed to set text domain codeset")
	}

	if textdomain(domain) == "" {
		return errors.New("failed to set text domain")
	}

	return nil
}

// BindI18n binds a text domain without setting it as the current domain. It runs the following C code:
//
//	setlocale(LC_ALL, "");
//	bindtextdomain(domain, dir);
//	bind_textdomain_codeset(domain, "UTF-8");
//
// This is useful for libraries that need their own translation domain. The library
// should use LocalDomain or the LD alias to look up strings in its specific domain.
// This does NOT change the global text domain used by Local/L.
func BindI18n(domain, dir string, logger *slog.Logger) error {
	if err := registerLibrary(); err != nil {
		return err
	}

	lcAll, err := getLCALL()
	if err != nil {
		return errors.Join(errors.New("could get LC_ALL value"), err)
	}

	if setlocale(lcAll, "") == "" {
		logger.Debug("failed to set locale, verify that system locale for this LC_ALL value is installed; ignoring and continuing", "LC_ALL", lcAll)
	}

	if bindtextdomain(domain, dir) == "" {
		return errors.New("failed to bind text domain")
	}

	if bindTextdomainCodeset(domain, "UTF-8") == "" {
		return errors.New("failed to set text domain codeset")
	}

	return nil
}

// Local localizes a string using gettext with the current text domain.
func Local(input string) string {
	return gettext(input)
}

// L is a shorthand for Local
var L = Local

// LocalDomain localizes a string using dgettext with a specific text domain.
// Use this in libraries that have their own translation domain.
func LocalDomain(domain, input string) string {
	return dgettext(domain, input)
}

// LD is a shorthand for LocalDomain
var LD = LocalDomain
