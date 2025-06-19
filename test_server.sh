#!/bin/bash

echo "Galaxy Game Server Test"
echo "======================="

SERVER_URL="http://localhost:8080"

echo "1. Testing server status..."
curl -s "$SERVER_URL/status" | jq '.' || echo "Server not running or jq not available"

echo -e "\n2. Testing player connection..."
curl -s -X POST "$SERVER_URL/connect" \
  -H "Content-Type: application/json" \
  -d '{"player_id": "player1"}' | jq '.' || echo "Raw response:"

echo -e "\n3. Testing order submission..."
curl -s -X POST "$SERVER_URL/orders" \
  -H "Content-Type: application/json" \
  -d '{
    "player_id": "player1",
    "order_type": "BUILD_FACILITY",
    "planet_id": "planet_player1_home",
    "parameters": {"facility_type": "MetalMine"},
    "priority": 5
  }' | jq '.' || echo "Raw response:"

echo -e "\n4. Testing player status..."
curl -s "$SERVER_URL/player/player1" | jq '.' || echo "Raw response:"

echo -e "\n5. Testing game state..."
curl -s "$SERVER_URL/game" | jq '.data.players' || echo "Raw response:"

echo -e "\nTest complete!"