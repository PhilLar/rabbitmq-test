package main

import (
	"fmt"
	"time"

	ws "github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

type Message struct {
	Type     string    `json:"type"`
	Channels []Channel `json:"channels"`
}

type Channel struct {
	Name       string   `json:"name"`
	ProductIDs []string `json:"product_ids"`
}

type Match struct {
	Type         string    `json:"type"`
	TradeID      int       `json:"trade_id,number"`
	Sequence     int64     `json:"sequence,number"`
	MakerOrderID string    `json:"maker_order_id"`
	TakerOrderID string    `json:"taker_order_id"`
	Time         time.Time `json:"time,string"`
	ProductID    string    `json:"product_id"`
	Size         string    `json:"size"`
	Price        string    `json:"price"`
	Side         string    `json:"side"`
}

func main() {

	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed.pro.coinbase.com", nil)
	if err != nil {
		println(err.Error())
	}

	subscribe := Message{
		Type: "subscribe",
		Channels: []Channel{
			Channel{
				Name:       "matches",
				ProductIDs: []string{"BTC-USD"},
			},
		},
	}

	if err := wsConn.WriteJSON(subscribe); err != nil {
		println(err.Error())
	}

	// rabbit
	conn, err := amqp.Dial("amqp://root:root@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		panic(err)
	}

	response := new(Match)
	for {
		if err := wsConn.ReadJSON(response); err != nil {
			println(err.Error())
			break
		}
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(fmt.Sprint(response)),
			})
	}
}
