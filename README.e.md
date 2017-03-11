# {{.Name}}
{{template "badge/travis" .}}{{template "badge/appveyor" .}}{{template "badge/godoc" .}}
{{pkgdoc}}

## Install
{{template "gh/releases" .}}

### Go
{{template "go/install" .}}

### Chocolatey
{{template "choco/install" .}}

### linux rpm/deb repository
{{template "linux/gh_src_repo" .}}

### linux rpm/deb standalone package
{{template "linux/gh_pkg" .}}

## Usage

{{cli "rclone-json" "-help"}}

### Cli examples

```sh
rclone-json stats src/ dst/
```

# Recipes

### Release the project

```sh
gump patch -d # check
gump patch # bump
```
