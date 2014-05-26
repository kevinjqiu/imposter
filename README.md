```
    ____                           __           
   /  _/___ ___  ____  ____  _____/ /____  _____
   / // __ `__ \/ __ \/ __ \/ ___/ __/ _ \/ ___/
 _/ // / / / / / /_/ / /_/ (__  ) /_/  __/ /    
/___/_/ /_/ /_/ .___/\____/____/\__/\___/_/     
             /_/                                
```

A server that pretends to be another service.


Start the Server
================

Start the server:

```bash
$ go run imposter.go
[martini] listening on :3000
```

Create a Preset
===============

```bash
curl -XPOST -H"Content-Type:application/json" http://localhost:3000/p -d'{
    "matcher": {
       "method": "POST",
       "endpoint": "/foo",
       "headers": {
         "Accept": "application/json"
       },
       "body": "{}"
     },
     "response": {
       "status_code": 400,
       "headers": {
         "Content-Type": "application/json"
       },
       "body": "{\"message\": \"nope\"}"
      }
    }'
```

This creates a preset for the endpoint `/m/foo`. When a request to `localhost:3000/m/foo` with the header `Accept: application/json` and request body `{}`, the response given will be `400 Bad Request` with the header `Content-Type: application/json` and the body `{"message": "nope"}`.

