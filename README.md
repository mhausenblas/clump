# clump

A configurable cluster dump tool for collecting system status information about nodes in a cluster, using SSH.

## Install

### Linux

    $ curl -s -L https://github.com/mhausenblas/clump/releases/download/v0.1.1/linux-clump -o clump
    $ sudo mv clump /usr/local/bin
    $ sudo chmod +x /usr/local/bin/clump

### MacOS

    $ curl -s -L https://github.com/mhausenblas/clump/releases/download/v0.1.1/osx-clump -o clump
    $ sudo mv clump /usr/local/bin
    $ sudo chmod +x /usr/local/bin/clump

## Use

You can use clump as is or via one of the recipes supplied:

- [Recipe for DC/OS](/recipes/dcos)
- ... 

Generic usage is as follows:

    $ clump -u $USERNAME -pk $PRIVATESSHKEY -nl $NODES -cmds $COMMANDS

with:

- `$USERNAME` … username to use for SSH connection
- `$PRIVATESSHKEY` … filename (with relative or absolute filepath) of the private SSH key to use
- `$NODES` … filename (with relative or absolute filepath) of a text file listing the target nodes, one IP address per line
- `$COMMANDS` … filename (with relative or absolute filepath) of a text file listing the commands to be executed with one entry per line; if an entry is prefixed with `LOCAL:` it will be executed locally, if `REMOTE:` then on each of the target nodes

For example:

    $ cat clusternodes
    35.160.157.251
    
    $ cat snapshot.cmds
    LOCAL:id
    REMOTE:hostname -f
    REMOTE:timedatectl
    
    $ clump -u core -pk /Users/mhausenblas/.ssh/test -nl clusternodes -cmds snapshot.cmds
    Trying to establish node list from clusternodes
    Got 1 target node(s)
    Executing 1 command(s) locally ...
    Executing 2 command(s) remotely ...
    Attempting to ssh into core@35.160.157.251
    Attempting to ssh into core@35.160.157.251
    
    $ tree
    .
    └── 35_160_157_251
        ├── hostname_-f
        └── timedatectl

Kudos go out to [Svett Ralchev](http://blog.ralch.com/tutorial/golang-ssh-connection/) for the seed code base around the SSH client.

## Disclaimer

> THIS SOFTWARE IS PROVIDED "AS IS" AND ANY EXPRESSED OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
