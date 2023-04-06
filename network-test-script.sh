# This is a sample script which is to run on the router itself. On USG it can be set up as follows:

INTERFACE=eth0                            # Name of the interface we expect to fail - i.e. not the failover one
HOSTS="google.com yahoo.com facebook.com" # List of hosts to ping
COUNT=10                                  # How many times to ping each host
MAX_LATENCY=100                           # Maximum average latency to expect across all hosts

# Let's determine the source IP we'll be using to ping the hosts from $INTERFACE provided
SOURCE_IP=$(ip -f inet addr show $INTERFACE | sed -En -e 's/.*inet ([0-9.]+).*/\1/p')
echo "Source IP: $SOURCE_IP"

# Run the program
/config/scripts/network-status -source-ip $SOURCE_IP -max-latency $MAX_LATENCY -count $COUNT $HOSTS

# exit with return code of the last command
exit $?