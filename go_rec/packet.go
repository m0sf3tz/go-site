package main

func create_packet(b []byte) []byte {
	b[0] = 1
	b[1] = 1
	b[2] = 125
	b[3] = 1
	b[4] = 15

	return b
}
