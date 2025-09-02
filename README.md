## go-streamer 🚀

Lightweight, WIP library for building P2P-style TCP transports in Go. It provides primitives to accept peers, run a customizable handshake, and decode incoming messages via a pluggable `Decoder` interface. Use it as a reusable library to embed networking into your own services and CLIs. ⚙️📦🧪

---

### Table of contents 📚
- **Overview**
- **Features**
- **Project layout**
- **Installation**
- **Library usage**
- **Configuration**
- **Architecture**
- **Extending the decoder**
- **Handshake customization**
- **Examples**
- **Troubleshooting**
- **FAQ**
- **Roadmap**
- **Contributing**
- **License**

---

### Overview 🧭

This repository provides a minimal TCP transport abstraction with a simple accept loop and hook points for handshakes and message decoding.

It is primarily intended to be imported as a library into your applications. A small `main.go` is included as an example harness for local testing.

Core ideas:
- A transport (`TCPTransport`) that listens and accepts connections
- A peer wrapper (`TCPPeer`) for each connection
- A handshake function (default: no-op) to validate/authenticate
- A `Decoder` interface for reading messages from the wire

> Note: The current implementation includes a skeleton `Decoder` and handshake. You can plug in your own protocol and message types.

---

### Features ✨
- **Simple TCP listener** with a clean start/accept loop
- **Pluggable handshake** via `handshakeFunc`
- **Pluggable decoder** via `Decoder` interface
- **Thread-safe peer map (scaffold)** ready for peer tracking
- **Tiny codebase** that’s easy to read and extend

---

### Project layout 🗂️

```
/home/aaaaa/Public/gooooo/go-streamer/
  ├─ go.mod
  ├─ main.go
  ├─ Makefile
  └─ p2p/
     ├─ decoder.go
     ├─ handshake.go
     ├─ tcp_transport.go
     └─ trasnport.go   ← interface definitions (filename contains a typo)
```

- `main.go`: boots a `TCPTransport` on `:30090` and blocks
- `p2p/tcp_transport.go`: `TCPTransport`, `TCPPeer`, accept loop, connection handler
- `p2p/decoder.go`: `Decoder` interface
- `p2p/handshake.go`: `handshakeFunc` + `NOPHandShake`
- `p2p/trasnport.go`: `Transport` interface (note the filename typo)

---

### Requirements 📦
- **Go**: module targets `go 1.24.5` (use the latest stable Go if 1.24.x is not available on your system)
- Unix-like OS recommended; examples use `zsh`/`bash`

---

### Installation 📥

Add the module to your project:

```bash
go get github.com/bidhan948/go-streamer@latest
```

Import the packages you need:

```go
import (
    "github.com/bidhan948/go-streamer/p2p"
)
```

---

### Library usage 🧩

Create and start a transport inside your service:

```go
package main

import (
    "log"
    "github.com/bidhan948/go-streamer/p2p"
)

func main() {
    tr := p2p.NewTCPTransport(":30090")
    // Optionally: set a custom handshake and decoder (see sections below)
    if err := tr.ListenAndAccept(); err != nil {
        log.Fatal(err)
    }
    select {} // block while your app runs other work
}
```

Notes:
- You’ll typically integrate the transport into an existing app and coordinate lifecycle with `context.Context` and proper shutdown hooks.
- Expose setters or exported fields for `handshakeFunc` and `decoder` as needed (see Roadmap).

---

### Makefile commands 🛠️

For developing this library locally, a simple Makefile is provided:

```bash
make build     # build example harness to ./bin/go-streamer
make run       # go run . (example harness)
make run-bin   # build then run ./bin/go-streamer
make clean     # clean artifacts
```

---

### Configuration ⚙️

The listen address is currently hardcoded in `main.go` when constructing `TCPTransport`:

```go
tr := p2p.NewTCPTransport(":30090")
```

To change the port, update the string above (for example `":40000"`). For production-grade apps, consider reading from env vars or flags.

---

### Architecture 🏗️

- **`Transport` interface (`p2p/trasnport.go`)**
  - Declares `ListenAndAccept() error`

- **`TCPTransport` (`p2p/tcp_transport.go`)**
  - Fields: `listenAddress`, `listener`, `handshakeFunc`, `decoder`, `peers` (with `mu`)
  - `ListenAndAccept()` binds and starts the accept loop in a goroutine
  - `startAcceptLoop()` accepts new connections and spawns `handleConnection`
  - `handleConnection()` wraps a `net.Conn` as a `TCPPeer`, performs handshake, then continuously decodes messages using the configured `Decoder`

- **`TCPPeer`**
  - Holds the connection and whether it’s outbound or inbound (currently used for bookkeeping)

- **Handshake (`p2p/handshake.go`)**
  - Type: `type handshakeFunc func(Peer) error`
  - Default: `NOPHandShake` (no-op)

- **Decoder (`p2p/decoder.go`)**
  - Interface: `decode(io.Reader, any) error`
  - Provide your implementation and assign it to `TCPTransport.decoder`

---

### Extending the decoder 🔌

Implement the interface in `p2p/decoder.go`:

```go
type Decoder interface {
    decode(io.Reader, any) error
}
```

Example: a minimal line-based decoder using `bufio.Scanner` (pseudo-code sketch):

```go
package p2p

import (
    "bufio"
    "io"
)

type LineDecoder struct{}

func (LineDecoder) decode(r io.Reader, v any) error {
    // Expect v to be *string
    s, ok := v.(*string)
    if !ok {
        return fmt.Errorf("expected *string")
    }
    scanner := bufio.NewScanner(r)
    if scanner.Scan() {
        *s = scanner.Text()
        return nil
    }
    return scanner.Err()
}
```

Then, wire it up before calling `ListenAndAccept()`:

```go
tr := p2p.NewTCPTransport(":30090")
tr.Decoder = LineDecoder{} // expose/set a public field or setter accordingly
```

Note: in the current code, `decoder` is an unexported field. You can either export it (e.g., `Decoder`) or add a setter method to inject your decoder.

---

### Handshake customization 🤝

Replace the default `NOPHandShake` with your own function to validate peers, perform version checks, or exchange keys:

```go
tr := p2p.NewTCPTransport(":30090")
tr.SetHandshakeFunc(func(p p2p.Peer) error {
    // perform validation
    return nil
})
```

Note: add a `SetHandshakeFunc` method to `TCPTransport` to support this cleanly.

---

### Examples 💡

- **Start the server**

```bash
go run .
```

- **Connect with netcat**

```bash
nc 127.0.0.1 30090
```

You should see a log line in the server like:

```text
NEW INCOMING CONNECTION : { ... }
```

> Tip: The format string currently uses `/n` instead of `\n` in the log. You may want to change it to `"\n"` for proper newlines.

---

### Troubleshooting 🧰
- **Port already in use**: change the listen address (e.g., `":0"` for an ephemeral port) or stop the conflicting process.
- **Decoder errors**: ensure your `Decoder` matches the wire format and the type of the destination message.
- **Handshake failures**: log and carefully validate what your handshake expects; ensure both sides agree on protocol/version.
- **Typos**: file `p2p/trasnport.go` is spelled with an `s` after `tra`—rename if desired.

---

### FAQ ❓
- **Q: Does this implement a complete P2P protocol?**
  - A: No. It’s a transport scaffold with hooks for you to build on.
- **Q: How do I send messages to peers?**
  - A: Extend `TCPPeer` with send/write helpers and track peers in the `peers` map.
- **Q: Is there TLS?**
  - A: Not yet. You can replace `net.Listen`/`net.Conn` with `tls` equivalents.

---

### Roadmap 🗺️
- Export/setter for `decoder`
- Public setter for `handshakeFunc`
- Proper peer management (add/remove, heartbeat, IDs)
- Outbound dialing (connect to remote peers)
- Graceful shutdown and context handling
- Structured logging
- Optional TLS support
- Example decoders (JSON, length-prefixed, protobuf)

---

### Contributing 👋
Contributions are welcome! Please:
- Open an issue to discuss significant changes
- Keep PRs small and focused
- Add clear commit messages and comments where reasoning isn’t obvious

---



