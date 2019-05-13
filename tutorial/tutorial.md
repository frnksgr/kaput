# How to break things

## Install required tools

- netcat
- curl
- python3 (optional)
- wireshark (optional)
- tmux (optional)

```bash
    apt-get update
    apt-get install netcat-openbsd curl wireshark tmux python3
```

## Warm Up

### Playing around with netcat

Run a simple http client and server just with netcat on the comnmandline.

Comandline is cool.

```bash
curl wttr.in/walldorf
```

### See what's going on

Start _wireshark_ capturing all packets relevant in this scenario.

```bash
wireshark -i loopback -f "tcp port 4711 or tcp port 4712" -k
```

### Start with a simple TCP chat

Run a simple tcp chat by running netcat in two different terminals.

```bash
# Start server in terminal-1
nc -l localhost 4711 # open server listening on localhost:4711
```

Actually this is not a _server_ rahter just a listening socket.

```bash
# Start client in terminal-2
nc localhost 4711
```

Keep server listening by adding `-k` to netcat server. 

Get some info on _stderr_ by adding `-v` to netcat.

Nice coloring of _stderr_ (fd 2) to distinguish from output on stdout (fd 1).

```bash
# Start server in terminal-1
nc -lkv localhost 4711 \
    2> >( \
            while read line; \
                do echo -e "\e[01;31m$line\e[0m" >&2; \
            done \
        )
```

```bash
# Start client in terminal-2
nc -v localhost 4711 \
    2> >( \
            while read line; \
                do echo -e "\e[01;31m$line\e[0m" >&2; \
            done \
        )
```
[bash-redirection-cheat-sheet] (https://catonmat.net/ftp/bash-redirections-cheat-sheet.pdf)

### Let's talk HTTP

Open a third terminal for running some other tools. Let's see a client
send a request.

```bash
# send http request to our netcat server using curl
curl localhost:4711 /path
```

In terminal-1 we see a proper http request created by curl.

```bash
...
GET /path HTTP/1.1
Host: localhost:4711
User-Agent: curl/7.58.0
Accept: */*

...
```

 Looks like

- first line is HTTP method, path, protcol/version separated by a space
- then come headers, line by line
- terminated by an empty line
- no body yet

Add an extra request header using `-H` with curl command,
and send some data using the `-d` switch. NOTE: we need to encapsulate
data in `$'...'` to prevent bash interpreting special chars.

```bash
# add headers and body with curl
curl -v -H "My-header: value" \
    -d $'the body\n' \
    http://localhost:4711/path
```

```bash
...
POST /path HTTP/1.1
Host: localhost:4711
User-Agent: curl/7.58.0
Accept: */*
My-header: value
Content-Length: 9
Content-Type: application/x-www-form-urlencoded

the body
...
```

NOTE: _curl_ handling HTTP protocol under the hood.
Method changed from _GET_ to _POST_, some additional headers.

Show a valid HTTP response with curl calling a real server.
Therefore run a small webserver on _localhost:4712_ with python3.

```bash
# Run a simple webserver browsing current directory
python3 -m http.server  -b localhost 4712
```
Check with browser on http://localhost:4712 

```bash
# send request to server to see HTTP response
curl -v http://localhost:4712/ >/dev/null
```

```bash
...
< HTTP/1.0 200 OK
< Server: SimpleHTTP/0.6 Python/3.6.7
< Date: Sun, 12 May 2019 16:31:24 GMT
< Content-type: text/html; charset=utf-8
< Content-Length: 5202
<
...
```

Looks like

- response line showing protocol/version and response code separated by spaces
- then come headers line by line
- terminated by an empty line
- followed by body

Let's try that with netcat. I.e. create a valid HTTP request and HTTP response

```bash
# create a calid HTTP request with netcat client
cat <<EOF | nc localhost 4712
GET / HTTP/1.0

EOF
```

```bash
# See what happens if we create an invalid request
cat <<EOF | nc localhost 4712
GET HTTP/1.0 /

EOF
```

Creating a valid HTTP response in netcat server terminal.

```bash
# Send a valid request with curl
curl -v http://localhost:4711/foo
```

Create a response interactively.

```bash
...
HTTP/1.0 200 OK


Damn it
```

Seems the client doesn't find the end of body, hm.
Help with some HTTP protocol header.

```bash
...
HTTP/1.0 200 OK
Content-Length: 0

```

Looks much better.

Stop guessing about HTTP protocol. Study its definition and always refer to it on questions?
[rfc7230](https://tools.ietf.org/html/rfc7230)
