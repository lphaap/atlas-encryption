#!/bin/sh

### --- LUKS NETWORK UNLOCK SCRIPT --- ###
## Fetches LUKS decryption key from Atlas Encryption API during boot
## This script is designed to run in initramfs environment
## Usage: Called by cryptsetup via /etc/crypttab keyscript

set -e

### --- CONFIGURATION --- ###
# These should be customized for your environment

API_HOST="192.168.1.100"  # IP address of the encryption API server
API_PORT="8001"
ENDPOINT="atlas"  # or sentryone, sentrytwo
ENCRYPTED_KEY="your_encrypted_luks_key_here"  # The encrypted LUKS key

API_URL="http://${API_HOST}:${API_PORT}"
MAX_RETRIES=10
RETRY_DELAY=3
CONNECT_TIMEOUT=5

### --- LOGGING --- ###

log() {
    echo "[LUKS-UNLOCK] $*" >&2
}

### --- NETWORK SETUP --- ###

wait_for_network() {
    log "Waiting for network..."

    # Configure network interfaces if not already done
    if ! ip addr show | grep -q "inet "; then
        log "Configuring network via DHCP..."

        # Try to bring up all network interfaces
        for iface in /sys/class/net/*; do
            iface_name=$(basename "$iface")

            # Skip loopback
            [ "$iface_name" = "lo" ] && continue

            log "Trying interface: $iface_name"
            ip link set "$iface_name" up || true
        done

        # Wait a moment for interfaces to come up
        sleep 2

        # Try DHCP on available interfaces
        for iface in /sys/class/net/*; do
            iface_name=$(basename "$iface")

            # Skip loopback
            [ "$iface_name" = "lo" ] && continue

            # Check if interface is up and has carrier
            if [ -f "/sys/class/net/$iface_name/operstate" ]; then
                state=$(cat "/sys/class/net/$iface_name/operstate")
                if [ "$state" = "up" ] || [ "$state" = "unknown" ]; then
                    log "Getting DHCP on $iface_name..."
                    ipconfig -t 10 -d "$iface_name" || true

                    # Check if we got an IP
                    if ip addr show "$iface_name" | grep -q "inet "; then
                        log "Network configured on $iface_name"
                        return 0
                    fi
                fi
            fi
        done
    fi

    # Verify we have network
    if ip addr show | grep -q "inet "; then
        log "Network is ready"
        return 0
    fi

    log "Failed to configure network"
    return 1
}

### --- API COMMUNICATION --- ###

check_api_status() {
    local status_url="$API_URL/status"
    local attempt=1

    log "Checking API availability at $API_URL..."

    while [ $attempt -le $MAX_RETRIES ]; do
        log "API check attempt $attempt/$MAX_RETRIES"

        if wget -q -O /dev/null -T "$CONNECT_TIMEOUT" "$status_url" 2>/dev/null; then
            log "API is ready"
            return 0
        fi

        if [ $attempt -lt $MAX_RETRIES ]; then
            log "API not ready, waiting ${RETRY_DELAY}s..."
            sleep "$RETRY_DELAY"
        fi

        attempt=$((attempt + 1))
    done

    log "ERROR: API is not responding after $MAX_RETRIES attempts"
    return 1
}

fetch_decryption_key() {
    local endpoint_url="$API_URL/$ENDPOINT/key"
    local temp_file="/tmp/luks-response.json"

    log "Fetching decryption key from $endpoint_url..."

    # Create JSON payload
    cat > /tmp/luks-payload.json <<EOF
{"encrypted":"$ENCRYPTED_KEY"}
EOF

    # Make API request using wget (available in initramfs)
    if ! wget -q -O "$temp_file" \
        --header="Content-Type: application/json" \
        --post-file=/tmp/luks-payload.json \
        -T "$CONNECT_TIMEOUT" \
        "$endpoint_url" 2>/dev/null; then
        log "ERROR: Failed to fetch key from API"
        rm -f /tmp/luks-payload.json "$temp_file"
        return 1
    fi

    # Extract decrypted value (using basic shell tools available in initramfs)
    decrypted=$(grep -o '"decrypted":"[^"]*"' "$temp_file" | sed 's/"decrypted":"\(.*\)"/\1/')

    # Clean up
    rm -f /tmp/luks-payload.json "$temp_file"

    if [ -z "$decrypted" ]; then
        log "ERROR: Failed to extract decrypted key from response"
        return 1
    fi

    log "Successfully retrieved decryption key"

    # Output the key (this goes to cryptsetup)
    echo -n "$decrypted"
    return 0
}

### --- MAIN --- ###

log "Starting LUKS network unlock process..."

# Wait for network to be available
if ! wait_for_network; then
    log "ERROR: Network unavailable - falling back to manual unlock"
    exit 1
fi

# Check if API is reachable
if ! check_api_status; then
    log "ERROR: Cannot reach API - falling back to manual unlock"
    exit 1
fi

# Fetch and output the decryption key
if ! fetch_decryption_key; then
    log "ERROR: Failed to fetch decryption key - falling back to manual unlock"
    exit 1
fi

log "LUKS network unlock successful"
exit 0
