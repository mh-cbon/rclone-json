# {{.Name}}
{{template "badge/travis" .}}{{template "badge/appveyor" .}}{{template "badge/godoc" .}}
{{pkgdoc}}

## Install

### Go
{{template "go/install" .}}

## API example

{{file "main_example.go"}}

## Usage

{{cli "rclone-json" "-help"}}

### Cli examples

```sh
rclone-json sync src/ dst/
```

# Recipes

### Release the project

```sh
gump patch -d # check
gump patch # bump
```
