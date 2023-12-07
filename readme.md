# SIMAP

SIMAP is a Simple IMAP Server, designed as a proof of concept command-and-control server to receive a connection via IMAP (with no security controls) and process a base64 encoded string.

# Usage

Set environment variables as follows, substituting in your desired username and password in plaintext. This will allow you to login to the IMAP server from a client.

```
[System.Environment]::SetEnvironmentVariable("simap_poc_username", "exampleUser", [System.EnvironmentVariableTarget]::User)
[System.Environment]::SetEnvironmentVariable("simap_poc_password", "examplePassword", [System.EnvironmentVariableTarget]::User)
```