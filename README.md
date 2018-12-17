# TVT

Reserve tickets for TV Shows on [tvtickets](https://www.tvtickets.com)

## Install 

Install with `go get -u github.com/fedemengo/tvt`

## Run

The tools support a few simples options


```
$ tvt

tvt - Reserve ticket for www.tvtickets.com
Usage:	 tvt option

Options:	ls - List all available tv shows
        	rs - Reserve ticket
```

### Examples

Parameter from command line

`tvt rs --show SHOW_NAME --fist FIRST --last LAST ...`

Parameter from file

`tvt rs --file FILE_PATH`

Structure of `FILE_PATH`

```
SHOW_NAME1
- FIRST LAST N PHONE EMAIL
- FIRST LAST N PHONE EMAIL
SHOW_NAME2
- FIRST LAST N PHONE EMAIL
....
```
