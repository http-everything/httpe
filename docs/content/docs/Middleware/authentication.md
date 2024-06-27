---
weight: 501
title: Authentication
description: ""
date: "2024-03-02T14:17:51+01:00"
lastmod: "2024-03-02T14:17:51+01:00"
draft: false
toc: true
---

## HTTP Basic Auth

You can protect any rule from unauthorized access by adding http basic authentication.

Remember that http basic authentication is not recommended on unencrypted connections. Consider activating TLS.

You can specify a list of username and password pairs per rule. Passwords are stored inside the rules yaml file, either
clear text or as sha265 or sha512 hash.

### Example

```yaml
---
rules:
  - name: Auth1
    on:
      path: /auth/1
    answer.content: I'm in
    with:
      auth_basic:
        - username: john
          password: abc123
        - username: jane
          password: 123abc

  - name: Auth2
    on:
      path: /auth/2
    answer.content: I'm in
    with:
      auth_basic:
        - username: john
          password: 6ca13d52ca70c883e0f0bb101e425a89e8624de51db2d2392593af6a84118090
        - username: jane
          password: dd130a849d7b29e5541b05d2f7f86a4acd4f1ec598c1c9438783f56bc4f0ff80
      auth_hashing: sha256

  - name: Auth3
    on:
      path: /auth/3
    answer.content: I'm in
    with:
      auth_basic:
        - username: john
          password: c70b5dd9ebfb6f51d09d4132b7170c9d20750a7852f00680f65658f0310e810056e6763c34c9a00b0e940076f54495c169fc2302cceb312039271c43469507dc
        - username: jane
          password: 7b6ad79b346fb6951275343948e13c1b4ebca82a5452a6c5d15684377f096ca927506a23a847e6e046061399631b16fc2820c8b0e02d0ea87aa5a203a77c2a7e
      auth_hashing: sha512
```

To create the hashes use:
```shell
echo -n <PASSWORD>|sha256sum
echo -n <PASSWORD>|sha512sum
```

Remember to always use `-n`. Otherwise, the newline character will be part of the password.

On Windows, you can use the below function to calculate the hash of a string.
```powershell
function Get-SHA256Hash($string) {
  $bytes = [System.Text.Encoding]::UTF8.GetBytes($string)
  $sha256 = New-Object System.Security.Cryptography.SHA256Managed
  $hash = $sha256.ComputeHash($bytes)
  return [System.BitConverter]::ToString($hash).Replace("-", "").ToLower()
}

Get-SHA256Hash "123abc"
dd130a849d7b29e5541b05d2f7f86a4acd4f1ec598c1c9438783f56bc4f0ff80
```

## DRY with anchors

Don't repeat yourself. If you want to protect multiple rules with the same credentials, use yaml anchors as shown in the
example.

```yaml
---
define:
  with: &auth
    auth_basic:
      - username: john
        password: "1234"

rules:
  - name: Test 1
    on:
      path: /test1
    answer.content: test1
    with:
      <<: *auth

  - name: Test
    on:
      path: /test2
    answer.content: test2
    with:
      <<: *auth
```