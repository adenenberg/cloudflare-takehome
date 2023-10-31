# Cloudflare Takehome: URL Shortener

## Design
I chose to create a HTTP service with Golang and to use MongoDB as storage. I chose Mongo because I felt like a NoSQL db was better suited to this problem. It would allow the service to scale easily and the data doesn't have any kind of relationship that needs to be managed. I chose Mongo specifically because I felt like the document storage paradigm would work nicely for storing the URL data and statistics.

I have used Docker so that each of the pieces could be containerized and run independently. In the future, if we choose to manage the service with Kubernetes, it will allow us to easily scale the containers independently.

## Run
**Docker is required to run**

From the root of this project, the following applies:

In order to run all the containers with docker, please run
```docker-compose up```

To stop all running containers
```docker-compose down```

## Testing
With the containers up and running, please change line 17 in ```api_test.go``` to ```mongodb://root:cloudflare@localhost:27017/``` and run ```go test ./...``` for automated tests.

You can use curl to manually call different endpoints in the service
| Endpoint | Command |
| --- | --- |
| Ping | ```curl localhost:8080/ping``` |
| Create Short URL | ```curl -X POST -d '{"original_url": "cloudflare.com", "expiration_date": "2023-10-30T15:04:05Z"}' localhost:8080/create``` |
| Delete Short URL | ```curl -X DELETE localhost:8080/delete/<short_url_key>``` |
| Get URL Stats | ```curl localhost:8080/stats/<short_url_key>``` |

You can visit a short link in the browser by following: ```localhost:8080/go/<short_url_key>```

## Viewing data
You can view the data stored in MongoDB by visiting ```localhost:8081``` in your browser and accessing the ```cloudflare``` DB.

## Future

The current implementation of this service would be fine for small scale usage. When this service needs to scale to support millions of users, there are some changes that would need to be made.

1. Add a new service for key generation. Currently, the API handles generating the short url key when a new one is created. This would be a separate service that would exist for the sole purpose of generating keys and storing them in a database, then every time a new short url is created, the API would read from this database instead of generating a key itself. The key generation service would run independently of the API and periodically check to make sure a number of keys are available for the API to use and when the number is below a certain threshold, it would add more keys. KGS (key generation service) will be also be responsible for maintaining the set of used keys. This is to ensure duplicated keys are not used and to help speed up the APIs
2. Add a message queue (kafka) to record stats. Currently, whenever a short URL is hit, a new record is inserted into the stats table. If lots of URLs are being hit constantly, the number of inserts going on could lead to a bottleneck. In order to prevent this, I would add a durable queue so that every time a URL is hit, a message is sent to the queue indicating which URL and the timestamp. A worker would be responsible for draining the queue and recording the data to the DB. This would take some pressure off the API for handling the inserts into the stats DB, by giving the responsibility to a worker running simultaneously. The downside of this is that now that stats data becomes eventually consistent instead of strongly, but I believe this is something that can be lived with because the latency for the API would be greatly improved.
3. Add a cache (redis) for popular short URLs. Currently, when a short URL is requested, the API must reach out to the DB to look up the original URL. For a service with a small scale, this could be ok, but as we scale with more URLs and more users, having to look up a URL each time could lead to slow downs. We can take some load off the DB by introducing a redis cache. We can keep a set of the most commonly accessed URLs in the cache so that the API isn't always having to query the DB for URL data.
4. Add a load balancer between the client and API and API and DBs. This will help to ensure that traffic is evenly distributed between our servers to make sure that any one particular server does not get overloaded.
5. User Kubernetes to manage containers. We can use Kubernetes to monitor load on our containers and once a certain threshold is reached, automatically add more containers to help with the load. Subsequently, it can also tear down those containers when the load has returned to normal.