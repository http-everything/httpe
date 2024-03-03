---
weight: 306
title: "Render Buttons"
description: ""
icon: "article"
date: "2024-03-02T14:17:51+01:00"
lastmod: "2024-03-02T14:17:51+01:00"
draft: false
toc: true
---

HTTPE can render a simple responsive website for you with some buttons to trigger actions. 

Currently, buttons can only trigger actions with a GET request.

## Example

```yaml
---
rules:
  - name: Play Music
    on:
      path: /music/play
    do:
      run.script: |
        nohup afplay /Users/thorsten/Swound-2023-03-4.mp3 >/dev/null 2>&1 &
        echo "Music now playing"

  - name: Stop
    on:
      path: /music/stop
    do:
      run.script: killall afplay

  - name: Ping
    on:
      path: /ping
    do:
      run.script: ping -c 4 {{ .Input.Params.Tgt }}
      args:
        timeout: 10

  - name: Some buttons
    on:
      path: /
    do:
      render.buttons:
        - name: Send ping to Google
          url: /ping?Tgt=8.8.8.8
        - name: Send ping to Quad9
          url: /ping?Tgt=9.9.9.9
        - name: ▶️ Play Music
          url: /music/play
          classes: btn-lg btn-outline-warning
        - name: ⏹️ Stop Music
          url: /music/stop
          classes: btn-lg btn-dark

```

{{< figure
src="buttons1.png"
caption="Some buttons created from the rules."
>}}

The first three rules define actions. The fourth rule defines buttons to trigger these actions.

You can change the style of the buttons by adding [button classes](https://getbootstrap.com/docs/5.0/components/buttons/)
from Bootstrap5.

{{% alert icon="🫥" context="warning" %}}
Currently, buttons will load required JS and CSS from public CDNs. From a privacy perspective, this isn't ideal.
Future versions will have all files embedded.
{{% /alert %}}