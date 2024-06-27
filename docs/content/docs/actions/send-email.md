---
weight: 307
title: "Send email"
description: ""
icon: "article"
date: "2024-06-02T13:08:55+01:00"
lastmod: "2024-06-02T13:08:55+01:00"
draft: false
toc: true
---

## Preface
The `send.email` action allows sending email on incoming http requests. An SMTP server and a `[smtp]` configuration
block in the main `httpe.conf` file is required.

## Server configuration

Add a block like the below to your server configuration file. The authentication parameters are optional.
Whether your SMTP server requires TLS or not is automatically detected.

The email-from can be overwritten on a per-message basis.

```toml
[smtp]
## Defines smtp server, optional.
## Required if you want to use the send.email action
server = "smtp.example.com"
port = 587

## An optional pair of username and password, if server requires authentication
#username = ""
#password = ""

## An optional from address, from address specified by the rule will have precedence
from = "info@example.com"
```

## Rules examples

### Example

```yaml
---
rules:
  - name: Send email
    on:
      path: /email
    send.email:
      to: recipient@example.com
      from: sender@example.com
      body: "{{ .Input.Form.text }}"
      subject: "{{ .Input.Form.subject }}"
```

The above email can be triggered with a request like

```shell
curl localhost:3000/email -F text="Have a nice day" -F subject="Nice $(date)"
```

Of course, you can take the `to` and `from`, and any other field, from the request data using the 
[built-in templating](/docs/templating).

If the SMTP server doesn't return errors, HTTP status 200 and a message 'email sent' is returned.

Here is another example in which the email recipient is retrieved from the URL.

```yaml
---
rules:
  - name: Send email
    on:
      path: /email/{to}
      methods:
        - post
    send.email:
      to: "{{ .Input.URLPlaceholders.to }}"
      from: registration@example.com
      subject: Your registration key
      body: |
        Welcome {{ .Input.Form.name }},
        Thanks for joining us.
        Your registration key is {{ .Input.Form.key }}.
        Have a nice day.
    respond:
      on_success:
        body: ""
        http_status: 204
```

The below example would trigger the above email action.

```shell
curl localhost:3000/email/user@example.com -F name=Thorsten -F key=4711
```
