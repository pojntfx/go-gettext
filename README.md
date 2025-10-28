# go-gettext

Go `gettext` bindings based on `purego`.

![Go Version](https://img.shields.io/badge/go%20version-%3E=1.25-61CFDD.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/pojntfx/go-gettext.svg)](https://pkg.go.dev/github.com/pojntfx/go-gettext)

## Overview

go-gettext provides simple bindings to the [GNU`gettext`](https://en.wikipedia.org/wiki/Gettext) internationalization system based on [`purego`](https://github.com/ebitengine/purego). Thanks to `purego`, this means it requires neither a full reimplementation of `gettext` and the associated complexities, nor does it require CGo.

## Installation

You can add go-gettext to your Go project by running the following:

```shell
$ go get github.com/pojntfx/go-gettext/...@latest
```

Please note that a gettext library (usually named `libintl`) needs to be installed on your system. On Linux this is almost always the case, but on Windows you might want to ship the relevant DLL manually, while macOS requires that you install it with Homebrew or MacPorts. Only Linux is a tested platform at this time.

## Tutorial

> TL;DR: Extract strings, initialize the i18n system, then call `i18n.Local`

### 1. Extract Strings

Just like in any `gettext`-based project, you'll start by extracting strings from your source code. For Go, this works:

```shell
find .. -name '*.go' | xgettext --language=C++ --keyword=_ --keyword=i18n.Local --keyword=Local --omit-header -o default.pot --files-from=
```

The resulting `.pot` file can then be translated. The standard `gettext` toolchain can now be used; for a full example (including building and installing the `.mo` files), see [pojntfx/sessions/po](https://github.com/pojntfx/sessions/tree/main/po).

### 2. Setting up the Internalization System

Next, in an `init` function or elsewhere, import and setup go-gettext:

```go
import "github.com/pojntfx/go-gettext/pkg/i18n"

const (
	gettextPackage = "sessions"
	localeDir      = "/usr/share/locale"
)

func init() {
	if err := i18n.InitI18n(gettextPackage, localeDir); err != nil {
		panic(err)
	}
}
```

Adjust `gettextPackage` and `localeDir` to match your local environment. If you're using Meson, see [pojntfx/senbara/senbara-gtk/src/config.go.in](https://github.com/pojntfx/senbara/blob/981fb805eab9c91c56985c92c62dbf4835178c90/senbara-gtk/src/config.go.in) for an example of how to get those dynamically.

### 3. Getting a Localized String

Now that everything is set up, getting a localized string is as easy as calling `i18n.Local`:

```go
i18n.Local("Session finished")
```

The translated string for "Session finished" should be returned by `i18n.Local`, e.g. "Sitzung beendet" in German.

🚀 That's it! We hope go-gettext helps you with internationalizing your app.

## Acknowledgements

- [jwijenbergh/purego](https://github.com/jwijenbergh/purego) allows us to call functions from `gettext` without the need for CGo.
- [jwijenbergh/puregotk](https://github.com/jwijenbergh/puregotk) is what I usually use with this library, and was very helpful for learning how to use purego.
- [diamondburned/gotk4](https://github.com/diamondburned/gotk4) was the inspiration for how the `InitI18n` function should work.
- [GNU gettext](https://en.wikipedia.org/wiki/Gettext) is the most commonly used implementation of gettext and what go-gettext is usually used with.

## Contributing

To contribute, please use the [GitHub flow](https://guides.github.com/introduction/flow/) and follow our [Code of Conduct](./CODE_OF_CONDUCT.md).

## License

go-gettext (c) 2025 Felicitas Pojtinger and contributors

SPDX-License-Identifier: Apache-2.0
