# Recipe for DC/OS

To use `clump` for dumping the state of a DC/OS cluster, clone this repo and use the templates provided here.

In order to work, you will need to have a DC/OS cluster 1.8 or above running as well as the DC/OS CLI installed and authenticated.
If you want low-level namespaces and cgroups info you need to have [cinf](https://github.com/mhausenblas/cinf/) installed and available from the home directory, otherwise comment out the last line in `snapshot.cmds`, that is, `# REMOTE:sudo ./cinf`.

## Prepare

Init `clusternodes` with the IP of the Master(s) manually and then do the following to populate it with the agents:

    $ echo ""  >> clusternodes &&  dcos node | tail -n +2 | awk '{print $2}' >> clusternodes
    $ cat clusternodes
    35.160.66.81
    10.0.2.25
    10.0.5.173

## Use

Now you run `clump` as so (assuming a cluster deployed in AWS):
    
    $ cd recipes/dcos/
    $ clump -u core -pk ~/.ssh/awskey -nl clusternodes -cmds snapshot.cmds