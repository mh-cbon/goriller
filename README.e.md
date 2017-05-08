---
License: MIT
LicenseFile: LICENSE
LicenseColor: yellow
---
# {{.Name}}

{{template "badge/travis" .}} {{template "badge/appveyor" .}} {{template "badge/goreport" .}} {{template "badge/godoc" .}} {{template "license/shields" .}}

{{pkgdoc}}

# {{toc 5}}

# Install
{{template "glide/install" .}}

## Usage

#### $ {{exec "goriller" "-help" | color "sh"}}

## Cli examples

```sh
# Create a goriller binder of Tomate to MyTomate
goriller tomate_gen.go Tomate:MyTomate
```

# API example

Following example demonstates a program using it to generate a goriller binder of a type.

#### Anootations

`goriller` reads and interprets annotations on `struct` and `methods`.

The `struct` annotations are used as default for the `methods` annotations.

| Name | Description |
| --- | --- |
| @route | The route path such as `/{param}` |
| @name | The route name `name` |
| @host | The route name `host` |
| @methods | The route methods `GET,POST,PUT` |
| @schemes | The route methods `http, https` |

#### > {{cat "demo/main.go" | color "go"}}

Following code is the generated implementation of the goriller binder.

#### > {{cat "demo/goriller_vegetables_gen.go" | color "go"}}


# Recipes

#### Release the project

```sh
gump patch -d # check
gump patch # bump
```

# History

[CHANGELOG](CHANGELOG.md)
