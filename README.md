# clump

A configurable cluster dump tool for collecting system status information about nodes in a cluster, using SSH.

## Usage

To dump a cluster state:

    $ clump -u $USERNAME -pk $PRIVATESSHKEY -nl $NODES -cmds $COMMANDS

with:

- `$USERNAME` … username to use for SSH connection  
- `$PRIVATESSHKEY` … filename (with relative or absolute filepath) of the private SSH key to use
- `$NODES` … filename (with relative or absolute filepath) of a text file listing the target nodes, one IP address or FQDN per line
- `$COMMANDS` … filename (with relative or absolute filepath) of a text file listing the commands to be executed with one entry per line; if an entry is prefixed with `LOCAL:` it will be executed locally, if `REMOTE:` then on each of the target nodes

For example:

    $ cat mynodes
    1.2.3.4
    5.6.7.8
    
    $ cat mycommands
    LOCAL:dcos node diagnostics create all dcosreport
    LOCAL:dcos node diagnostics download dcosreport
    REMOTE:hostname -f
    REMOTE:cat /proc/version
    
    $ clump -u core -pk /Users/mhausenblas/.ssh/test -nl test/clusternodes -cmds test/snapshot.cmds
    Trying to establish node list from test/clusternodes
    Got 1 target nodes
    Trying to establish list of commands from test/snapshot.cmds
    Got 2 commands to execute
    Attempting to ssh into core@35.160.157.251
    Attempting to ssh into core@35.160.157.251

Kudos go out to [Svett Ralchev](http://blog.ralch.com/tutorial/golang-ssh-connection/) for the seed code base around the SSH client.

## Disclaimer

> THIS SOFTWARE IS PROVIDED "AS IS" AND ANY EXPRESSED OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
