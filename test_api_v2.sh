#!/bin/bash
# Comprehensive API Test Script for NetConfigHub v2
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

    local resp_body=""
    local actual_status=""
    local tmpfile=$(mktemp)
    
    if [ "$method" = "DELETE" ] || ([ "$method" = "POST" ] && [ -z "$data" ]); then
        if [ -n "$data" ]; then
            actual_status=$(curl -s -o "$tmpfile" -w "%{http_code}" -X "$method" "$url" -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d "$data" 2>/dev/null)
        else
            actual_status=$(eval curl -s -o "$tmpfile" -w '"%{http_code}"' -X "$method" '"$url"' $extra_args 2>/dev/null)
        fi
    else
        actual_status=$(eval curl -s -o "$tmpfile" -w '"%{http_code}"' -X "$method" '"$url"' $extra_args 2>/dev/null)
    fi
    
    resp_body=$(head -c 300 "$tmpfile" 2>/dev/null)
    rm -f "$tmpfile"
    
    actual_status=$(echo "$actual_status" | tr -d '\n\r')
    
    if [ "$actual_status" = "$expected_status" ]; then
        echo "  PASS: $method $url -> $actual_status (expected $expected_status)"
        echo "  RESP: $resp_body"
        PASS=$((PASS + 1))
    else
        echo "  FAIL: $method $url -> $actual_status (expected $expected_status)"
        echo "  RESP: $resp_body"
        FAIL=$((FAIL + 1))
        FAILURES="$FAILURES\n  $method $url [$desc]: got $actual_status, expected $expected_status"
    fi
    echo ""
}

extract_field() {
    python3 -c "
import sys,json
try:
    d=json.load(sys.stdin)
    if isinstance(d,dict):
        # Try data -> field
        dd = d.get('data',{})
        if isinstance(dd,dict):
            print(dd.get('$1',''))
        elif isinstance(dd,list) and dd:
            print(dd[0].get('$1',''))
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
    -d '{"username":"admin","password":"admin"}' | extract_field "token")

if [ -z "$TOKEN" ]; then
    echo "FATAL: Could not obtain auth token. Aborting."
    exit 1
fi
echo "Token obtained: ${TOKEN:0:30}..."
echo ""

############################################
# Health
############################################
echo ">>> 1. HEALTH"
echo "============"
test_endpoint "GET" "$BASE/api/v1/health" "200" "none" "" "Health check"

############################################
# Auth
############################################
echo ">>> 2. AUTH"
echo "=========="
test_endpoint "POST" "$BASE/api/v1/auth/login" "200" "none" '{"username":"admin","password":"admin"}' "Login correct"
test_endpoint "POST" "$BASE/api/v1/auth/login" "401" "none" '{"username":"admin","password":"wrongpass"}' "Login wrong pass"
test_endpoint "POST" "$BASE/api/v1/auth/refresh" "200" "token" "" "Refresh token"
test_endpoint "POST" "$BASE/api/v1/auth/logout" "200" "token" "" "Logout"

# Re-login for subsequent tests
TOKEN=$(curl -s -X POST "$BASE/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin"}' | extract_field "token")

############################################
# Admins
############################################
echo ">>> 3. ADMINS"
echo "============"
test_endpoint "GET" "$BASE/api/v1/admins" "200" "token" "" "List admins"
test_endpoint "POST" "$BASE/api/v1/admins" "200" "token" '{"username":"testuser1","password":"test123"}' "Create admin"
test_endpoint "GET" "$BASE/api/v1/admins/1" "200" "token" "" "Get admin 1"
test_endpoint "PUT" "$BASE/api/v1/admins/1" "200" "token" '{"username":"admin-updated"}' "Update admin 1"
test_endpoint "DELETE" "$BASE/api/v1/admins/2" "200" "token" "" "Delete admin 2"
# Re-create test admin
test_endpoint "POST" "$BASE/api/v1/admins" "200" "token" '{"username":"testuser2","password":"test123"}' "Re-create admin"

############################################
# Tokens
############################################
echo ">>> 4. API TOKENS"
echo "================"
test_endpoint "POST" "$BASE/api/v1/tokens" "200" "token" '{"name":"test-token"}' "Create API token"

# Capture API token for ext tests
API_TOKEN_RAW=$(curl -s -X POST "$BASE/api/v1/tokens" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"name":"ext-token"}' | extract_field "token")
echo "  API Token (ext): ${API_TOKEN_RAW:-N/A}"

test_endpoint "GET" "$BASE/api/v1/tokens" "200" "token" "" "List tokens"

# Get first token ID for delete
TOKEN_LIST=$(curl -s -X GET "$BASE/api/v1/tokens" -H "Authorization: Bearer $TOKEN")
TOKEN_ID=$(echo "$TOKEN_LIST" | python3 -c "
import sys,json
d=json.load(sys.stdin)
items=d.get('data',d)
if isinstance(items,dict) and 'items' in items:
    items=items['items']
if isinstance(items,list) and items:
    print(items[0].get('id',''))
" 2>/dev/null)

if [ -n "$TOKEN_ID" ] && [ "$TOKEN_ID" != "null" ]; then
    test_endpoint "DELETE" "$BASE/api/v1/tokens/$TOKEN_ID" "200" "token" "" "Delete token"
else
    echo "  SKIP: Could not determine token ID"
    echo ""
fi

############################################
# Devices (using new field names)
############################################
echo ">>> 5. DEVICES"
echo "============="
# Create device (correct field names)
test_endpoint "POST" "$BASE/api/v1/devices" "200" "token" \
    '{"name":"core-sw-01","ip":"192.168.1.100","vendor":"cisco","model":"ios","protocol":"ssh","port":22,"username":"admin","password":"cisco123"}' \
    "Create device"

test_endpoint "GET" "$BASE/api/v1/devices" "200" "token" "" "List devices"
test_endpoint "GET" "$BASE/api/v1/devices?search=core" "200" "token" "" "Search devices"
test_endpoint "GET" "$BASE/api/v1/devices?vendor=cisco" "200" "token" "" "Filter by vendor"
test_endpoint "GET" "$BASE/api/v1/devices/1" "200" "token" "" "Get device 1"
test_endpoint "PUT" "$BASE/api/v1/devices/1" "200" "token" \
    '{"name":"core-sw-01-updated","vendor":"huawei","model":"vrp"}' \
    "Update device 1"
test_endpoint "PUT" "$BASE/api/v1/devices/1/enable" "200" "token" "" "Enable device 1"
test_endpoint "PUT" "$BASE/api/v1/devices/1/disable" "200" "token" "" "Disable device 1"
test_endpoint "PUT" "$BASE/api/v1/devices/1/enable" "200" "token" "" "Re-enable device 1"

# Trigger backup
test_endpoint "POST" "$BASE/api/v1/devices/1/backup" "200" "token" "" "Trigger backup device 1"

# Config and versions: the dummy device cannot be collected in smoke tests,
# so no stored config should be reported as 404.
test_endpoint "GET" "$BASE/api/v1/devices/1/config" "404" "token" "" "Get device 1 config"
test_endpoint "GET" "$BASE/api/v1/devices/1/versions" "404" "token" "" "Get device 1 versions"
test_endpoint "GET" "$BASE/api/v1/devices/1/versions?limit=10" "404" "token" "" "Get device 1 versions limit"
test_endpoint "GET" "$BASE/api/v1/devices/1/deviations" "200" "token" "" "Get device 1 deviations"

# Diff
test_endpoint "GET" "$BASE/api/v1/devices/1/diff-latest" "200" "token" "" "Get dev 1 latest diff"
test_endpoint "GET" "$BASE/api/v1/devices/1/diff?from=abc&to=def" "200" "token" "" "Get dev 1 diff"
test_endpoint "GET" "$BASE/api/v1/devices/1/diff-download?from=abc&to=def" "200" "token" "" "Download diff"

# Export
test_endpoint "GET" "$BASE/api/v1/devices/export" "200" "token" "" "Export devices CSV"

# Import
echo "  --- Device Import (multipart/form-data) ---"
cat > /tmp/test_import.csv << 'CSVEOF'
name,ip,vendor,model,protocol,port,username,password,group_name
imported-dev-1,10.10.10.10,cisco,ios,ssh,22,admin,cisco123,
imported-dev-2,10.10.10.11,huawei,vrp,telnet,23,admin,huawei123,
CSVEOF
IMPORT_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/api/v1/devices/import" \
    -H "Authorization: Bearer $TOKEN" \
    -F "file=@/tmp/test_import.csv" 2>/dev/null)
if [ "$IMPORT_STATUS" = "200" ]; then
    echo "  PASS: POST /api/v1/devices/import -> $IMPORT_STATUS"
    PASS=$((PASS + 1))
else
    echo "  FAIL: POST /api/v1/devices/import -> $IMPORT_STATUS (expected 200)"
    FAIL=$((FAIL + 1))
    FAILURES="$FAILURES\n  POST /api/v1/devices/import: got $IMPORT_STATUS, expected 200"
fi
echo ""

# Validation: duplicate name
test_endpoint "POST" "$BASE/api/v1/devices" "400" "token" \
    '{"name":"core-sw-01-updated","ip":"192.168.1.200","vendor":"cisco","model":"ios","protocol":"ssh","port":22,"username":"admin","password":"cisco123"}' \
    "Create device (duplicate name -> 400)"
# Validation: bad IP
test_endpoint "POST" "$BASE/api/v1/devices" "400" "token" \
    '{"name":"bad-ip-device","ip":"not-an-ip","vendor":"cisco","model":"ios","protocol":"ssh","port":22,"username":"admin","password":"cisco123"}' \
    "Create device (bad IP -> 400)"

############################################
# Groups
############################################
echo ">>> 6. GROUPS"
echo "============"
test_endpoint "POST" "$BASE/api/v1/groups" "200" "token" \
    '{"name":"test-group","description":"A test group"}' "Create group"
test_endpoint "GET" "$BASE/api/v1/groups" "200" "token" "" "List groups"
test_endpoint "GET" "$BASE/api/v1/groups/1" "200" "token" "" "Get group 1"
test_endpoint "PUT" "$BASE/api/v1/groups/1" "200" "token" \
    '{"name":"test-group-updated","description":"Updated desc"}' "Update group 1"
test_endpoint "POST" "$BASE/api/v1/groups/1/backup" "400" "token" "" "Trigger empty group backup"
test_endpoint "DELETE" "$BASE/api/v1/groups/1" "200" "token" "" "Delete group 1"

############################################
# Backups
############################################
echo ">>> 7. BACKUPS"
echo "============="
test_endpoint "GET" "$BASE/api/v1/backups" "200" "token" "" "List backups"
test_endpoint "GET" "$BASE/api/v1/backups?device_id=1" "200" "token" "" "List backups device 1"
test_endpoint "GET" "$BASE/api/v1/backups?status=failed" "200" "token" "" "List backups failed"

# Get backup ID
BACKUP_RESP=$(curl -s -X GET "$BASE/api/v1/backups" -H "Authorization: Bearer $TOKEN")
BACKUP_ID=$(echo "$BACKUP_RESP" | extract_field "id")

# Also try with jobs endpoint
test_endpoint "GET" "$BASE/api/v1/jobs" "200" "token" "" "List jobs"

if [ -n "$BACKUP_ID" ] && [ "$BACKUP_ID" != "null" ]; then
    test_endpoint "GET" "$BASE/api/v1/backups/$BACKUP_ID" "200" "token" "" "Get backup $BACKUP_ID"
    test_endpoint "GET" "$BASE/api/v1/backups/$BACKUP_ID/config" "200" "token" "" "Get backup $BACKUP_ID config"
else
    echo "  SKIP: No backup tasks found"
    echo ""
fi

############################################
# Baselines
############################################
echo ">>> 8. BASELINES"
echo "==============="
# Create baseline for device 1
test_endpoint "POST" "$BASE/api/v1/baselines" "200" "token" \
    '{"scope":"device","content":"hostname router1\ninterface eth0\n ip address 10.0.0.1/24","device_id":1}' "Create baseline"
test_endpoint "GET" "$BASE/api/v1/baselines" "200" "token" "" "List baselines"
test_endpoint "GET" "$BASE/api/v1/baselines/1" "200" "token" "" "Get baseline 1"
test_endpoint "PUT" "$BASE/api/v1/baselines/1" "200" "token" \
    '{"content":"hostname router1-updated"}' "Update baseline 1"
test_endpoint "DELETE" "$BASE/api/v1/baselines/1" "200" "token" "" "Delete baseline 1"

############################################
# Deviations
############################################
echo ">>> 9. DEVIATIONS"
echo "================"
test_endpoint "GET" "$BASE/api/v1/deviations" "200" "token" "" "List deviations"
test_endpoint "GET" "$BASE/api/v1/deviations?device_id=1" "200" "token" "" "List deviations device 1"
test_endpoint "GET" "$BASE/api/v1/deviations/export" "200" "token" "" "Export deviations"

############################################
# Alerts
############################################
echo ">>> 10. ALERTS"
echo "============"
test_endpoint "GET" "$BASE/api/v1/alerts" "200" "token" "" "List alerts"
test_endpoint "GET" "$BASE/api/v1/alerts/unread-count" "200" "token" "" "Unread count"

# Find alert ID
ALERT_RESP=$(curl -s -X GET "$BASE/api/v1/alerts" -H "Authorization: Bearer $TOKEN")
ALERT_ID=$(echo "$ALERT_RESP" | extract_field "id")
if [ -n "$ALERT_ID" ] && [ "$ALERT_ID" != "null" ]; then
    test_endpoint "PUT" "$BASE/api/v1/alerts/$ALERT_ID/read" "200" "token" "" "Mark alert read"
    test_endpoint "DELETE" "$BASE/api/v1/alerts/$ALERT_ID" "200" "token" "" "Delete alert"
else
    echo "  SKIP: No alerts found"
    echo ""
fi

test_endpoint "PUT" "$BASE/api/v1/alerts/read-all" "200" "token" "" "Mark all alerts read"

############################################
# Hooks
############################################
echo ">>> 11. HOOKS"
echo "==========="
test_endpoint "POST" "$BASE/api/v1/hooks" "200" "token" \
    '{"name":"test-webhook","type":"webhook","config":"http://example.com/hook","events":"backup_completed"}' \
    "Create webhook"
test_endpoint "GET" "$BASE/api/v1/hooks" "200" "token" "" "List hooks"
test_endpoint "GET" "$BASE/api/v1/hooks/1" "200" "token" "" "Get hook 1"
test_endpoint "PUT" "$BASE/api/v1/hooks/1" "200" "token" \
    '{"name":"updated-webhook","type":"webhook","config":"http://example.com/hook2","events":"backup_completed","enabled":false}' \
    "Update hook 1"
test_endpoint "DELETE" "$BASE/api/v1/hooks/1" "200" "token" "" "Delete hook 1"

############################################
# Configs
############################################
echo ">>> 12. CONFIGS"
echo "============="
test_endpoint "GET" "$BASE/api/v1/configs" "200" "token" "" "List configs"
test_endpoint "PUT" "$BASE/api/v1/configs" "200" "token" '{"key":"test-key","value":"test-value"}' "Set config"
test_endpoint "GET" "$BASE/api/v1/configs/test-key" "200" "token" "" "Get config test-key"

############################################
# Failed Devices
############################################
echo ">>> 13. FAILED DEVICES"
echo "===================="
test_endpoint "GET" "$BASE/api/v1/failed-devices" "200" "token" "" "List failed devices"

############################################
# Audit Logs
############################################
echo ">>> 14. AUDIT LOGS"
echo "================"
test_endpoint "GET" "$BASE/api/v1/audit-logs" "200" "token" "" "List audit logs"

############################################
# EXT API (API Token Auth)
############################################
echo ">>> 15. EXT API (API Token Auth)"
echo "============================="
if [ -n "$API_TOKEN_RAW" ] && [ "$API_TOKEN_RAW" != "null" ]; then
    test_endpoint "GET" "$BASE/api/v1/ext/devices" "200" "api_token" "" "Ext list devices"
    test_endpoint "GET" "$BASE/api/v1/ext/backups" "200" "api_token" "" "Ext list backups"
    test_endpoint "GET" "$BASE/api/v1/ext/baselines" "200" "api_token" "" "Ext list baselines"
    test_endpoint "GET" "$BASE/api/v1/ext/deviations" "200" "api_token" "" "Ext list deviations"
    test_endpoint "GET" "$BASE/api/v1/ext/alerts" "200" "api_token" "" "Ext list alerts"
    test_endpoint "GET" "$BASE/api/v1/ext/failed-devices" "200" "api_token" "" "Ext list failed devices"
    test_endpoint "POST" "$BASE/api/v1/ext/devices/1/backup" "200" "api_token" "" "Ext trigger backup"
    test_endpoint "GET" "$BASE/api/v1/ext/devices/1/config" "404" "api_token" "" "Ext get config"
else
    echo "  SKIP: No API token available"
    echo ""
fi

############################################
# SECURITY TESTS
############################################
echo ">>> 16. SECURITY TESTS"
echo "===================="
echo "--- Without auth (expect 401) ---"
test_endpoint "GET" "$BASE/api/v1/admins" "401" "none" "" "No auth - admins"
test_endpoint "GET" "$BASE/api/v1/devices" "401" "none" "" "No auth - devices"
test_endpoint "GET" "$BASE/api/v1/groups" "401" "none" "" "No auth - groups"

echo "--- Invalid token (expect 401) ---"
test_endpoint "GET" "$BASE/api/v1/admins" "401" "bad_token" "" "Bad token - admins"
test_endpoint "GET" "$BASE/api/v1/devices" "401" "bad_token" "" "Bad token - devices"

echo "--- Password not in responses ---"
PWD_RESULT=$(curl -s -X GET "$BASE/api/v1/admins/1" -H "Authorization: Bearer $TOKEN" 2>/dev/null)
if echo "$PWD_RESULT" | grep -qi '"password"'; then
    if echo "$PWD_RESULT" | grep -q '"password":"\*\*\*\*\*\*"'; then
        echo "  PASS: Password masked in admin response"
        PASS=$((PASS + 1))
    else
        echo "  FAIL: Real password found in response!"
        echo "  RESP: $PWD_RESULT"
        FAIL=$((FAIL + 1))
        FAILURES="$FAILURES\n  SECURITY: Password found in admin response"
    fi
else
    echo "  PASS: No password in admin response"
    PASS=$((PASS + 1))
fi
echo ""

############################################
# Edge Cases
############################################
echo ">>> 17. EDGE CASES"
echo "=================="
# Get non-existent device
test_endpoint "GET" "$BASE/api/v1/devices/99999" "404" "token" "" "Get non-existent device"
# Update non-existent device
test_endpoint "PUT" "$BASE/api/v1/devices/99999" "400" "token" '{"name":"xyz"}' "Update non-existent device"
# Get non-existent group
test_endpoint "GET" "$BASE/api/v1/groups/99999" "404" "token" "" "Get non-existent group"

############################################
# Summary
############################################
echo "=============================================="
echo "  COMPREHENSIVE API TEST SUMMARY"
echo "=============================================="
echo "  PASS: $PASS"
echo "  FAIL: $FAIL"
echo "  TOTAL: $((PASS + FAIL))"
if [ $((PASS + FAIL)) -gt 0 ]; then
    echo "  PASS RATE: $(echo "scale=1; $PASS * 100 / ($PASS + $FAIL)" | bc)%"
fi
echo ""

if [ "$FAIL" -gt 0 ]; then
    echo "=============================================="
    echo "  FAILURES:"
    echo "=============================================="
    echo -e "$FAILURES"
    echo ""
fi

echo "Test completed at: $(date)"
