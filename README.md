# kahva

Web client for rTorrent.

The backend is used to deserialize XML-RPC requests from rTorrent to more browser friendly JSON presentation. It also serves the static Vue frontend so a dedicated HTTP server is not required.

---

## Configuration

### Environment variables

- `SERVER_ADDRESS` HTTP server bind address, it is `0.0.0.0:8080` by default. 
- `XMLRPC_URL` Remote XML-RPC server URL, this will be the URL exposed by your nginx (or similar) web server (e.g. `https://yourdomain.tld/rpc`)
- `XMLRPC_USERNAME` Optional basic authentication username
- `XMLRPC_PASSWORD` Optional basic authentication password
- `CORS_ORIGIN` CORS origin if the frontend runs on a different path
- `CORS_AGE` CORS age if the frontend runs on a different path

### Remote server setup

#### rTorrent

rTorrent should be configured to expose the XML-RPC interface through a UNIX socket.This can be achieved by adding the following lines to your `.rtorrent.rc` file. A regular SCGI socket should _never_ be used unless you are for certain that it cannot be accessed from outside your network.

```
...
scgi_local = /path/to/some/directory/xmlrpc.socket
schedule = scgi_permission,0,0,"execute.nothrow=chmod,\"g+w,o=\",/path/to/some/directory/xmlrpc.socket"
...
```

**Make sure that the socket is writable by the user that runs nginx (e.g. www-data)**


If you have an excessive amount of metadata you might need to increase the default XML-RPC limit. This is almost never required.

```
network.xmlrpc.size_limit.set = 10M
```

#### nginx

Create a nginx virtual host that serves the XML-RPC socket. Basic authentication is optional but recommended (read: a must) if you are accessing the server remotely. Note that the XML-RPC interface can be used to execute shell commands remotely.

You can use the snippet below as an example.

```
server {
        root /usr/share/nginx/html;
        index index.html;
        server_name yourdomain.tld;

        server_tokens off;
        autoindex off;
        auth_basic_user_file /path/to/some/directory/.htpasswd;
        auth_basic "super secret";

        location /rpc {
                include    scgi_params;
                scgi_pass  unix:/path/to/some/directory/xmlrpc.socket;
        }
}
```

A configuration like this will serve the socket at `http://yourdomain.tld/rpc`.

### Building

Easiest way to run the application is with Docker. The server will bind to `0.0.0.0` and use port `8080` by default. 

Build and run the image using `docker compose up --build`. Remember to change the environment variables to match your server configuration.

Alternatively the frontend and backend can be compiled separately.

#### Backend

Install the Go compiler and run `go build -v -o ./kahva ./cmd`. This should result in a single `kahva` binary.

#### Frontend

Install Node and the required dependencies. Run `npm run build`. This should result in a `dist/` directory.

Finally move the dist directory on the same level as the backend executable and rename the folder to `www/`.

The directory structure should look rougly like this:

```
./kahva
./www/
     /index.html
     /assets/
```

Run the binary with `SERVER_ADDRESS=0.0.0.0:8080 OTHER_ENV_VARAIBLES=... ./kahva`

## Problems?

Check the backend stdout for clues and always verify that your environment variables are correct and point to the right place.

If the server is reporting deserialization issues check the nginx error log and enable XML-RPC logging  in `.rtorrent.rc` by adding the line `log.xmlrpc = "/path/to/somewhere/xmlrpc.log"`.

If you are still experiencing issues and you are absolutely sure that the problem is not with your server and/or configuration, create an issue in this repository.