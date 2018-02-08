# Authentication-proxy

[![Build status](https://ci.appveyor.com/api/projects/status/u0tbng5wockjgd97?svg=true)](https://ci.appveyor.com/project/IvoBellinSalarin/authentication-proxy)

A simple Windows Authentication proxy. This application injects the current session's Kerberos token in the communication between a client software and a corporate proxy. The client software is unable to perform the Negotiate authentication exchanges and the corporate proxy only accepts authenticated users. Usually, the corporate proxy accepts NTLM and Negotiate as authentication protocols.

## Which is the use case, exactly ?

Many package managers and source control managers are not able to perform a Negotiate exchange to authenticate the communication. This means that npm, git, docker, bower and so on will be unable to pass through a corporate proxy.

Some tools, like CNTLM, allow you to pass your NTLM token to the proxy. This is a different protocol, less secure than Negotiate. CNTLM shall be configured writing your username/domain/encrypted password to a file.
A patch for CNTLM allows you to use the Negotiate protocol (avoiding the need for a password saved in a file), but no binary is available nowadays. Moreover, in my personal  environment, CNTLM is *slow*: it won't be able to follow the rhythm of the exchanges between npm and the npm registry.

## Building

The following command will build the application.

```
go get github.com/qichaozhao/authentication-proxy
```

## User manual

	:: authentication-proxy proxies a corporate proxy on port 80
	:: this proxy supports Negotiate or NTLM
    > .\authentication-proxy.exe http://a_corporate_proxy:80

    :: authentication-proxy listens on 3128
    > set HTTPS_PROXY=http://127.0.0.1:3128

    :: one of the following ..
    > npm install
    > git clone https://github.com/..
    > go get github.com/..

## Notes

* You must run the application as user having a valid kerberos session ticket (i.e. log in with your corporate identity, then start this software)
* Does not reply to mutual authentication request, but it's probably somewhat rare to bump into with web applications.
* 64-bit platforms should still offer 32-bit compatible library/API so the application should compile and work. There's afaik no reason for which the application should be 64-bit.
* The application does not add proxy headers, or manipulate any other headers besides Www-Authenticate/Authorization intentionally.
* Works only on Windows (because of the syscalls being called)
