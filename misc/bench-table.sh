#!/bin/bash
# Run Go benchmarks and display results as a table

# Run benchmarks and capture output
echo "Running benchmarks..."
BENCH_OUTPUT=$(go test -bench=. -benchmem ./rand/ 2>&1)

# Check if benchmarks ran successfully
if [ $? -ne 0 ]; then
    echo "Error running benchmarks"
    echo "$BENCH_OUTPUT"
    exit 1
fi

echo ""
echo "═══════════════════════════════════════════════════════════════════════════"
echo "  Benchmark Results Summary"
echo "═══════════════════════════════════════════════════════════════════════════"
echo ""

# Parse benchmark results using awk
echo "$BENCH_OUTPUT" | awk '
BEGIN {
    # Initialize arrays
    split("", results)
    split("", rng_order)
    rng_count = 0
    split("", test_order)
    test_count = 0

    # Track which RNGs and tests we have seen
    split("", seen_rng)
    split("", seen_test)
}

/^Benchmark/ {
    # Extract benchmark name and ns/op
    benchmark = $1
    ns_per_op = $3

    # Parse benchmark name: BenchmarkRNG_Test-cores
    # Remove "Benchmark" prefix
    sub(/^Benchmark/, "", benchmark)

    # Remove -XX suffix (core count)
    sub(/-[0-9]+$/, "", benchmark)

    # Split on underscore to get RNG and Test
    split(benchmark, parts, "_")
    rng = parts[1]
    test = parts[2]

    # Store result
    results[rng, test] = ns_per_op

    # Track unique RNGs (in order)
    if (!(rng in seen_rng)) {
        seen_rng[rng] = 1
        rng_order[rng_count++] = rng
    }

    # Track unique tests (in order)
    if (!(test in seen_test)) {
        seen_test[test] = 1
        test_order[test_count++] = test
    }
}

END {
    # Define desired order of RNGs
    split("MathRand MathRandV2 R30R2 CryptoRand", ordered_rngs)
    ordered_count = 4

    # Print table header
    printf "%-15s", "Algorithm"
    for (i = 0; i < test_count; i++) {
        test = test_order[i]
        printf " | %12s", test
    }
    printf "\n"

    # Print separator
    printf "---------------"
    for (i = 0; i < test_count; i++) {
        printf "-|-------------"
    }
    printf "\n"

    # Print data rows in specified order
    for (i = 1; i <= ordered_count; i++) {
        rng = ordered_rngs[i]
        printf "%-15s", rng

        for (j = 0; j < test_count; j++) {
            test = test_order[j]
            ns = results[rng, test]

            # Format the value
            if (ns != "") {
                printf " | %9.2f ns", ns
            } else {
                printf " | %12s", "N/A"
            }
        }
        printf "\n"
    }

    print ""
    print "═══════════════════════════════════════════════════════════════════════════"
    print ""

    # Calculate and show speed comparison (MathRand as baseline)
    if ("MathRand" in seen_rng) {
        print "Speed Comparison (vs MathRand baseline = 1.00x):"
        print ""
        printf "%-15s", "Algorithm"
        for (i = 0; i < test_count; i++) {
            test = test_order[i]
            printf " |%12s", test
        }
        printf "\n"

        printf "---------------"
        for (i = 0; i < test_count; i++) {
            printf "-|------------"
        }
        printf "\n"

        for (i = 1; i <= ordered_count; i++) {
            rng = ordered_rngs[i]
            printf "%-15s", rng

            for (j = 0; j < test_count; j++) {
                test = test_order[j]
                ns = results[rng, test]
                baseline = results["MathRand", test]

                if (ns != "" && baseline != "") {
                    # Speed = baseline_time / actual_time
                    # > 1.0 = faster, < 1.0 = slower
                    speed = baseline / ns
                    printf " | %10.2fx", speed
                } else {
                    printf " | %12s", "N/A"
                }
            }
            printf "\n"
        }
        print ""
    }
}
'

echo "═══════════════════════════════════════════════════════════════════════════"
echo ""
