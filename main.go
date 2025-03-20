package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Result struct {
	IP           string `json:"ip"`
	StartAddress string `json:"start_address"`
	EndAddress   string `json:"end_address"`
	OrgName      string `json:"organization"`
	Country      string `json:"country"`
	City         string `json:"city"`
	Holder       string `json:"holder"`
}

func fetchRDAP(ctx context.Context, url string, resultChan chan<- map[string]interface{}) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return
	}

	select {
	case resultChan <- data:
	case <-ctx.Done():
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <IP>")
		return
	}

	ip := os.Args[1]

	rdapURLs := []string{
		fmt.Sprintf("https://rdap.arin.net/registry/ip/%s", ip),
		fmt.Sprintf("https://rdap.db.ripe.net/ip/%s", ip),
		fmt.Sprintf("https://rdap.lacnic.net/rdap/ip/%s", ip),
		fmt.Sprintf("https://rdap.apnic.net/ip/%s", ip),
		fmt.Sprintf("https://rdap.afrinic.net/rdap/ip/%s", ip),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultChan := make(chan map[string]interface{})

	for _, url := range rdapURLs {
		go fetchRDAP(ctx, url, resultChan)
	}

	select {
	case rdapData := <-resultChan:
		cancel()
		// fmt.Println("IP found in RDAP servers:", ip, "\n", rdapData)
		output := parseRDAP(rdapData, ip)
		jsonOut, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(jsonOut))
	case <-time.After(5 * time.Second):
		fmt.Println(`{"error": "Timeout or IP not found in RDAP servers."}`)
	}
}

func parseRDAP(rdapData map[string]interface{}, ip string) Result {
	result := Result{IP: ip}

	if start, ok := rdapData["startAddress"].(string); ok {
		result.StartAddress = start
	}
	if end, ok := rdapData["endAddress"].(string); ok {
		result.EndAddress = end
	}
	// Fallback country at root (some servers provide this)
	if country, ok := rdapData["country"].(string); ok {
		result.Country = country
	}

	if entities, ok := rdapData["entities"].([]interface{}); ok {
		for _, e := range entities {
			entity := e.(map[string]interface{})

			// Try to find holder handle
			if handle, ok := entity["handle"].(string); ok && result.Holder == "" {
				result.Holder = handle
			}

			// vCard parsing
			if vcardArray, ok := entity["vcardArray"].([]interface{}); ok && len(vcardArray) >= 2 {
				vcard := vcardArray[1].([]interface{})
				for _, item := range vcard {
					entry := item.([]interface{})
					key := entry[0].(string)

					if key == "fn" {
						if orgName, ok := entry[3].(string); ok && result.OrgName == "" {
							result.OrgName = orgName
						}
					}

					if key == "adr" {
						// Address field
						if adrMap, ok := entry[1].(map[string]interface{}); ok {
							// Try to extract label (full address block)
							if label, ok := adrMap["label"].(string); ok {
								lines := splitLines(label)
								// Usually: Street, City, Region, Zip, Country
								if len(lines) >= 2 {
									result.City = lines[len(lines)-3]    // e.g., Mountain View
									result.Country = lines[len(lines)-1] // e.g., United States
								}
							}
						}
					}
				}
			}
		}
	}

	return result
}

func splitLines(s string) []string {
	var lines []string
	current := ""
	for _, r := range s {
		if r == '\n' || r == '\r' {
			if current != "" {
				lines = append(lines, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}
