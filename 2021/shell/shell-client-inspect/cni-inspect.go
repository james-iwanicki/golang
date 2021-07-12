package main

import (
//	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/cni/pkg/version"
)

const (
	logPath     = "/var/log/ntct.log"
	defaultSock = "/opt/cni/bin/ntct.sock"
// different sock location for openshift
//	defaultSock = "/var/lib/cni/bin/ntct.sock"
)

// PluginConf is whatever you expect your configuration json to be. This is whatever
// is passed in on stdin. Your plugin may wish to expose its functionality via
// runtime args, see CONVENTIONS.md in the CNI spec.
type PluginConf struct {
	// This is the previous result, when called in the context of a chained
	// plugin. Because this plugin supports multiple versions, we'll have to
	// parse this in two passes. If your plugin is not chained, this can be
	// removed (though you may wish to error if a non-chainable plugin is
	// chained.
	// If you need to modify the result before returning it, you will need
	// to actually convert it to a concrete versioned struct.
	RawPrevResult *map[string]interface{} `json:"prevResult"`
	PrevResult    *current.Result         `json:"-"`
	types.NetConf
	Unix string `json:"unix"`
	IP   string `json:"ip"`
}

// parseConfig parses the supplied configuration (and prevResult) from stdin.
func parseConfig(stdin []byte) (*PluginConf, error) {
	conf := PluginConf{}

	if err := json.Unmarshal(stdin, &conf); err != nil {
		return nil, fmt.Errorf("failed to parse network configuration: %v", err)
	}

	// Parse previous result. Remove this if your plugin is not chained.
	if conf.RawPrevResult != nil {
		resultBytes, err := json.Marshal(conf.RawPrevResult)
		if err != nil {
			return nil, fmt.Errorf("could not serialize prevResult: %v", err)
		}
		res, err := version.NewResult(conf.CNIVersion, resultBytes)
		if err != nil {
			return nil, fmt.Errorf("could not parse prevResult: %v", err)
		}
		conf.RawPrevResult = nil
		conf.PrevResult, err = current.NewResultFromResult(res)
		if err != nil {
			return nil, fmt.Errorf("could not convert result to current version: %v", err)
		}
	}
	// End previous result parsing

	if conf.Unix == "" && conf.IP == "" {
		conf.Unix = defaultSock
	}

	return &conf, nil
}

func appendLog(d string) error {
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("Failed opening file: %s", err)
	}

	_, err = file.WriteString(d)
	return err
}

func getClient(conf *PluginConf) http.Client {
	httpc := http.Client{
		Timeout: 3 * time.Second, // The whole connection should take less than 3 seconds
	}

	if conf.Unix != "" {
		httpc.Transport = &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				dialer := net.Dialer{}
				return dialer.DialContext(ctx, "unix", conf.Unix)
			},
		}
	}

	return httpc
}

func getURL(conf *PluginConf, uri string) string {
	if uri == "" {
		uri = "/"
	}

	url := "http://unix" + uri
	if conf.IP != "" {
		url = "http://" + conf.IP + uri
	}
	return url
}

type cmdType int

const (
	typeAdd = iota
	typeDel
	typeCheck
)

func send(conf *PluginConf, args *skel.CmdArgs, cmd cmdType) error {
	//httpc := getClient(conf)

	buf, err := json.Marshal(args)
	if err != nil {
		return err
	}

	var method string
	switch cmd {
	case typeAdd:
		method = http.MethodPut
		appendLog(fmt.Sprintf("%s %s\n", method, buf))
	case typeDel:
		method = http.MethodDelete
	case typeCheck:
		method = http.MethodPost
	default:
		return fmt.Errorf("Unknown cmd type: %v", cmd)
	}

/*	request, err := http.NewRequest(method, getURL(conf, "/pod"), bytes.NewReader(buf))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	resp, err := httpc.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with: %d", resp.StatusCode)
	}
*/
	return nil
}

func genericCmd(args *skel.CmdArgs, cmd cmdType) error {
	conf, err := parseConfig(args.StdinData)
	if err != nil {
		return err
	}

	defer func() {
		// Just in case we need to recover
		_ = recover()

		// Pass through the result for the next plugin
		// ignore any errors since we can't do anything about them anyways
		if conf.PrevResult == nil {
			res := &current.Result{}
			_ = types.PrintResult(res, conf.CNIVersion)
		} else {
			_ = types.PrintResult(conf.PrevResult, conf.CNIVersion)
		}
	}()

	return send(conf, args, cmd)
}

// cmdAdd is called for ADD requests
func cmdAdd(args *skel.CmdArgs) error {
	return genericCmd(args, typeAdd)
}

// cmdDel is called for DELETE requests
func cmdDel(args *skel.CmdArgs) error {
	return genericCmd(args, typeDel)
}

func cmdCheck(args *skel.CmdArgs) error {
	return genericCmd(args, typeCheck)
}

func main() {
	err := skel.PluginMainWithError(cmdAdd, cmdCheck, cmdDel, version.All, "NetScout Prototype v0.0.2")
	if err != nil {
		_ = appendLog(fmt.Sprintf("%s\n", err.Error()))
	}
}
