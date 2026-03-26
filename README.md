# rtmp-go

Go implementation of RTMP, focusing on ease of use and performance. This project does not aim to be complete -- it implements those portions of RTMP that are in use by modern RTMP clients and servers. For example, AMF0 object references are not implemented, and AMF3 is not currently sent or received. Aggregate message handling is also not completed.

# Examples

Three trivial example programs are provided:

- **server** — An RTMP server that accepts connections on port 1935. When a client sends a publish command, the server logs the received media. When a client sends a play command, the server streams blank video and silent audio.
- **client-publish** — An RTMP client that connects to a server, publishes a stream, and sends blank media to a server.
- **client-play** — An RTMP client that connects to a server, sends a play command, and logs received media.

## Building

```sh
go build ./examples/server/
go build ./examples/client-publish/
go build ./examples/client-play/
```

## Usage

Each example can be paired with ffmpeg/ffplay or with the other examples. The examples listen on or connect to `localhost:1935`.

### Publishing to the server

The server logs media it receives from a publishing client.

#### With ffmpeg as the publish client

```sh
# Terminal 1: Start the RTMP server
./server

# Terminal 2: Publish a test stream to the server
ffmpeg -re -f lavfi -i testsrc=duration=60:size=1280x720:rate=30 \
  -f lavfi -i "sine=frequency=1000:duration=60" \
  -shortest -f flv rtmp://localhost/live
```

#### With client-publish as the publish client

```sh
# Terminal 1: Start the RTMP server
./server

# Terminal 2: Publish blank media to the server
./client-publish
```

### Playing from the server

The server sends blank video and silent audio to a playing client.

#### With ffplay as the play client

```sh
# Terminal 1: Start the RTMP server
./server

# Terminal 2: Play the stream from the server
ffplay rtmp://localhost/live
```

#### With client-play as the play client

```sh
# Terminal 1: Start the RTMP server
./server

# Terminal 2: Play the stream
./client-play
```

### Client-publish with ffplay as the server

Use ffplay in listen mode as a simple RTMP server, then publish to it:

```sh
# Terminal 1: Start ffplay as a listening RTMP server
ffplay -listen 1 rtmp://localhost/live

# Terminal 2: Publish blank media to ffplay
./client-publish
```

### Client-play with ffmpeg as the server

Use ffmpeg in listen mode as a simple RTMP server that generates test media, then play from it:

```sh
# Terminal 1: Start ffmpeg as a listening RTMP server with test content
ffmpeg -re -f lavfi -i testsrc=duration=60:size=1280x720:rate=30 \
  -f lavfi -i "sine=frequency=1000:duration=60" \
  -shortest -f flv -listen 1 rtmp://localhost/live

# Terminal 2: Play the stream
./client-play
```

# Status

## Initial Release

This library is currently a work in progress. The planned initial release consists of the following high-level tasks:

🗹 AMF0 Library

🗹  RTMP Message Library

🗹 Chunk Stream Implementation

☐ Connection Implementation

☐ Examples

## Fast Follow

These will be added shortly after the 1.0 release is complete

☐ RTMPS (TLS) examples

☐ RTMP Aggregate Message type

## Not Currently Planned

The following features are not currently planned, but may be candidates for consideration in future versions if there are use cases to motivate them.

☒ AMF0 Object References

# Future Plans

We're tracking the progress of [enhanced RTMP](https://github.com/veovera/enhanced-rtmp), and intend to add support in post-1.0 version of the library
