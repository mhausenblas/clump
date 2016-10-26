package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	VERSION = "0.1.0"
)

var (
	user                  string
	fprivatekey           string
	fnodes                string
	fcmds                 string
	IPv4_24bit_block_low  net.IP
	IPv4_24bit_block_high net.IP
	IPv4_20bit_block_low  net.IP
	IPv4_20bit_block_high net.IP
	IPv4_16bit_block_low  net.IP
	IPv4_16bit_block_high net.IP
	CONNECTION_TIMEOUT    time.Duration
)

type SSHClient struct {
	Config *ssh.ClientConfig
	Client *ssh.Client
	Host   string
	Port   int
}

type Commands struct {
	local  []string
	remote []string
}

func init() {
	// What follows is the IPv4 private address space
	// as of https://tools.ietf.org/html/rfc1918
	IPv4_24bit_block_low = net.ParseIP("10.0.0.0")
	IPv4_24bit_block_high = net.ParseIP("10.255.255.255")
	IPv4_20bit_block_low = net.ParseIP("172.16.0.0")
	IPv4_20bit_block_high = net.ParseIP("172.31.255.255")
	IPv4_16bit_block_low = net.ParseIP("192.168.0.0")
	IPv4_16bit_block_high = net.ParseIP("192.168.255.255")

	// how long to wait for the SSH connection to establish (TODO: make configurable via env variable)
	CONNECTION_TIMEOUT = 5 * time.Second

	flag.StringVar(&user, "u", "", "user name to use for SSH connection")
	flag.StringVar(&fprivatekey, "pk", "", "path to private key to use for SSH connection")
	flag.StringVar(&fnodes, "nl", "", "path to file listing the target nodes, one IP address per line")
	flag.StringVar(&fcmds, "cmds", "", "path to file listing the commands to be executed, one entry per line")
	flag.Usage = func() {
		about()
		fmt.Println("\nUsage: clump -u $USERNAME -pk $PRIVATESSHKEY -nl $NODES -cmds $COMMANDS")
		fmt.Println("\nExample: clump -u core -pk /Users/mhausenblas/.ssh/test -nl clusternodes -cmds snapshot")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func about() {
	fmt.Println(fmt.Sprintf("This is clump version %s", VERSION))
}

///////////////////////////////////////////////////////////////////////////////
// utils

func privateIP(ip string) (bool, error) {
	trial := net.ParseIP(ip)
	isprivate := false
	if trial.To4() == nil { // can't parse as IP4
		return isprivate, errors.New("Can't parse IP address")
	}
	if bytes.Compare(trial, IPv4_24bit_block_low) >= 0 && bytes.Compare(trial, IPv4_24bit_block_high) <= 0 {
		isprivate = true
	} else {
		if bytes.Compare(trial, IPv4_20bit_block_low) >= 0 && bytes.Compare(trial, IPv4_20bit_block_high) <= 0 {
			isprivate = true
		} else {
			if bytes.Compare(trial, IPv4_16bit_block_low) >= 0 && bytes.Compare(trial, IPv4_16bit_block_high) <= 0 {
				isprivate = true
			}
		}
	}
	return isprivate, nil
}

func nodes() ([]string, error) {
	fmt.Println(fmt.Sprintf("Trying to establish node list from %s", fnodes))
	n := make([]string, 0)
	if nl, err := ioutil.ReadFile(fnodes); err == nil {
		lines := strings.Split(string(nl), "\n")
		for _, l := range lines {
			if !strings.HasPrefix(l, "#") {
				n = append(n, strings.TrimSpace(l))
			}
		}
	} else {
		return nil, err
	}
	return n, nil
}

func commands() (*Commands, error) {
	fmt.Println(fmt.Sprintf("Trying to establish list of commands from %s", fcmds))
	c := &Commands{}
	if cmdl, err := ioutil.ReadFile(fcmds); err == nil {
		lines := strings.Split(string(cmdl), "\n")
		for _, l := range lines {
			scope := strings.Split(l, ":")[0]
			cmd := strings.Split(l, ":")[1]
			if !strings.HasPrefix(l, "#") {
				if scope == "LOCAL" {
					c.local = append(c.local, cmd)
				} else {
					c.remote = append(c.remote, cmd)
				}
			}
		}
	} else {
		return nil, err
	}
	return c, nil
}

func nexec(nodes []string, commands []string) {
	for _, node := range nodes {
		mexec(node, commands)
	}
}

func mexec(node string, commands []string) {
	if isprivate, err := privateIP(node); err == nil {
		resultdir := strings.Replace(node, ".", "_", -1)
		if _, err := os.Stat(resultdir); os.IsNotExist(err) {
			os.Mkdir(resultdir, 0700)
		}
		if isprivate { // node has an IP that is in the private address space
			fmt.Println(fmt.Sprintf("%s is an IP address in the private address space", node))
			for _, cmd := range commands {
				sshConfig := &ssh.ClientConfig{
					User: user,
					Auth: []ssh.AuthMethod{
						publickey(fprivatekey)},
					Timeout: CONNECTION_TIMEOUT,
				}
				client := &SSHClient{
					Config: sshConfig,
					Host:   "35.160.66.81",
					Port:   22,
				}
				remote := &SSHClient{
					Config: sshConfig,
					Host:   node,
					Port:   22,
				}

				var bin, bout bytes.Buffer
				buf := bufio.NewReadWriter(bufio.NewReader(&bin), bufio.NewWriter(&bout))

				if err := remote.run(cmd, "", buf); err != nil {
					fmt.Println(fmt.Sprintf("Remote target executing %s on %s failed ", cmd, remote.Host, err))
				}
				if err := client.run(cmd, resultdir, buf); err != nil {
					fmt.Println(fmt.Sprintf("Jump host executing %s on %s failed ", cmd, client.Host, err))
				}

			}
		} else {
			for _, cmd := range commands {
				rexec(node, cmd, resultdir)
			}
		}
	} else {
		fmt.Println(fmt.Sprintf("Skipping %s since it's not a valid IPv4 address", node))
	}
}

func rexec(node string, command string, resultdir string) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			publickey(fprivatekey)},
		Timeout: CONNECTION_TIMEOUT,
	}
	fmt.Println(fmt.Sprintf("Attempting to ssh into %s@%s to execute %s", user, node, command))
	client := &SSHClient{
		Config: sshConfig,
		Host:   node,
		Port:   22,
	}
	if err := client.run(command, resultdir, nil); err != nil {
		fmt.Println(fmt.Sprintf("Executing %s on %s failed ", command, client.Host, err))
	}
}

///////////////////////////////////////////////////////////////////////////////
// SSH connection stuff

func (client *SSHClient) run(command string, resultdir string, sink *bufio.ReadWriter) error {
	s := &ssh.Session{}
	err := errors.New("")

	if s, err = client.create(); err != nil {
		return err
	}
	defer s.Close()

	so, _ := s.StdoutPipe()
	if resultdir != "" {
		resultfile := strings.Replace(command, " ", "_", -1)
		resultfile = strings.Replace(resultfile, "/", "-", -1)
		resultfile = strings.Replace(resultfile, ".", "", -1)
		rfname := filepath.Join(resultdir, resultfile)
		rf := &os.File{}
		if rf, err = os.Create(rfname); err != nil {
			return err
		}
		defer rf.Close()
		go io.Copy(rf, so)
		err = s.Run(command)
		return err
	} else {
		go io.Copy(sink, so)
	}
	return nil
}

func (client *SSHClient) create() (*ssh.Session, error) {
	c, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port), client.Config)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial: %s", err)
	}
	client.Client = c
	s, err := c.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %s", err)
	}
	return s, nil
}

func publickey(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func main() {
	if user != "" &&
		fprivatekey != "" &&
		fnodes != "" &&
		fcmds != "" {
		if n, nerr := nodes(); nerr != nil {
			fmt.Println(fmt.Sprintf("Problem reading node list ", nerr))
			os.Exit(1)
		} else {
			fmt.Println(fmt.Sprintf("Got %d target node(s)", len(n)))
			if c, cerr := commands(); cerr != nil {
				fmt.Println(fmt.Sprintf("Problem reading commands ", cerr))
				os.Exit(2)
			} else {
				fmt.Println(fmt.Sprintf("Executing %d command(s) locally ...", len(c.local)))
				for _, c := range c.local {
					cmd := &exec.Cmd{Path: strings.Fields(c)[0], Args: strings.Fields(c)[1:]}
					out, err := cmd.Output()
					if err != nil {
						fmt.Println(fmt.Sprintf("Problem executing \"%s\" %s", c, err))
						os.Exit(3)
					}
					fmt.Printf(string(out))
				}
				fmt.Println(fmt.Sprintf("Executing %d command(s) remotely ...", len(c.remote)))
				nexec(n, c.remote)
			}
		}
	} else {
		flag.Usage()
		os.Exit(4)
	}
	os.Exit(0)
}
