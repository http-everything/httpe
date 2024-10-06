## HTTP Everything

Make everything available over HTTP. 
Execute shell scripts on remote machines via an HTTP request.

To execute a command or script on a remote system the most common approach is using SSH. While SSH is super versatile,
it has some drawbacks. Usually you expose a login shell to the network. For many use cases this is too much. 

So why not exposing a single command via HTTP(s) and call it with curl? 
The httpe server makes it easy and secure. Just define a route and a command and start the server.

HTTPE can do more for you:
- Execute commands via HTTP request
- Make requests and commands dynamic by using templates
- Respond with static content
- Service directories
- Send emails via HTTP requests
- Execute commands asynchronously

[Read more](https://httpe.io/)