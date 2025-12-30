#!/usr/bin/env bash

# Test entropy for all three RNGs at multiple data sizes
# Compatible with bash 3.x (macOS default)

set -e

# Test sizes: 1MB, 10MB, 100MB
SIZES=(1048576 10485760 104857600)
SIZE_LABELS=("1MB" "10MB" "100MB")
SEED=12345

# Store results in indexed arrays (bash 3.x compatible)
# Format: RNG_METRIC_SIZE
declare -a ENTROPY_RULE30 ENTROPY_MATH ENTROPY_CRYPTO
declare -a CHI_RULE30 CHI_MATH CHI_CRYPTO
declare -a MEAN_RULE30 MEAN_MATH MEAN_CRYPTO
declare -a PI_RULE30 PI_MATH PI_CRYPTO
declare -a CORR_RULE30 CORR_MATH CORR_CRYPTO

echo "═══════════════════════════════════════════════════════════"
echo "  Entropy Testing Suite"
echo "═══════════════════════════════════════════════════════════"
echo ""
echo "Testing three data sizes: 1MB, 10MB, 100MB"
echo "This shows how metrics converge with larger datasets..."
echo ""

# Check if ent is installed
if ! command -v ent &> /dev/null; then
    echo "Error: 'ent' is not installed."
    echo "Install with: brew install ent (macOS) or apt-get install ent (Linux)"
    exit 1
fi

# Build binaries if needed
if [ ! -f ./rule30 ]; then
    echo "Building rule30..."
    make rule30 > /dev/null 2>&1
fi

if [ ! -f ./stdlib-rng ]; then
    echo "Building stdlib-rng..."
    go build -o stdlib-rng stdlib-rng.go
fi

# Temporary files
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# Function to parse ent output
parse_ent() {
    local file=$1
    local output=$(ent "$file")

    # Extract values
    local entropy=$(echo "$output" | grep "Entropy" | awk '{print $3}')
    local chi_square=$(echo "$output" | grep "Chi square" | awk '{print $8}' | tr -d ',')
    local mean=$(echo "$output" | grep "Arithmetic mean" | awk '{print $8}')
    local pi=$(echo "$output" | grep "Monte Carlo" | awk '{print $7}')
    local correlation=$(echo "$output" | grep "Serial correlation" | awk '{print $5}')

    echo "$entropy|$chi_square|$mean|$pi|$correlation"
}

# Generate and test data for each size
for i in "${!SIZES[@]}"; do
    BYTES=${SIZES[$i]}
    LABEL=${SIZE_LABELS[$i]}

    echo "Testing $LABEL ($BYTES bytes)..."

    # Generate test files
    RULE30_FILE="$TMP_DIR/rule30_${LABEL}.bin"
    MATH_FILE="$TMP_DIR/math_${LABEL}.bin"
    CRYPTO_FILE="$TMP_DIR/crypto_${LABEL}.bin"

    ./rule30 --seed=$SEED --bytes=$BYTES > "$RULE30_FILE" 2>/dev/null
    ./stdlib-rng --type=math --seed=$SEED --bytes=$BYTES > "$MATH_FILE" 2>/dev/null
    ./stdlib-rng --type=crypto --bytes=$BYTES > "$CRYPTO_FILE" 2>/dev/null

    # Parse results and store in indexed arrays
    IFS='|' read -r ent chi mean pi corr <<< "$(parse_ent "$RULE30_FILE")"
    ENTROPY_RULE30[$i]=$ent
    CHI_RULE30[$i]=$chi
    MEAN_RULE30[$i]=$mean
    PI_RULE30[$i]=$pi
    CORR_RULE30[$i]=$corr

    IFS='|' read -r ent chi mean pi corr <<< "$(parse_ent "$MATH_FILE")"
    ENTROPY_MATH[$i]=$ent
    CHI_MATH[$i]=$chi
    MEAN_MATH[$i]=$mean
    PI_MATH[$i]=$pi
    CORR_MATH[$i]=$corr

    IFS='|' read -r ent chi mean pi corr <<< "$(parse_ent "$CRYPTO_FILE")"
    ENTROPY_CRYPTO[$i]=$ent
    CHI_CRYPTO[$i]=$chi
    MEAN_CRYPTO[$i]=$mean
    PI_CRYPTO[$i]=$pi
    CORR_CRYPTO[$i]=$corr
done

echo ""
echo "═══════════════════════════════════════════════════════════"
echo "  Shannon Entropy (bits/byte) - Higher is better"
echo "═══════════════════════════════════════════════════════════"
echo ""

printf "%-15s" "RNG"
for label in "${SIZE_LABELS[@]}"; do
    printf " │ %10s" "$label"
done
echo ""
echo "────────────────┼────────────┼────────────┼────────────"

printf "%-15s" "Rule30RNG"
for i in "${!SIZE_LABELS[@]}"; do
    printf " │ %10s" "${ENTROPY_RULE30[$i]}"
done
echo ""

printf "%-15s" "math/rand"
for i in "${!SIZE_LABELS[@]}"; do
    printf " │ %10s" "${ENTROPY_MATH[$i]}"
done
echo ""

printf "%-15s" "crypto/rand"
for i in "${!SIZE_LABELS[@]}"; do
    printf " │ %10s" "${ENTROPY_CRYPTO[$i]}"
done
echo ""

echo ""
echo "═══════════════════════════════════════════════════════════"
echo "  Chi-Square Distribution - Should be ~255 (200-310)"
echo "═══════════════════════════════════════════════════════════"
echo ""

printf "%-15s" "RNG"
for label in "${SIZE_LABELS[@]}"; do
    printf " │ %10s" "$label"
done
echo ""
echo "────────────────┼────────────┼────────────┼────────────"

printf "%-15s" "Rule30RNG"
for i in "${!SIZE_LABELS[@]}"; do
    printf " │ %10s" "${CHI_RULE30[$i]}"
done
echo ""

printf "%-15s" "math/rand"
for i in "${!SIZE_LABELS[@]}"; do
    printf " │ %10s" "${CHI_MATH[$i]}"
done
echo ""

printf "%-15s" "crypto/rand"
for i in "${!SIZE_LABELS[@]}"; do
    printf " │ %10s" "${CHI_CRYPTO[$i]}"
done
echo ""

echo ""
echo "═══════════════════════════════════════════════════════════"
echo "  Arithmetic Mean - Should be 127.5"
echo "═══════════════════════════════════════════════════════════"
echo ""

printf "%-15s" "RNG"
for label in "${SIZE_LABELS[@]}"; do
    printf " │ %10s" "$label"
done
echo ""
echo "────────────────┼────────────┼────────────┼────────────"

printf "%-15s" "Rule30RNG"
for i in "${!SIZE_LABELS[@]}"; do
    printf " │ %10s" "${MEAN_RULE30[$i]}"
done
echo ""

printf "%-15s" "math/rand"
for i in "${!SIZE_LABELS[@]}"; do
    printf " │ %10s" "${MEAN_MATH[$i]}"
done
echo ""

printf "%-15s" "crypto/rand"
for i in "${!SIZE_LABELS[@]}"; do
    printf " │ %10s" "${MEAN_CRYPTO[$i]}"
done
echo ""

echo ""
echo "═══════════════════════════════════════════════════════════"
echo "  Monte Carlo π Approximation - Should be ~3.141593"
echo "═══════════════════════════════════════════════════════════"
echo ""

printf "%-15s" "RNG"
for label in "${SIZE_LABELS[@]}"; do
    printf " │ %10s" "$label"
done
echo ""
echo "────────────────┼────────────┼────────────┼────────────"

printf "%-15s" "Rule30RNG"
for i in "${!SIZE_LABELS[@]}"; do
    printf " │ %10s" "${PI_RULE30[$i]}"
done
echo ""

printf "%-15s" "math/rand"
for i in "${!SIZE_LABELS[@]}"; do
    printf " │ %10s" "${PI_MATH[$i]}"
done
echo ""

printf "%-15s" "crypto/rand"
for i in "${!SIZE_LABELS[@]}"; do
    printf " │ %10s" "${PI_CRYPTO[$i]}"
done
echo ""

echo ""
echo "═══════════════════════════════════════════════════════════"
echo "  Serial Correlation - Should be ~0.0 (uncorrelated)"
echo "═══════════════════════════════════════════════════════════"
echo ""

printf "%-15s" "RNG"
for label in "${SIZE_LABELS[@]}"; do
    printf " │ %10s" "$label"
done
echo ""
echo "────────────────┼────────────┼────────────┼────────────"

printf "%-15s" "Rule30RNG"
for i in "${!SIZE_LABELS[@]}"; do
    printf " │ %10s" "${CORR_RULE30[$i]}"
done
echo ""

printf "%-15s" "math/rand"
for i in "${!SIZE_LABELS[@]}"; do
    printf " │ %10s" "${CORR_MATH[$i]}"
done
echo ""

printf "%-15s" "crypto/rand"
for i in "${!SIZE_LABELS[@]}"; do
    printf " │ %10s" "${CORR_CRYPTO[$i]}"
done
echo ""

echo ""
echo "═══════════════════════════════════════════════════════════"
echo ""
echo "Interpretation:"
echo "  • Entropy closer to 8.000000 = better randomness"
echo "  • Chi-square within 200-310 = good (95% confidence)"
echo "  • Mean closer to 127.5 = more uniform distribution"
echo "  • Monte π closer to 3.141593 = better randomness"
echo "  • Serial correlation closer to 0.0 = less predictable"
echo ""
echo "Note: Metrics should converge/stabilize with larger datasets."
echo ""
