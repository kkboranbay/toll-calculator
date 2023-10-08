package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kkboranbay/toll-calculator/types"
)

var SendInterval = time.Second * 5
var wsEndpint = "ws://127.0.0.1:30000/ws"

func getLatLong() (float64, float64) {
	return generateCoordinate(), generateCoordinate()
}

func generateCoordinate() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()

	return n + f
}

func generateOBUIds(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}

	return ids
}

func main() {
	obuIDS := generateOBUIds(20)

	// The github.com/gorilla/websocket package in Go provides a Dialer type
	// that can be used to connect to WebSocket servers.
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(wsEndpint, nil)
	if err != nil {
		log.Fatal(err)
	}

	for {
		for i := 0; i < len(obuIDS); i++ {
			lat, long := getLatLong()
			data := types.OBUData{
				OBUID: obuIDS[i],
				Lat:   lat,
				Long:  long,
			}
			fmt.Printf("%+v\n", data)
			if err := conn.WriteJSON(data); err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(SendInterval)
	}
}
