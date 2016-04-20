# Authentication-proxy

Simple, performing authentication proxy. It injects the current session Kerberos token to the communication between a client (unable to perform the negotiate authentication scheme) and a corporate proxy accessible thtough Negotiate.

## Which is the problem, exactly ?

Many package managers, source control managers are not able to perform a Negotiate exchange to authenticate the communication. This means that npm, git, docker, bower and so on will be unable to pass through a corporate proxy.

Some tools, like CNTLM, allow you to pass your NTLM token to the proxy. This is a different protocol, less secure than Negotiate. A patch for CNTLM allows you to use the Negotiate protocol, but no binary is available nowadays. Moreover, in my personal  environment, CNTLM is slow. It won't be able to follow the rhythm of the exchanges between npm and the npm registry.

# Building

The following command should build the application. It is a little bit large, but it should not require any dependencies from target the systems.

```
go get github.com/nilleb/authentication-proxy
go build github.com/nilleb/authentication-proxy
```

## User manual

	:: authentication-proxy proxies a corporate proxy on port 80
	:: this proxy supports Negotiate or NTLM
    $ .\authentication-proxy.exe http://a_corporate_proxy:80

    :: authentication-proxy listens on 8080
    $ set HTTPS_PROXY-http://127.0.0.1:8080
    :: one of the following ..
    $ npm install
    $ git clone https://github.com/..
    $ go get github.com/..

# Notes

* You must run the application as user that has valid kerberos login and tickets. 
* Does not reply to mutual authentication request, but it's probably somewhat rare to bump into with web applications.
* 64-bit platforms should still offer 32-bit compatible library/API so the application should compile and work. There's afaik no reason why the application should be 64-bit.
* The application does not add proxy headers, or manipulate any other headers besides Www-Authenticate/Authorization intentionally.
* Works only on Windows (because the syscalls being called)