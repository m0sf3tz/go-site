package main

/**********************************************************
 *        Used to define type of incomming packet
 *********************************************************/
const DATA_PACKET = (0)         // Server -> Device
const CMD_PACKET = (1)          // Server -> Device
const INTERNAL_ACK_PACKET = (2) // Internal
const DEVICE_ACK_PACKET = (3)   // Device -> Server
const LOGIN_PACKET = (4)        // Device -> Server

const TYPE_SIZE = (1)
const TRANSACTION_ID_SIZE = (2)
const SMALL_PAYLOAD_SIZE = (8)
const MEDIUM_PAYLOAD_SIZE = (256)
const LARGE_PLAYLOAD_SIZE = (512)
const REASON_SIZE = (1)
const HOST_ACK_REQ_SIZE = (1)

const PACKET_TYPE_OFFSET = (0)
const PACKET_TRANSACTION_ID_OFFSET = (1)
const PACKET_HOST_ACK_REQ_OFFSET = (3)
const PACKET_ACK_REASON_OFFSET = (4)

const DATA_PACKET_SIZE = (TYPE_SIZE + TRANSACTION_ID_SIZE + HOST_ACK_REQ_SIZE + LARGE_PLAYLOAD_SIZE)
const CMD_PACKET_SIZE = (TYPE_SIZE + TRANSACTION_ID_SIZE + HOST_ACK_REQ_SIZE + SMALL_PAYLOAD_SIZE)
const ACK_PACKET_SIZE = (TYPE_SIZE + TRANSACTION_ID_SIZE + HOST_ACK_REQ_SIZE + REASON_SIZE) //Same format for device/server acks
const LOGIN_PACKET_SIZE = (TYPE_SIZE + TRANSACTION_ID_SIZE + HOST_ACK_REQ_SIZE + MEDIUM_PAYLOAD_SIZE)

const PACKET_LEN_MAX = DATA_PACKET_SIZE

const MAX_OUTSTANDING_TCP_CORE_SEND = (16)

const QUEUE_LEN_ONE = (1)
const QUEUE_MIN_SIZE = (4) //sizeof int32
const DONT_WAIT_QUEUE = (0)

/**********************************************************
 *                    ACK/NAK REASONS
 *********************************************************/
const ACK_GOOD = (0)
const NAK_TCP_DOWN = (1)
const SERVER_ACK_TIMED_OUT = (2)

/**********************************************************
 *               CONSUMER ACK REQUIREMENTS
 *********************************************************/
const CONSUMER_ACK_NOT_NEEDED = (0)
const CONSUMER_ACK_REQUIRED = (1)

const INTERNALLY_ACKED = (1)
const PENDING_ACK = (0)
