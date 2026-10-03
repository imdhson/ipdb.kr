## 2023-10-25 - Optimize getIP Function
**Learning:** Using standard library functions like `strings.Cut` instead of `strings.Split` helps avoid unnecessary slice allocations when retrieving split substrings. Replacing custom port parsing on IP addresses using `net.SplitHostPort` handles bracketed IPv6 and IPv4 out-of-the-box natively and efficiently.
**Action:** In Go services reading repetitive fields per HTTP request, ensure memory string allocations are minimized by leveraging string index cutting natively when possible.
