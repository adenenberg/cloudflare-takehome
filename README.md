# Cloudflare Takehome: URL Shortener

## Design
I chose to create a HTTP service with Golang and to use MongoDB as storage. I chose Mongo because I felt like a NoSQL db was better suited to this problem. It would allow the service to scale easily and the data doesn't have any kind of relationship that needs to be managed. I chose Mongo specifically because I felt like the document storage paradigm would work nicely for storing the URL data and statistics.

I implemented...

I have used Docker so that each of the pieces could be containerized and run independently. In the future, if we choose to manage the service with Kubernetes, it will allow us to easily scale the containers independently.

## Run
**Docker is required to run**

From the root of this project, the following applies:

In order to run all the containers with docker, please run
```docker-compose up```

To stop all running containers
```docker-compose down```

## Testing
Automated tests can be run with ```go test ./...```

You can use curl to manually call different endpoints in the service
| Endpoint | Command |
| --- | --- |
| Ping | ```curl localhost:8080/ping``` |
| Create Short URL | ```curl -X POST -d '{"original_url": "cloudflare.com", "expiration_date": "2023-10-30T15:04:05Z"}' localhost:8080/create``` |
| Delete Short URL | ```curl -X DELETE localhost:8080/delete/<short_url_key>``` |
| Get URL Stats | ```curl localhost:8080/stats/<short_url_key>``` |

You can visit a short link in the browser by following: ```localhost:8080/go/<short_url_key>```