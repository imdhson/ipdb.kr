#!/bin/bash

# Exit on any error
set -e

APP_NAME="ipdb"
PID_FILE="app.pid"

echo "=== Building $APP_NAME ==="
go build -o $APP_NAME main.go
echo "Build successful."

echo "=== Deploying $APP_NAME ==="

# Stop the running application if it exists
if [ -f "$PID_FILE" ]; then
    OLD_PID=$(cat "$PID_FILE")
    if ps -p $OLD_PID > /dev/null; then
        echo "Stopping existing instance (PID: $OLD_PID)..."
        kill $OLD_PID
        sleep 2
    else
        echo "Stale PID file found. Removing..."
    fi
    rm "$PID_FILE"
else
    # Try to find existing process just in case
    OLD_PID=$(pgrep -f "./$APP_NAME" || echo "")
    if [ ! -z "$OLD_PID" ]; then
        echo "Stopping existing instance found by pgrep (PID: $OLD_PID)..."
        kill $OLD_PID
        sleep 2
    fi
fi

# Run the application in the background
echo "Starting new instance..."
nohup ./$APP_NAME > app.log 2>&1 &
NEW_PID=$!
echo $NEW_PID > "$PID_FILE"

echo "Deployment complete! Application is running (PID: $NEW_PID)."
echo "Check app.log for output."
