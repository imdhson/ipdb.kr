## 2024-05-23 - Prevent Log Injection and XSS via Spoofed IP Headers
**Vulnerability:** The application was reading user IP addresses from `CF-Connecting-IP` and `X-Forwarded-For` headers and storing/displaying them without any validation. Since HTTP headers can be easily spoofed, a malicious user could inject arbitrary text, scripts, or newline characters.
**Learning:** This could lead to Log Injection (corrupting `iphistory.txt`) and potentially Cross-Site Scripting (XSS) if the template escaping wasn't perfectly secure or used elsewhere without escaping. The application was trusting header input completely.
**Prevention:** Always validate parsed network data using established standard library functions like `net.ParseIP` to ensure it only contains structurally valid IP addresses before storing or using the data.
## 2026-10-04 - Broken Access Control in History Removal
**Vulnerability:** The '/sendremove' endpoint called 'clearHistory()' which truncated the entire 'iphistory.txt' file, allowing any user to delete the IP history of all users without authorization.
**Learning:** Insecure Direct Object Reference (IDOR) can occur when an action applies globally instead of filtering by the user's identity/scope.
**Prevention:** Ensure that state-modifying actions (like deletions) validate the requester's context (e.g. matching their IP to the entries being deleted) to limit the impact of the action to only data owned by the requester.
