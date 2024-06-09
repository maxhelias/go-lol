package process

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

const (
	LeagueClientUx = `LeagueClientUx`
	AuthTokenRegex = `--remoting-auth-token=(\S+)` // #nosec G101
	PortRegex      = `--app-port=(\d+)`
)

type LcuProcessNotFoundError struct{}

func (e LcuProcessNotFoundError) Error() string {
	return `no process found with the name ` + LeagueClientUx
}

type LcuConnectInfo struct {
	ProcessID int32
	// Indicates the port of the LCU
	Port int
	// Indicates the authentication token of the LCU
	AuthToken string
}

func FindLcuConnectInfo() (*LcuConnectInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	for _, p := range processes {
		name, err := p.Name()

		if err != nil {
			continue
		}

		if strings.Contains(name, LeagueClientUx) {
			cmdLine, err := p.Cmdline()
			if err != nil {
				continue
			}

			authToken, authTokenErr := extractAuthToken(cmdLine)
			port, portErr := extractPort(cmdLine)

			if (authTokenErr != nil) || (portErr != nil) {
				continue
			}

			return &LcuConnectInfo{
				ProcessID: p.Pid,
				Port:      port,
				AuthToken: authToken,
			}, nil
		}
	}

	return nil, &LcuProcessNotFoundError{}
}

func extractAuthToken(cmdLine string) (string, error) {
	return extractRegex(cmdLine, AuthTokenRegex)
}

func extractPort(cmdLine string) (int, error) {
	str, err := extractRegex(cmdLine, PortRegex)
	if nil != err {
		return 0, err
	}

	num, err := strconv.Atoi(str)

	return num, err
}

func extractRegex(input, regex string) (string, error) {
	regexCompiled := regexp.MustCompile(regex)

	match := regexCompiled.FindStringSubmatch(input)
	if match != nil {
		return match[1], nil
	}

	return ``, errors.New("no match found")
}
