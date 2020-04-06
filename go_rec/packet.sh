go build server.go tcp_core.go constants.go ipc_core.go packet_helper.go server_config.go

if [ $? != 0 ]; then
  exit 
fi

if [ ! -z "$1" ]; then
  if [ "$1" == "-r" ]; then 
    ./server
  fi
fi
