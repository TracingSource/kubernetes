package remote

import (
	"flag"
	"fmt"
	"os/exec"
	"os/user"
	"strings"
	"sync"

	"k8s.io/klog"
)

var sshOptions = flag.String("ssh-options", "", "Commandline options passed to ssh.")
var sshEnv = flag.String("ssh-env", "", "Use predefined ssh options for environment.  Options: gce")
var sshKey = flag.String("ssh-key", "", "Path to ssh private key.")
var sshUser = flag.String("ssh-user", "", "Use predefined user for ssh.")

var sshOptionsMap map[string]string
var sshDefaultKeyMap map[string]string

func init() {
	usr, err := user.Current()
	if err != nil {
		klog.Fatal(err)
	}
	sshOptionsMap = map[string]string{
		"gce": "-o UserKnownHostsFile=/dev/null -o IdentitiesOnly=yes -o CheckHostIP=no -o StrictHostKeyChecking=no -o ServerAliveInterval=30 -o LogLevel=ERROR",
	}
	sshDefaultKeyMap = map[string]string{
		"gce": fmt.Sprintf("%s/.ssh/google_compute_engine", usr.HomeDir),
	}
}

var hostnameIpOverrides = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

func AddHostnameIp(hostname, ip string) {
	hostnameIpOverrides.Lock()
	defer hostnameIpOverrides.Unlock()
	hostnameIpOverrides.m[hostname] = ip
}

// GetHostnameOrIp converts hostname into ip and apply user if necessary.
func GetHostnameOrIp(hostname string) string {
	hostnameIpOverrides.RLock()
	defer hostnameIpOverrides.RUnlock()
	host := hostname
	if ip, found := hostnameIpOverrides.m[hostname]; found {
		host = ip
	}
	if *sshUser != "" {
		host = fmt.Sprintf("%s@%s", *sshUser, host)
	}
	return host
}

// getSSHCommand handles proper quoting so that multiple commands are executed in the same shell over ssh
func getSSHCommand(sep string, args ...string) string {
	return fmt.Sprintf("'%s'", strings.Join(args, sep))
}

// SSH executes ssh command with runSSHCommand as root. The `sudo` makes sure that all commands
// are executed by root, so that there won't be permission mismatch between different commands.
func SSH(host string, cmd ...string) (string, error) {
	return runSSHCommand("ssh", append([]string{GetHostnameOrIp(host), "--", "sudo"}, cmd...)...)
}

// SSHNoSudo executes ssh command with runSSHCommand as normal user. Sometimes we need this,
// for example creating a directory that we'll copy files there with scp.
func SSHNoSudo(host string, cmd ...string) (string, error) {
	return runSSHCommand("ssh", append([]string{GetHostnameOrIp(host), "--"}, cmd...)...)
}

// runSSHCommand executes the ssh or scp command, adding the flag provided --ssh-options
func runSSHCommand(cmd string, args ...string) (string, error) {
	if *sshKey != "" {
		args = append([]string{"-i", *sshKey}, args...)
	} else if key, found := sshDefaultKeyMap[*sshEnv]; found {
		args = append([]string{"-i", key}, args...)
	}
	if env, found := sshOptionsMap[*sshEnv]; found {
		args = append(strings.Split(env, " "), args...)
	}
	if *sshOptions != "" {
		args = append(strings.Split(*sshOptions, " "), args...)
	}
	output, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command [%s %s] failed with error: %v", cmd, strings.Join(args, " "), err)
	}
	return string(output), nil
}
