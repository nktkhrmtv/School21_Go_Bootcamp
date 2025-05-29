package main

import (
	"context"
	"flag"
	"log"
	"math"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"team00/message"
)

type anomalyMessage struct {
	SessionID string  `gorm:"column:session_id"`
	Frequency float64 `gorm:"column:frequency"`
	Timestamp int64   `gorm:"column:timestamp"`
}

func connectToDB() (*gorm.DB, error) {
	dsn := "host=localhost user=meteoriw dbname=godb password=11111 port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func predictedValues(frequencies []float64) (float64, float64) {
	sum := 0.0
	for _, freq := range frequencies {
		sum += freq
	}
	mean := sum / float64(len(frequencies))

	variance := 0.0
	for _, freq := range frequencies {
		variance += math.Pow(freq-mean, 2)
	}
	variance /= float64(len(frequencies) - 1)
	std := math.Sqrt(variance)

	return mean, std
}

func anomalyDetection(msg *message.Message, db *gorm.DB, k float64, meanVal float64, devSTD float64) {
	currentFrequency := msg.GetFrequency()
	if currentFrequency > meanVal+k*devSTD || currentFrequency < meanVal-k*devSTD {
		anomaly := anomalyMessage{
			SessionID: msg.SessionId,
			Frequency: msg.Frequency,
			Timestamp: msg.Timestamp,
		}
		if err := db.Create(&anomaly).Error; err != nil {
			log.Printf("Failed to save anomaly to database: %v", err)
		} else {
			log.Printf("Anomaly saved: SessionID=%s, Frequency=%.2f, Timestamp=%d", msg.SessionId, msg.Frequency, msg.Timestamp)
		}
	}
}

func main() {
	k := flag.Float64("k", 0., "STD anomaly coefficient")
	flag.Parse()

	conn, err := grpc.NewClient("localhost:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()
	client := message.NewMessageServiceClient(conn)
	stream, err := client.StreamMessages(context.Background(), &message.StreamMessagesRequest{})
	if err != nil {
		log.Fatalf("Failed to stream messages: %v", err)
	}
	freqIter, devSTD, meanVal := 0, 0., 0.
	var frequencies []float64

	db, err := connectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&anomalyMessage{})

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive message: %v", err)
		}

		if freqIter < 50 {
			frequencies = append(frequencies, msg.GetFrequency())
			freqIter++
			if freqIter%10 == 0 {
				meanVal, devSTD = predictedValues(frequencies)
				log.Printf("meanVal:%.3f devSTD:%.3f for %d values\n", meanVal, devSTD, freqIter)
			}
		} else {
			anomalyDetection(msg, db, *k, meanVal, devSTD)
		}
	}
}
