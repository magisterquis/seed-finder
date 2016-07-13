# Seed finder

Go tool that finds math/rand seed number for given word. Right now it's limited to [a-z] but can be easily modified for wider ascii range.

## Usage

```sh
go get github.com/heppu/seed-finder
go install seed-finder
seed-finder -w hello
```

## Example

```go
package main

import (
    "fmt"
    "math/rand"
)

func main() {
    fmt.Println(getRandStr(-4611686018426513711, 5))
    fmt.Println(getRandStr(13143826, 5))
}

func getRandStr(seed int64, len int) (s string) {
    rand.Seed(seed)
    for i := 0; i < len; i++ {
        s += string(rand.Intn(26) + 97)
    }
    return
}

```
