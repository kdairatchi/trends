package trends

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// RSS Structs
type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title   string `xml:"title"`
	GUID    string `xml:"guid"`
	PubDate string `xml:"pubDate"`
}

// Constants
const maxTitleLength = 65

// Fetch RSS Feed with Timeout
func fetchRSSFeed(url string) (*RSS, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error fetching URL %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Received status %d from %s", resp.StatusCode, url)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response from %s: %v", url, err)
	}

	var rss RSS
	err = xml.Unmarshal(data, &rss)
	if err != nil {
		return nil, fmt.Errorf("Error parsing XML from %s: %v", url, err)
	}

	return &rss, nil
}

// Extract Feed Name from URL
func extractFeedName(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

// Check if a Date is Today
func isToday(pubDate, currentDate string) string {
	pubTime, err := time.Parse(time.RFC1123, pubDate)
	if err != nil {
		fmt.Printf("Error parsing date %s: %v\n", pubDate, err)
		return ""
	}
	if pubTime.Format("Mon, 02 Jan 2006") == currentDate {
		return "Yes"
	}
	return ""
}

// Sanitize Titles for Markdown
func sanitizeTitle(title string) string {
	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.ReplaceAll(title, "\r", " ")
	title = strings.ReplaceAll(title, "|", "\\|")
	title = strings.ReplaceAll(title, "[", "\\[")
	title = strings.ReplaceAll(title, "]", "\\]")
	if len(title) > maxTitleLength {
		title = title[:maxTitleLength] + "..."
	}
	return title
}

// Write Report to Markdown
func WriteReport(entries []map[string]string) {
	file, err := os.Create("trends.md")
	if err != nil {
		fmt.Println("Error creating trends.md:", err)
		return
	}
	defer file.Close()

	fmt.Fprintln(file, "| Time | Title | Feed | IsNew | IsToday |")
	fmt.Fprintln(file, "|------|-------|------|-------|--------|")

	for _, entry := range entries {
		title := sanitizeTitle(entry["title"])
		fmt.Fprintf(file, "| %s | [%s](%s) | %s | %s | %s |\n",
			entry["pubDate"], title, entry["guid"], entry["feeds"], entry["isNew"], entry["isToday"])
	}
}

// Core Function to Fetch & Process Trends
func TrackTrends() {
	urls := []string{
		"https://medium.com/feed/tag/bug-bounty",
		"https://medium.com/feed/tag/security",
		"https://medium.com/feed/tag/vulnerability",
	}

	entries := make([]map[string]string, 0)

	for _, url := range urls {
		rss, err := fetchRSSFeed(url)
		if err != nil {
			fmt.Println(err)
			continue
		}

		feedName := extractFeedName(url)
		for _, item := range rss.Channel.Items {
			entry := map[string]string{
				"title":   item.Title,
				"guid":    item.GUID,
				"pubDate": item.PubDate,
				"feeds":   fmt.Sprintf("[%s](%s)", feedName, url),
				"isNew":   "Yes",
				"isToday": isToday(item.PubDate, time.Now().Format("Mon, 02 Jan 2006")),
			}
			entries = append(entries, entry)
		}
		time.Sleep(2 * time.Second)
	}

	WriteReport(entries)
}
