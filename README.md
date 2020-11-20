# scjalliance/smb-http-proxy

```
docker run --name smb_something --restart=unless-stopped -p 0.0.0.0:5432:80 --rm -it -e "URLPREFIX=/some/url/prefix/" -e "DOMAIN=yourdomain" -e "USERNAME=somebody" -e "PASSWORD=secreteating" -e "SOURCE=//flange.example.com/best-laid-plans/" scjalliance/smb-http-proxy
```
