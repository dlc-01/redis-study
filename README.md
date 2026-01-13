Redis-compatible Server (Go)

A Redis-compatible server implemented in Go, built as a learning project focused on infrastructure-level backend development, concurrency, and protocol design.

The project implements a subset of Redis features with emphasis on Streams and blocking reads.

⚠️ Work in progress — new commands and improvements are added incrementally.

⸻

Features

Core
	•	RESP2 protocol support
	•	In-memory storage
	•	Concurrent access with proper synchronization
	•	Expiration (TTL) support

Implemented Redis Commands

Strings
	•	GET
	•	SET (with PX)
	•	TYPE

Lists
	•	LPUSH
	•	RPUSH
	•	LPOP
	•	LLEN
	•	LRANGE
	•	BLPOP

Streams
	•	XADD
	•	XRANGE
	•	XREAD
	•	Multiple streams
	•	Blocking reads (BLOCK <ms>)
	•	Infinite blocking (BLOCK 0)
	•	$ semantics (read only new entries)

⸻

Streams Implementation Details
	•	Stream entries are stored sorted by ID
	•	Each stream maintains an indexed LastID to avoid reparsing
	•	Efficient range and read operations using binary search
	•	Blocking XREAD implemented via per-stream waiters
	•	Correct Redis semantics:
	•	XREAD is exclusive (> comparison)
	•	$ resolves once at subscription time
	•	Blocking returns null array only on timeout

⸻

Architecture Overview

The project follows a clean layered architecture:

parser        → RESP decoding, command parsing
application  → command definitions & handlers
domain       → core types (Stream, StreamID, errors)
storage      → in-memory engine, concurrency & blocking logic

Key Design Choices
	•	Stream ID parsing isolated in streamidcodec
	•	Domain-level streamid.ID used instead of raw strings
	•	Storage layer owns concurrency and blocking semantics
	•	Handlers are thin and protocol-focused

⸻

Running the Server

go run ./cmd/server

Then connect using redis-cli:

redis-cli -p 6379


⸻

Example

XADD mystream * temperature 25
XREAD BLOCK 0 STREAMS mystream $


⸻

Goals of the Project
	•	Understand Redis internals (Streams, blocking, IDs)
	•	Practice concurrent system design in Go
	•	Build infrastructure-grade abstractions
	•	Prepare for backend / infrastructure interviews

⸻

Roadmap
	•	Persistence (AOF / RDB-like)
	•	Consumer groups (XGROUP, XREADGROUP)
	•	Memory optimizations
	•	Better error compatibility with Redis
	•	Benchmarks

⸻

Disclaimer

This is not a production-ready Redis replacement.
It is a learning-oriented implementation designed to explore how Redis works internally.
