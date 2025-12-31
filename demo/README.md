# HTML/JS Web Application Demo

Demonstration of how you can run gonrg in a web application.
You must have a running gonrg server which is reachable from
your client.


## Prepare

First, edit `meter.html` and set the `server-select` options
values to your gonrg server(s), e.g. `http://<your-gonrg-server>:8080`.

You cannot run `meter.html` from disk, because of CORS rules,
you must host it somewhere first. For example, run a simple
local web server e.g. with

```bash
# cd to this dir here, e.g. cd demo first
python3 -m http.server 8081
```

Now, if not already done, you must allow your server hosting the
web page (`http://localhost:8081` in this example) to be a valid
origin in your gonrg application. For this, add the server to the
`alloworigins` list in the gonrg server conf and restart it.


## Run application

Open the html file `http://localhost:8081/meter.html` in your browser
and enjoy.


## Reverse proxy

You can run gonrg behind a reverse proxy to enable authentication,
enforce authorization, and provide encryption (TLS).

Example Apache2 reverse proxy configuration:
```
<IfModule mod_ssl.c>
    <VirtualHost _default_:443>
        ServerAdmin webmaster@localhost

        DocumentRoot /var/www/html

        [ ... ]

        ProxyPass "/meterapi/ws" "ws://internalserver:8123/ws"
        ProxyPassReverse "/meterapi/ws" "ws://internalserver:8123/ws"
        ProxyPass "/meterapi" "http://internalserver:8123"
        ProxyPassReverse "/meterapi" "http://internalserver:8123"

	</VirtualHost>
</IfModule>
```
