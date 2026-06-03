#!/bin/bash
# Comprehensive API Test Script for NetConfigHub
# Tests ALL API endpoints on http://localhost:8080

BASE="http://localhost:8080"
PASS=0
FAIL=0
FAILURES=""

test_endpoint() {
    local method="$1"
    local url="$2"
    local expected_status="$3"
    local auth="$4"
    local data="$5"
    local desc="$6"
    
    local extra_args=""
    if [ "$method" = "POST" ] || [ "$method" = "PUT" ]; then
        if [ -n "$data" ]; then
            extra_args="-H \"Content-Type: application/json\" -d '$data'"
        fi
    fi
    
    if [ "$auth" = "token" ]; then
        extra_args="$extra_args -H \"Authorization: Bearer $TOKEN\""
    elif [ "$auth" = "api_token" ]; then
        extra_args="$extra_args -H \"Authorization: Bearer $API_TOKEN_RAW\""
    elif [ "$auth" = "bad_token" ]; then
        extra_args="$extra_args -H \"Authorization: Bearer invalid_token_12345\""
    fi
    
    local curl_cmd="curl -s -o /tmp/api_response.txt -w \"%{http_code}\" -X $method \"$url\" $extra_args"
    
    local actual_status=$(eval $curl_cmd 2>/dev/null)
    local resp_body=$(cat /tmp/api_response.txt 2>/dev/null)
    local resp_trunc=$(echo "$resp_body" | head -c 250)
    
    if [ "$actual_status" -eq "$expected_status" ]; then
        echo "  PASS: $method $url -> $actual_status (expected $expected_status)"
        echo "  RESP: $resp_trunc"
        PASS=$((PASS + 1))
    else
        echo "  FAIL: $method $url -> $actual_status (expected $expected_status)"
        echo "  RESP: $resp_trunc"
        FAIL=$((FAIL + 1))
        FAILURES="$FAILURES\n  $method $url: got $actual_status, expected $expected_status"
    fi
    echo ""
}

extract_json() {
    python3 -c "
import sys,json
try:
    d=json.load(sys.stdin)
    # Try data.token, data.id, then direct
    if isinstance(d,dict):
        if 'data' in d and isinstance(d['data'],dict):
            print(d['data'].get('$1',''))
        else:
            print(d.get('$1',''))
    else:
        print('')
except:
    print('')
" 2>/dev/null
}

echo "=============================================="
echo "  NetConfigHub Comprehensive API Test Report"
echo "  Started: $(date)"
echo "=============================================="
echo ""

############################################
# STEP 1: Login and get token
############################################
echo ">>> STEP 1: LOGIN"
echo "=================="

TOKEN=$(curl -s -X POST "$BASE/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin"}' | extract_json "token")

if [ -z "$TOKEN" ]; then
    echo "FATAL: Could not obtain auth token. Aborting."
    exit 1
fi
echo "Token obtained: ${TOKEN:0:30}..."
echo ""

############################################
# Health
############################################
echo ">>> HEALTH"
echo "============"
test_endpoint "GET" "$BASE/api/v1/health" "200" "none" "" "" "Health check"

############################################
# Auth
############################################
echo ">>> AUTH"
echo "=========="
test_endpoint "POST" "$BASE/api/v1/auth/login" "200" "none" '{"username":"admin","password":"admin"}' "" "Login correct"
test_endpoint "POST" "$BASE/api/v1/auth/login" "401" "none" '{"username":"admin","password":"wrongpass"}' "" "Login wrong pass"

# Refresh - try with Bearer token
test_endpoint "POST" "$BASE/api/v1/auth/refresh" "200" "token" "" "" "Refresh token (Bearer)"

test_endpoint "POST" "$BASE/api/v1/auth/logout" "200" "token" "" "" "Logout"
# Re-login for subsequent tests
TOKEN=$(curl -s -X POST "$BASE/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin"}' | extract_json "token")

############################################
# Admins
############################################
echo ">>> ADMINS"
echo "============"
test_endpoint "GET" "$BASE/api/v1/admins" "200" "token" "" "" "List admins"
test_endpoint "POST" "$BASE/api/v1/admins" "201" "token" '{"username":"test","password":"test123"}' "" "Create admin test"
test_endpoint "GET" "$BASE/api/v1/admins/1" "200" "token" "" "" "Get admin 1"
test_endpoint "PUT" "$BASE/api/v1/admins/1" "200" "token" '{"username":"updated"}' "" "Update admin 1"
test_endpoint "DELETE" "$BASE/api/v1/admins/2" "200" "token" "" "" "Delete admin 2"

# Re-create test admin (was deleted)
test_endpoint "POST" "$BASE/api/v1/admins" "201" "token" '{"username":"test2","password":"test123"}' "" "Re-create admin test2"

############################################
# Tokens
############################################
echo ">>> TOKENS"
echo "============"
test_endpoint "POST" "$BASE/api/v1/tokens" "201" "token" '{"name":"test-token"}' "" "Create API token #1"

# Capture API token for ext tests
API_TOKEN_RAW=$(curl -s -X POST "$BASE/api/v1/tokens" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"name":"ext-token"}' | extract_json "token")
echo "  API Token (ext): ${API_TOKEN_RAW:-N/A}"

test_endpoint "GET" "$BASE/api/v1/tokens" "200" "token" "" "" "List tokens"

# Get first token ID for delete
TOKEN_RESP=$(curl -s -X GET "$BASE/api/v1/tokens" -H "Authorization: Bearer $TOKEN")
TOKEN_ID=$(echo "$TOKEN_RESP" | extract_json "id")
# Try array format
if [ -z "$TOKEN_ID" ]; then
    TOKEN_ID=$(echo "$TOKEN_RESP" | python3 -c "
import sys,json
d=json.load(sys.stdin)
items=d.get('data',d)
if isinstance(items,list) and items:
    print(items[0].get('id',''))
elif isinstance(items,dict) and 'tokens' in items:
    t=items['tokens']
    print(t[0].get('id','')) if t else print('')
" 2>/dev/null)
fi
echo "  Token ID to delete: ${TOKEN_ID:-N/A}"
if [ -n "$TOKEN_ID" ] && [ "$TOKEN_ID" != "null" ]; then
    test_endpoint "DELETE" "$BASE/api/v1/tokens/$TOKEN_ID" "200" "token" "" "" "Delete token $TOKEN_ID"
else
    echo "  SKIP: Could not determine token ID"
    echo ""
fi

############################################
# Devices
############################################
echo ">>> DEVICES"
echo "============="
test_endpoint "POST" "$BASE/api/v1/devices" "201" "token" '{"name":"test-device","ip_address":"192.168.1.100","vendor":"cisco","model":"Catalyst","conn_protocol":"ssh","port":22,"username":"admin","password":"cisco123"}' "" "Create device"
test_endpoint "GET" "$BASE/api/v1/devices" "200" "token" "" "" "List devices"
test_endpoint "GET" "$BASE/api/v1/devices?search=test" "200" "token" "" "" "Search devices"
test_endpoint "GET" "$BASE/api/v1/devices/1" "200" "token" "" "" "Get device 1"
test_endpoint "PUT" "$BASE/api/v1/devices/1" "200" "token" '{"name":"test-device-updated","vendor":"juniper"}' "" "Update device 1"
test_endpoint "POST" "$BASE/api/v1/devices/1/backup" "200" "token" "" "" "Trigger backup device 1"
test_endpoint "GET" "$BASE/api/v1/devices/1/config" "200" "token" "" "" "Get device 1 config"
test_endpoint "GET" "$BASE/api/v1/devices/1/versions" "200" "token" "" "" "Get device 1 versions"
test_endpoint "GET" "$BASE/api/v1/devices/1/deviations" "200" "token" "" "" "Get device 1 deviations"
test_endpoint "GET" "$BASE/api/v1/devices/1/diff?from=abc&to=def" "200" "token" "" "" "Get dev 1 diff"
test_endpoint "GET" "$BASE/api/v1/devices/export" "200" "token" "" "" "Export devices"

# Import test
echo "  --- Device Import (multipart/form-data) ---"
cat > /tmp/test_import.csv << 'CSVEOF'
name,ip_address,vendor,model,conn_protocol,port,username,password
imported-dev,10.10.10.1,cisco,2950,ssh,22,admin,cisco123
CSVEOF
IMPORT_STATUS=$(curl -s -o /tmp/api_response.txt -w "%{http_code}" -X POST "$BASE/api/v1/devices/import" \
    -H "Authorization: Bearer $TOKEN" \
    -F "file=@/tmp/test_import.csv" 2>/dev/null)
IMP_RESP=$(head -c 250 /tmp/api_response.txt)
echo "  IMPORT: Status=$IMPORT_STATUS"
echo "  RESP: $IMP_RESP"
if [ "$IMPORT_STATUS" = "200" ] || [ "$IMPORT_STATUS" = "201" ]; then
    echo "  PASS: POST /api/v1/devices/import -> $IMPORT_STATUS"
    PASS=$((PASS + 1))
else
    echo "  FAIL: POST /api/v1/devices/import -> $IMPORT_STATUS (expected 200 or 201)"
    FAIL=$((FAIL + 1))
fi
echo ""

test_endpoint "DELETE" "$BASE/api/v1/devices/1" "200" "token" "" "" "Delete device 1"

############################################
# Groups
############################################
echo ">>> GROUPS"
echo "============"
test_endpoint "POST" "$BASE/api/v1/groups" "201" "token" '{"name":"test-group","description":"A test group"}' "" "Create group"
test_endpoint "GET" "$BASE/api/v1/groups" "200" "token" "" "" "List groups"
test_endpoint "GET" "$BASE/api/v1/groups/1" "200" "token" "" "" "Get group 1"
test_endpoint "PUT" "$BASE/api/v1/groups/1" "200" "token" '{"name":"test-group-updated","description":"Updated desc"}' "" "Update group 1"
test_endpoint "POST" "$BASE/api/v1/groups/1/backup" "200" "token" "" "" "Trigger group 1 backup"
test_endpoint "DELETE" "$BASE/api/v1/groups/1" "200" "token" "" "" "Delete group 1"

############################################
# Backups
############################################
echo ">>> BACKUPS"
echo "============="

# Create a device and trigger backup first for backup data
curl -s -X POST "$BASE/api/v1/devices" -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" \
    -d '{"name":"backup-dev","ip_address":"10.0.0.1","vendor":"cisco","conn_protocol":"ssh","port":22,"username":"admin","password":"cisco"}' > /dev/null 2>&1
curl -s -X POST "$BASE/api/v1/devices/2/backup" -H "Authorization: Bearer $TOKEN" > /dev/null 2>&1
# Small sleep to let backup "complete"
sleep 1

test_endpoint "GET" "$BASE/api/v1/backups" "200" "token" "" "" "List backups"
test_endpoint "GET" "$BASE/api/v1/backups/1" "200" "token" "" "" "Get backup 1"
test_endpoint "GET" "$BASE/api/v1/backups/1/config" "200" "token" "" "" "Get backup 1 config"

############################################
# Baselines
############################################
echo ">>> BASELINES"
echo "==============="
test_endpoint "POST" "$BASE/api/v1/baselines" "201" "token" '{"scope":"device","content":"hostname router1\ninterface eth0\n ip address 10.0.0.1/24","device_id":2}' "" "Create baseline"
test_endpoint "GET" "$BASE/api/v1/baselines" "200" "token" "" "" "List baselines"
test_endpoint "GET" "$BASE/api/v1/baselines/1" "200" "token" "" "" "Get baseline 1"
test_endpoint "PUT" "$BASE/api/v1/baselines/1" "200" "token" '{"content":"hostname router1-updated\ninterface eth0\n ip address 10.0.0.2/24"}' "" "Update baseline 1"
test_endpoint "DELETE" "$BASE/api/v1/baselines/1" "200" "token" "" "" "Delete baseline 1"

############################################
# Deviations
############################################
echo ">>> DEVIATIONS"
echo "================"
test_endpoint "GET" "$BASE/api/v1/deviations" "200" "token" "" "" "List deviations"
test_endpoint "GET" "$BASE/api/v1/deviations/1" "200" "token" "" "" "Get deviation 1 (may 404)"
test_endpoint "GET" "$BASE/api/v1/deviations/export" "200" "token" "" "" "Export deviations"

############################################
# Alerts
############################################
echo ">>> ALERTS"
echo "============"
test_endpoint "GET" "$BASE/api/v1/alerts" "200" "token" "" "" "List alerts"
test_endpoint "GET" "$BASE/api/v1/alerts/unread-count" "200" "token" "" "" "Unread count"

# Try to find alert ID
ALERT_RESP=$(curl -s -X GET "$BASE/api/v1/alerts" -H "Authorization: Bearer $TOKEN")
ALERT_ID=$(echo "$ALERT_RESP" | python3 -c "
import sys,json
d=json.load(sys.stdin)
items=d.get('data',d)
if isinstance(items,list) and items:
    print(items[0].get('id',''))
" 2>/dev/null)
if [ -n "$ALERT_ID" ] && [ "$ALERT_ID" != "null" ]; then
    test_endpoint "PUT" "$BASE/api/v1/alerts/$ALERT_ID/read" "200" "token" "" "" "Mark alert $ALERT_ID read"
    test_endpoint "DELETE" "$BASE/api/v1/alerts/$ALERT_ID" "200" "token" "" "" "Delete alert $ALERT_ID"
else
    echo "  SKIP: No alerts found"
    echo ""
fi

test_endpoint "PUT" "$BASE/api/v1/alerts/read-all" "200" "token" "" "" "Mark all alerts read"

############################################
# Hooks
############################################
echo ">>> HOOKS"
echo "==========="
test_endpoint "POST" "$BASE/api/v1/hooks" "201" "token" '{"name":"test-webhook","type":"webhook","config":{"url":"http://example.com/hook"},"events":["backup_completed"],"enabled":true}' "" "Create webhook"
test_endpoint "GET" "$BASE/api/v1/hooks" "200" "token" "" "" "List hooks"
test_endpoint "GET" "$BASE/api/v1/hooks/1" "200" "token" "" "" "Get hook 1"
test_endpoint "PUT" "$BASE/api/v1/hooks/1" "200" "token" '{"name":"updated-webhook","enabled":false}' "" "Update hook 1"
test_endpoint "DELETE" "$BASE/api/v1/hooks/1" "200" "token" "" "" "Delete hook 1"

############################################
# Configs
############################################
echo ">>> CONFIGS"
echo "============="
test_endpoint "GET" "$BASE/api/v1/configs" "200" "token" "" "" "List configs"
test_endpoint "PUT" "$BASE/api/v1/configs" "200" "token" '{"key":"test","value":"val"}' "" "Set config test=val"
test_endpoint "GET" "$BASE/api/v1/configs/test" "200" "token" "" "" "Get config test"

############################################
# Failed Devices
############################################
echo ">>> FAILED DEVICES"
echo "===================="
test_endpoint "GET" "$BASE/api/v1/failed-devices" "200" "token" "" "" "List failed devices"

############################################
# Ext API (use API token)
############################################
echo ">>> EXT API (API Token Auth)"
echo "============================="
if [ -n "$API_TOKEN_RAW" ] && [ "$API_TOKEN_RAW" != "null" ]; then
    test_endpoint "GET" "$BASE/api/v1/ext/devices" "200" "api_token" "" "" "Ext list devices"
    test_endpoint "GET" "$BASE/api/v1/ext/backups" "200" "api_token" "" "" "Ext list backups"
    test_endpoint "GET" "$BASE/api/v1/ext/baselines" "200" "api_token" "" "" "Ext list baselines"
else
    echo "  SKIP: No API token available"
    echo ""
fi

############################################
# Security Tests
############################################
echo ">>> SECURITY TESTS"
echo "===================="
echo "--- Without auth (expect 401) ---"
test_endpoint "GET" "$BASE/api/v1/admins" "401" "none" "" "" "No auth - admins"
test_endpoint "GET" "$BASE/api/v1/devices" "401" "none" "" "" "No auth - devices"
test_endpoint "GET" "$BASE/api/v1/groups" "401" "none" "" "" "No auth - groups"

echo "--- Invalid token (expect 401) ---"
test_endpoint "GET" "$BASE/api/v1/admins" "401" "bad_token" "" "" "Bad token - admins"
test_endpoint "GET" "$BASE/api/v1/devices" "401" "bad_token" "" "" "Bad token - devices"

echo "--- Password not in responses ---"
PWD_RESULT=$(curl -s -X GET "$BASE/api/v1/admins/1" -H "Authorization: Bearer $TOKEN" 2>/dev/null)
if echo "$PWD_RESULT" | grep -qi "password"; then
    echo "  FAIL: Password found in response!"
    echo "  RESP: $PWD_RESULT"
    FAIL=$((FAIL + 1))
    FAILURES="$FAILURES\n  SECURITY: Password found in admin response"
else
    echo "  PASS: No password in admin response"
    PASS=$((PASS + 1))
fi
echo ""

############################################
# Summary
############################################
echo "=============================================="
echo "  COMPREHENSIVE API TEST SUMMARY"
echo "=============================================="
echo "  PASS: $PASS"
echo "  FAIL: $FAIL"
echo "  TOTAL: $((PASS + FAIL))"
echo "  PASS RATE: $(echo "scale=1; $PASS * 100 / ($PASS + $FAIL)" | bc)%"
echo ""

if [ "$FAIL" -gt 0 ]; then
    echo "=============================================="
    echo "  FAILURES:"
    echo "=============================================="
    echo -e "$FAILURES"
    echo ""
fi

echo "Test completed at: $(date)"