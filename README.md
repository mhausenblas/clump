# clump

A configurable cluster dump tool for collecting system status information about nodes in a cluster, using SSH.

## Usage

To dump a cluster state:

    $ clump -u $USERNAME -pk $PRIVATESSHKEY -nl $NODES -cmds $COMMANDS

with:

- `$USERNAME` … username to use for SSH connection  
- `$PRIVATESSHKEY` … filename (with relative or absolute filepath) of the private SSH key to use
- `$NODES` … filename (with relative or absolute filepath) of a text file listing the target nodes, one IP address per line
- `$COMMANDS` … filename (with relative or absolute filepath) of a text file listing the commands to be executed with one entry per line; if an entry is prefixed with `LOCAL:` it will be executed locally, if `REMOTE:` then on each of the target nodes

For example:

    $ cat test/clusternodes
    35.160.157.251
    
    $ cat test/snapshot.cmds
    REMOTE:hostname -f
    REMOTE:timedatectl
    #REMOTE:cat /proc/version
    #REMOTE:sudo ps faux
    
    $ clump -u core -pk /Users/mhausenblas/.ssh/test -nl test/clusternodes -cmds test/snapshot.cmds
    Trying to establish node list from test/clusternodes
    Got 1 target nodes
    Trying to establish list of commands from test/snapshot.cmds
    Got 2 commands to execute
    Attempting to ssh into core@35.160.157.251
    Attempting to ssh into core@35.160.157.251
    
    $ tree
    .
    └── 35_160_157_251
        ├── hostname_-f
        └── timedatectl

Kudos go out to [Svett Ralchev](http://blog.ralch.com/tutorial/golang-ssh-connection/) for the seed code base around the SSH client.

To do:

- fix DC/OS recipe
- credential forwarding?


## Disclaimer

> THIS SOFTWARE IS PROVIDED "AS IS" AND ANY EXPRESSED OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

