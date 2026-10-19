## 2023-10-25 - Optimize getIP Function
**Learning:** Using standard library functions like `strings.Cut` instead of `strings.Split` helps avoid unnecessary slice allocations when retrieving split substrings. Replacing custom port parsing on IP addresses using `net.SplitHostPort` handles bracketed IPv6 and IPv4 out-of-the-box natively and efficiently.
**Action:** In Go services reading repetitive fields per HTTP request, ensure memory string allocations are minimized by leveraging string index cutting natively when possible.

## 2024-10-04 - Avoid string slice allocations in hot paths
**Learning:** Using `strings.Split` or `strings.SplitN` allocates string slices, which creates unnecessary garbage and CPU overhead when only parsing one or two values. In the Go HTTP handlers, replacing `strings.SplitN(line, "|", 2)` with `strings.Cut(line, "|")` and `strings.Split(ip, ",")` with `strings.IndexByte(ip, ',')` reduced string manipulation benchmark times by ~85% in `Cut` vs `SplitN` and ~88% in `IndexByte` vs `Split` (avoiding slice heap allocations entirely).
**Action:** When extracting the first component or dividing strings into simple pairs, prefer `strings.Cut` or `strings.IndexByte` + string slicing over `strings.Split` to avoid unnecessary memory allocations.

## 2024-10-05 - Avoid fmt.Sprintf for simple string concatenation
**Learning:** Using `fmt.Sprintf` for simple string building incurs significant runtime reflection and parsing overhead compared to direct string concatenation. Benchmarks showed concatenation is >3x faster (~70ns vs ~230ns) for simple string joining like `dateStr + "|" + ip + "\n"`.
**Action:** Prefer direct string concatenation using `+` over `fmt.Sprintf` when joining a few strings together, especially in hot paths like logging, formatting, or file writing logic.

## 2024-05-24 - Optimize File Reading for Tiny Log Files
**Learning:** `bufio.Scanner` is an anti-pattern for reading very small, artificially constrained files (like `iphistory.txt` which is capped at 1024 bytes). The 4KB buffer allocation in `bufio` is overkill and dominates the performance.
**Action:** When working with explicit small file caps (< 4KB), prefer reading the entire file into memory with `os.ReadFile` and parsing with `strings.Cut` over allocating readers/scanners. This reduces allocations dramatically (from ~30 down to ~6) and improves execution time. Remember to handle `\r` carriage returns.
## 2024-10-07 - [Zero-allocation time formatting]
**Learning:** Using `time.Now().AppendFormat` onto a fixed-size stack array (e.g. `var buf [64]byte`) eliminates heap allocations entirely for time formatting and string concatenation in Go, offering a ~35% speedup.
**Action:** Always prefer `time.AppendFormat` over `time.Format` and direct slice appends over string concatenation using `+` when the max length of the output string is predictable and small, especially in hot paths.

## 2026-10-08 - Fast Line-by-Line Processing of Small Files
**Learning:** For small files (e.g. < 1024 bytes), parsing line-by-line with `bufio.Scanner` is unnecessarily expensive in Go due to multiple tiny allocations per loop.
**Action:** Use `os.ReadFile` with `strings.Cut(s, "\n")` and `strings.Builder` on the full string for processing without multiple allocations per line. This reduces allocations dramatically and improves execution speed by ~10x.

## 2023-10-09 - Pre-allocate slice capacity for text file lines
**Learning:** Pre-allocating slice capacity by using `strings.Count(s, "\n")` before parsing lines of text can significantly reduce heap allocations (from ~8 to ~2) and improve speed (by ~35%), even for small files loaded entirely into memory.
**Action:** Always count lines or items upfront to determine required slice capacity before looping and appending, particularly when parsing delimited strings or files.

## 2023-10-25 - [In-place byte slice filtering panic]
**Learning:** When performing in-place byte slice filtering/modification (e.g., dropping lines from a file loaded into a `[]byte`), directly indexing the slice via `data[outLen] = '\n'` is dangerous. If the final length equals the original capacity (e.g., when retaining a final line that lacked a trailing newline in the source), this manual assignment will cause an index out of bounds panic.
**Action:** Always use `append` for in-place slice modifications (e.g., `data = append(data[:outLen], line...)` followed by `data = append(data, '\n')`). `append` safely handles length boundaries and will automatically expand capacity if ever necessary, preventing out-of-bounds panics while remaining allocation-free if capacity permits.

## 2024-10-10 - Eliminate heap allocations in static HTTP responses
**Learning:** Using `w.Write([]byte("constant string"))` in an `http.ResponseWriter` forces a heap allocation for the byte slice conversion. Replacing this with `io.WriteString(w, "constant string")` eliminates the allocation entirely (drops from 1 to 0 per call) and reduces execution time significantly (~70% faster in microbenchmarks).
**Action:** When writing static string content directly to an `http.ResponseWriter`, always prefer `io.WriteString` over `w.Write([]byte(...))` to minimize heap allocations and improve performance.
## 2026-10-11 - [Use netip for allocation-free IP parsing]
 **Learning:** Using `net.ParseIP` and `net.SplitHostPort` allocates strings and slices on the heap. Trying to optimize `net.SplitHostPort` manually with `strings.LastIndexByte` is dangerous for IPv6 parsing and harms readability. The `net/netip` package (introduced in Go 1.18) provides `netip.ParseAddrPort` and `netip.ParseAddr` which solve this by parsing addresses without heap allocations while remaining safe and idiomatic.
 **Action:** Prefer `net/netip` over `net.ParseIP` and manual string manipulation for all new IP parsing code.

## 2024-10-13 - [Zero-allocation IP port parsing]
**Learning:** Using `strings.Split(ip, ":")` to split an IP and a port forces a heap allocation (e.g. 32 bytes) because it allocates a new slice. By rearranging the logic to use `strings.Cut(ip, ":")` when the input is not an IPv6 string, we can eliminate the slice allocation entirely and reduce processing time significantly (from ~70ns down to ~10ns).
**Action:** Prefer `strings.Cut` over `strings.Split` whenever only the first delimited element is required or when processing simple two-part strings. Avoid blindly calling `strings.Split` before ensuring the split parts will actually be used.

## 2023-10-27 - [Zero-allocation string conversion for read-only parsing]
**Learning:** When parsing small files (like `iphistory.txt`) into structs where strings only need to be sequentially read (using `strings.Count`, `strings.Cut` etc.), standard `string(bytes)` casts cause unnecessary heap allocations by copying the byte slice.
**Action:** Use `unsafe.String(unsafe.SliceData(data), len(data))` instead of `string(data)` to perform zero-allocation string casts when the underlying byte slice will not be modified during parsing.

## $(date +%Y-%m-%d) - [Zero-allocation URL host extraction for CSRF validation]
**Learning:** `url.Parse` is surprisingly heavy for simple operations in hot paths (like middleware or CSRF checks), allocating ~144 bytes and taking ~470ns per call to simply extract a hostname from `Origin` or `Referer` headers.
**Action:** Use string slicing with `strings.Index`, `strings.IndexAny`, and `strings.IndexByte` to manually extract the host for simple URL validations. It reduces allocations to 0 and cuts execution time down to ~47ns (10x faster).
