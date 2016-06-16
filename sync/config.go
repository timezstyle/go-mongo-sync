package sync

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

// Config implemention that parse from command arguments.
type Config struct {
	GoroutineNumber int
	From            string
	To              string
	Database        string
	Query           string
	Log             string
	StartOptime     int
	SrcHost         string
	SrcPort         int
	DstHost         string
	DstPort         int
	Sleep		int
	RetryCodes	string
	RetryCodeList   []int
}

// load and parse command-line flags
func (p *Config) Load() error {
	flag.StringVar(&p.From, "from", "", "source, should be a member of replica-set")
	flag.StringVar(&p.To, "to", "", "destination, should be a mongos or mongod instance")
	flag.StringVar(&p.Database, "db", "", "database to sync")
	flag.IntVar(&p.GoroutineNumber, "c", 20, "goroutine number")
	flag.IntVar(&p.StartOptime, "start_optime", -1, "start optime, -1 means no specify.")
	flag.IntVar(&p.Sleep, "sleep", 100, "sleep ms, when move snapshot error!")
	flag.StringVar(&p.RetryCodes, "retryCodes", "10058,", "oplog retry when meet retryCodes, use ',' to seperate error codes.")
	flag.Parse()
	var retryCodeList = []int{}
	retryCodesStr := strings.Split(p.RetryCodes, ",")
	for _, v := range retryCodesStr {
		trimV := strings.TrimSpace(v)
		if trimV == "" {
			continue
		}
		i, err := strconv.Atoi(trimV)
		if err != nil {
			panic(err)
		}
		retryCodeList = append(retryCodeList, i)
	}
	p.RetryCodeList = retryCodeList
	if err := p.validate(); err != nil {
		return err
	}
	p.print()
	return nil
}

// validate command-line flags
func (p *Config) validate() error {
	var err error
	p.SrcHost, p.SrcPort, err = parse_host_port(p.From)
	if err != nil {
		return errors.New("from error: " + err.Error())
	}
	p.DstHost, p.DstPort, err = parse_host_port(p.To)
	if err != nil {
		return errors.New("to error: " + err.Error())
	}
	return nil
}

// print config
func (p *Config) print() {
	fmt.Printf("from: %s:%d\n", p.SrcHost, p.SrcPort)
	fmt.Printf("to:   %s:%d\n", p.DstHost, p.DstPort)
}

// parse hostportstr
func parse_host_port(hostportstr string) (host string, port int, err error) {
	s := strings.Split(hostportstr, ":")
	if len(s) != 2 {
		return host, port, errors.New("invalid hostportstr")
	}
	host = s[0]
	port, err = strconv.Atoi(s[1])
	if err != nil {
		return host, port, errors.New("invalid port")
	}
	if port < 0 || port > 65535 {
		return host, port, errors.New("invalid port")
	}
	return host, port, nil
}
