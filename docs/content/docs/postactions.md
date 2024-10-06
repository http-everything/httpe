---
weight: 500
title: "Running asynchronous post actions"
description: ""
icon: "article"
date: "2024-08-18T15:09:30+01:00"
lastmod: "2024-08-18T15:09:30+01:00"
draft: false
toc: true
---

## Preface

Since version 0.0.6 you can perform additional actions after the http request has been answered.
This is called a post (http request) action.
These actions continue to run in the background. The output of these actions is written to a JSON
file in the data directory.

Each execution of a post action triggers a cleanup of the data directory. Files older than the specified
retention are deleted.

Supported post actions are:
1. `run.script` to execute a script
2. `send.email` to send a message

Multiple post actions are executed in the same routine. Scripts are run first and messages sent last.
The results of all post actions are stored in the same file.

## Example

```yaml
---
rules:
  - name: Single Line
    on:
      path: /content/1
      methods:
        - post
    answer.content: Hello World
    postaction:
      run.script: |
        date
        echo {{ .Input.Form.text }}
      send.email:
        subject: Your configuration has changed to {{ .Input.Form.text }}
        to: user@example.com
        from: configalerts@example.com
        body: |
          This is a confirmation about your latest config changes.
          From now on the following configuration is active:
          
          {{ .Input.Form.text }}

```

The example above responds to the http request with some static text. Once the request has been fully answered 
and the connection is closed, two actions are performed in the background.

Remember: Scripts are executed first, then emails are sent. This order cannot be changed.

Running the above rule via an HTTP request like `curl http://localhost:3000/content/1 -F text="some new values"`,
will create a file like `httpe-postrun-2024-08-18-10-29-19-736771776767.json`
in the data directory with a content like this:

```json
[
  {
    "action_type": "run.script",
    "success_body": "Sun Aug 18 12:20:03 CEST 2024\nsome new values\n",
    "error_body": "",
    "code": 0,
    "internal_error": ""
  },
  {
    "action_type": "send.email",
    "success_body": "email sent",
    "error_body": "",
    "code": 0,
    "internal_error": ""
  }
]
```
