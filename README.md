pingoo
======
A simple, quick (and hacky) wrapper lib for the go-fastping package.

NOTE: You'll probably have to run this (and tests) via `sudo` or as a privileged user due to raw socket usage.

## Usage
```go
package main

import (
    "fmt"
    "os"
    "time"

    "github.com/dselans/pingoo"
)

func main() {
    host := "google.com"

    myPinger, err := pingoo.New(host)
    if err != nil {
        fmt.Printf("ERROR: %v\n", err)
        os.Exit(1)
    }

    count := 5

    fmt.Printf("Attempting to ping %v...\n", host)

    received, err := myPinger.Ping(count, time.Duration(100)*time.Millisecond)
    if err != nil {
        fmt.Printf("ERROR: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Received %v/%v responses!\n", received, count)
}
```

Example in `example/example.go`
