# DDD-GO

Golang DDD 实践框架，可以快速创建DDD项目以及实体相关代码，使开发更简单，更快捷。

## Usage

```bash
$ ddd-go -h
ddd-go is a library for golang DDD development. 
It provides some useful tools for DDD development. 
It could also help users generating entity, etc.

Usage:
  ddd-go [command]

Available Commands:
  entity      Generate entity related files
  help        Help about any command

Flags:
  -h, --help   help for ddd-go

Use "ddd-go [command] --help" for more information about a command.
```

创建实体

```bash
$ ddd-go entity -h
Generate entity related files

Usage:
  ddd-go entity [flags]

Flags:
  -d, --directory string   project directory, default is '.' (default ".")
  -h, --help               help for entity
  -l, --layout string      project layout, default is 'kratos' (default "kratos")
  -m, --module string      project module name
  -n, --name string        entity name
  -s, --short string       entity short name

ddd-go entity -m myproj -d ../myproj -e User -s usr
```
