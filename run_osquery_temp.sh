#!/bin/bash

# Configuration - Adjust these as needed
SOCKET_PATH="/tmp/osquery.$(whoami).$RANDOM.em"
OSQUERYD_BIN="/opt/osquery/lib/osquery.app/Contents/MacOS/osqueryd"
CONFIG_PATH="/var/osquery/osquery.conf" # Or use a custom path
LOG_DIR="/tmp/osquery_logs_$(whoami)"

# Cleanup previous runs
killall osqueryd 2>/dev/null
rm -f "$SOCKET_PATH"
mkdir -p "$LOG_DIR"

# Copy deamon configs
sudo cp /var/osquery/osquery.example.conf /var/osquery/osquery.conf
sudo cp /var/osquery/io.osquery.agent.plist /Library/LaunchDaemons
sudo launchctl load /Library/LaunchDaemons/io.osquery.agent.plist

# Start osqueryd with temporary settings
echo "Starting temporary osqueryd..."
echo "  - Socket: $SOCKET_PATH"
echo "  - Logs: $LOG_DIR"

"$OSQUERYD_BIN" \
  --database_path=/tmp/osquery_temp.db \
  --logger_path="$LOG_DIR" \
  --extensions_socket="$SOCKET_PATH" \
  --config_path="$CONFIG_PATH" \
  --disable_database \
  --ephemeral \
  --disable_audit \
  --disable_events &

# Wait for socket creation
MAX_WAIT=10
for ((i = 1; i <= $MAX_WAIT; i++)); do
  if [ -S "$SOCKET_PATH" ]; then
    break
  fi
  sleep 1
  echo -n "."
done

if [ ! -S "$SOCKET_PATH" ]; then
  echo -e "\n❌ Failed to create socket!"
  exit 1
fi

echo -e "\n✅ Ready! Use this socket path in your Go app:"
echo "$SOCKET_PATH"

# Cleanup on exit
trap "cleanup" EXIT
cleanup() {
  killall osqueryd 2>/dev/null
  rm -f "$SOCKET_PATH"
  echo "Cleaned up temporary osqueryd"
}

# Keep running until Ctrl+C
echo "Press Ctrl+C to stop osqueryd"
wait
