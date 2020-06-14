package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go"
)

//Member status
type Member struct {
	ID      string `json:"id,omitempty"`
	Counter int    `json:"counter,omitempty"`
	DtTm    string `json:"dttm,omitempty"`
	HHmm    string `json:"hhmm,imitempty"`
	Status  string `json:"status,omitempty"`
	City    string `json:"city,omitempty"`
}

//Partition helps with reading through the partitions
func Partition(conn string, partitionid string) {
	hub, err := eventhub.NewHubFromConnectionString(conn)
	if err != nil {
		LogIt("ERR:Could not connect to Event Hub", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Hour)
	defer cancel()
	_, err = hub.Receive(ctx, partitionid, handler, eventhub.ReceiveWithLatestOffset())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("INFO:I am listening to partition %s...\n", partitionid)
}

var entityPath string

func main() {
	connStr := os.Getenv("MSG_EVENTHUBNS")
	entityPath = os.Getenv("MSG_EVENTHUB")
	connStr = fmt.Sprintf("%s;EntityPath=%s", connStr, entityPath)

	hub, err := eventhub.NewHubFromConnectionString(connStr)
	if err != nil {
		LogIt("ERR:Could not connect to Event Hub", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	info, err := hub.GetRuntimeInformation(ctx)
	if err != nil {
		log.Fatalf("failed to get runtime info: %s\n", err)
	}
	log.Printf("got partition IDs: %s\n", info.PartitionIDs)

	for _, partitionID := range info.PartitionIDs {
		go Partition(connStr, partitionID)
	}

	//Wait for a signal to quit
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
	hub.Close(context.Background())
}

//LogIt the Logger
func LogIt(pipeline string, err error) {
	fmt.Printf("%s,%s\n", pipeline, err)
}

//Handle EventHub Messages
func handler(c context.Context, event *eventhub.Event) error {

	var member Member
	json.Unmarshal(event.Data, &member)
	if member.City != "" {
		keyStr := fmt.Sprintf("%s-%d", member.City, member.Counter)
		enqTime := fmt.Sprintf("%s", event.SystemProperties.EnqueuedTime)
		fmt.Printf("%s:%s\t%s\n", entityPath, keyStr, enqTime)
	}

	return nil
}
