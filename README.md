# caddy-supervisor [![Build Status](https://travis-ci.org/lucaslorentz/caddy-supervisor.svg?branch=master)](https://travis-ci.org/lucaslorentz/caddy-supervisor)

## Introduction
This plugin enables caddy to run and supervise background processes.

## How it works
For every **supervisor** caddy directive a command is executed in background and killed when caddy stops.

You can use **supervisor** plugin as an http directive or as a server type.

## Supervisor http directive
You can activate a supervisor directly from your web caddyfile using:
```
supervisor command arg1 arg2 arg3
```

Or using a block for more control
```
supervisor {
  command command
  args arg1 arg2 arg3
  dir directory
  env VARIABLEA VALUEA
  env VARIABLEB VALUEB
  redirect_stdout file
  redirect_stderr file
  restart_policy policy
  termination_grace_period period
  replicas replicas
}
```

## Supervisor server type
You can also use a supervisor server type using `-type` CLI option:
```
caddy -type supervisor
```

The Caddyfile syntax for supervisor server type is:
```
name {
  command command
  args arg1 arg2 arg3
  dir directory
  env VARIABLEA VALUEA
  env VARIABLEB VALUEB
  redirect_stdout file
  redirect_stderr file
  restart_policy policy
  termination_grace_period period
  replicas replicas
}
...
```

## Options description

- **command**: the command or executable name to be executed. Supports template.
- **args**: args provided to the command, separated by whitespace. Supports template.
- **dir**: the working directory the command should be executed in. Supports template.
- **env**: declare environment variable that should be passed to command. This property can be repeated. Supports template.
- **redirect_stdout**: redirect command stdout to a file. Use "stdout" to redirect to caddy stdout
- **redirect_stderr**: redirect command stderr to a file. Use "stderr" to redirect to caddy stderr
- **restart_policy**: define under which conditions the command should be restarted after exit. Valid values:
  - **never**: do not restart the command
  - **on_failure**: restart if exit code is not 0
  - **always**: always restart
- **termination_grace_period**: amount of time to wait for application graceful termination before killing it. Ex: 10s
- **replicas**: number of instances that should be executed. Default: 1.

On windows **termination_grace_period** is ignored and the command is killed immediatelly due to lack of signals support.

## Templates
To enable different configuration per replica, you can use go templates on the fields marked with Supports template".

The following information are available to templates:
- **Replica**: the index of the current replica, starting from 0

Templates also supports all functions from http://masterminds.github.io/sprig/

Example:
```
supervisor myapp --port "{{add 8000 .Replica}}" {
  replicas 5
}
proxy / localhost:8000-8004
```

## Exponential backoff
To avoid spending too many resources on a crashing application, this plugin makes use of exponential backoff.

That means that when the command fail, it will be restarted with a delay of 0 seconds. If it fails again it will be restarted with a delay of 1 seconds, then on every sucessive failure the delay time doubles, with a max limit of 5 minutes.

If the command runs stable for at least 10 minutes, the restart delay is reset to 0 seconds.

## Examples
AspNet Core application on windows:
```
example.com {
  run {
    env ASPNETCORE_URLS http://localhost:5000
    command dotnet ./MyApplication.dll
    dir "C:\MyApplicationFolder"
    redirect_stdout stdout
    redirect_stderr stderr
    restart_policy always
  }
  proxy / localhost:5000 {
    transparent
  }
}
```

Php fastcgi on windows:
```
example.com {
  run {
    command ./php-cgi.exe
    args -b 9800
    dir C:/php/
    redirect_stdout stdout
    redirect_stderr stderr
    restart_policy always
  }
  root C:/Site
  fastcgi / localhost:9800 php
}
```

## Building it
Build from caddy repository and import  **caddy-supervisor** plugin on file https://github.com/mholt/caddy/blob/master/caddy/caddymain/run.go :
```
import (
  _ "github.com/lucaslorentz/caddy-supervisor/httpplugin"
  _ "github.com/lucaslorentz/caddy-supervisor/servertype"
)
```
