#!/bin/bash -x

if [ -z "$USERNAME" -o -z "$PASSWORD" -o -z "$UNCPATH" ]; then
  echo "This requires that you provide USERNAME, PASSWORD, and UNCPATH environment variables."
  echo "You might also consider using the DOMAIN environment variable."
  echo "Also maybe you'd find URLPREFIX useful, which is used for HTTP URL prefix."
  exit 1
fi

if [ ! -z "$DOMAIN" ]; then
  DOMAIN="domain=$DOMAIN,"
fi

if [ ! -z "$USERNAME" ]; then
  USERNAME="user=$USERNAME,"
fi

if [ ! -z "$PASSWORD" ]; then
  PASSWORD="password=$PASSWORD,"
fi

#mount -t cifs -o ro,${USERNAME}${PASSWORD}${DOMAIN}uid=0,gid=0,forceuid,forcegid,vers=2.1,sec=ntlm "$UNCPATH" /mnt/smb || exit 1
mount -t cifs -o ro,${USERNAME}${PASSWORD}${DOMAIN}uid=0,gid=0,forceuid,forcegid,vers=2.1 "$UNCPATH" /mnt/smb || exit 1

/smb-http-proxy

umount /mnt/smb
