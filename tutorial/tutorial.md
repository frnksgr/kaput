# How to break things

## Install required tools

- wireshark
- netcat
- curl
- some terminals or a terminal multiplexer like _tmux_

```bash
    apt-get update
    apt-get install netcat-openbsd curl wireshark tmux
```

## Warm Up

### Playing around with netcat

We will basically use the command line running tools like _curl_ and
_netcat_. Comandline and Web are pretty cool...

```bash
curl wttr.in/walldorf
```

### See what's going on

Start _wireshark_ capturing packets on _loopback_ device filtering for
_tcp port 4711_

```bash
wireshark -i loopback -f "tcp port 4711 or tcp port 4712" -k
```

### Start with a simple tcp chat

Run a simple tcp chat by running netcat in two different terminals.

```bash
# Start server in terminal-1
nc -l localhost 4711 # open server listening on localhost:4711
# Start client in terminal-2
nc localhost 4711
```

Now we can chat between terminal-1 and terminal-2.

Shut down client with _CTRL+C_ in terminal-2. Inspect the tcp session
in wireshark.

Repeat it. But now, shutdown server with _CTRL+C_ in terminal-1. See
the differences in _wireshark_.

Now, keep server listening by adding `-k` to netcat server.

Optionaly get some out of band info on _stderr_ by adding `-v` to
netcat.

Optionaly coloring _stderr_ (fd 2) to distinguish from stdout (fd 1).

```bash
# Start server in terminal-1
nc -lkv localhost 4711 \
    2> >( \
            while read line; \
                do echo -e "\e[01;31m$line\e[0m" >&2; \
            done \
        )

# Start client in terminal-2
nc -v localhost 4711 \
    2> >( \
            while read line; \
                do echo -e "\e[01;31m$line\e[0m" >&2; \
            done \
        )
```

NOTE: `2> >(...) >&2` redirect _stderr_ to anonymous fifo and back to
_stderr_.
[bash-redirection-cheat-sheet] (https://catonmat.net/ftp/bash-redirections-cheat-sheet.pdf)

### Let's talk http

Open a third terminal for running some other tools. Let's see a client
send a request.

```bash
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

- first line is HTTP method /path protcol/version 
- then come headers, line by line (optional)
- terminated by an mepty line (somehow)
- no body yet

Now, lets add an extra request header using `-H` with curl command,
and send some data using the `-d` switch. NOTE: we need to encapsulate
data in `$'...'` to prevent bash interpreting special chars.

```bash
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

NOTE: the client (_curl_) serving the http protocol. Method has been
changed to _POST_ and some additional header fields have been added.

Now, let's reengineer a valid http response sending a request to a
webserver. Therefore we run a small webserver on _localhost:4712_.

```bash
python3 -m http.server  -b localhost 4712
```

This [server](http://localhost:4712]) let you browse your current
directory.

Let's call it with netcat and ignore the response body sending it to
_/dev/null_ for now. We are only interested in the protocol header
which we get with `-v` switch in curl.

```bash
curl -v http://localhost:4712/ >/dev/null
```

We see the protocol header of the response in the terminal running
curl.

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

Ok, looks like we first have

- a response line showing protocol/version and response code
- then come header lines
- terminated by an empty line (somehow)
- followed by body

Let's try that with our netcat server/client.

We open a connection with _netcat_ to our server running on _port
4712_ and send a request and send a (hopefully) valid request. We do
that by pypint a _here-document_ to stdin of our _netcat_ client
command.

```bash
cat <<EOF | nc localhost 4712
GET / HTTP/1.0

EOF
```

Now, let's do something wrong.

```bash
cat <<EOF | nc localhost 4712
GET HTTP/1.0 /

EOF
```

See how the server reacts on client misbehavior.

In our _netcat_ server we try to respond to a valid request from curl.

```bash
curl -v http://localhost:4711/foo
```

In our server terminal we create a response interactively

```bash
...
HTTP/1.0 200 OK



Damn it
```

Seems the client doesn't find the end of our response, hm. We could
help with some http protocol header.

```bash
...
HTTP/1.0 200 OK
Content-Length: 0

```

Looks much better.

Finally, do a HTTP/1.0 request/response with your netcat client and
server.

Time to stop guessing about HTTP protocol. Where is it defined?
Searching for HTTP e.g. in wikipedia we find
[rfc7230](https://tools.ietf.org/html/rfc7230)
That's the only truth (as long as you are not fooled by some client,
server, intermediary or hacker in between)
