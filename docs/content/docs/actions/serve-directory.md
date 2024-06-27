---
weight: 305
title: "Serve Directory"
description: ""
icon: "article"
date: "2024-03-02T14:17:51+01:00"
lastmod: "2024-03-02T14:17:51+01:00"
draft: false
toc: true
---

HTTPE comes with a simple static file server.

## Example 

```yaml
---
rules:
  - name: Before
    on:
      path: /dir/some-file.txt
    answer.content: Access denied
    respond:
      on_success:
        http_status: 403

  - name: Server Dir
    on:
      path: /dir
    serve.directory: /tmp/
```

On accessing `http://localhost:3000/dir/` you get a list of files and folders inside the `/tmp` directory.
The directory listing is rendered if no `index.html` file is found.

{{% alert icon="💁‍♂️" context="primary" %}}
On the rule definition do not put a trailing slash at the end of `path`.  
Good: `path: /dir`  
Bad: `path: /dir/`

However, when you want the directory listing getting rendered or the index.html returned, the URL must end with a slash.
{{% /alert %}}

The rules are processed from top to bottom. If you specify a rule with a path that's within the path of a directory, and
if the rule is above the `serve.directory` rule, the rule has precedence. Note that `path` always refers to the URL,
not to the filesystem.