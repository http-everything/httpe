---
weight: 301
title: "Answer Content"
description: ""
icon: "article"
date: "2024-02-12T13:09:34+01:00"
lastmod: "2024-02-12T13:09:34+01:00"
draft: false
toc: true
---

## Example

```yaml
---
rules:
  - name: Single Line
    on:
      path: /content/1
      methods:
        - get
    do:
      answer.content: Hello World

  - name: Multi-Line
    on:
      path: /content/2
      methods:
        - post
    do:
      answer.content: |
        Line 1
        Line 2
        {{ .Input.Form.Text }}
        
```

To see it in action, execute requests as follows:

```shell
$ curl localhost:3000/content/1
Hello World

$ curl localhost:3000/content/2 -F Text="Hello World"
Line 1
Line 2
Hello World
```

Read more about how to use [placeholders](/docs/templating/).