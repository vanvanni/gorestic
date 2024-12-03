#!/bin/bash

# API KEY
API_KEY="e20b98bf9144ad5ef093fdcdb439dedd"

# Send test data
curl -X POST "http://localhost:2010/api/stats" \
    -H "X-API-Key: ${API_KEY}" \
    -H "Content-Type: application/json" \
    -d '{
        "total_size": 636249527,
        "total_file_count": 1842,
        "snapshots_count": 1
    }'

echo -e "\nTest data sent!"