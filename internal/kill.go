package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func ParseCSVInts(s string) []int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var out []int
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if v, err := strconv.Atoi(p); err == nil {
			out = append(out, v)
		}
	}
	return out
}

// KillByPorts finds PIDs for given ports and kills them. Returns status messages.
func KillByPorts(entries []PortEntry, ports []string) []string {
	var pids []int
	for _, port := range ports {
		port = strings.TrimSpace(port)
		if port == "" {
			continue
		}
		for _, e := range entries {
			if e.LocalPort == port && e.PID > 0 {
				pids = append(pids, e.PID)
			}
		}
	}
	if len(pids) == 0 {
		return []string{"no matching PIDs for given ports"}
	}
	return KillPIDs(pids)
}

// KillPIDs sends SIGTERM to each PID (unique). Returns status messages.
func KillPIDs(pids []int) []string {
	seen := make(map[int]struct{})
	var msgs []string

	for _, pid := range pids {
		if pid <= 0 {
			continue
		}
		if _, ok := seen[pid]; ok {
			continue
		}
		seen[pid] = struct{}{}

		prefix := fmt.Sprintf("PID %d:", pid)
		proc, err := os.FindProcess(pid)
		if err != nil {
			msgs = append(msgs, prefix+" "+err.Error())
			continue
		}
		if err := proc.Signal(syscall.SIGTERM); err != nil {
			msgs = append(msgs, prefix+" SIGTERM failed: "+err.Error())
		} else {
			msgs = append(msgs, prefix+" terminated")
		}
	}

	if len(msgs) == 0 {
		msgs = []string{"no valid PIDs to kill"}
	}
	return msgs
}
