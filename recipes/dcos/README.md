# Recipe for DC/OS

This recipe is for dumping the state of a DC/OS cluster.

Prerequisites:

- A running DC/OS cluster 1.8 or above
- The DC/OS CLI installed and authenticated (`dcos auth login`)
- If you want low-level namespaces and cgroups info you need to have [cinf](https://github.com/mhausenblas/cinf/) installed and available from the home directory of all you nodes, otherwise comment out the last line in `snapshot.cmds`, that is do, `# REMOTE:sudo ./cinf`.

## Prepare

To use `clump` for dumping the state of a DC/OS cluster you need to do the following steps as a preparation.

### Get the recipe

You'll need the recipe here so first clone this repo and change to the respective directory:

    $ git clone git://github.com/mhausenblas/clump.git
    $ cd clump/recipes/dcos

### Establish the `clusternode` list

Init `clusternodes` with the IP of the Master(s) manually and then do the following to populate it with the agents:

    $ echo ""  >> clusternodes &&  dcos node | tail -n +2 | awk '{print $2}' >> clusternodes
    $ cat clusternodes
    35.160.66.81
    10.0.2.25
    10.0.5.173

### Make agents accessible via VPN

In order to gather information from the agents, which is especially necessary for the private ones since they have IP addresses from the private address space as of [RFC1918](https://tools.ietf.org/html/rfc1918), you need to be on the same network as the DC/OS cluster. For this, use the DC/OS CLI to create a VPN following the instructions for [tunneling](https://dcos.io/docs/1.8/administration/access-node/tunnel/). In a nutshell, this means:

- We need the DC/OS CLI `tunnel` subcommand so do: `dcos package install tunnel-cli --cli`
- install an OpenVPN client on your machine; in case of Tunnelblick (on MacOS) launch it like so: `sudo dcos tunnel vpn --client=/Applications/Tunnelblick.app/Contents/Resources/openvpn/openvpn-*/openvpn`
- add the DNS servers IP addresses as shown by the OpenVPN client; in case of MacOS see instructions [from Apple support](https://support.apple.com/kb/PH18499?locale=en_US) and make sure the three IPs are the first in the list (once applied you can also check via `/etc/resolv.conf` if they are set correctly)
- to verify the setup you can do `curl leader.mesos` and you should see some HTML content, if not, either your OpenVPN client is not properly working or the DNS servers are not on top of the resolver list

## Use

Once you've completed all the preparation steps above, you run `clump` as so (note that in the following, I'm assuming a cluster deployed in AWS, change the user name with `-u XXX` respectively if you have a different setup; same goes for the private key, change it to yours respectively with `-pk ZZZ`):
    
    $ clump -u core -pk ~/.ssh/awskey -nl clusternodes -cmds snapshot.cmds
    Executing 1 command(s) locally ...
    Diagnostics bundle downloaded to /Users/mhausenblas/tmp/clump/bundle-2016-10-26T16:17:08-645940292.zip
    Executing 12 command(s) remotely ...
    Attempting to ssh into core@35.160.66.81 to execute /bin/hostname -f
    Attempting to ssh into core@35.160.66.81 to execute timedatectl
    Attempting to ssh into core@35.160.66.81 to execute cat /proc/version
    Attempting to ssh into core@35.160.66.81 to execute sudo ps faux
    Attempting to ssh into core@35.160.66.81 to execute cat /etc/passwd
    Attempting to ssh into core@35.160.66.81 to execute df -h
    Attempting to ssh into core@35.160.66.81 to execute mount
    Attempting to ssh into core@35.160.66.81 to execute sudo ss -lptn
    Attempting to ssh into core@35.160.66.81 to execute sudo ip route
    Attempting to ssh into core@35.160.66.81 to execute cat /etc/resolv.conf
    Attempting to ssh into core@35.160.66.81 to execute systemctl list-unit-files --all
    Attempting to ssh into core@35.160.66.81 to execute sudo docker --version
    10.0.2.25 is an IP address in the private address space
    Attempting to ssh into core@10.0.2.25 to execute /bin/hostname -f
    Attempting to ssh into core@10.0.2.25 to execute timedatectl
    Attempting to ssh into core@10.0.2.25 to execute cat /proc/version
    Attempting to ssh into core@10.0.2.25 to execute sudo ps faux
    Attempting to ssh into core@10.0.2.25 to execute cat /etc/passwd
    Attempting to ssh into core@10.0.2.25 to execute df -h
    Attempting to ssh into core@10.0.2.25 to execute mount
    Attempting to ssh into core@10.0.2.25 to execute sudo ss -lptn
    Attempting to ssh into core@10.0.2.25 to execute sudo ip route
    Attempting to ssh into core@10.0.2.25 to execute cat /etc/resolv.conf
    Attempting to ssh into core@10.0.2.25 to execute systemctl list-unit-files --all
    Attempting to ssh into core@10.0.2.25 to execute sudo docker --version
    10.0.5.173 is an IP address in the private address space
    Attempting to ssh into core@10.0.5.173 to execute /bin/hostname -f
    Attempting to ssh into core@10.0.5.173 to execute timedatectl
    Attempting to ssh into core@10.0.5.173 to execute cat /proc/version
    Attempting to ssh into core@10.0.5.173 to execute sudo ps faux
    Attempting to ssh into core@10.0.5.173 to execute cat /etc/passwd
    Attempting to ssh into core@10.0.5.173 to execute df -h
    Attempting to ssh into core@10.0.5.173 to execute mount
    Attempting to ssh into core@10.0.5.173 to execute sudo ss -lptn
    Attempting to ssh into core@10.0.5.173 to execute sudo ip route
    Attempting to ssh into core@10.0.5.173 to execute cat /etc/resolv.conf
    Attempting to ssh into core@10.0.5.173 to execute systemctl list-unit-files --all
    Attempting to ssh into core@10.0.5.173 to execute sudo docker --version

The results are available now in the directory you executed `clump`:

     $ tree
    .
    ├── 10_0_2_25
    │   ├── -bin-hostname_-f
    │   ├── cat_-etc-passwd
    │   ├── cat_-etc-resolvconf
    │   ├── cat_-proc-version
    │   ├── df_-h
    │   ├── mount
    │   ├── sudo_docker_--version
    │   ├── sudo_ip_route
    │   ├── sudo_ps_faux
    │   ├── sudo_ss_-lptn
    │   ├── systemctl_list-unit-files_--all
    │   └── timedatectl
    ├── 10_0_5_173
    │   ├── -bin-hostname_-f
    │   ├── cat_-etc-passwd
    │   ├── cat_-etc-resolvconf
    │   ├── cat_-proc-version
    │   ├── df_-h
    │   ├── mount
    │   ├── sudo_docker_--version
    │   ├── sudo_ip_route
    │   ├── sudo_ps_faux
    │   ├── sudo_ss_-lptn
    │   ├── systemctl_list-unit-files_--all
    │   └── timedatectl
    ├── 35_160_66_81
    │   ├── -bin-hostname_-f
    │   ├── cat_-etc-passwd
    │   ├── cat_-etc-resolvconf
    │   ├── cat_-proc-version
    │   ├── df_-h
    │   ├── mount
    │   ├── sudo_docker_--version
    │   ├── sudo_ip_route
    │   ├── sudo_ps_faux
    │   ├── sudo_ss_-lptn
    │   ├── systemctl_list-unit-files_--all
    │   └── timedatectl
    ├── bundle-2016-10-26T16:17:08-645940292.zip

