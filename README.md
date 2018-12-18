# TVT

Reserve tickets for TV Shows on [tvtickets](https://www.tvtickets.com)

## Install 

Install with `go get -u github.com/fedemengo/tvt`

## Run

The tools support a few simples options


```
$ tvt --help
usage: tvt [<flags>] <command> [<args> ...]

tvt - Reserve ticket for www.tvtickets.com

Flags:
      --help     Show context-sensitive help (also try --help-long and --help-man).
  -f, --force    Force reserving/creating a ticket
  -v, --verbose  Verbose output of what is happening

Commands:
  help [<command>...]
    Show help.

  ls
    List all available tv shows

  rs [<flags>]
    Reserve ticket
```

### Examples

Parameter from command line

`tvt rs --show SHOW_NAME --fist FIRST --last LAST ...`

Parameter from file

`tvt rs --config FILE_PATH`

The file configuration file is just a json file that contains all data necessary (the `config` file show the structure)

