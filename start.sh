#!/bin/bash

# start the SIMAP server in the background
./main &

# keep the container running by starting a Bash shell
exec /bin/bash