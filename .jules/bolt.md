## 2023-10-25 - Optimize getIP Function
**Learning:** Using standard library functions like `strings.Cut` instead of `strings.Split` helps avoid unnecessary slice allocations when retrieving split substrings. Replacing custom port parsing on IP addresses using `net.SplitHostPort` handles bracketed IPv6 and IPv4 out-of-the-box natively and efficiently.
**Action:** In Go services reading repetitive fields per HTTP request, ensure memory string allocations are minimized by leveraging string index cutting natively when possible.

## 2024-10-04 - Avoid string slice allocations in hot paths
**Learning:** Using `strings.Split` or `strings.SplitN` allocates string slices, which creates unnecessary garbage and CPU overhead when only parsing one or two values. In the Go HTTP handlers, replacing `strings.SplitN(line, "|", 2)` with `strings.Cut(line, "|")` and `strings.Split(ip, ",")` with `strings.IndexByte(ip, ',')` reduced string manipulation benchmark times by ~85% in `Cut` vs `SplitN` and ~88% in `IndexByte` vs `Split` (avoiding slice heap allocations entirely).
**Action:** When extracting the first component or dividing strings into simple pairs, prefer `strings.Cut` or `strings.IndexByte` + string slicing over `strings.Split` to avoid unnecessary memory allocations.

## 2024-10-05 - Avoid fmt.Sprintf for simple string concatenation
**Learning:** Using `fmt.Sprintf` for simple string building incurs significant runtime reflection and parsing overhead compared to direct string concatenation. Benchmarks showed concatenation is >3x faster (~70ns vs ~230ns) for simple string joining like `dateStr + "|" + ip + "\n"`.
**Action:** Prefer direct string concatenation using `+` over `fmt.Sprintf` when joining a few strings together, especially in hot paths like logging, formatting, or file writing logic.
