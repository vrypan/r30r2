# TestU01 Statistical Tests for ring30mix RNG

This directory contains setup and programs for testing ring30mix RNG with TestU01, a comprehensive statistical test suite.

## Installation

### macOS

```bash
# Install dependencies
brew install gsl

# Download and build TestU01
cd /tmp
curl -O http://simul.iro.umontreal.ca/testu01/TestU01.zip
unzip TestU01.zip
cd TestU01-1.2.3

# Configure and install
./configure --prefix=/usr/local
make
sudo make install
```

### Linux (Ubuntu/Debian)

```bash
sudo apt-get install libgsl-dev

# Download and build TestU01
cd /tmp
wget http://simul.iro.umontreal.ca/testu01/TestU01.zip
unzip TestU01.zip
cd TestU01-1.2.3

./configure --prefix=/usr/local
make
sudo make install
```

## Verify Installation

```bash
# Check if libraries are installed
ls /usr/local/lib/libtestu01*
ls /usr/local/include/TestU01/

# You should see:
# /usr/local/lib/libtestu01.a
# /usr/local/include/TestU01/*.h
```

## Running Tests

### 1. SmallCrush (Quick Test - ~1 minute)

```bash
make smallcrush
./test-smallcrush
```

### 2. Crush (Medium Test - ~1 hour)

```bash
make crush
./test-crush
```

### 3. BigCrush (Comprehensive Test - ~8 hours)

```bash
make bigcrush
./test-bigcrush
```

## Test Programs

- `test-smallcrush.c` - Runs SmallCrush battery (10 tests)
- `test-crush.c` - Runs Crush battery (96 tests)
- `test-bigcrush.c` - Runs BigCrush battery (106 tests)

All programs read random data from the ring30mix RNG via a pipe.

## Understanding Results

TestU01 will print:
- **PASS**: p-value in reasonable range (typically 0.001 to 0.999)
- **FAIL**: p-value very close to 0 or 1 (< 0.001 or > 0.999)

A good RNG should pass all or nearly all tests. A few failures out of hundreds of tests can be acceptable due to random chance.

## References

- [TestU01 Home Page](http://simul.iro.umontreal.ca/testu01/tu01.html)
- [TestU01 Guide](http://simul.iro.umontreal.ca/testu01/guideshorttestu01.pdf)
