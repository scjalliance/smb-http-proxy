# scjalliance/smb-http-proxy

```
docker run -n smb_something  --restart=always-p 0.0.0.0:5432:80 --security-opt apparmor:unconfined --cap-add=SYSDMIN --cap-add=DAC_READ_SEARCH --rm -it -e "URLPREFIX=/some/url/prefix/" -e "DOMAIN=yourdomain" -e "USERNAME=somebody" -e "PASSWORD=secreteating" -e "UNCPATH=//flange.example.com/best-laid-plans/" scjalliance/smb-http-proxy
```
