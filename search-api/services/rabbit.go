package services

import (
	"encoding/json"
	"fmt"
	"log"
	"search-api/dtos"

	"github.com/streadway/amqp"
)

const (
	rabbitMQURL = "amqp://guest:guest@localhost:5672/" // Reemplázalo con la configuración correcta
	queueName   = "search_queue"
)

// ConsumeMessages escucha los mensajes de RabbitMQ y los indexa en Solr
func ConsumeMessages() {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatal("❌ Error al conectar con RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("❌ Error al abrir canal en RabbitMQ:", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatal("❌ Error al declarar la cola:", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("❌ Error al consumir mensajes de la cola:", err)
	}

	go func() {
		for msg := range msgs {
			fmt.Println("📩 Recibido mensaje:", string(msg.Body))

			// Convertir mensaje JSON a estructura HotelDTO
			var hotelDTO dtos.HotelDTO
			if err := json.Unmarshal(msg.Body, &hotelDTO); err != nil {
				log.Println("❌ Error al parsear JSON:", err)
				continue
			}

			// Indexar en Solr
			if err := AddHotel(hotelDTO); err != nil {
				log.Println("❌ Error al indexar hotel en Solr:", err)
				continue
			}

			fmt.Println("✅ Hotel indexado en Solr:", hotelDTO.Name)
		}
	}()

	select {} // Mantener el consumidor activo
}
