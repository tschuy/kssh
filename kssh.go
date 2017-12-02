package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

var helpText = `kssh: the Kubernetes SSH config helper

Automatically generate an ssh config file for every node on a Kubernetes cluster
with passthrough through a bastion host.

For use with the Include directive of ssh 7.3+. See README for setup instructions.
`

func getNodes() []string {
	out, err := exec.Command("kubectl", "get", "nodes", "-o=custom-columns=:.metadata.name").Output()
	if err != nil {
		log.Fatal(err)
	}
	nodeStr := strings.TrimSpace(string(out))
	return strings.Split(nodeStr, "\n")
}

func main() {
	var bastion = flag.String("bastion", "", "hostname of the bastion server")
	var listNodes = flag.Bool("nodes", false, "just list nodes, don't render config")
	var cssh = flag.String("cssh", "", "output a ClusterSSH configuration line")
	flag.Parse()
	// show usage and exit if necessary
	if !*listNodes && *bastion == "" && *cssh == "" {
		fmt.Print(helpText)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	nodes := getNodes()
	if *cssh != "" {
		// output a cssh config line
		fmt.Printf("%s ", *cssh)
		for _, n := range nodes {
			fmt.Printf("%s ", n)
		}
	} else if *listNodes {
		// a convoluted way to get the output of the source kubectl command
		for _, n := range nodes {
			fmt.Printf("%s\n", n)
		}
	} else {
		renderTemplate(*bastion, nodes)
	}
}

func renderTemplate(bastion string, hosts []string) {
	const sshConfig = `# this file managed by kssh -- changes may be overwritten

Host {{.Bastion}}
	User core
{{ range $host := .Hosts }}
Host {{$host}}
	User core
	ProxyCommand ssh -W %h:%p {{$.Bastion}}
{{ end -}}
`

	// Prepare some data to insert into the template.
	type SSHConfig struct {
		Bastion string
		Hosts   []string
	}
	var config = SSHConfig{bastion, hosts}

	t := template.Must(template.New("sshConfig").Parse(sshConfig))
	_ = t.Execute(os.Stdout, config) // check err
}
