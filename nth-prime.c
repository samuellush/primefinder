// CS 131: nth-prime in C
//                   
// Names:
//
// Compile with:
//      make nth-prime
//
// Run with
//    ./nth-prime N
//
// to find the Nth prime

#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <inttypes.h>
#include <stdbool.h>
#include <limits.h>
#include <assert.h>
#include <math.h>
#include <string.h> // for memset

// str = ordinal(n)
//   -- Function for putting cute "st", "nd", "rd", or "th" on numbers.

const char* ordinal(uint64_t num) 
{
    int digit = num % 10UL;
    switch (digit) {
        case 1:  return "st";
        case 2:  return "nd";
        case 3:  return "rd";
        default: return "th";
    }
}

// p = primeUnderestimate(nth)
//   -- Returns a number p <= the nth prime.

double primeUnderestimate(double n)
{
    if (n < 5.0)
        return 2.0;
    if (n < 6077.0)
        return n*(log(n) + log(log(n)) - 1.0);
    return n*(log(n) + log(log(n)) - 0.9718);   // 0.97174375151226449448
}

// p = primeOverestimate(nth)
//   -- Returns a number p > the nth prime.

double primeOverestimate(double n)
{
    if (n < 6.0)
        return 12.0;
    if (n < 8602.0) //  according to one source, it's 7022.0, but that's wrong
        return n*log(n) + n*log(log(n));
    return n*(log(n) + log(log(n)) - 0.9385);
}

// count = primeCountOverestimate(n)
//   -- Returns an overestimate of the number of primes < n.  

double primeCountOverestimate(double n)
{
    if (n < 2.0)
        return 0.0;
    return 1.25506*n/log(n);
}


// t = primesTable(n, &c)
//   -- returns a dynamically allocated array of the first n primes,
//      and a count of how many there is stored in the table.
//      We employ the classic Sieve of Eratosthenes algorithm, with the
//      usual optimization of avoiding even numbers.
//      The code limits itself to primes < 2^32.

uint32_t* primesTable(uint32_t limit, uint32_t* countPtr)
{
    const uint32_t maxPrimes = (uint32_t) primeCountOverestimate(limit);
    uint32_t* primeNum = calloc(maxPrimes, sizeof(uint32_t));
    uint32_t numPrimes = 0;
 
    // The primes we'll be using will all be less than 2^32 because that's
    // sufficient to enumerate all primes up to 2^64. 

    // Allocate our array and initially declare that everything is prime.
    uint8_t* const isPrime  = malloc(limit/2);
    memset(isPrime, true, limit/2);

    primeNum[numPrimes++] = 2;  // never used, but there for completeness
                                // and proper numbering

    for (uint32_t n = 3; n < limit; n += 2) {
        if (isPrime[n/2] == false)
            continue;
        primeNum[numPrimes++] = n;
        
        if (n >= (1<<16))       // if n*n would overflow a 32 bit int, we're
            continue;           // already done with crossing out, and the
                                // test might not catch it
        
        // Cross off all the composites
        for (uint32_t i = n*n; i < limit; i += 2*n)
            isPrime[i/2] = false;
    }

    free(isPrime);
    
    *countPtr = numPrimes;
    return primeNum;
}


// count = countPrimes(from, to, nth, &nth_prime)
//    -- Counts how many primes there are in the range [from,to); if
//       nth_prime is not NULL, it will also store the nth prime it found
//       in *nth_prime (if there aren't that many, it'll store nothing)

uint64_t countPrimes(uint64_t rangeStart, uint64_t rangeEnd, uint64_t nth, 
             uint64_t* thePrime)
{
    if (thePrime != NULL)
        printf("- Counting the primes in %" PRIu64 "..%" PRIu64
               " and finding the %" PRIu64 "%s ...\n",
               rangeStart, rangeEnd, nth, ordinal(nth));
    else
        printf("- Counting the primes in %" PRIu64 "..%" PRIu64
               " ...\n", rangeStart, rangeEnd);

    uint32_t  crossingMax = (uint32_t) ceil(sqrt((double)rangeEnd));
    uint32_t  crossingPrimesCount;
    uint32_t* crossingPrimes = primesTable(crossingMax, &crossingPrimesCount);

    printf("\tUsing %d primes (%d..%d) for crossing off\n", crossingPrimesCount,
            crossingPrimes[0], crossingPrimes[crossingPrimesCount-1]);

    uint64_t primeCount = 0;

    // Handling 2 is a pain, but necessary...
    if (rangeStart < 3) {
        ++primeCount;
        --nth;
        if (thePrime != NULL && nth == 0)
            *thePrime = 2;
        rangeStart = 3;
    }

    // Adjust the start of the range up to be an odd number.
    rangeStart = rangeStart + 1 - (rangeStart % 2);
    
    // Bail if there's nothing to do
    if (rangeStart >= rangeEnd)
        return primeCount;

    // Create our array for crossing off...
    uint64_t range = rangeEnd - rangeStart;
    uint8_t* const isPrime  = malloc((range+1)/2);
    memset(isPrime, true, (range+1)/2);

    // Cross off all the composites...
    for (uint32_t i = 1; i < crossingPrimesCount; ++i) {
        uint64_t prime = crossingPrimes[i];
        
        // We need to adjust the crossing-off prime to be the first prime
        // in the range.  First, move it into the range.
        uint64_t rangeStartAdjusted =
            rangeStart + (prime - (rangeStart % prime)) % prime;
        if (rangeStartAdjusted % 2 == 0)        // It can't be even.
            rangeStartAdjusted += prime;
        if (rangeStartAdjusted < prime*prime)   // Only cross off the things
            rangeStartAdjusted = prime*prime;   // we need to.

        // Actually cross off the composites.
        for (uint64_t j = rangeStartAdjusted; j < rangeEnd; j += 2*prime)
            isPrime[(j-rangeStart)/2] = false;
    }
    
    // Now, actually count how many primes we're left with...
    for (uint64_t j = rangeStart; j < rangeEnd; j += 2) {
        if (isPrime[(j-rangeStart)/2]) {
            ++primeCount;
            if (thePrime != NULL && --nth == 0)
                *thePrime = j;
        }
    }

    printf("\t%" PRIu64 " primes found\n", primeCount);

    free(crossingPrimes);       // Get rid of our array...
    free(isPrime);

    return primeCount;
}

int main(int argc, char **argv)
{
    if (argc != 2) {
	printf("Usage: %s nth\n", argv[0]);
	exit(1);
    }

    uint64_t nth  = strtoull(argv[1], NULL, 10);

    uint64_t rangeStart = 2UL;
    uint64_t rangeEnd   = primeUnderestimate(nth);

    assert(rangeStart <= rangeEnd);
    assert(rangeStart > 0);
    assert(rangeEnd > 1);
    
    uint64_t result = 0;
    
    uint64_t pieces   = 10;
    printf("+ Will break counting into %" PRIu64 " pieces...\n", pieces);

    uint64_t range    = rangeEnd - rangeStart;
    uint64_t step     = (range + pieces -1) / pieces;   // round up
    
    uint64_t count = 0;

    for (uint64_t r = rangeStart; r < rangeEnd; r += step) {
        uint64_t rnext = r+step;
        if (rnext > rangeEnd)
            rnext = rangeEnd;
        count += countPrimes(r, rnext, 0, NULL);
    }

    printf("+ Counted %" PRIu64 " noncandidate primes, now finishing up...\n",
           count);
    
    assert(count < nth);
    
    rangeStart = rangeEnd;
    rangeEnd   = primeOverestimate(nth);

    count += countPrimes(rangeStart, rangeEnd, nth-count, &result);
    
    assert(result != 0);

    printf("Prime #%" PRIu64 " = %" PRIu64 "\n", nth, result);
    printf("\t(%" PRIu64 " primes calculated)\n", count);
           
    return 0;
}
