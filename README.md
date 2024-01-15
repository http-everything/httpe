## HTTP Everything

Make everything available over HTTP. 
Execute shell scripts on remote machines via an HTTP request.

To execute a command or script on a remote system the most common approach is using SSH. While SSH is super versatile,
it has some drawbacks. Usually you expose a login shell to the network. For many use cases this is too much. 

So why not exposing a single command via HTTP(s) and call it with curl? 
The httpe server makes it easy and secure. Just define a route and a command and start the server.

```yaml
---
# a routes.yaml file for httpe
- name: Execute a shell command
  on:
    path: /commands/log
    method: get
  do:
    script:
      interpreter: /bin/bash
      exec: |
        echo "$(date) -- /command/logs has been called" >> /tmp/command.log
        echo "httpe did the job"
```

Start the httpe server with
```shell
./httpe --rules ./rules.yaml
```

Calling `curl http://localhost:8080/commands/log` will give you `httpe did the job` and the defined script has been
executed.