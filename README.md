# info
docker pull defreitas/dns-proxy-server

docker run --rm --hostname dns.mageddo -v /var/run/docker.sock:/var/run/docker.sock -v /etc/resolv.conf:/etc/resolv.conf defreitas/dns-proxy-server


POST advert

curl --location 'http://ad-serving:8080/api/v1/advert' \
--header 'Content-Type: application/json' \
--data '{
    "id":"some_id",
    "title":"some_title",
    "link":"some_link"
}'

or 

curl --location 'http://ad-serving:8080/api/v1/advert' \
--header 'Content-Type: application/json' \
--data '{
    "id":"some_id",
    "title":"some_title",
    "description": "some_description",
    "link":"some_link"
}'


GET advert

curl --location 'http://ad-serving:8080/api/v1/advert/some_id'


GET url with impression number 

curl --location 'http://ad-serving:8080/api/v1/advert/some_id/impression'

