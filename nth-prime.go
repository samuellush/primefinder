// Compile with:
//
// To compile (on Macs or Knuth)
//    make nth-prime-go
//
// Run with
//    ./nth-prime-go N
//
// to find the Nth prime

package main

import (
    "fmt"
    "math"
    "os"
    "strconv"
    "runtime"
)

// str = ordinal(n)
//   -- Function for putting cute "st", "nd", "rd", or "th" on numbers.
func ordinal(num uint64) string {
    digit := math.Mod(float64(num), 10)
    switch digit {
    case 1:
        return "st"
    case 2:
        return "nd"
    case 3:
        return "rd"
    default:
    }
    return "th"
}

// p = primeUnderestimate(nth)
//   -- Returns a number p <= the nth prime.

func primeUnderestimate(n float64) float64 {
    if n < 5.0 {
        return 2.0
    }
    if n < 6077.0 {
        return n * (math.Log(n) + math.Log(math.Log(n)) - 1.0)
    }
    return n * (math.Log(n) + math.Log(math.Log(n)) - 0.9718) // 0.97174375151226449448
}

// p = primeOverestimate(nth)
//   -- Returns a number p > the nth prime.

func primeOverestimate(n float64) float64 {
    if n < 6.0 {
        return 12.0
    }
    if n < 8602.0 { //  according to one source, it's 7022.0, but that's wrong
        return n*math.Log(n) + n*math.Log(math.Log(n))
    }
    return n * (math.Log(n) + math.Log(math.Log(n)) - 0.9385)
}

// count = primeCountOverestimate(n)
//   -- Returns an overestimate of the number of primes < n.

func primeCountOverestimate(n float64) float64 {
    if n < 2.0 {
        return 0.0
    }
    return 1.25506 * n / math.Log(n)
}

// t = primesTable(n, &c)
//   -- returns a dynamically allocated array of the first n primes,
//      and a count of how many there is stored in the table.
//      We employ the classic Sieve of Eratosthenes algorithm, with the
//      usual optimization of avoiding even numbers.
//      The code limits itself to primes < 2^32.

func primesTable(limit uint32) ([]uint32, uint32) {

    maxPrimes := uint32(primeCountOverestimate(float64(limit)))
    primeNum := make([]uint32, maxPrimes)
    numPrimes := uint32(0)

    // The primes we'll be using will all be less than 2^32 because that's
    // sufficient to enumerate all primes up to 2^64.

    // Allocate our array and initially declare that everything is prime.
    isPrime := make([]bool, limit/2)
    for i := range isPrime {
        isPrime[i] = true
    }

    // never used, but there for completeness and proper numbering
    primeNum[numPrimes] = 2
    numPrimes++

    for n := uint32(3); n < limit; n += 2 {
        if isPrime[n/2] == false {
            continue
        }
        primeNum[numPrimes] = n
        numPrimes++

        if n >= (1 << 16) { // if n*n would overflow a 32 bit int, we're
            continue // already done with crossing out, and the
        }   // test might not catch it

        // Cross off all the composites
        for i := n * n; i < limit; i += 2 * n {
            isPrime[i/2] = false
        }
    }

    return primeNum, numPrimes
}

// count = countPrimes(from, to, nth, &nth_prime)
//    -- Counts how many primes there are in the range [from,to); if
//       nth_prime is not nil, it will also store the nth prime it found
//       in *nth_prime (if there aren't that many, it'll store nothing)

func countPrimes(rangeStart uint64, rangeEnd uint64,
    nth uint64) (uint64, uint64) {
    thePrime := uint64(0)
    if nth != 0 {
        fmt.Printf("- Counting the primes in %d..%d and finding the %d%s ...\n",
            rangeStart, rangeEnd, nth, ordinal(nth))
    } else {
        fmt.Printf("- Counting the primes in %d..%d ...\n",
            rangeStart, rangeEnd)
    }

    crossingMax := math.Ceil(math.Sqrt(float64(rangeEnd)))
    crossingPrimes, crossingPrimesCount := primesTable(uint32(crossingMax))

    fmt.Printf("\tUsing %d primes (%d..%d) for crossing off\n",
        crossingPrimesCount, crossingPrimes[0],
        crossingPrimes[crossingPrimesCount-1])

    primeCount := uint64(0)

    // Handling 2 is a pain, but necessary...
    if rangeStart < 3 {
        primeCount += 1
        nth -= 1
        if nth == 0 {
            thePrime = 2
        }
        rangeStart = 3
    }

    // Adjust the start of the range up to be an odd number.
    rangeStart = rangeStart + 1 - uint64(math.Mod(float64(rangeStart), 2))

    // Bail if ther's nothing to do
    if rangeStart >= rangeEnd {
        return primeCount, thePrime
    }

    // Create our array for crossing off...
    thisRange := rangeEnd - rangeStart
    isPrime := make([]bool, (thisRange+1)/2)
    for i := range isPrime {
        isPrime[i] = true
    }

    // Cross off all the composites...
    for i := uint32(1); i < crossingPrimesCount; i += 1 {
        prime := uint64(crossingPrimes[i])

        // We need to adjust the crossing-off prime to be the first prime
        // in the range.  First, move it into the range.
        rangeStartAdjusted := rangeStart
        rangeStartAdjusted += uint64(math.Mod((float64(prime) -
            math.Mod(float64(rangeStart), float64(prime))), float64(prime)))

        if math.Mod(float64(rangeStartAdjusted), 2) == 0 { // It can't be even.
            rangeStartAdjusted += prime
        }
        if rangeStartAdjusted < prime*prime { // Only cross off the things
            rangeStartAdjusted = prime * prime // we need to.
        }

        // Actually cross off the composites.
        for j := rangeStartAdjusted; j < rangeEnd; j += 2 * prime {
            isPrime[(j-rangeStart)/2] = false
        }
    }

    // Now, actually count how many primes we're left with...
    for j := rangeStart; j < rangeEnd; j += 2 {
        if isPrime[(j-rangeStart)/2] {
            primeCount += 1
            nth -= 1
            if nth == 0 {
                thePrime = j
            }
        }
    }

    fmt.Printf("\t%d primes found\n", primeCount)

    return primeCount, thePrime
}

func uintMax(a, b uint64) uint64 {
    if a < b {
        return b
    }
    return a
}

func main() {
    argv := os.Args
    argc := len(argv)

    if argc != 2 {
        fmt.Printf("Usage: %s nth\n", argv[0])
        os.Exit(1)
    }

    nth, _ := strconv.ParseUint(argv[1], 10, 64)

    var rangeStart uint64 = 2
    rangeEnd := uint64(primeUnderestimate(float64(nth)))

    if rangeStart > rangeEnd {
        fmt.Printf("Assertion failed: %d..%d is not a valid range!\n",
            rangeStart, rangeEnd)
        os.Exit(-1)
    }
    if rangeStart <= 0 {
        fmt.Printf("Assertion failed: %d is not a valid starting point!\n",
            rangeStart)
        os.Exit(-1)
    }
    if rangeEnd <= 1 {
        fmt.Printf("Assertion failed: %d is not a valid end point!\n",
            rangeEnd)
        os.Exit(-1)
    }

    thisRange := rangeEnd - rangeStart

    pieces := uint64(runtime.GOMAXPROCS(runtime.NumCPU()))
    fmt.Printf("+ Will break counting into %d pieces...\n", pieces)

    step := uint64((thisRange + pieces - 1) / pieces) // round up

    count := uint64(0)
    countChannel := make (chan uint64)
    forCounter := uint64(0) // Supposedly there's a bug for certain values of N
                            // so we're going to run the second for-loop this
                            // many times rather than pieces

    for r := uint64(rangeStart); r < rangeEnd; r += step {
        rnext := r + step
        if rnext > rangeEnd {
            rnext = rangeEnd
        }
        go (func(countChannel1 chan uint64, r1 uint64, count1 uint64){
          numFound, _ := countPrimes(r1, rnext, 0)
          countChannel1 <- numFound
        })(countChannel, r, count)

        forCounter += uint64(1)
    }

    for i := uint64(0); i < forCounter; i += uint64(1){
      numFound1 := <- countChannel
      count += numFound1
    }


    fmt.Printf("+ Counted %d noncandidate primes, now finishing up...\n", count)

    if !(count < nth) {
        fmt.Printf("Assertion failed: we found %d primes instead of %d!\n", count, nth)
        os.Exit(-1)
    }

    rangeStart = rangeEnd
    rangeEnd = uint64(primeOverestimate(float64(nth)))

    fcount, result := countPrimes(rangeStart, rangeEnd, nth - count)
    count += fcount

    if !(result != 0) {
        fmt.Println("Assertion failed: result is zero!")
        os.Exit(-1)
    }

    fmt.Printf("Prime #%d = %d\n", nth, result)
    fmt.Printf("\t(%d primes calculated)\n", count)
}
