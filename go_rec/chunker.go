package main

import "fmt"

const DATA_PACKET = 0
const CMD_PACKET = 1

const SMALL_PAYLOAD_SIZE = 8
const LARGE_PLAYLOAD_SIZE = 512

const DATA_PACKET_SIZE = (1 + LARGE_PLAYLOAD_SIZE) // extra one for "type"
const CMD_PACKET_SIZE = (1 + SMALL_PAYLOAD_SIZE)   // extra one for "type"
const MES_LEN_MAX = (DATA_PACKET_SIZE)

const MAX_Q_LEN = 10

var curr_buff int
var pckt_size int
var current_parse_type = -1
var chunk_buff []byte

var c chan []byte

func chunker(rx []byte, lenght int) int {
	if len(rx) == 0 || lenght == 0 {
		return -1
	}

	var rx_processed = 0
	for rx_processed != lenght {
		if current_parse_type == -1 {
			current_parse_type = int(rx[rx_processed]) //the type is always the first byte in any packet

			switch current_parse_type {
			case DATA_PACKET:
				pckt_size = DATA_PACKET_SIZE
			case CMD_PACKET:
				pckt_size = CMD_PACKET_SIZE
			default:
				return -1
			}
		}

		var rx_left = lenght - rx_processed
		read_len := 0

		if rx_left <= (pckt_size - curr_buff) {
			read_len = rx_left
		} else {
			read_len = pckt_size - curr_buff
		}

		chunk_buff = append(chunk_buff, rx[rx_processed:rx_processed+read_len]...)
		curr_buff += read_len
		rx_processed += read_len

		if curr_buff == pckt_size {
			//do something
			fmt.Println("Parsed a slice!")
			c <- chunk_buff
			chunk_buff = nil
			curr_buff = 0
			current_parse_type = -1
		}
	}
	return 1
}

/*
func main() {
	c = make(chan []byte, MAX_Q_LEN)

	b := make([]byte, 9)
	d := make([]byte, 9)

	b[0] = 1

	go func() {
		chunker(b[0:8], 8)
		chunker(b[8:9], 1)
		chunker(b, 9)
	}()

	i := 0
	for {
		d = <-c
		fmt.Println(d)
		i++
		if i == 2 {
			break
		}
	}
}
*/

//int chunker(const char * rx, const int len)
//{
//  if(rx == NULL || len == 0)
//  {
//    return -1;
//  }
//
//  // curr_buff: tracks the position inside chunk_buffer
//  // pckt_size: size of current packet being processed
//  // rx_processed : marks the position within the rx buffer;
//  static int curr_buff;
//  static int pckt_size;
//  static int current_parse_type = -1;
//  static char chunk_buff[DATA_PACKET_SIZE];
//  int rx_proccessed = 0;
//
//  while(rx_proccessed != len)
//  {
//    if (current_parse_type == -1)
//    {
//      current_parse_type = *(uint8_t*)(rx+rx_proccessed);
//
//      switch(current_parse_type)
//      {
//        case DATA_PACKET:
//          pckt_size = DATA_PACKET_SIZE;
//          break;
//        case CMD_PACKET:
//          pckt_size = CMD_PACKET_SIZE;
//          break;
//        default:
//          ESP_LOGE(TAG, "unknown command parse type recieved?");
//          return -1;
//      }
//    }
//
//    // Only read within the next message boundary based on the current
//    // Packet size being processed.
//    int rx_left  = len - rx_proccessed;
//    int read_len = (rx_left <= (pckt_size - curr_buff)) ? rx_left : (pckt_size - curr_buff);
//
//    memcpy(chunk_buff + curr_buff, rx + rx_proccessed, read_len);
//    curr_buff               += read_len;
//    rx_proccessed           += read_len;
//
//    if(curr_buff > pckt_size)
//    {
//      ESP_LOGE(TAG, "currentBuff > pckt_size - huge error");
//      return -1;
//    }
//
//    if (curr_buff == pckt_size)
//    {
//      ESP_LOGI(TAG, "Adding a packet to the RX_LL");
//      // Reset the state machine internal to this function
//      curr_buff = 0;
//      current_parse_type = -1;
//      ll_add_node(RX_LL, (void*)chunk_buff, pckt_size);
//      const char foo[] = "hi";
//      BaseType_t xStatus = xQueueSendToBack(xQueue_RX, (void * const) foo, 0 );
//    }
//  }
//  return 0;
//}
