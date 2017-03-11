# rclone-json
[![travis Status](https://travis-ci.org/mh-cbon/rclone-json.svg?branch=master)](https://travis-ci.org/mh-cbon/rclone-json)[![appveyor Status](https://ci.appveyor.com/api/projects/status/github/mh-cbon/rclone-json?branch=master&svg=true)](https://ci.appveyor.com/project/mh-cbon/rclone-json)
[![GoDoc](https://godoc.org/github.com/mh-cbon/rclone-json?status.svg)](http://godoc.org/github.com/mh-cbon/rclone-json)

Package rclone-json streams an rclone sync activity as a json object stream.


## Install

### Go

```sh
go get github.com/mh-cbon/rclone-json
```


## Usage


###### $ rclone-json -help
```sh
rclone-json - 0.0.0
Usage of rclone-json:
  -bwlimit string
    	
  -checkers string
    	
  -help
    	Show help
  -rclone string
    	 (default "rclone")
  -stats string
    	
  -transfers string
    	
  -version
    	Show version
```

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
