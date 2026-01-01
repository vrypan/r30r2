# Rule 30 RNG Examples

This directory contains example programs demonstrating how to use the Rule 30 RNG library in your Go projects.

## Running the Examples

```bash
# Basic usage example
go run basic-usage.go

# Monte Carlo π estimation
go run monte-carlo.go
```

## Examples

### basic-usage.go

Demonstrates all the basic RNG functions:
- Integer generation (Uint32, Uint64, Int)
- Bounded integers (Intn)
- Float generation (Float32, Float64)
- Statistical distributions (NormFloat64, ExpFloat64)
- Byte generation (io.Reader interface)

### monte-carlo.go

Real-world example using Rule 30 RNG for Monte Carlo simulation to estimate π.
Shows how the estimate converges as sample size increases.

## Importing the Library

In your own projects:

```go
import "github.com/vrypan/r30r2/rand"
```

Then install the dependency:

```bash
go get github.com/vrypan/r30r2
```
