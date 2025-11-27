# Go Routine Patterns: Concurrency, Worker Pools, Graceful Shutdown & Retry Logic

## Overview
This repository provides idiomatic Go patterns for concurrency, including graceful shutdown, worker pools, retry logic, and error handling. Each example is paired with a **Mermaid sequence diagram** for clarity.

## Features
- Graceful shutdown with `context` and `WaitGroup`
- Signal handling for clean exits
- Worker pool implementation with dynamic tasks
- Retry logic with safe channel operations
- Timeout-based shutdown
- Error handling and logging
- Visual sequence diagrams with channels as actors

## Notes

- For unbuffered channel
  - Send  to channel `ch <- true`, _must have_ a **matching receive** `<-ch` **happening at the same time**. If not it will be blocked on send. Because it doesn't have **storage capacity**.
- 