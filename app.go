package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/vnykmshr/gowsay/api"
	"github.com/vnykmshr/gowsay/cow"
)

var version = "devel"

func main() {
	// Check if first argument is "serve"
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		runServer()
		return
	}

	// Run CLI mode
	runCLI()
}

func runCLI() {
	// Define CLI flags
	var (
		cowName    = flag.String("c", "default", "Cow name to use")
		mood       = flag.String("m", "", "Mood (borg, dead, greedy, paranoid, stoned, wired, young)")
		think      = flag.Bool("t", false, "Think instead of say")
		list       = flag.Bool("l", false, "List available cows and moods")
		random     = flag.Bool("r", false, "Random cow and mood")
		columns    = flag.Int("w", 40, "Column width for text wrapping")
		showVer    = flag.Bool("v", false, "Show version")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "gowsay %s - cowsay implementation in Go\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  gowsay [options] [message...]\n")
		fmt.Fprintf(os.Stderr, "  echo \"message\" | gowsay [options]\n")
		fmt.Fprintf(os.Stderr, "  gowsay serve                    Start HTTP server\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Show version
	if *showVer {
		fmt.Printf("gowsay %s\n", version)
		os.Exit(0)
	}

	// List cows and moods
	if *list {
		cows := cow.List()
		sort.Strings(cows)
		fmt.Println("Available cows:")
		for _, c := range cows {
			fmt.Printf("  %s\n", c)
		}
		fmt.Println("\nAvailable moods:")
		moods := cow.ListMoods()
		sort.Strings(moods)
		for _, m := range moods {
			fmt.Printf("  %s\n", m)
		}
		os.Exit(0)
	}

	// Get message text
	var text []string
	args := flag.Args()

	if len(args) > 0 {
		// Use command line arguments
		text = args
	} else {
		// Read from stdin
		text = readStdin()
		if len(text) == 0 {
			flag.Usage()
			os.Exit(1)
		}
	}

	// Apply random if requested
	if *random {
		*cowName = cow.RandomCow()
		*mood = cow.RandomMood()
	}

	// Validate cow exists
	if !cow.Exists(*cowName) {
		fmt.Fprintf(os.Stderr, "Error: cow '%s' not found\n", *cowName)
		os.Exit(1)
	}

	// Validate mood if specified
	if *mood != "" && !cow.MoodExists(*mood) {
		fmt.Fprintf(os.Stderr, "Error: mood '%s' not found\n", *mood)
		os.Exit(1)
	}

	// Determine action
	action := cow.ActionSay
	if *think {
		action = cow.ActionThink
	}

	// Render and output (without markdown code blocks for CLI)
	output := cow.Render(text, *cowName, *mood, action, *columns)
	// Strip markdown code blocks (```) for CLI output
	output = strings.TrimPrefix(output, "```\n")
	output = strings.TrimSuffix(output, "\n```\n")
	fmt.Print(output)
}

func readStdin() []string {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil
	}

	// Check if stdin is a pipe or has data
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		var lines []string
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		return lines
	}

	return nil
}

func runServer() {
	// Setup structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	m := api.NewModule()

	// Legacy Slack endpoint (backward compatibility)
	http.HandleFunc("/say", m.Gowsay)

	// New API endpoints (with CORS)
	http.HandleFunc("/api/moo", api.CORS(m.APIMoo))
	http.HandleFunc("/api/cows", api.CORS(m.APICows))
	http.HandleFunc("/api/moods", api.CORS(m.APIMoods))
	http.HandleFunc("/health", api.CORS(api.Health(version)))

	// Web UI - serve at root
	http.Handle("/", api.ServeWeb())

	fmt.Println(api.GetBanner(version))
	slog.Info("routes registered",
		"endpoints", []string{"/", "/say", "/api/moo", "/api/cows", "/api/moods", "/health"})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		slog.Info("starting server", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	<-stop
	slog.Info("shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
