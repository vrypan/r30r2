#!/bin/bash
# Compare Rule30 RNG vs /dev/urandom throughput

set -e

# Build rule30 if needed
if [ ! -f ./rule30 ]; then
    echo "Building rule30..."
    make rule30 > /dev/null 2>&1
fi

echo "═══════════════════════════════════════════════════════════"
echo "  Rule30 vs /dev/urandom - Throughput Comparison"
echo "═══════════════════════════════════════════════════════════"
echo ""

# Test sizes: 10MB, 100MB, 1GB
SIZES=(10485760 104857600 1073741824)
SIZE_LABELS=("10 MB" "100 MB" "1 GB")

# Store results
declare -a RULE30_TIMES
declare -a URANDOM_TIMES

for i in "${!SIZES[@]}"; do
    BYTES=${SIZES[$i]}
    LABEL=${SIZE_LABELS[$i]}

    echo "Testing $LABEL..."

    # Test Rule30
    START=$(date +%s.%N)
    ./rule30 --bytes=$BYTES > /dev/null 2>&1
    END=$(date +%s.%N)
    RULE30_TIME=$(echo "$END - $START" | bc)
    RULE30_TIMES[$i]=$RULE30_TIME

    # Test /dev/urandom
    START=$(date +%s.%N)
    dd if=/dev/urandom of=/dev/null bs=1M count=$((BYTES/1048576)) 2>/dev/null
    END=$(date +%s.%N)
    URANDOM_TIME=$(echo "$END - $START" | bc)
    URANDOM_TIMES[$i]=$URANDOM_TIME

    echo "  Rule30:      ${RULE30_TIME}s"
    echo "  /dev/urandom: ${URANDOM_TIME}s"
    echo ""
done

echo "═══════════════════════════════════════════════════════════"
echo "  Summary Table"
echo "═══════════════════════════════════════════════════════════"
echo ""

# Table header
printf "%-12s │ %12s │ %12s │ %12s\n" "Size" "Rule30" "/dev/urandom" "Speedup"
printf "─────────────┼──────────────┼──────────────┼─────────────\n"

# Table rows
for i in "${!SIZES[@]}"; do
    LABEL=${SIZE_LABELS[$i]}
    RULE30_TIME=${RULE30_TIMES[$i]}
    URANDOM_TIME=${URANDOM_TIMES[$i]}

    # Calculate throughput (MB/s)
    BYTES=${SIZES[$i]}
    MB=$(echo "scale=2; $BYTES / 1048576" | bc)

    RULE30_MBPS=$(echo "scale=2; $MB / $RULE30_TIME" | bc)
    URANDOM_MBPS=$(echo "scale=2; $MB / $URANDOM_TIME" | bc)

    # Calculate speedup
    SPEEDUP=$(echo "scale=2; $URANDOM_TIME / $RULE30_TIME" | bc)

    printf "%-12s │ %9.0f MB/s │ %9.0f MB/s │ %10.2fx\n" \
        "$LABEL" "$RULE30_MBPS" "$URANDOM_MBPS" "$SPEEDUP"
done

echo ""
echo "═══════════════════════════════════════════════════════════"
echo ""
echo "Notes:"
echo "  • Rule30: Pure computational PRNG"
echo "  • /dev/urandom: Kernel CSPRNG (syscalls + context switches)"
echo "  • Speedup > 1.0 means Rule30 is faster"
echo ""
