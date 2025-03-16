# rtmp-go

Go implementation of RTMP, focusing on ease of use and performance. This project does not aim to be complete -- it implements those portions of RTMP that are in use by modern RTMP clients and servers. For example, AMF0 object references are not implemented, and AMF3 is not currently sent or received. Aggregate message handling is also not completed.

# Status

## Initial Release

This library is currently a work in progress. The planned initial release consists of the following high-level tasks:

🗹 AMF0 Library

🗹  RTMP Message Library

☐ Chunk Stream Implementation

☐ Connection Implementation

☐ Server Interface

☐ Client Interface

☐ Examples

## Fast Follow

These will be added shortly after the 1.0 release is complete

☐ RTMPS (TLS) support

☐ RTMP Aggregate Message type

## Not Currently Planned

The following features are not currently planned, but may be candidates for consideration in future versions if there are use cases to motivate them.

☒ AMF3 Library

☒ AMF0 Object References

# Future Plans

We're tracking the progress of [enhanced RTMP](https://github.com/veovera/enhanced-rtmp), and intend to add support in post-1.0 version of the library