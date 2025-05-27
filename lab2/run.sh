#!/bin/sh

SHARED_DIR="/shared"

LOCK_FILE="$SHARED_DIR/.lock"

CONTAINER_ID=$(hostname)

FILE_COUNTER=1

while true; do
  (
    flock -x 200

    for i in $(seq 1 999); do
      FILENAME=$(printf "%03d" $i)
      if [ ! -f "$SHARED_DIR/$FILENAME" ]; then
        break
      fi
    done

    FILEPATH="$SHARED_DIR/$FILENAME"

    echo "$CONTAINER_ID:$FILE_COUNTER" > "$FILEPATH"

    FILE_COUNTER=$((FILE_COUNTER + 1))

  ) 200>"$LOCK_FILE"

  sleep 1

  if [ -f "$FILEPATH" ]; then
    rm -f "$FILEPATH"
  fi

  sleep 1
done