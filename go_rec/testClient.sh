go build testClient.go constants.go server_config.go packet_helper.go

if [ $? != 0 ]; then
  exit 
fi

if [ ! -z "$1" ]; then
  if [ "$1" == "-r" ]; then 
    ./testClient
  fi
fi
