## HTTP Everything

Make everything available over HTTP. 
Execute shell scripts on remote machines via an HTTP request.

To execute a command or script on a remote system the most common approach is using SSH. While SSH is super versatile,
it has some drawbacks. Usually you expose a login shell to the network. For many use cases this is too much. 

So why not exposing a single command via HTTP(s) and call it with curl? 
The httpe server makes it easy and secure. Just define a route and a command and start the server.

[Read more](https://http-everything.github.io)