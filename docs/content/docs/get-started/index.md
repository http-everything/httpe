---
weight: 100
title: "Get Started"
description: "Download a pre-compiled version for your OS and run an example."
keywords: "foo bla"
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
{{% tab tabName="Linux & macOS" %}}

Download, unpack and start the server with the default example rules.

```shell
DOWNLOAD=httpe_{{% httpe-version %}}_$(uname -o|sed "s/GNU\///g")_$(uname -m).tar.gz
curl -LO https://github.com/http-everything/httpe/releases/download/{{% httpe-version %}}/${DOWNLOAD}
tar xzf $DOWNLOAD
./httpe -r example.rules.unix.yaml
```

## Execute the examples

On a second terminal, execute a request like this

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

Download, unpack and start the server with the default example rules.

```text
$download = "httpe_{{% httpe-version %}}_Windows_x86_64.zip"
iwr https://github.com/http-everything/httpe/releases/download/{{% httpe-version%}}/$download -OutFile $download
Expand-Archive $download
cd httpe_{{% httpe-version %}}_Windows_x86_64
.\httpe.exe -r .\example.rules.win.yaml
```

You might get a confirmation dialogue from Windows Defender Firewall because `httpe` wants to open a TCP port.
Click "Allow Access"

On a second terminal, execute a request like this
```text
(iwr http://localhost:3000/hello?Name=John).content
```

which will give you the following response:
```shell
Hello John
Have a lovely Sunday 😎
```

The weekday may vary.

{{% /tab %}}
{{< /tabs >}}

Run the other examples with `curl` or `Invoke-WebRequest`.

```shell
$ curl http://localhost:3000/date
Sun Feb 25 11:42:13 CET 2024

$ curl http://localhost:3000/hello-world
Hello World 

$ curl http://localhost:3000/hosts
##
# Host Database
#
# localhost is used to configure the loopback interface
# when the system is booting.  Do not change this entry.
##
127.0.0.1       localhost
255.255.255.255 broadcasthost
::1             localhost
```

Finally, point a browser to the URL `http://localhost:3000/button`. You will get a button called "Date".
Click on it. After the request has been executed, inside the green confirmation box, click on "Click to inspect" to see
how HTTPE has executed a script and send back the output.

{{< figure 
  src="buttons-date.png" 
  caption="A button created from the rules."
>}}

## Understand the rules

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
* A rule must have a name for identification. (Line 4)
* The `on` object defines the request matcher. In the shown example the rule takes action if the request goes to 
  the `hello`path. Because the `method` is not defined, this rules takes action on all request methods. (Lines 5-6)
* With the `do` object you define which action to execute if the `on` definition matches the request. The example 
  launches the `run.script` action. The script specified will be executed by the default shell. Stdout is returned as 
  http response. (Lines 8-10)
* `{{ .Input.Params.Name }}` is a template macro. HTTPE will replace it by the URL parameter `Name`
  before execution. (Line 9)

