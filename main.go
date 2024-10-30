package main

import (
	"bufio"
	"context"
	"encoding/json"
	// "fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	skip_raw := os.Getenv("SKIP")
	skip, _ := strconv.Atoi(skip_raw)

	if err := os.MkdirAll(filepath.Join("results"), 0755); err != nil {
		panic(err)
	}

	domainQueue := make(chan string, 1024)

	// start reading from file
	go func() {
		f, err := os.Open("./cloudflare-radar_top-100000-domains_20241021-20241028.csv")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)

		scanner.Scan() // skip first line

		// skip more lines
		for i := 0; i < skip; i++ {
			scanner.Scan()
		}

		for scanner.Scan() {
			line := scanner.Text()
			domainQueue <- strings.TrimSpace(line)
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	}()

	// print stats
	var processedDomainsCount atomic.Uint64
	go func() {
		t := time.Tick(15 * time.Second)
		for range t {
			log.Printf("processed %v domains", processedDomainsCount.Load())
		}
	}()

	// start workers
	numWorkers := runtime.NumCPU() * 2
	log.Printf("Starting %v workers ...", numWorkers)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for domain := range domainQueue {
				processDomain("https://" + strings.TrimSpace(domain))
				processedDomainsCount.Add(1)
			}
		}()
	}

	wg.Wait()
	log.Println("Done")
}

func processDomain(url string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "./dump-signedexchange", "-json", "-verify", "-uri", url)
	out, err := cmd.CombinedOutput()

	if err == nil && len(out) > 0 {
		data := make(map[string]any)
		if err := json.Unmarshal(out, &data); err != nil {
			log.Println(url, err)
			return
		}

		if data["Valid"].(bool) {
			if err := os.WriteFile(filepath.Join("results", sanitizePathElement(url)), out, 0644); err != nil {
				log.Println(url, err)
			}
			return
		}
	}

	// fmt.Fprintf(os.Stdout, "ERROR %v: %v: %s", url, err, out)
}

func sanitizePathElement(input string) string {
	unsafeChars := regexp.MustCompile(`[^a-zA-Z0-9_\-]`)
	sanitized := unsafeChars.ReplaceAllString(input, "_")
	sanitized = strings.TrimSpace(sanitized)
	return sanitized
}

