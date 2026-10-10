## 2024-05-23 - Prevent Log Injection and XSS via Spoofed IP Headers
**Vulnerability:** The application was reading user IP addresses from `CF-Connecting-IP` and `X-Forwarded-For` headers and storing/displaying them without any validation. Since HTTP headers can be easily spoofed, a malicious user could inject arbitrary text, scripts, or newline characters.
**Learning:** This could lead to Log Injection (corrupting `iphistory.txt`) and potentially Cross-Site Scripting (XSS) if the template escaping wasn't perfectly secure or used elsewhere without escaping. The application was trusting header input completely.
**Prevention:** Always validate parsed network data using established standard library functions like `net.ParseIP` to ensure it only contains structurally valid IP addresses before storing or using the data.
## 2026-10-04 - Broken Access Control in History Removal
**Vulnerability:** The '/sendremove' endpoint called 'clearHistory()' which truncated the entire 'iphistory.txt' file, allowing any user to delete the IP history of all users without authorization.
**Learning:** Insecure Direct Object Reference (IDOR) can occur when an action applies globally instead of filtering by the user's identity/scope.
**Prevention:** Ensure that state-modifying actions (like deletions) validate the requester's context (e.g. matching their IP to the entries being deleted) to limit the impact of the action to only data owned by the requester.
## 2025-02-18 - [Secure Terms Cookie Implementation]
**Vulnerability:** Terms acceptance cookies were set with `HttpOnly: false` and missing `SameSite` configurations, making them vulnerable to XSS and CSRF attacks.
**Learning:** Found two places (`handleSendAccept`, `handleSendReject`) managing a simple string state ("accept"/"reject") via cookies to handle ToS acceptance.
**Prevention:** Make sure cookies default to `HttpOnly: true` and `SameSite: http.SameSiteStrictMode` especially for simple state cookies not required directly by Javascript logic.
## 2024-05-18 - Missing timeouts in Go HTTP servers
**Vulnerability:** Default Go HTTP server configurations do not enforce timeouts (`ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout`), making the application vulnerable to slow-connection DoS attacks (e.g., Slowloris). Additionally, basic security headers were missing.
**Learning:** `http.ListenAndServe()` in Go uses `http.Server{}` without any connection timeouts by default.
**Prevention:** Always initialize `http.Server{}` explicitly with appropriate timeouts and use a middleware to add fundamental security headers like `X-Content-Type-Options`, `X-Frame-Options`, `X-XSS-Protection`, and `Strict-Transport-Security`.
## 2023-10-25 - Prevent CSRF on State-Changing Endpoints via GET Requests
**Vulnerability:** The `/sendremove` endpoint, which clears a user's IP history, was accessible via an HTTP GET request. This allowed Cross-Site Request Forgery (CSRF) attacks, where an attacker could trick a user into navigating to the endpoint (e.g., via a hidden image tag or malicious link) to delete their IP history without their consent.
**Learning:** Performing state-changing operations via GET requests violates HTTP semantics and bypasses basic protections that browsers and web frameworks apply to POST requests. In unauthenticated (or IP-based) systems, this makes such endpoints easily exploitable.
**Prevention:** Always require `http.MethodPost` (or PUT/DELETE) for endpoints that modify or delete data, and use an HTML `<form>` or XHR/fetch to submit to it.

## 2024-10-08 - IP Spoofing via Proxy Headers
**Vulnerability:** The application unconditionally trusted `CF-Connecting-IP` and `X-Forwarded-For` headers to determine the client's IP address. This allows any external client to easily spoof their IP address by injecting these headers into their requests.
**Learning:** Blindly trusting proxy headers without validating the source of the connection (e.g., checking if it's from a trusted load balancer or local proxy) is a critical security vulnerability.
**Prevention:** Always parse `r.RemoteAddr` first to verify the direct connection source. Only read proxy headers if `r.RemoteAddr` corresponds to a trusted proxy IP (e.g., local, loopback, or internal private IPs).

## 2024-10-24 - Prevent CSRF on Terms Acceptance via GET Requests
**Vulnerability:** The `/terms/sendaccept` and `/terms/sendreject` endpoints, which set terms acceptance cookies, were accessible via HTTP GET requests. This allowed Cross-Site Request Forgery (CSRF) attacks where an attacker could trick a user into accepting or rejecting terms without their consent.
**Learning:** State-changing operations via GET requests are inherently vulnerable to CSRF. While a previous entry documented this for `/sendremove`, other routes were missed, highlighting the need to audit *all* endpoints that modify state (including cookies).
**Prevention:** Ensure that all endpoints modifying application state or setting critical cookies enforce HTTP POST (or other appropriate state-modifying methods) and are integrated with proper frontend HTML forms or XHR submissions.

## 2024-10-25 - IP Spoofing via left-most X-Forwarded-For parsing
**Vulnerability:** Even when behind a trusted proxy, parsing the `X-Forwarded-For` header by taking the left-most IP address blindly trusts user input. If a user sets a spoofed XFF header, the trusted proxy will append the real IP, but the application would still read the spoofed left-most IP.
**Learning:** The left-most entries in `X-Forwarded-For` are controlled by the client and cannot be trusted. Only the right-most entries appended by trusted proxies are reliable.
**Prevention:** When parsing `X-Forwarded-For`, always iterate from right-to-left, checking each IP. Stop and use the first IP that is NOT part of your trusted proxy infrastructure (e.g., not an internal/private IP).
