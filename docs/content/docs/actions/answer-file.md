---
weight: 302
title: "Answer File"
description: ""
icon: "article"
date: "2024-02-27T15:19:59+01:00"
lastmod: "2024-02-27T15:19:59+01:00"
draft: false
toc: true
---

HTTPE can serve files from disk. Optionally these files are send through the template engine to do on-the-fly manipulation.

## Example

```yaml
---
rules:
  - name: GET File 1
    on:
      path: /file/1
      methods:
        - get
    do:
      answer.file: /etc/hosts

  - name: GET File 2
    on:
      path: /file/2
      methods:
        - get
    do:
      answer.file: /tmp/test.txt
      args:
        template: true
    respond:
      on_success:
        headers:
          My-Header: Super Dupa
      on_error:
        headers:
          My-Header: THIS IS AN ERROR

  - name: Music
    on:
      path: /file/3
      methods:
        - get
    do:
      answer.file: /tmp/music.mp3
      args:
        templating: false
    respond:
      on_success:
        headers:
          Content-Type: audio/mpeg
          Content-Disposition: filename="music.mp3"
```

To see it in action, create a file `/tmp/test.txt` with the following content:

```text
Hello {{ .Input.Params.Name }}!
How are you today?
```

Then fire some requests.
```shell
$ curl localhost:3000/file/1
127.0.0.1       localhost
255.255.255.255 broadcasthost
::1             localhost

$ curl localhost:3000/file/2?Name=John
Hello John!
How are you today?
```

When serving files from disk, templating is turned off by default. Read more about [templating](/docs/templating).

{{% alert context="danger" %}}
All files are read into memory before sending the response even if templating is tuned off.  
Consider using [serve.directory](/docs/actions/serve-directory) instead which implements a static file server that serves
files from disk.
{{% /alert %}}

## File Paths

On Windows you can you both, a backslash or a forward slash to specify a file path. The below examples are both valid.
```text
  - name: Hosts
    on:
       path: /hosts
    do:
      #answer.file: C:\Windows\System32\drivers\etc\hosts
      answer.file: C:/Windows/System32/drivers/etc/hosts
```

Paths containing blank spaces doesn't require escaping or quoting. Just type it in with blank space such as
```text
answer.file: C:/Users/thorsten/httpe_0.0.1_Windows_x86_64/an important folder/file.txt
```
