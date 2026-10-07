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
