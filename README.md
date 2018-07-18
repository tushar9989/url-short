# URL Shortener

## Requirements

- Cassandra >= 3.0.9
- Go >= 1.10
- Nginx >= 1.15.0

## System Logic

Components:
- load balancer
	- authenticate all requests and add headers
		- auth server(s)
	- for get requests
		- proxy to read-server(s)
	- for post requests
		- proxy to write-server(s)

- read-server
	- Look up db for link
	- Check if user is authorized in case of retricted url
	- Log view event
	- return data

- write-server 
	- Startup
		- load segment start, end and current unused id from db
	- Set id for request
		- auto incr number
		- while already used
			- incr
			- if boundary + 1 move to offset
	- Store into db
	- Periodically write current unused id for segment to db

## Setup

### Cassandra

- use `nodetool status` to get datacenter name
- start `cqlsh`
- Create the keyspace: 
    - `CREATE KEYSPACE IF NOT EXISTS urlshort WITH REPLICATION = { 'class' : 'NetworkTopologyStrategy', 'datacenter1' : 1 } AND DURABLE_WRITES = true;`
    - On local system: `ALTER KEYSPACE urlshort WITH replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };`
- Use created keyspace `use urlshort;`
- Create the links table: `CREATE TABLE links (
  id varint,
  expire_time timestamp,
  target_url text,
  allowed_emails set<text>,
  user_id text,
  PRIMARY KEY (id, user_id)
);`
- Create the write_server_meta table: `CREATE TABLE write_server_meta (
  name text PRIMARY KEY,
  start varint,
  end varint,
  current varint
);`
- Create user stats table: `CREATE TABLE views (
    user_id text,
    id varint,
    views counter,
    PRIMARY KEY(user_id, id)
) WITH CLUSTERING ORDER BY (id DESC);`
- Write default data for write servers. Depends on the number of write servers you want to run. 2 here:
    - `insert into write_server_meta (name, start, end, current) values ('0', 1, 1760807303104, 1);`
    - `insert into write_server_meta (name, start, end, current) values ('1', 1760807303105, 3521614606207, 1760807303105);`

### Nginx

- Copy `nginx-proxy.conf` to the http section of your `nginx.conf` or copy the file to the `servers` subdirectory of your nginx setup.
- Reload http server: `nginx -s reload`

### Go

- Copy project source code to `$GOHOME/go/src/github.com/tushar9989`
- For auth, read and write servers copy the example `config.json` file and modify as needed.
- Install 3rd party dependencies
    - `go get github.com/tkanos/gonfig` Used to load config from json file.
    - `go get github.com/gocql/gocql` Used to interact with Cassandra.

## Start Servers

- If on a unix environment `./start.sh` from the project's root directory.
- Otherwise start auth, read and write servers after modifying config files before starting each instance.

## Usage

- To login with Google visit `localhost:8000`
- To create a Short Link
    - Optional `CustomSlug`(max 7 characters [0-9a-zA-Z]) and `ValidEmails` (will only work if logged in)
    - Method: POST
    - URL: localhost:8000/shorten
    - e.g: 
        - Request: `{
            "TargetUrl": "http://google.com",
            "ExpireAt": "2006-01-02T15:04:05Z",
            "CustomSlug": "tushar",
            "ValidEmails": [
                "tushar9989@gmail.com"
            ]
        }`
        - Response: `{
            "data": {
                "slug": "tushar"
            },
            "status": "OK"
        }`

- To visit short link: `localhost:8000/s/tushar`.
- To check created url stats: `localhost:8000/stats`. Will only work if logged in.
