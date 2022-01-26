package eosconnect

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aristanetworks/goeapi"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Conn struct {
	Transport string
	Host      string
	Username  string
	Password  string
	Port      int
	Config    string
}

func (c *Conn) Connect() (*goeapi.Node, error) {
	connect, err := goeapi.Connect(c.Transport, c.Host, c.Username, c.Password, c.Port)
	if err != nil {
		fmt.Println(err)
	}
	return connect, nil
}

func (c *Conn) Compare() string {
	compare, err := c.Connect()
	if err != nil {
		fmt.Println(err)
	}
	runnconfig, err := compare.Enable([]string{"show running-config"})
	if err != nil {
		fmt.Println(err)
	}
	var config string
	for _, v := range runnconfig {
		for _, r := range v {
			config = r
		}
	}
	someslice := []string{}
	for _, line := range strings.Split(strings.TrimRight(config, "\n"), "\n") {
		someslice = append(someslice, line)
	}
	SliceLength := len(someslice)
	var trimmedconfig string
	if len(someslice) < 3 {
		log.Log.Info("Could not fetch the running config.  Try again next reconcile.")
	} else {
		a := strings.Join(someslice[3:SliceLength], "\n")
		trimmedconfig = a
	}
	return trimmedconfig
}

func (c *Conn) Configure(runningconfig string) string {
	year, month, day := time.Now().Date()
	hour := time.Now().Hour()
	minutes := time.Now().Minute()
	seconds := time.Now().Second()

	concatenated := strings.Join([]string{strconv.Itoa(year), strconv.Itoa(day), strconv.Itoa(int(month)), strconv.Itoa(hour), strconv.Itoa(minutes), strconv.Itoa(seconds)}, "")
	configure, err := c.Connect()
	if err != nil {
		fmt.Println(err)
	}
	splitData := strings.Split(runningconfig, "\n")
	cmds := []string{"configure session eapi-" + concatenated, "rollback clean-config", "copy terminal: session-config"}
	cmds = append(cmds, splitData...)
	configdevice, err := configure.RunCommands(cmds, "json")
	if err != nil {
		fmt.Println(err)
	}

	_, err = configure.RunCommands([]string{"configure session eapi-" + concatenated, "commit"}, "json")
	if err != nil {
		fmt.Println(err)
	}

	return configdevice.Jsonrpc
}
