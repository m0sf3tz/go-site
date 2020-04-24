go build example.go server_config.go constants.go packet_helper.go ipc_constants.go logger.go

if [ $? != 0 ]; then
  exit 
fi

if [ ! -z "$1" ]; then
  if [ "$1" == "-r" ]; then 
    ./example
  fi
fi
