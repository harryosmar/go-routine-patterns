```mermaid
sequenceDiagram
    participant Main
    participant SignalHandler as Signal Handler Goroutine
    participant Worker as Worker Goroutine

    Main->>Main: Create cancellable context
    Main->>Main: Create progressCh (buffer size 5)
    Main->>SignalHandler: Start signal listener goroutine
    SignalHandler->>SignalHandler: Wait for SIGINT/SIGTERM
    SignalHandler->>Main: cancel() on signal received

    Main->>Worker: Start Worker goroutine
    Worker->>Worker: Loop: select { ctx.Done() or work }
    Worker->>Worker: Send progressCh <- true
    Worker->>Worker: Print "start Working..."
    Worker->>Worker: Sleep 5s
    Worker->>Worker: Print "done Working..."
    Worker->>Worker: Receive <-progressCh
    Worker->>Worker: Repeat until ctx.Done()
    Worker->>Main: Exit when ctx.Done()

    Main->>Main: Wait for ctx.Done()
    Main->>Main: After cancellation, fill progressCh (cap times)
    Main->>Main: Print "Shutdown complete."
```

## ⚠️ Important Note:
This pattern works for one worker, but if you add multiple workers or dynamic tasks, it can deadlock or panic because:
• Hardcoded buffer size (cap(progressCh)) assumes max concurrency.
• Closing the channel while other goroutines send can cause send on closed channel.

for best
