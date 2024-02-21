---
weight: 001
title: "Run scripts"
description: ""
icon: "article"
date: "2024-02-12T13:08:55+01:00"
lastmod: "2024-02-12T13:08:55+01:00"
draft: false
toc: true
---

## Examples

### Simple

```yaml
rules:
  - name: Execute some commands
    on:
      path: /commands/1
      methods:
        - get
        - post
    do:
      run.script: |
        date
        id
        whoami
      args:
        interpreter: /bin/bash
        timeout: 2
```

### Advanced

```yaml
rules:
  - name: Execute some more commands
    on:
      path: /commands/2
    do:
      run.script: "{{ .Input.Form.Script }}"
      args:
        timeout: 3
    respond:
      on_error:
        body: |
          Your script '{{ .Input.Form.Script }}' has failed with: 
            Stderr: {{ .Action.ErrorBody }}
            Stdout: {{ .Action.SuccessBody }}
          Exit Code: {{ .Action.Code }}
      on_success:
        body: |
          Your script '{{ .Input.Form.Script }}' returns:
          {{ .Action.SuccessBody }}
          {{ .Action.ErrorBody }}

  - name: execute python
    on:
      path: /commands/3
    do:
      run.script: |
        import platform
        print(platform.python_version())
      args:
        interpreter: python3
```

## Supported Args

### `interpreter`

Define which interpreter shall execute the script.  
Defaults:
* On Linux, system default shell, usually `sh`.  
* On Windows: `powershell` taken from the path. Usually this resolves to PowerShell 5.
  Interpreters ending in `powershell` or `pwsh` are started with the options `-NoProfile -NonInteractive`.

### `timeout`

A timeout (integer, seconds) after which the script and eventually spawned child processes will be killed.

Default: `30` seconds.

{{% alert context="warning" %}}
`timeout` is currently not fully supported on Windows. After the timeout is exceeded, the httpe server process will 
decouple from the script and return a timeout exceeded error, but **the script continues running**. 
{{% /alert %}}

