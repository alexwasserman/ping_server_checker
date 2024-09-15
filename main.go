package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-ping/ping"
)

func main() {
	// Define the -file and -quiet flags
	fileFlag := flag.String("file", "", "Path to file containing server hostnames, one per line")
	quietFlag := flag.Bool("quiet", false, "Suppress non-error output")
	flag.Parse()

	var servers []string
	var err error

	// If the -file flag is provided, read the list of servers from the file
	if *fileFlag != "" {
		servers, err = readServersFromFile(*fileFlag)
		if err != nil {
			log.Fatalf("Error reading servers from file: %v\n", err)
		}
	} else {
		// Otherwise, read the list of servers from stdin
		servers, err = readServersFromStdin()
		if err != nil {
			log.Fatalf("Error reading servers from stdin: %v\n", err)
		}
	}

	// Log unresponsive servers
	unresponsiveServers := []string{}

	for _, server := range servers {
		if !*quietFlag {
			fmt.Printf("Pinging server: %s\n", server)
		}
		if !pingServer(server, *quietFlag) {
			unresponsiveServers = append(unresponsiveServers, server)
			logUnresponsiveServer(server)
		}
	}

	if len(unresponsiveServers) > 0 {
		if !*quietFlag {
			fmt.Println("Some servers are unresponsive. Check logs for details.")
		}
	} else if !*quietFlag {
		fmt.Println("All servers are responsive.")
	}
}

// pingServer pings a server and returns true if it's responsive, false otherwise
// The quietFlag controls whether non-error output should be suppressed
func pingServer(server string, quiet bool) bool {
	pinger, err := ping.NewPinger(server)
	if err != nil {
		log.Printf("Failed to create pinger for server %s: %v\n", server, err)
		return false
	}
	pinger.Count = 3
	pinger.Timeout = time.Second * 5

	err = pinger.Run()
	if err != nil {
		log.Printf("Ping failed for server %s: %v\n", server, err)
		return false
	}

	stats := pinger.Statistics()
	if stats.PacketsRecv == 0 {
		log.Printf("No response from server %s\n", server)
		return false
	}

	if !quiet {
		fmt.Printf("Server %s is responsive\n", server)
	}
	return true
}

// logUnresponsiveServer logs an alert when a server is unresponsive
func logUnresponsiveServer(server string) {
	// Log the unresponsive server to a file
	f, err := os.OpenFile("unresponsive_servers.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer f.Close()

	logger := log.New(f, "ALERT: ", log.LstdFlags)
	logger.Printf("Server %s is unresponsive\n", server)
}

// readServersFromFile reads server hostnames from a file
func readServersFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var servers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			servers = append(servers, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return servers, nil
}

// readServersFromStdin reads server hostnames from stdin
func readServersFromStdin() ([]string, error) {
	fmt.Println("Enter server hostnames, one per line (Press Ctrl+D when done):")
	scanner := bufio.NewScanner(os.Stdin)

	var servers []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			servers = append(servers, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading from stdin: %w", err)
	}

	return servers, nil
}
