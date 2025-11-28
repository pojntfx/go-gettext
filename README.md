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

> TL;DR: Extract strings, initialize the i18n system with `InitI18n` (for applications) or `BindI18n` (for libraries), then call `i18n.Local` or `i18n.LocalDomain`

### 1. Extract Strings

Just like in any `gettext`-based project, you'll start by extracting strings from your source code. For Go, this works:

```shell
find . -name '*.go' | xgettext --language=C++ --keyword=_ --keyword=i18n.Local --keyword=Local --keyword=i18n.LocalDomain:2 --keyword=LocalDomain:2 --omit-header -o default.pot --files-from=
```

Alternatively, if you're using the `L` and `LD` shorthands instead of `i18n.Local` and `i18n.LocalDomain`:

```shell
find . -name '*.go' | xgettext --language=C++ --keyword=_ --keyword=L --keyword=LD:2 --omit-header -o default.pot --files-from=
```

The resulting `.pot` file can then be translated. The standard `gettext` toolchain can now be used; for a full example (including building and installing the `.mo` files), see [pojntfx/sessions/po](https://github.com/pojntfx/sessions/tree/main/po).

### 2. Setting up the Internalization System

#### For Applications

If you're building an application, use `InitI18n` to initialize and set the global text domain:

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

#### For Libraries

If you're building a library that needs its own translation domain without changing the global text domain, use `BindI18n` instead:

```go
import "github.com/pojntfx/go-gettext/pkg/i18n"

const (
	gettextPackage = "mylib"
	localeDir      = "/usr/share/locale"
)

func init() {
	if err := i18n.BindI18n(gettextPackage, localeDir); err != nil {
		panic(err)
	}
}
```

When using `BindI18n`, you'll need to use `i18n.LocalDomain` (or `LD`) instead of `i18n.Local` (or `L`) to look up strings in your library's specific domain (see [3. Getting a Localized String](#3-getting-a-localized-string)).

#### Configuration

Adjust `gettextPackage` and `localeDir` to match your local environment. If you're using Meson, see [pojntfx/senbara/senbara-gtk/src/config.go.in](https://github.com/pojntfx/senbara/blob/981fb805eab9c91c56985c92c62dbf4835178c90/senbara-gtk/src/config.go.in) for an example of how to get those dynamically. Since `go-gettext` uses the system `gettext` library, using `go:embed` is a bit harder than usual; one (somewhat hacky) solution is to embed the generated `.mo` files and extract them to a temporary directory at runtime like this:

<details>
  <summary>Expand section</summary>

```go
//go:embed *
var FS embed.FS

// ...

i18t, err := os.MkdirTemp("", "")
if err != nil {
	panic(err)
}
defer os.RemoveAll(i18t)

if err := fs.WalkDir(po.FS, ".", func(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if d.IsDir() {
		if err := os.MkdirAll(filepath.Join(i18t, path), os.ModePerm); err != nil {
			return err
		}

		return nil
	}

	src, err := po.FS.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filepath.Join(i18t, path))
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}); err != nil {
	panic(err)
}

i18n.InitI18n("default", i18t)
```

If you're looking for a pure Go library that has support for `go:embed` out of the box, we recommend [leonelquinteros/gotext](https://github.com/leonelquinteros/gotext).

</details>

### 3. Getting a Localized String

#### For Applications (using InitI18n)

Now that everything is set up, getting a localized string is as easy as calling `i18n.Local`:

```go
i18n.Local("Session finished")
```

Alternatively you can also use the `L` shorthand like so:

```go
import . "github.com/pojntfx/go-gettext/pkg/i18n"

L("Session finished")
```

The translated string for "Session finished" should be returned by `i18n.Local` or `L`, e.g. "Sitzung beendet" in German.

#### For Libraries (using BindI18n)

If you're using `BindI18n` in a library, use `i18n.LocalDomain` to look up strings in your specific text domain:

```go
i18n.LocalDomain("mylib", "Operation completed")
```

Alternatively you can also use the `LD` shorthand like so:

```go
import . "github.com/pojntfx/go-gettext/pkg/i18n"

LD("mylib", "Operation completed")
```

This allows libraries to have their own translation domain without interfering with the application's global text domain.

🚀 That's it! We hope go-gettext helps you with internationalizing your app.

## Acknowledgements

- [jwijenbergh/purego](https://github.com/jwijenbergh/purego) allows us to call functions from `gettext` without the need for CGo.
- [jwijenbergh/puregotk](https://github.com/jwijenbergh/puregotk) is what is commonly used with this library, and was very helpful for learning how to use purego.
- [diamondburned/gotk4](https://github.com/diamondburned/gotk4) was the inspiration for how the `InitI18n` function should work.
- [GNU gettext](https://en.wikipedia.org/wiki/Gettext) is the most commonly used implementation of gettext and what go-gettext is usually used with.
- [leonelquinteros/gotext](https://github.com/leonelquinteros/gotext) is a great, pure Go gettext reimplementation.

## Contributing

To contribute, please use the [GitHub flow](https://guides.github.com/introduction/flow/) and follow our [Code of Conduct](./CODE_OF_CONDUCT.md).

## License

go-gettext (c) 2025 Felicitas Pojtinger and contributors

SPDX-License-Identifier: Apache-2.0
