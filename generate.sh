export FABRIC_CFG_PATH=${PWD}
./bin/cryptogen generate --config=./crypto-config.yaml
./bin/configtxgen -profile OneOrgsOrdererGenesis -channelID byfn-sys-channel -outputBlock ./channel-artifacts/genesis.block
./bin/configtxgen -profile OneOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID mychannel
./bin/configtxgen -profile OneOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP

