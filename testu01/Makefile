# TestU01 Tests Makefile

CC = gcc
CFLAGS = -O3 -Wall -Wextra -I/usr/local/include
LDFLAGS = -L/usr/local/lib -ltestu01 -lprobdist -lmylib -lm

TARGETS = test-smallcrush test-crush test-bigcrush
GO_TARGETS = mathrand-gen mathrandv2-gen

.PHONY: all clean smallcrush crush bigcrush help mathrand-smallcrush mathrand-crush mathrand-bigcrush mathrandv2-smallcrush mathrandv2-crush mathrandv2-bigcrush

all: $(TARGETS)

test-smallcrush: test-smallcrush.c
	$(CC) $(CFLAGS) -o $@ $< $(LDFLAGS)

test-crush: test-crush.c
	$(CC) $(CFLAGS) -o $@ $< $(LDFLAGS)

test-bigcrush: test-bigcrush.c
	$(CC) $(CFLAGS) -o $@ $< $(LDFLAGS)

# Build Go programs
mathrand-gen: mathrand-gen.go
	go build -o mathrand-gen mathrand-gen.go

mathrandv2-gen: mathrandv2-gen.go
	go build -o mathrandv2-gen mathrandv2-gen.go

# Run SmallCrush (quick test)
smallcrush: test-smallcrush ../ring30mix
	@echo "Running SmallCrush test (approx 1-2 minutes)..."
	@echo "Generating unlimited random data and piping to TestU01..."
	@LOGBASE=smallcrush; \
	if [ ! -f $$LOGBASE.log ]; then \
		LOGFILE=$$LOGBASE.log; \
	else \
		N=2; \
		while [ -f $$LOGBASE-$$N.log ]; do \
			N=$$((N+1)); \
		done; \
		LOGFILE=$$LOGBASE-$$N.log; \
	fi; \
	echo "Output will be saved to $$LOGFILE"; \
	../ring30mix --bytes=0 2>/dev/null | ./test-smallcrush | tee $$LOGFILE

# Run Crush (medium test)
crush: test-crush ../ring30mix
	@echo "Running Crush test (approx 1 hour)..."
	@echo "Generating unlimited random data and piping to TestU01..."
	@LOGBASE=crush; \
	if [ ! -f $$LOGBASE.log ]; then \
		LOGFILE=$$LOGBASE.log; \
	else \
		N=2; \
		while [ -f $$LOGBASE-$$N.log ]; do \
			N=$$((N+1)); \
		done; \
		LOGFILE=$$LOGBASE-$$N.log; \
	fi; \
	echo "Output will be saved to $$LOGFILE"; \
	../ring30mix --bytes=0 2>/dev/null | ./test-crush | tee $$LOGFILE

# Run BigCrush (comprehensive test)
bigcrush: test-bigcrush ../ring30mix
	@echo "Running BigCrush test (approx 8 hours)..."
	@echo "Generating unlimited random data and piping to TestU01..."
	@LOGBASE=bigcrush; \
	if [ ! -f $$LOGBASE.log ]; then \
		LOGFILE=$$LOGBASE.log; \
	else \
		N=2; \
		while [ -f $$LOGBASE-$$N.log ]; do \
			N=$$((N+1)); \
		done; \
		LOGFILE=$$LOGBASE-$$N.log; \
	fi; \
	echo "Output will be saved to $$LOGFILE"; \
	../ring30mix --bytes=0 2>/dev/null | ./test-bigcrush | tee $$LOGFILE

# Build ring30mix binary if needed
../ring30mix:
	@cd .. && $(MAKE) $(notdir $@)

# Run SmallCrush with math/rand (comparison)
mathrand-smallcrush: test-smallcrush mathrand-gen
	@echo "Running SmallCrush test with math/rand (approx 1-2 minutes)..."
	@echo "Generating unlimited random data from math/rand..."
	@LOGBASE=mathrand-smallcrush; \
	if [ ! -f $$LOGBASE.log ]; then \
		LOGFILE=$$LOGBASE.log; \
	else \
		N=2; \
		while [ -f $$LOGBASE-$$N.log ]; do \
			N=$$((N+1)); \
		done; \
		LOGFILE=$$LOGBASE-$$N.log; \
	fi; \
	echo "Output will be saved to $$LOGFILE"; \
	./mathrand-gen --bytes=0 2>/dev/null | ./test-smallcrush | tee $$LOGFILE

# Run Crush with math/rand
mathrand-crush: test-crush mathrand-gen
	@echo "Running Crush test with math/rand (approx 1 hour)..."
	@echo "Generating unlimited random data from math/rand..."
	@LOGBASE=mathrand-crush; \
	if [ ! -f $$LOGBASE.log ]; then \
		LOGFILE=$$LOGBASE.log; \
	else \
		N=2; \
		while [ -f $$LOGBASE-$$N.log ]; do \
			N=$$((N+1)); \
		done; \
		LOGFILE=$$LOGBASE-$$N.log; \
	fi; \
	echo "Output will be saved to $$LOGFILE"; \
	./mathrand-gen --bytes=0 2>/dev/null | ./test-crush | tee $$LOGFILE

# Run BigCrush with math/rand
mathrand-bigcrush: test-bigcrush mathrand-gen
	@echo "Running BigCrush test with math/rand (approx 8 hours)..."
	@echo "Generating unlimited random data from math/rand..."
	@LOGBASE=mathrand-bigcrush; \
	if [ ! -f $$LOGBASE.log ]; then \
		LOGFILE=$$LOGBASE.log; \
	else \
		N=2; \
		while [ -f $$LOGBASE-$$N.log ]; do \
			N=$$((N+1)); \
		done; \
		LOGFILE=$$LOGBASE-$$N.log; \
	fi; \
	echo "Output will be saved to $$LOGFILE"; \
	./mathrand-gen --bytes=0 2>/dev/null | ./test-bigcrush | tee $$LOGFILE

# Run SmallCrush with math/rand/v2 (comparison)
mathrandv2-smallcrush: test-smallcrush mathrandv2-gen
	@echo "Running SmallCrush test with math/rand/v2 (approx 1-2 minutes)..."
	@echo "Generating unlimited random data from math/rand/v2..."
	@LOGBASE=mathrandv2-smallcrush; \
	if [ ! -f $$LOGBASE.log ]; then \
		LOGFILE=$$LOGBASE.log; \
	else \
		N=2; \
		while [ -f $$LOGBASE-$$N.log ]; do \
			N=$$((N+1)); \
		done; \
		LOGFILE=$$LOGBASE-$$N.log; \
	fi; \
	echo "Output will be saved to $$LOGFILE"; \
	./mathrandv2-gen --bytes=0 2>/dev/null | ./test-smallcrush | tee $$LOGFILE

# Run Crush with math/rand/v2
mathrandv2-crush: test-crush mathrandv2-gen
	@echo "Running Crush test with math/rand/v2 (approx 1 hour)..."
	@echo "Generating unlimited random data from math/rand/v2..."
	@LOGBASE=mathrandv2-crush; \
	if [ ! -f $$LOGBASE.log ]; then \
		LOGFILE=$$LOGBASE.log; \
	else \
		N=2; \
		while [ -f $$LOGBASE-$$N.log ]; do \
			N=$$((N+1)); \
		done; \
		LOGFILE=$$LOGBASE-$$N.log; \
	fi; \
	echo "Output will be saved to $$LOGFILE"; \
	./mathrandv2-gen --bytes=0 2>/dev/null | ./test-crush | tee $$LOGFILE

# Run BigCrush with math/rand/v2
mathrandv2-bigcrush: test-bigcrush mathrandv2-gen
	@echo "Running BigCrush test with math/rand/v2 (approx 8 hours)..."
	@echo "Generating unlimited random data from math/rand/v2..."
	@LOGBASE=mathrandv2-bigcrush; \
	if [ ! -f $$LOGBASE.log ]; then \
		LOGFILE=$$LOGBASE.log; \
	else \
		N=2; \
		while [ -f $$LOGBASE-$$N.log ]; do \
			N=$$((N+1)); \
		done; \
		LOGFILE=$$LOGBASE-$$N.log; \
	fi; \
	echo "Output will be saved to $$LOGFILE"; \
	./mathrandv2-gen --bytes=0 2>/dev/null | ./test-bigcrush | tee $$LOGFILE

clean:
	rm -f $(TARGETS)
	rm -f $(GO_TARGETS)
	rm -f *.o

help:
	@echo "TestU01 Tests for ring30mix RNG and math/rand comparison"
	@echo ""
	@echo "ring30mix Targets:"
	@echo "  smallcrush           Run SmallCrush with ring30mix (1-2 min)"
	@echo "  crush                Run Crush with ring30mix (1 hour)"
	@echo "  bigcrush             Run BigCrush with ring30mix (8 hours)"
	@echo ""
	@echo "math/rand Comparison Targets:"
	@echo "  mathrand-smallcrush  Run SmallCrush with math/rand (1-2 min)"
	@echo "  mathrand-crush       Run Crush with math/rand (1 hour)"
	@echo "  mathrand-bigcrush    Run BigCrush with math/rand (8 hours)"
	@echo ""
	@echo "math/rand/v2 Comparison Targets:"
	@echo "  mathrandv2-smallcrush Run SmallCrush with math/rand/v2 (1-2 min)"
	@echo "  mathrandv2-crush      Run Crush with math/rand/v2 (1 hour)"
	@echo "  mathrandv2-bigcrush   Run BigCrush with math/rand/v2 (8 hours)"
	@echo ""
	@echo "Build Targets:"
	@echo "  all                  Build all test programs"
	@echo "  clean                Remove built programs"
	@echo ""
	@echo "Example:"
	@echo "  make smallcrush           # Test ring30mix"
	@echo "  make mathrand-smallcrush  # Test math/rand for comparison"
	@echo "  make mathrandv2-smallcrush # Test math/rand/v2 for comparison"
