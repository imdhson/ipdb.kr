package main

import (
	"bufio"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const historyFile = "iphistory.txt"

var fileMutex sync.Mutex

type HistoryEntry struct {
	Date string
	IP   string
}

// readHistory reads the IP history from the file
func readHistory() ([]HistoryEntry, error) {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	var history []HistoryEntry

	file, err := os.Open(historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return history, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if date, ip, found := strings.Cut(line, "|"); found {
			history = append(history, HistoryEntry{
				Date: date,
				IP:   ip,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return history, nil
}

// appendHistory appends a new entry to the IP history file and truncates if it gets too large
func appendHistory(ip string) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// Check file size
	info, err := os.Stat(historyFile)
	if err == nil && info.Size() > 1024 {
		// Truncate file
		err = os.WriteFile(historyFile, []byte(""), 0644)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(historyFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	dateStr := time.Now().Format("2006-01-02 15:04:05")
	// Optimized: Replace fmt.Sprintf with faster string concatenation
	entry := dateStr + "|" + ip + "\n"
	_, err = file.WriteString(entry)
	return err
}

// clearUserHistory removes entries matching the given IP from the history file
// SECURITY: Prevents a user from deleting other users' history
func clearUserHistory(userIP string) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file, err := os.Open(historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Optimized: Avoid slice allocation in string parsing
		_, ipPart, found := strings.Cut(line, "|")
		// Only drop the line if it's a valid entry and the IP matches the user's IP.
		if !found || ipPart != userIP {
			lines = append(lines, line)
		}
	}
	file.Close()

	if err := scanner.Err(); err != nil {
		return err
	}

	// Re-write the file with the remaining lines
	content := ""
	if len(lines) > 0 {
		content = strings.Join(lines, "\n") + "\n"
	}
	return os.WriteFile(historyFile, []byte(content), 0644)
}

var (
	indexTmpl *template.Template
	termsTmpl *template.Template
)

func init() {
	indexTmpl = template.Must(template.ParseFiles("templates/index.html"))
	termsTmpl = template.Must(template.ParseFiles("templates/terms.html"))
}

func getIP(r *http.Request) string {
	// 1. Cloudflare IP header
	ip := r.Header.Get("CF-Connecting-IP")
	if ip != "" {
		if parsedIP := net.ParseIP(ip); parsedIP != nil {
			return parsedIP.String()
		}
	}

	// 2. X-Forwarded-For header
	// Optimized: Avoid strings.Split allocation when parsing X-Forwarded-For
	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		if idx := strings.IndexByte(ip, ','); idx != -1 {
			ip = strings.TrimSpace(ip[:idx])
		} else {
			ip = strings.TrimSpace(ip)
		}
		if parsedIP := net.ParseIP(ip); parsedIP != nil {
			return parsedIP.String()
		}
	}

	// 3. Fallback to RemoteAddr
	ip = r.RemoteAddr
	if host, _, err := net.SplitHostPort(ip); err == nil {
		ip = host
	} else {
		// Fallback for cases where port might not be present or format is odd
		if strings.Contains(ip, ":") {
			parts := strings.Split(ip, ":")
			if strings.HasPrefix(ip, "[") && strings.Contains(ip, "]:") {
				idx := strings.LastIndex(ip, ":")
				ip = strings.Trim(ip[:idx], "[]")
			} else if len(parts) == 2 {
				ip = parts[0]
			}
		}
	}

	if parsedIP := net.ParseIP(ip); parsedIP != nil {
		return parsedIP.String()
	}

	return "Unknown"
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/accept" {
		http.NotFound(w, r)
		return
	}

	cookie, err := r.Cookie("terms")
	if err != nil || cookie.Value != "accept" {
		http.Redirect(w, r, "/terms", http.StatusFound)
		return
	}

	ip := getIP(r)
	if err := appendHistory(ip); err != nil {
		log.Printf("Failed to append history: %v", err)
	}

	history, err := readHistory()
	if err != nil {
		log.Printf("Failed to read history: %v", err)
	}

	data := struct {
		IP      string
		History []HistoryEntry
	}{
		IP:      ip,
		History: history,
	}

	if err := indexTmpl.Execute(w, data); err != nil {
		log.Printf("Template execution failed: %v", err)
	}
}

func handleTerms(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r)
	data := struct {
		IP string
	}{
		IP: ip,
	}

	if err := termsTmpl.Execute(w, data); err != nil {
		log.Printf("Template execution failed: %v", err)
	}
}

func handleSendAccept(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "terms",
		Value:    "accept",
		Path:     "/",
		MaxAge:   60 * 60 * 24,
		HttpOnly: false,
	}
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<meta http-equiv='refresh' content='0;url=/' />\n<h1>Wait...</h1>"))
}

func handleSendReject(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "terms",
		Value:    "reject",
		Path:     "/",
		MaxAge:   60 * 60 * 24,
		HttpOnly: false,
	}
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<meta http-equiv='refresh' content='0;url=/' />\n<h1>Wait...</h1>"))
}

func handleSendRemove(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r)
	if err := clearUserHistory(ip); err != nil {
		log.Printf("Failed to clear history: %v", err)
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<meta http-equiv='refresh' content='0;url=/terms/sendreject' />\n<h1>Wait...</h1>"))
}

func main() {
	// Static files
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	// Routes
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/accept", handleIndex) // Handled same as root
	http.HandleFunc("/terms", handleTerms)
	http.HandleFunc("/terms/sendaccept", handleSendAccept)
	http.HandleFunc("/terms/sendreject", handleSendReject)
	http.HandleFunc("/sendremove", handleSendRemove)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
