# LUKS Network Unlock Setup Guide

This guide explains how to configure automatic LUKS disk encryption unlocking using the Atlas Encryption API when the machine is on a local network.

## Overview

The `luks-network-unlock.sh` script runs during boot in the initramfs environment, fetches the LUKS decryption key from your encryption API over the network, and automatically unlocks your encrypted disk.

## Prerequisites

1. A working Atlas Encryption API server on your local network
2. A LUKS-encrypted disk/partition
3. Root access to configure the system
4. Network connectivity during boot (DHCP or static IP)

## Installation Steps

### 1. Configure the Script

Edit `luks-network-unlock.sh` and update these variables:

```bash
API_HOST="192.168.1.100"      # IP address of your encryption API server
API_PORT="8001"               # Port of your encryption API
ENDPOINT="atlas"              # Endpoint to use (atlas, sentryone, or sentrytwo)
ENCRYPTED_KEY="your_encrypted_luks_key_here"  # Your encrypted LUKS key
```

**Important:** The `ENCRYPTED_KEY` should be your actual LUKS key, encrypted using the encryption API's `/crypto/encrypt` endpoint.

### 2. Encrypt Your LUKS Key

First, get your current LUKS key or create a new one:

```bash
# If you don't have a key file, you can create one from your passphrase
# (Replace /dev/sdX1 with your encrypted partition)
echo "your-current-passphrase" | cryptsetup luksAddKey /dev/sdX1 /root/luks-key.txt

# Or generate a random key
dd if=/dev/urandom of=/root/luks-key.txt bs=512 count=1
cryptsetup luksAddKey /dev/sdX1 /root/luks-key.txt
```

Then encrypt it using your API:

```bash
# Read the key and encrypt it
LUKS_KEY=$(cat /root/luks-key.txt)

# Encrypt using the API
curl -X POST http://YOUR_API:8001/crypto/encrypt \
  -H "Content-Type: application/json" \
  -d "{\"data\":\"$LUKS_KEY\"}" | jq -r '.encrypted'

# Copy the encrypted value into the script
```

**Security Note:** Keep `/root/luks-key.txt` secure and delete it after setup if you only want network-based unlock.

### 3. Copy Script to System

```bash
sudo cp luks-network-unlock.sh /usr/local/sbin/
sudo chmod +x /usr/local/sbin/luks-network-unlock.sh
```

### 4. Add Script to initramfs Hooks

Create a hook to include the script in initramfs:

```bash
sudo mkdir -p /etc/initramfs-tools/hooks
sudo tee /etc/initramfs-tools/hooks/luks-network-unlock > /dev/null << 'EOF'
#!/bin/sh
PREREQ=""
prereqs()
{
    echo "$PREREQ"
}

case $1 in
prereqs)
    prereqs
    exit 0
    ;;
esac

. /usr/share/initramfs-tools/hook-functions

# Copy the unlock script
copy_exec /usr/local/sbin/luks-network-unlock.sh /bin/luks-network-unlock

# Ensure required tools are available
copy_exec /usr/bin/wget /bin/wget
copy_exec /bin/grep /bin/grep
copy_exec /bin/sed /bin/sed

exit 0
EOF

sudo chmod +x /etc/initramfs-tools/hooks/luks-network-unlock
```

### 5. Configure Network in initramfs

Enable network during boot:

```bash
# For DHCP (most common)
echo "IP=dhcp" | sudo tee -a /etc/initramfs-tools/initramfs.conf

# Or for static IP:
# echo "IP=192.168.1.50::192.168.1.1:255.255.255.0:hostname:eth0:off" | sudo tee -a /etc/initramfs-tools/initramfs.conf
```

Add network tools to initramfs:

```bash
sudo tee -a /etc/initramfs-tools/initramfs.conf > /dev/null << 'EOF'
DEVICE=eth0
EOF
```

### 6. Update crypttab

Find your encrypted device UUID:

```bash
sudo blkid | grep crypto_LUKS
```

Edit `/etc/crypttab` and add the keyscript:

```bash
sudo nano /etc/crypttab
```

Change the line from:
```
cryptroot UUID=your-uuid-here none luks,discard
```

To:
```
cryptroot UUID=your-uuid-here none luks,discard,keyscript=/bin/luks-network-unlock
```

### 7. Update initramfs

```bash
sudo update-initramfs -u -k all
```

### 8. Test (Optional but Recommended)

Before rebooting, you can test if the script works in your current environment:

```bash
sudo /usr/local/sbin/luks-network-unlock.sh
# Should output your decrypted LUKS key
```

### 9. Reboot and Test

```bash
sudo reboot
```

The system should:
1. Boot into initramfs
2. Configure network via DHCP
3. Connect to your encryption API
4. Fetch and decrypt the LUKS key
5. Automatically unlock the disk
6. Continue booting normally

## Troubleshooting

### Check Logs

After boot, check the logs:

```bash
journalctl -b | grep LUKS-UNLOCK
```

### Manual Unlock Fallback

If network unlock fails, the system will fall back to manual password entry. The original passphrase still works.

### Common Issues

1. **Network not available during boot**
   - Ensure DHCP is configured in initramfs
   - Check if your network interface is recognized
   - Verify network cable is connected

2. **API not reachable**
   - Ensure the API server is running before the client boots
   - Check firewall rules
   - Verify the API_HOST IP address is correct

3. **Wrong encrypted key**
   - Re-encrypt your LUKS key using the correct endpoint
   - Verify the ENDPOINT variable matches what you used to encrypt

4. **Script not found in initramfs**
   - Verify the hook was created correctly
   - Run `sudo update-initramfs -u` again
   - Check `lsinitramfs /boot/initrd.img-$(uname -r) | grep luks-network-unlock`

### Debug Mode

To see detailed logs during boot:

1. Edit `/etc/default/grub`
2. Remove `quiet splash` from `GRUB_CMDLINE_LINUX_DEFAULT`
3. Run `sudo update-grub`
4. Reboot

## Security Considerations

1. **Network Security**: The API communication is over HTTP. Consider:
   - Using VPN/WireGuard for the network
   - Adding TLS support to the API
   - Using network segmentation/isolation

2. **Key Storage**: The encrypted key is stored in plaintext in the script
   - The script is in initramfs, which is readable by root
   - Consider additional encryption of the initramfs itself

3. **Fallback**: Always keep a backup method to unlock:
   - The original passphrase still works
   - Keep a recovery USB with the key file

4. **API Security**:
   - Ensure the encryption API is not exposed to the internet
   - Use firewall rules to restrict access to local network only
   - Monitor API access logs

## Removal

To remove network unlock and return to manual unlock:

1. Edit `/etc/crypttab` and remove `,keyscript=/bin/luks-network-unlock`
2. Run `sudo update-initramfs -u`
3. Reboot

The system will prompt for manual password entry again.
