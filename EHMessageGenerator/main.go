package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go"
	uuid "github.com/nu7hatch/gouuid"
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

var citycode = []string{"DFW", "ORD", "MIA", "LAX", "PHI", "AUS", "SFO", "DBX", "MAA", "KOC", "LHR", "XYZ"}

//https://notes.shichao.io/gopl/ch9/

var (
	mu          sync.Mutex // guards memberCount
	memberCount int
)

func main() {

	mu.Lock()
	memberCount = 616657
	mu.Unlock()
	for _, cityCODE := range citycode {
		go GenMember(cityCODE)
	}

	//Wait for a signal to quit
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

}

//GenMember - get the city code for filtering
func GenMember(city string) {
	var member Member

	connStr := os.Getenv("MSG_EVENTHUBNS")
	entityPath := os.Getenv("MSG_EVENTHUB")
	connStr = fmt.Sprintf("%s;EntityPath=%s", connStr, entityPath)
	hub, err := eventhub.NewHubFromConnectionString(connStr)
	if err != nil {
		fmt.Printf("ERR:Error connecting to .\n")
	}

	var bytesRepresentation []byte

	for x := 0; x < 10000000000; x++ {

		//var printStr string
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		mu.Lock()
		memberCount++
		member.Counter = memberCount
		mu.Unlock()

		uuid, err := uuid.NewV4()
		if err != nil {
			log.Fatal(err)
		}
		member.DtTm = time.Now().String()
		member.ID = uuid.String()
		member.Status = "New"
		member.City = city

		bytesRepresentation, _ = json.Marshal(member)
		event := eventhub.NewEvent(bytesRepresentation)
		event.PartitionKey = &member.City

		start := time.Now()
		err = hub.Send(ctx, event)
		if err != nil {
			fmt.Println(err)
			return
		}
		/*t := time.Now()
		elapsed := t.Sub(start)
		if err != nil {
			fmt.Println(err)
			return
		}*/

		rand.Seed(time.Now().UnixNano())
		//k := rand.Intn(10)
		k := 2 //Second Wait between Member Generation

		dispInfo := fmt.Sprintf("%s-%d-%v", city, member.Counter, start.UTC())
		//fmt.Printf("%s\t%s\tElapsed Time:%v\tNext Member in:%d seconds\n", dispInfo, printStr, elapsed, k)

		fmt.Printf("SENT:%s\n", dispInfo)
		time.Sleep(time.Duration(k) * time.Second)
	}
}
