/*
 * TestU01 Crush test for ring30mix RNG
 * Reads random data from stdin (pipe from Go program)
 */

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <unistd.h>
#include "TestU01.h"

#define BUFFER_SIZE 8192

static uint8_t buffer[BUFFER_SIZE];
static size_t buffer_pos = 0;
static size_t buffer_len = 0;

/* Refill buffer from stdin */
static int refill_buffer(void) {
    buffer_len = fread(buffer, 1, BUFFER_SIZE, stdin);
    buffer_pos = 0;

    if (buffer_len == 0) {
        fprintf(stderr, "Error: Failed to read from stdin\n");
        return -1;
    }

    return 0;
}

/* TestU01 generator function - returns uint32 */
static unsigned int stdin_gen(void) {
    if (buffer_pos + 4 > buffer_len) {
        if (refill_buffer() < 0) {
            return 0;
        }
    }

    uint32_t value = 0;
    value |= (uint32_t)buffer[buffer_pos++];
    value |= (uint32_t)buffer[buffer_pos++] << 8;
    value |= (uint32_t)buffer[buffer_pos++] << 16;
    value |= (uint32_t)buffer[buffer_pos++] << 24;

    return value;
}

int main(void) {
    unif01_Gen *gen;

    printf("═══════════════════════════════════════════════════════════\n");
    printf("  TestU01 Crush - ring30mix RNG\n");
    printf("═══════════════════════════════════════════════════════════\n");
    printf("\n");
    printf("Reading random data from stdin...\n");
    printf("This test will take approximately 1 hour.\n");
    printf("\n");

    /* Create generator */
    gen = unif01_CreateExternGenBits("ring30mix via stdin", stdin_gen);

    /* Run Crush battery */
    bbattery_Crush(gen);

    /* Cleanup */
    unif01_DeleteExternGenBits(gen);

    printf("\n");
    printf("═══════════════════════════════════════════════════════════\n");
    printf("  Crush Complete\n");
    printf("═══════════════════════════════════════════════════════════\n");

    return 0;
}
