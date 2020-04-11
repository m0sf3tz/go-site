go build testCore.go constants.go server_config.go packet_helper.go logger.go ipc_constants.go

if [ $? != 0 ]; then
  exit 
fi

if [ ! -z "$1" ]; then
  if [ "$1" == "-r" ]; then 
    ./testCore
  fi
fi
