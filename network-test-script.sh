# This is a sample script which is to run on the router itself. On USG it can be set up as follows:

INTERFACE=eth0
HOSTS="google.com yahoo.com facebook.com"
COUNT=10
MAX_LATENCY=100

SOURCE_IP=$(ip -f inet addr show $INTERFACE | sed -En -e 's/.*inet ([0-9.]+).*/\1/p')
echo "Source IP: $SOURCE_IP"

/config/scripts/network-status -source-ip $SOURCE_IP -max-latency $MAX_LATENCY -count $COUNT $HOSTS

# exit with return code of the last command
exit $?