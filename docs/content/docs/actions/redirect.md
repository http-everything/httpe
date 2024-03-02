---
weight: 304
title: "Redirect"
description: ""
icon: "article"
date: "2024-02-12T13:12:42+01:00"
lastmod: "2024-02-12T13:12:42+01:00"
draft: false
toc: true
---

## Example

```yaml
---
rules:
  - name: redirect perm
    on:
      path: /redirect/google
    do:
      redirect.permanent: https://www.google.com

  - name: redirect temp
    on:
      path: /redirect/{new_loc}
    do:
      redirect.temporary: https://{{ .Input.URLPlaceholders.new_loc }}
```

To see it in action execute the following requests:

```shell
$ curl localhost:3000/redirect/google -v
*   Trying [::1]:3000...
* Connected to localhost (::1) port 3000
> GET /redirect/google HTTP/1.1
> Host: localhost:3000
> User-Agent: curl/8.4.0
> Accept: */*
> 
< HTTP/1.1 301 Moved Permanently
< Location: https://www.google.com

$ curl localhost:3000/redirect/www.example.com -v
*   Trying [::1]:3000...
* Connected to localhost (::1) port 3000
> GET /redirect/www.example.com HTTP/1.1
> Host: localhost:3000
> User-Agent: curl/8.4.0
> Accept: */*
> 
< HTTP/1.1 302 Found
< Location: https://www.example.com
```