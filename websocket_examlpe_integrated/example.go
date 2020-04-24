package main

import (
	"bitbucket.org/avd/go-ipc/mq"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//var dmq_from_site_to_core, dmq_from_core_to_site *mq.LinuxMessageQueue
var dmq_from_site_to_core *mq.LinuxMessageQueue

func init_ipc() {

	var err error

	//dmq_from_site_to_core, err = mq.CreateLinuxMessageQueue("smq_from_site_to_core", os.O_RDWR, IPC_QUEUE_PERM, IPC_QUEUE_DEPTH, PACKET_LEN_MAX)
	dmq_from_site_to_core, err = mq.OpenLinuxMessageQueue("smq_from_site_to_core", os.O_RDWR)
	if err != nil {
		logger(PRINT_FATAL, "Could not create linux smq_client - (did yu increase the limit in /proc/sys/fs/mqueue/msg_max?  error: ", err)
	}
	/*
		dmq_from_core_to_site, err := mq.OpenLinuxMessageQueue("smq_from_core_to_site", os.O_RDWR)
		if err != nil {
			logger(PRINT_FATAL, "Could not create linux smq_client - (did yu increase the limit in /proc/sys/fs/mqueue/msg_max?  error: ", err)
		}
	*/
}

var ti int16

func hack_ipc() {
	fmt.Println(dmq_from_site_to_core)
	fmt.Println("sendin to core!")
	p := Packet{}
	p.Packet_type = CMD_PACKET
	p.Transaction_id = ti
	p.Consumer_ack_req = CONSUMER_ACK_NOT_NEEDED
	ti++

	c := Cmd_add_user{0, []byte("cool boy saman")}
	p.Data = packet_pack_cmd(c)

	packed := packet_pack(p)
	fmt.Println(packed)

	dmq_from_site_to_core.Send(packed)

}

var upgrader = websocket.Upgrader{CheckOrigin: CO} // use default options

func CO(r *http.Request) bool {
	fmt.Println(r.Header.Get("Origin"))
	return true
}

type Person struct {
	FirstName string
	LastName  string
}

func upgrade(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()

		hack_ipc()

		p := Person{}
		er := json.Unmarshal(message, &p)
		fmt.Println(er)
		fmt.Println("first: ", p.FirstName, " Last: ", p.LastName)
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}

}

func home(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadFile("./static/button.html")
	w.Write(b)
}

func main() {
	init_ipc()
	http.HandleFunc("/upgrade", upgrade)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
