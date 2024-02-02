---
title: "Quick Start"
description: "Install and run HTTPE in no time"
summary: ""
draft: false
menu:
  docs:
    parent: ""
    identifier: "get-started-c177b0c56f8225d55374a8ca1ac5fd9f2893"
weight: 1
toc: true
seo:
  title: "" # custom title (optional)
  description: "Learn how to install and run the HTTPE server" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

## Install
Currently, you must install from the sources. Very shorty we will provide binaries and RPM/DEB packages.
```bash
git clone https://github.com/http-everything/httpe/
cd httpe
go build -o httpe main.go
./httpe --help
```

## First Test
Start the server with included example rules by executing

    ./httpe -r example.rules.yaml

and fire a request like that

    curl http://localhost:3000/hello-world

which will return `hello world`.

Open the `example.rules.yaml` file with an editor or [read it online](https://github.com/http-everything/httpe/blob/main/example.rules.yaml)
to get a first impression of what HTTPE is able to do.

