# SIMAP

SIMAP is a Simple IMAP Server, designed as a proof of concept command-and-control server to receive a connection via IMAP (with no security controls) and process a base64 encoded string.

# Usage

You will need your own copy of https://github.com/0xflux/ZestyChips (my dotnet stealer this C2 is built for), and that executable will need to be accessible by the webserver so it can be served over HTTP.

Set environment variables as follows, substituting in your desired username and password in plaintext. This will allow you to login to the IMAP server from a client.

```
[System.Environment]::SetEnvironmentVariable("simap_poc_username", "exampleUser", [System.EnvironmentVariableTarget]::User)
[System.Environment]::SetEnvironmentVariable("simap_poc_password", "examplePassword", [System.EnvironmentVariableTarget]::User)
```

In the dockerfile, edit:

```
ENV simap_poc_username=defaultUsername
ENV simap_poc_password=defaultPassword
```

With your chosen username and password, must match the environment of where you run ZestyChips.

If you would prefer this hardcoded, please submit an issue request on Git (or Tweet me).