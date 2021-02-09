
# generate crypto material for orderer
  
./bin/cryptogen generate --config=./crypto-config.yaml
if [ "$?" -ne 0 ]; then
  echo "Failed to generate crypto material..."
  exit 1
fi

echo "=================== Replace _sk files  =========================================================="
find $PWD -type f -name *_sk -execdir mv {} key \;

echo "================ Check crypto-config Ffolder  ==================================================="

export FABRIC_CFG_PATH=$PWD

# generate genesis block 

echo "================== Generating genesis block =========================================="

./bin/configtxgen -profile SampleMultiNodeEtcdRaft -channelID byfn-sys-channel -outputBlock ./channel-artifacts/genesis.block
if [ "$?" -ne 0 ]; then
  echo "Failed to generate genesis block..."
  exit 1
fi



#generate channel tx

echo "================== Generating channel Tx =============================================="

export CHANNEL_NAME=mychannel  

./bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
if [ "$?" -ne 0 ]; then
  echo "Failed to generate channel tx..."
  exit 1
fi


echo "================== Generating anchor peer update for Org1  ============================="

export CHANNEL_NAME=mychannel 
./bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP

if [ "$?" -ne 0 ]; then
  echo "Failed to update sun anchor peer update ..."
  exit 1
fi


#update anchor peer of Org2

echo "================== Generating anchor peer update for Org2  ============================="

export CHANNEL_NAME=mychannel  
./bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP

if [ "$?" -ne 0 ]; then
  echo "Failed to generate Org2 anchor peer update..."
  exit 1
fi

