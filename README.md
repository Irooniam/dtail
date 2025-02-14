# Dtail

Dtail is a tool to monitor changes in a database table (inserts, updates, deletes) and broadcasts them via websockets.

Dtail comes with two parts.
1. console setup utility
2. table monitor and websocket server

### How to setup utility
1. Run `go run cmd/setup/run.go`
2. Update the connection string and click `Connect DB`
3. Choose which table you want to monitor and click `Choose Table`
4. Repeat step 4 for as many tables you want to monitor
5. Click `Save` button and it will write a `.env` file to the root of the project
6. Verify the DB connection string is correct in `.env` file

