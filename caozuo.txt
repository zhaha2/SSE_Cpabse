在~/fabric-samples/chaincode-docker-devmode下打开三个终端

第一个：
	sudo docker-compose -f docker-compose-simple.yaml up

第二个：
	sudo docker exec -it chaincode bash
	
	配置PBC和我自己的cpabe：
	cd
	apt-get update
	apt-get install wget
	apt-get install tar
	sudo apt-get install libgmp-dev
	sudo apt-get install build-essential flex bison
	wget https://crypto.stanford.edu/pbc/files/pbc-0.5.14.tar.gz
	tar -xzvf pbc-0.5.14.tar.gz
	cd pbc-0.5.14
	./configure
        make
        sudo make install
	ldconfig
	cd /opt/gopath/src/github.com/
	mkdir Nik-U
	cd Nik-U
	git clone https://github.com/Nik-U/pbc.git
	cd /opt/gopath/src
	
	从本机往docker内传cpabse包
	docker cp cpabse a8b907129bfc（容器码）:/opt/gopath/src
	
	cd /opt/gopath/src/chaincode/my
	go build
	CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=my:0 ./my

第三个：
	sudo docker exec -it cli bash
	cd
	apt-get update
	apt-get install wget
	apt-get install tar
	sudo apt-get install libgmp-dev
	sudo apt-get install build-essential flex bison
	wget https://crypto.stanford.edu/pbc/files/pbc-0.5.14.tar.gz
	tar -xzvf pbc-0.5.14.tar.gz
	cd pbc-0.5.14
	./configure
        make
        sudo make install
	ldconfig
	cd /opt/gopath/src
	mkdir github.com
	cd github.com
	mkdir Nik-U
	cd Nik-U
	git clone https://github.com/Nik-U/pbc.git
	cd /opt/gopath/src
	
	从本机往docker内传cpabse包
	docker cp cpabse a8b907129bfc（容器码）:/opt/gopath/src
	
	peer chaincode install -p chaincodedev/chaincode/my -n my -v 0
	peer chaincode instantiate -n my -v 0 -c '{"Args":["a","b"]}' -C myc
