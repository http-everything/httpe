---
title : "Welcome to HTTP Everything"
description: ""
lead: "Wrap everything into HTTP requests"
draft: false
seo:
 title: "" # custom title (optional)
 description: "Create powerful APIs with low code" # custom description (recommended)
 canonical: "" # custom canonical URL (optional)
 noindex: false # false (default) or true
---
## Build HTTP APIs with low code

HTTPE is a web server that maps URLs to actions.

```yaml {title=rules.yaml}
- rules:
  on:
    path: /date
  do:
    run.script: date
```
Start the server with `./httpe -r rules.yaml` and fire a requests with
`curl http://localhost:3000/date`.
The command will be executed and the output is sent back.

HTTPE has templating on board allowing you to customize almost every part of the request and the response.

```yaml {title=rules.yaml}
- rules:
  on:
    path: /dir/{dir}
  do:
    run.script: ls {{ .Input.URLPlaceholder.dir }}
  args:
    cwd: /home/john.doe
```

Actions are not limited to script execution. You can serve static and dynamic content, directories, and redirects.
More actions coming soon such as email sending, storing and retrieving data and much more.

HTTPE has TLS and authentication built-in.

{{< callout context="tip" title="Caution" icon="alert-triangle-filled" >}}
The project is in a very early stage. The development has just begun (Feb 2024).
{{< /callout >}}


HTTPE is open-source released under the MIT license.
