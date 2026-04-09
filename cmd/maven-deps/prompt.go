package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func Confirm(w io.Writer, r io.Reader, msg string) (bool, error) {
	if _, err := fmt.Fprintf(w, "%s [y/N]: ", msg); err != nil {
		return false, err
	}
	reader := bufio.NewReader(r)
	line, err := reader.ReadString('\n')
	if err != nil && line == "" {
		return false, nil
	}
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
		return true, nil
	default:
		return false, nil
	}
}
