#!/bin/sh

IP="127.0.0.1"
PORT="8080"

curl -k -i -X POST -H "Content-Type: application/json" -d @wardrobe.json http://$IP:$PORT/allitems


#curl http://localhost:8080/allitems
#curl -k -i -X DELETE http://localhost:8080/allitems
