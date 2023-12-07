# seraphim
CLI ToolBelt

## Resources
- GoLang
 - Cobra (for commands and subcommands)
 - Viper (for configuration)
 - BubbleTea (for TUI)
 - <X> (for SQL)

## Usefull commands
Create new command using cobra
```bash
$ cobra-cli add <command name>
```
Create new subcommand
```bash
$ cobra-cli <subcommand name> -p '<parent command name in lowercase>Cmd'
$ #e.g. cobra-cli add sub -p 'parentCmd'
```

## Basic structure
1. Commands are defined using Cobra
2. Configuration is handled using Viper
3. Command operations are handled using BubbleTea to make use of the TUI framework capabilities 

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/IMNOTIKE/seraphim.git)