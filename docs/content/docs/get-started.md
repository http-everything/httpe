---
weight: 001
title: "Get Started"
description: ""
icon: "article"
date: "2024-02-12T13:06:04+01:00"
lastmod: "2024-02-12T13:06:04+01:00"
draft: false
slug: "get-started"
toc: true
---

## Download and run

Just download a pre-compiled version for your OS.

{{< tabs tabTotal="2">}}
{{% tab tabName="Linux & MacOS" %}}

Download and unpack
```shell
curl -LO https://github.com/http-everything/httpe/releases/{{% httpe-version %}}/http-Linux-{{% httpe-version %}}tar.gz
```

Once downloaded and unpacked, start the server with the default example rule.
```shell
./httpe -r rules.yaml
```

Execute a request like this
```shell
curl localhost:3000/hello?Name=John
```

which will give you the following response:

```shell
Hello John.
Have a lovely Monday! 😎
```

The weekday may vary ;-).

{{% /tab %}}
{{% tab tabName="Windows" %}}

Example content specific to **Mac** operating systems

{{% /tab %}}
{{< /tabs >}}

Let's look at the rule file.

```yaml {linenos=true}
---
## Rules for the HTTP application server
rules:
  - name: Hello World
    on:
      path: /hello
    do:
      run.script: |
        echo "Hello {{ .Input.Params.Name }}.
        Have a lovely $(date +%A)! 😎"

```

* All rules must be child items of the `rules` object. (Line 3)
* Rules must be defined as a list, hence each rule must start with a dash. (Line 4)
* A rule can have a name for better identification. (Line 4)
* The `on` object defines the request matcher. In the shown example the rule takes action if the request goes to 
  the `hello`path. Because the `method` is not defined, this rules takes action on all request methods. (Lines 5-6)
* With the `do` object you define which action to execute if the `on` definition matches the request. The example 
  launches the `run.script` action. The script specified will be executed by the default shell. Stdout is returned as 
  http response. (Lines 8-10)
* `{{ .Input.Params.Name }}` is a template macro. HTTPE will replace it by the URL parameter `Name`
  before execution. (Line 9)

