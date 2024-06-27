---
weight: 502
title: Request Limit
description: ""
date: "2024-03-02T14:17:51+01:00"
lastmod: "2024-03-02T14:17:51+01:00"
draft: false
toc: true
---


## Max request body

By default, a request can send a maximum of 512KB. You can change this limit per rule as shown in the following
example.

```yaml
---
rules:
  - name: small limit
    on:
      path: /small
    rule.answer.content: Hello
    with:
      max_request_body: 2B

  - name: large limit
    on:
      path: /large
    rule.answer.content: Hello
    with:
      max_request_body: 10MB
```

On exceeding the limit, the request will be answered with `HTTP/1.1 413 Request Entity Too Large`