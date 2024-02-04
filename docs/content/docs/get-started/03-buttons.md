---
title: "Buttons"
description: "Trigger actions with the built-in buttons UI"
summary: ""
draft: false
menu:
  docs:
    parent: ""
    identifier: "03-buttons-a4580ccf2a8774c4340f4d902761f78b"
weight: 03
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---
## Trigger actions with the built-in buttons UI

HTTPE comes with the option to create a responsive UI with buttons to execute requests to the URLs of other
actions.

```yaml
---
rules:
  - name: Wake up the PCs
    on:
      path: /wakeup/{mac}
    do:
      run.script: wakeonlan {{ .Input.URLPlaceholders.mac }}
  - name: Ping the PC
    on:
      path: /ping/{ip}
    do:
      run.script: ping -c 2 -t 2 {{ .Input.URLPlaceholders.ip }} 2>&1
    respond:
      on_error:
        body: "{{.Action.SuccessBody }}"
  - name: My Network
    on:
      path: /
    do:
      render.buttons:
        - name: Wakeup PC1
          url: /wakeup/c0:3f:d5:6a:4e:83
        - name: Wakeup PC2
          url: /wakeup/00:4E:01:C4:16:8C
        - name: Ping PC 1
          url: /ping/192.168.178.91
          classes: btn-lg btn-dark
        - name: Ping PC 2
          url: /ping/192.168.178.70
          classes: btn-lg btn-dark
```

The above example creates the following UI:
{{< figure
src="images/buttons-iphone-1.png"
alt="Button UI rendered by HTTPE's built-in button function"
caption="Button UI rendered by HTTPE's built-in button function"
>}}

