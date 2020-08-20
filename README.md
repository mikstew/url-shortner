# url-shortner

## Run
- go run main.go
- localhost:8080/shorten?url=<longUrl> to request a shortened URL
- localhost:8080/<token> to request the long URL


## Explanation
 - Create and serve shorten requests at the endpoint dol.ly/shorten?url=long-url.com
 - Create and serve expand requests at the endpoint dol.ly/token
 - Shorten: Parse the user submitted input parameter for a URL. Write that URL to the DB, which has an auto-incrementing ID as the primary key. The DB insert responds with the ID. The ID returned from the DB is encoded as a base62 token which is used as the unique identifier for the shortened URL.
 - Expand: Parse the user submitted path variable for a token. Decode the token from base62 to base10, which will give the id (index in slice) where the long URL is in the DB. Read the long URL from the DB and return to the user.

## How would you deploy this service?
- A service like this would be easy to containerize and deploy to a k8s cluster behind a horizontal pod autoscaler to scale up/down with load. I would use a NoSQL DB such as MongoDB or Cassandra since the nature of the use will be submitting a URL to shorten and then a potentially high number of reads of that value as the shortened URL is used. 

## What problems might you anticipate at 1million users vs 10million?
- With the current implementation, there is no cleanup of old URLs. This means the DB will continue to grow and the base62 encoded values will get longer. It would eventually get to the point that the URLs aren't shortened much, if at all. A duration could be added to the request, which will add a column to the DB entry for a cleanup script to purge when the duration expires. At that point, you could insert to the lowest available id.
- As use of the service increases the DB could start to have trouble keeping up with requests. Adding a caching layer in to quickly serve requests should make load on the DB more manageable. 

## How would you go about measuring the performance of this system?
- Average Response Time: This will be used to help identify a degradation in the service and general satisfaction of the users. If response times start to creep up, an investigation could be started to find out why and make adjustments before a problem occurs.
- Aggregate response codes: This will be used to identify issues with the service. If you know customers are having a bad user experience an investigation into the cause can be started, ideally before too many customers are impacted.
- Shorten/Expand request counts: These values will be used to identify adoption of the service and to make decisions on scaling.
- Pod count: This will be used to track the number of pods necessary over the course of the day to help identify use trends. This will provide good information for cost analysis as well as identifying good times to perform work that may impact service.
- DB CPU/Memory usage: Used to monitor DB health to identify problems. It could also be used to idenify if the DB needs to be scaled up/down.