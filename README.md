# Dtail

Dtail is a tool to monitor changes in a database table (inserts, updates, deletes) and broadcasts them via websockets.

Dtail comes with two parts.
1. console setup utility
2. table monitor and websocket server

### How to setup Dtail
1. Run `go run cmd/setup/run.go`
2. Update the connection string and click `Connect DB`
3. Choose which table you want to monitor and click `Choose Table`
4. Repeat step 4 for as many tables you want to monitor
5. Click `Save` button and it will write a `.env` file to the root of the project
6. Verify the DB connection string is correct in `.env` file

### How to run Dtail server
1. Run `go run cmd/notify/server.go
2. Insert a row into one of the tables you are monitoring
3. You should see see the message printed in the terminal `received notification:  {"id": 22928, "op": "INSERT", "date": "2025-03-08", "table": "freddie_30y", "value": 1.98, "created_at": "2025-02-13T19:53:51.023876"}`
4. You can also use the demo websocket client by pointing your browser to [localhost](http://localhost:9999/demo)
5. You should now also see the message printed on the webpage.


### How to connect to Dtail server
You can connect to Dtail via websockets.  Point your websocket client to the Dtail servre and wait for events to stream in.  If you want to get basics, look at the html of the demo page under `internal/server/ws_client.html`
