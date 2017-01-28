[![Build Status](https://travis-ci.org/JoeReid/fastTCP.svg)](https://travis-ci.org/JoeReid/fastTCP)
[![GoDoc](https://godoc.org/github.com/JoeReid/fastTCP?status.svg)](http://godoc.org/github.com/JoeReid/fastTCP)
[![Go Report](https://goreportcard.com/badge/github.com/JoeReid/fastTCP)](https://goreportcard.com/report/github.com/JoeReid/fastTCP)

# fastTCP
A highly optimised TCP listener/server system for golang taking inspiration from the performance tweaks in
[tcplisten](https://github.com/valyala/tcplisten) and the go-blocking server in [volley](https://github.com/jonhoo/volley).

## Performance tunings
The default standard library [net.Listener](https://golang.org/pkg/net/#Listener) listens for new TCP connections on one thread.
However on Linux kernel versions > 3.9 the SO_REUSEPORT port flag can be set to allow multiple threads to serve connections simultaneously.
This is done with the help of [this](https://github.com/valyala/tcplisten) project.

This alternative TCP listener also supports use of the TCP_DEFER_ACCEPT flag which can be used to aid performance
on servers with a high rate of new connections. This flag causes a delay in accepting new TCP connections until data is available
from the client. For this reason this option should not be used unless the client writes to the connection first.
See [here](http://man7.org/linux/man-pages/man7/tcp.7.html) for details.

Another performance option supported is the TCP_FASTOPEN flag which will set up the socket to send data before the ACK is received from the.
Client. See [here](https://lwn.net/Articles/508865/) for details.

