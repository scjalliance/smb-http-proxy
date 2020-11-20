# scjalliance/smb-http-proxy

```
Usage: smb-http-proxy --url-prefix=STRING --source=STRING --user=STRING --password=STRING --domain=STRING --conn-timeout=5s

Serves files from an SMB share over HTTP.

Flags:
  -h, --help                 Show context-sensitive help.
      --url-prefix=STRING    ($URL_PREFIX)
      --source=STRING        ($SOURCE)
      --user=STRING          ($USER)
      --password=STRING      ($PASSWORD)
      --domain=STRING        ($DOMAIN)
      --conn-timeout=5s      ($CONN_TIMEOUT)
```

```
docker run --name smb_something --restart=unless-stopped -p 0.0.0.0:5432:80 --rm -it -e "URL_PREFIX=/some/url/prefix/" -e "SOURCE=//flange.example.com/best-laid-plans/" -e "USER=somebody" -e "PASSWORD=secreteating" -e "DOMAIN=yourdomain" scjalliance/smb-http-proxy
```
