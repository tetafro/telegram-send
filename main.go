// Tiny CLI tool for sending messages to Telegram.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Version is the current git tag, injected on build.
var Version = "devel"

type config struct {
	chat  string
	token string
}

func main() {
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Println(Version)
		return
	}

	conf, err := readConfig()
	if err != nil {
		fatalf("Failed to read config: %w", err)
	}

	if len(os.Args) == 1 {
		fatalf("Message text is missing")
	}
	text := os.Args[1]

	if err := sendRequest(conf, text); err != nil {
		fatalf("Telegram API request failed: %w", err)
	}
	fmt.Println("Success")
}

func readConfig() (config, error) {
	exec, err := os.Executable()
	if err != nil {
		fatalf("Failed to get exec path: %w", err)
	}
	file := filepath.Join(filepath.Dir(exec), "telegram-send.ini")

	data, err := os.ReadFile(file) //nolint:gosec
	if err != nil {
		return config{}, fmt.Errorf("read file: %w", err)
	}

	var conf config
	for _, line := range strings.Split(string(data), "\n") {
		// Ignore empty lines and comments
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "chat":
			conf.chat = value
		case "token":
			conf.token = value
		default:
			return config{}, fmt.Errorf("unknown parameter: %s", key)
		}
	}
	return conf, nil
}

func sendRequest(conf config, text string) error {
	addr := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", conf.token)
	data := url.Values{}
	data.Set("chat_id", conf.chat)
	data.Set("text", text)

	//nolint:gosec,noctx
	resp, err := http.Post(
		addr,
		"application/x-www-form-urlencoded",
		bytes.NewBufferString(data.Encode()),
	)
	if err != nil {
		fatalf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		msg := fmt.Sprintf("invalid response code: %d", resp.StatusCode)
		if body, _ := io.ReadAll(resp.Body); len(body) != 0 {
			msg += fmt.Sprintf(", body: %s", string(body))
		}
		return errors.New(msg)
	}
	return nil
}

func fatalf(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
	os.Exit(1)
}
