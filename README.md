# dssgolib

This is a collection of various Go packages.

## Packages

| Package | Description |
|---------|-------------|
| `bitfield` | Marshal/unmarshal structs with bitfields to/from bytes |
| `brokers` | Generic pub/sub broker pattern implementation |
| `btree` | In-memory B-Tree implementation |
| `crc` | CRC16/CRC32 calculations with various polynomials |
| `debounce` | Debounce utilities for rate limiting |
| `debug` | Debug utilities |
| `digest` | Digest and hashing utilities |
| `endianness` | Runtime endianness detection |
| `filepos` | File position tracking |
| `fuzzy` | Fuzzy string matching (Levenshtein distance) |
| `i18n` | Internationalization with relative date formatting |
| `jsonc` | Parse JSON with comments (JSONC) |
| `kvcfg` | Key-value config file parsing |
| `leb128` | Little-endian base 128 encoding |
| `llrb` | Left-Leaning Red-Black tree (2-3 balanced BST) |
| `logrot` | Log rotation with privilege dropping |
| `mapx` | Extended map operations |
| `option` | Functional options pattern |
| `poolx` | Generic object pooling |
| `postgresql` | PostgreSQL name quoting and keyword detection |
| `randshort` | Generate random short strings |
| `ring` | Ring buffer / circular buffer |
| `semver` | Semantic versioning parsing |
| `set` | Set and bitset data structures |
| `sh` | Shell command utilities |
| `sysmon` | Cross-platform system monitoring (processes) |
| `tcpreader` | TCP data reader with buffering |
| `txt` | Text buffer utilities |
| `uptime` | System uptime detection |
| `utils` | Common utilities (strings, math, time, slices, files, etc.) |
| `workermanager` | Worker pool management |

## Installation

```bash
go get github.com/dssutg/dssgolib
```

## Requirements

- Go 1.26+

## Usage

Example of using the bitfield package:

```go
package main

import (
    "fmt"
    "github.com/dssutg/dssgolib/bitfield"
)

type Header struct {
    Version uint8  `bitfield:"4"`
    Flags   uint8  `bitfield:"4"`
    Length  uint16 `bitfield:"16"`
}

func main() {
    h := Header{Version: 1, Flags: 2, Length: 100}
    data, err := bitfield.Marshal(h)
    fmt.Printf("%x, err: %v\n", data, err)
}
```

## License

See individual packages for license information. Main code is under various licenses (see LICENSE files in subdirectories).
