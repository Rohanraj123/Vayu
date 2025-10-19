#!/bin/bash

# Exit immediately on error
set -e

# Colors for clarity
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW} Starting mock backend servers....${NC}"

# Start mock services in backend
go run tests/test-routing/backend-users/test_server.go > tests/test-routing/users.log 2>&1 & USERS_PID=$!
go run tests/test-routing/backend-orders/test_server.go > tests/test-routing/orders.log 2>&1 & ORDERS_PID=$!

# Wait a bit for servers to start
sleep 2

echo -e "${YELLOW}🚀 Starting Vayu API Gateway...${NC}"
go run cmd/vayu/main.go tests/test-routing/test_config.yaml > tests/test-routing/vayu.log 2>&1 &
VAYU_PID=$!

# Wait for gateway to start
sleep 2

echo -e "${YELLOW}📡 Running test requests...${NC}"

# Test endpoints
echo -e "${GREEN}→ /users${NC}"
curl -s -H "sajlksf: sajlksf" -o /dev/null -w "Status: %{http_code}, Time: %{time_total}s\n" http://localhost:8087/users

echo -e "${GREEN}-> /orders${NC}"
curl -s -o /dev/null -w "Status: %{http_code}, Time: %{time_total}s\n" http://localhost:8087/orders

echo -e "${YELLOW}🧾 Gateway logs:${NC}"
tail -n 10 tests/test-routing/vayu.log

echo -e "${YELLOW}✅ All tests completed.${NC}"

# Cleanup on exit
trap "echo -e '\n🧹 Cleaning up...'; kill $USERS_PID $ORDERS_PID $VAYU_PID 2>/dev/null || true" EXIT

# Removal of temp files
# Remove below lines if you want to check the logs in the temp files. 
rm tests/test-routing/orders.log
rm tests/test-routing/vayu.log
rm tests/test-routing/users.log