package main

import (
	"errors"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	VERSION = "0.1.0"
)

var (
	user        string
	fprivatekey string
	fnodes      string
	fcmds       string
)

type SSHClient struct {
	Config *ssh.ClientConfig
	Host   string
	Port   int
}

type Commands struct {
	local  []string
	remote []string
}

func init() {
	flag.StringVar(&user, "u", "", "user name to use for SSH connection")
	flag.StringVar(&fprivatekey, "pk", "", "path to private key to use for SSH connection")
	flag.StringVar(&fnodes, "nl", "", "path to file listing the target nodes, one IP address or FQDN per line")
	flag.StringVar(&fcmds, "cmds", "", "path to file listing the commands to be executed, one entry per line")
	flag.Usage = func() {
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
	resultdir := strings.Replace(node, ".", "_", -1)
	if _, err := os.Stat(resultdir); os.IsNotExist(err) {
		os.Mkdir(resultdir, 0700)
	}
	for _, cmd := range commands {
		rexec(node, cmd, resultdir)
	}
}

func rexec(node string, command string, resultdir string) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			publickey(fprivatekey)},
	}
	fmt.Println(fmt.Sprintf("Attempting to ssh into %s@%s to execute %s", user, node, command))
	client := &SSHClient{
		Config: sshConfig,
		Host:   node,
		Port:   22,
	}
	if err := client.run(command, resultdir); err != nil {
		fmt.Println(fmt.Sprintf("Executing %s on %s failed ", command, client.Host, err))
		os.Exit(3)
	}
}

///////////////////////////////////////////////////////////////////////////////
// SSH connection stuff

func (client *SSHClient) run(command string, resultdir string) error {
	s := &ssh.Session{}
	err := errors.New("")
	if s, err = client.create(); err != nil {
		return err
	}
	defer s.Close()
	so, _ := s.StdoutPipe()

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
}

func (client *SSHClient) create() (*ssh.Session, error) {
	c, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port), client.Config)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial: %s", err)
	}
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

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
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
