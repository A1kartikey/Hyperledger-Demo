docker-compose -f orderer-peer.yaml up -d 

docker-compose -f org1-peer.yaml -f docker-compose-couch.yaml up -d 

docker-compose -f org2-peer.yaml up -d 

docker-compose -f ca-org.yaml up -d

docker-compose -f cli.yaml up -d 

export FABRIC_START_TIMEOUT=5

echo "sleeping for ${FABRIC_START_TIMEOUT} secs ... "

sleep ${FABRIC_START_TIMEOUT}

docker ps -a