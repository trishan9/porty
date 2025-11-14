package internal

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

type PortEntry struct {
	Proto       string `json:"proto"`
	State       string `json:"state"`
	LocalAddr   string `json:"local_addr"`
	LocalPort   string `json:"local_port"`
	PID         int    `json:"pid"`
	ProcessName string `json:"process"`
	UserName    string `json:"user"`
	Tag         string `json:"tag"` // USER / SYSTEM / UNKNOWN / SELF
}

// ListPorts scans /proc for TCP/UDP sockets and maps them to processes.
func ListPorts() ([]PortEntry, error) {
	inodeToPID := buildInodePIDMap()

	curUser, _ := user.Current()
	curUID := ""
	if curUser != nil {
		curUID = curUser.Uid
	}

	var entries []PortEntry

	// tcp / tcp6
	entries = append(entries, parseNetFile("/proc/net/tcp", "tcp", inodeToPID, curUID)...)
	entries = append(entries, parseNetFile("/proc/net/tcp6", "tcp", inodeToPID, curUID)...)

	// udp / udp6
	entries = append(entries, parseNetFile("/proc/net/udp", "udp", inodeToPID, curUID)...)
	entries = append(entries, parseNetFile("/proc/net/udp6", "udp", inodeToPID, curUID)...)

	return entries, nil
}

// ------------------------------------------------------------
// /proc/<pid>/fd -> socket inode -> pid map
// ------------------------------------------------------------

func buildInodePIDMap() map[string]int {
	result := make(map[string]int)

	procEntries, err := os.ReadDir("/proc")
	if err != nil {
		return result
	}

	for _, e := range procEntries {
		if !e.IsDir() {
			continue
		}
		pidStr := e.Name()
		pid, err := strconv.Atoi(pidStr)
		if err != nil || pid <= 0 {
			continue
		}

		fdDir := filepath.Join("/proc", pidStr, "fd")
		fdEntries, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}

		for _, fd := range fdEntries {
			link, err := os.Readlink(filepath.Join(fdDir, fd.Name()))
			if err != nil {
				continue
			}
			// socket:[12345]
			if strings.HasPrefix(link, "socket:[") && strings.HasSuffix(link, "]") {
				inode := link[len("socket:[") : len(link)-1]
				result[inode] = pid
			}
		}
	}

	return result
}

// ------------------------------------------------------------
// /proc/net/{tcp,udp} parsing
// ------------------------------------------------------------

func parseNetFile(path, proto string, inodeToPID map[string]int, curUID string) []PortEntry {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	isIPv6 := strings.HasSuffix(path, "6")

	var entries []PortEntry

	sc := bufio.NewScanner(file)
	firstLine := true
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		// skip header
		if firstLine {
			firstLine = false
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		localField := fields[1] // local_address
		stateHex := fields[3]   // hex state
		inode := fields[9]      // inode

		localAddr, localPort := parseIPPort(localField, isIPv6)

		state := decodeState(proto, stateHex)

		// We care mostly about listening / unconnected (like btop).
		if proto == "tcp" && state != "LISTEN" {
			continue
		}

		pid := inodeToPID[inode]

		// -------------------------
		// Kernel-owned sockets:
		// inode is present but no PID maps to it
		// -------------------------
		if pid == 0 {
			entries = append(entries, PortEntry{
				Proto:       proto,
				State:       state,
				LocalAddr:   localAddr,
				LocalPort:   localPort,
				PID:         0,
				ProcessName: "<kernel>",
				UserName:    "kernel",
				Tag:         "KERNEL",
			})
			continue
		}

		pname := getProcessNameFromPID(pid)
		uname, uid := getUserFromPID(pid)
		tag := classifyEntry(uid, curUID, pid)

		entries = append(entries, PortEntry{
			Proto:       proto,
			State:       state,
			LocalAddr:   localAddr,
			LocalPort:   localPort,
			PID:         pid,
			ProcessName: pname,
			UserName:    uname,
			Tag:         tag,
		})
	}

	return entries
}

// local_field looks like "0100007F:1F90" (IPv4) or "0000000000000000FFFFFFFF00000000:0035" (IPv6)
func parseIPPort(field string, isIPv6 bool) (string, string) {
	parts := strings.Split(field, ":")
	if len(parts) != 2 {
		return field, ""
	}
	ipHex := parts[0]
	portHex := parts[1]

	// port
	portVal, err := strconv.ParseUint(portHex, 16, 16)
	if err != nil {
		return field, ""
	}
	port := strconv.Itoa(int(portVal))

	if isIPv6 {
		// keep IPv6 in compact form (don't overcomplicate for CLI)
		ip := shortenIPv6(ipHex)
		return ip, port
	}

	// IPv4: hex is little-endian, 4 bytes
	if len(ipHex) != 8 {
		return field, port
	}
	b1, _ := strconv.ParseUint(ipHex[6:8], 16, 8)
	b2, _ := strconv.ParseUint(ipHex[4:6], 16, 8)
	b3, _ := strconv.ParseUint(ipHex[2:4], 16, 8)
	b4, _ := strconv.ParseUint(ipHex[0:2], 16, 8)
	ip := fmt.Sprintf("%d.%d.%d.%d", b1, b2, b3, b4)
	return ip, port
}

func shortenIPv6(hex string) string {
	// hex is 32 chars; group into 8 segments
	if len(hex) != 32 {
		return hex
	}
	parts := make([]string, 0, 8)
	for i := 0; i < 32; i += 4 {
		part := hex[i : i+4]
		// strip leading zeros
		part = strings.TrimLeft(part, "0")
		if part == "" {
			part = "0"
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, ":")
}

func decodeState(proto, hexState string) string {
	hexState = strings.ToUpper(hexState)
	if proto == "tcp" {
		switch hexState {
		case "01":
			return "ESTAB"
		case "02":
			return "SYN-SENT"
		case "03":
			return "SYN-RECV"
		case "04":
			return "FIN-WAIT1"
		case "05":
			return "FIN-WAIT2"
		case "06":
			return "TIME-WAIT"
		case "07":
			return "CLOSE"
		case "08":
			return "CLOSE-WAIT"
		case "09":
			return "LAST-ACK"
		case "0A":
			return "LISTEN"
		case "0B":
			return "CLOSING"
		default:
			return "UNKNOWN"
		}
	}
	// UDP states are less meaningful; 07 is usually "UNCONN"
	if hexState == "07" {
		return "UNCONN"
	}
	return "UNKNOWN"
}

// ------------------------------------------------------------
// process / user helpers (mostly unchanged from before)
// ------------------------------------------------------------

func getProcessNameFromPID(pid int) string {
	if pid <= 0 {
		return "?"
	}

	// /proc/<pid>/comm
	commPath := filepath.Join("/proc", strconv.Itoa(pid), "comm")
	if data, err := os.ReadFile(commPath); err == nil {
		name := strings.TrimSpace(string(data))
		if name != "" {
			return name
		}
	}

	// /proc/<pid>/exe symlink
	exePath := filepath.Join("/proc", strconv.Itoa(pid), "exe")
	if link, err := os.Readlink(exePath); err == nil && link != "" {
		return filepath.Base(link)
	}

	// /proc/<pid>/cmdline
	cmdPath := filepath.Join("/proc", strconv.Itoa(pid), "cmdline")
	if data, err := os.ReadFile(cmdPath); err == nil {
		parts := strings.Split(string(data), "\x00")
		if len(parts) > 0 && parts[0] != "" {
			return filepath.Base(parts[0])
		}
	}

	return "?"
}

func getUserFromPID(pid int) (string, string) {
	if pid <= 0 {
		return "?", ""
	}

	statusPath := filepath.Join("/proc", strconv.Itoa(pid), "status")
	data, err := os.ReadFile(statusPath)
	if err != nil {
		return "?", ""
	}

	var uidLine string
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "Uid:") {
			uidLine = line
			break
		}
	}
	if uidLine == "" {
		return "?", ""
	}

	parts := strings.Fields(uidLine)
	if len(parts) < 2 {
		return "?", ""
	}
	uid := parts[1]

	u, err := user.LookupId(uid)
	if err != nil {
		return "uid=" + uid, uid
	}
	return u.Username, uid
}

func classifyEntry(uid, curUID string, pid int) string {
	if pid == os.Getpid() {
		return "SELF"
	}
	if uid == "" || pid == 0 {
		return "SYSTEM"
	}
	if uid == "0" {
		return "SYSTEM"
	}
	if curUID != "" && uid == curUID {
		return "USER"
	}
	if v, err := strconv.Atoi(uid); err == nil && v < 1000 {
		return "SYSTEM"
	}
	return "USER"
}
