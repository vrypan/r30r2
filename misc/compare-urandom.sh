#!/bin/bash
# Compare Rule30 RNG vs /dev/urandom throughput using dd
# Both sources pipe to dd for fair comparison

set -e

# Change to script directory, then go to parent
cd "$(dirname "$0")"
cd ..

# Build rule30 if needed
if [ ! -f ./rule30 ]; then
    echo "Building rule30..."
    make rule30 > /dev/null 2>&1
fi

echo "═══════════════════════════════════════════════════════════"
echo "  Rule30 vs /dev/urandom - Throughput Comparison (dd)"
echo "═══════════════════════════════════════════════════════════"
echo ""

# Test sizes in MB: 10MB, 100MB, 1GB
SIZES_MB=(10 100 1024)
SIZE_LABELS=("10 MB" "100 MB" "1 GB")

# Store results
declare -a RULE30_TIMES
declare -a URANDOM_TIMES

for i in "${!SIZES_MB[@]}"; do
    SIZE_MB=${SIZES_MB[$i]}
    LABEL=${SIZE_LABELS[$i]}
    SIZE_BYTES=$((SIZE_MB * 1024 * 1024))

    echo "Testing $LABEL..."

    # Test Rule30 (fixed bytes mode piped to dd)
    START=$(date +%s.%N)
    ./rule30 --bytes=$SIZE_BYTES 2>/dev/null | dd of=/dev/null bs=1m 2>/dev/null
    END=$(date +%s.%N)
    RULE30_TIME=$(echo "$END - $START" | bc)
    RULE30_TIMES[$i]=$RULE30_TIME

    # Test /dev/urandom
    START=$(date +%s.%N)
    dd if=/dev/urandom of=/dev/null bs=1m count=$SIZE_MB 2>/dev/null
    END=$(date +%s.%N)
    URANDOM_TIME=$(echo "$END - $START" | bc)
    URANDOM_TIMES[$i]=$URANDOM_TIME

    echo "  Rule30:       ${RULE30_TIME}s"
    echo "  /dev/urandom: ${URANDOM_TIME}s"
    echo ""
done

echo "═══════════════════════════════════════════════════════════"
echo "  Summary Table"
echo "═══════════════════════════════════════════════════════════"
echo ""

# Table header
printf "%-10s │ %15s │ %15s │ %10s\n" "Size" "Rule30" "/dev/urandom" "Speedup"
printf "───────────┼─────────────────┼─────────────────┼───────────\n"

# Table rows
for i in "${!SIZES_MB[@]}"; do
    LABEL=${SIZE_LABELS[$i]}
    RULE30_TIME=${RULE30_TIMES[$i]}
    URANDOM_TIME=${URANDOM_TIMES[$i]}

    # Calculate throughput (MB/s)
    MB=${SIZES_MB[$i]}

    RULE30_MBPS=$(echo "scale=0; $MB / $RULE30_TIME" | bc)
    URANDOM_MBPS=$(echo "scale=0; $MB / $URANDOM_TIME" | bc)

    # Calculate speedup
    SPEEDUP=$(echo "scale=2; $URANDOM_TIME / $RULE30_TIME" | bc)

    printf "%-10s │ %10s MB/s │ %10s MB/s │ %9.2fx\n" \
        "$LABEL" "$RULE30_MBPS" "$URANDOM_MBPS" "$SPEEDUP"
done

echo ""
echo "═══════════════════════════════════════════════════════════"
echo ""
echo "Notes:"
echo "  • Both sources pipe to dd for fair comparison (bs=1m)"
echo "  • Rule30: Fixed-size mode (--bytes=N) | dd"
echo "  • /dev/urandom: Kernel CSPRNG with dd reader"
echo "  • Speedup > 1.0 means Rule30 is faster"
echo ""
echo "Command format:"
echo "  rule30 --bytes=\$((SIZE * 1024 * 1024)) | dd of=file.data bs=1m"
echo ""
