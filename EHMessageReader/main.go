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

//CycleEvent to Register
type CycleEvent struct {
	Event  string `json:"event"`
	Flight struct {
		Event      string `json:"event"`
		Version    string `json:"version"`
		TrackingID string `json:"trackingID"`
		DataTime   struct {
			SourceTimeStamp string `json:"sourceTimeStamp"`
			FltHubTimeStamp string `json:"fltHubTimeStamp"`
		} `json:"dataTime"`
		Key struct {
			AirlineCode struct {
				IATA string `json:"IATA"`
				ICAO string `json:"ICAO"`
			} `json:"airlineCode"`
			FltNum       string `json:"fltNum"`
			FltOrgDate   string `json:"fltOrgDate"`
			DepSta       string `json:"depSta"`
			DupDepStaNum string `json:"dupDepStaNum"`
		} `json:"key"`
		FosPartition string `json:"fosPartition"`
		Leg          struct {
			Stations struct {
				Arr          string `json:"arr"`
				DupArrStaNum string `json:"dupArrStaNum"`
			} `json:"stations"`
			Costs struct {
				TimesCancelled int `json:"timesCancelled"`
			} `json:"costs"`
			DepartureArrival struct {
				DepGate     string `json:"depGate"`
				DepTerminal string `json:"depTerminal"`
				ArrTerminal string `json:"arrTerminal"`
			} `json:"departureArrival"`
			Equipment struct {
				SkdEquipType             string `json:"skdEquipType"`
				AssignedEquipType        string `json:"assignedEquipType"`
				SkdNumericEquipType      string `json:"skdNumericEquipType"`
				AssignedNumericEquipType string `json:"assignedNumericEquipType"`
				MaintTurnCode            string `json:"maintTurnCode"`
				EquipReq                 string `json:"equipReq"`
				SasEquipCode             string `json:"sasEquipCode"`
			} `json:"equipment"`
			Times struct {
				STD              string `json:"STD"`
				STA              string `json:"STA"`
				LTD              string `json:"LTD"`
				LTA              string `json:"LTA"`
				SkdTaxiOut       string `json:"skdTaxiOut"`
				SkdTaxiIn        string `json:"skdTaxiIn"`
				DepGMTAdjustment string `json:"depGMTAdjustment"`
				ArrGMTAdjustment string `json:"arrGMTAdjustment"`
				OGT              string `json:"OGT"`
				LegDepDOM        string `json:"legDepDOM"`
				LegArrDOM        string `json:"legArrDOM"`
			} `json:"times"`
			Linkage struct {
				NextLegFltNum      string `json:"nextLegFltNum"`
				NextLegOrgDate     string `json:"nextLegOrgDate"`
				NextLegFltDupCode  string `json:"nextLegFltDupCode"`
				PriorLegFltNum     string `json:"priorLegFltNum"`
				PriorLegOrgDate    string `json:"priorLegOrgDate"`
				PriorLegFltDupCode string `json:"priorLegFltDupCode"`
			} `json:"linkage"`
			Planners struct {
				DispDesk    string `json:"dispDesk"`
				LoadPlanner string `json:"loadPlanner"`
			} `json:"planners"`
			Status struct {
				Dep       string `json:"dep"`
				Arr       string `json:"arr"`
				PaxStatus string `json:"paxStatus"`
			} `json:"status"`
			Type struct {
				CallSign     string `json:"callSign"`
				MealCode     string `json:"mealCode"`
				IntOrDom     string `json:"intOrDom"`
				KickOffFlt   string `json:"kickOffFlt"`
				AntiIceInd   string `json:"antiIceInd"`
				DblProvision string `json:"dblProvision"`
				EROWneeded   string `json:"EROWneeded"`
				EROWcomplete string `json:"EROWcomplete"`
			} `json:"type"`
			CodeShareInfo struct {
				MarketCode string `json:"marketCode"`
				OperatedBy struct {
					IATA      string `json:"IATA"`
					FlightNbr string `json:"flightNbr"`
				} `json:"operatedBy"`
				PartnerFlight []struct {
					IATA      string `json:"IATA"`
					FlightNbr string `json:"flightNbr"`
				} `json:"partnerFlight"`
			} `json:"codeShareInfo"`
		} `json:"leg"`
		LUSInd string `json:"LUSInd"`
	} `json:"flight"`
}

//EtdEvent
type EtdEvent struct {
	Event  string `json:"event"`
	Flight struct {
		Event      string `json:"event"`
		Version    string `json:"version"`
		TrackingID string `json:"trackingID"`
		DataTime   struct {
			SourceTimeStamp string `json:"sourceTimeStamp"`
			FltHubTimeStamp string `json:"fltHubTimeStamp"`
		} `json:"dataTime"`
		Key struct {
			AirlineCode struct {
				IATA string `json:"IATA"`
				ICAO string `json:"ICAO"`
			} `json:"airlineCode"`
			FltNum       string `json:"fltNum"`
			FltOrgDate   string `json:"fltOrgDate"`
			DepSta       string `json:"depSta"`
			DupDepStaNum string `json:"dupDepStaNum"`
		} `json:"key"`
		FosPartition string `json:"fosPartition"`
		Leg          struct {
			Stations struct {
				Arr          string `json:"arr"`
				DupArrStaNum string `json:"dupArrStaNum"`
			} `json:"stations"`
			Times struct {
				STD              string `json:"STD"`
				STA              string `json:"STA"`
				LTD              string `json:"LTD"`
				LTA              string `json:"LTA"`
				PTD              string `json:"PTD"`
				PTA              string `json:"PTA"`
				LatestTaxiOut    string `json:"latestTaxiOut"`
				LatestTaxiIn     string `json:"latestTaxiIn"`
				DepGMTAdjustment string `json:"depGMTAdjustment"`
				ArrGMTAdjustment string `json:"arrGMTAdjustment"`
				AutoETDAccumMins string `json:"autoETDAccumMins"`
			} `json:"times"`
			Status struct {
				Dep string `json:"dep"`
				Arr string `json:"arr"`
			} `json:"status"`
			Type struct {
				CaptQuals    string `json:"captQuals"`
				FltRelStatus string `json:"fltRelStatus"`
				AntiIceInd   string `json:"antiIceInd"`
			} `json:"type"`
			Reason struct {
				Information string `json:"information"`
			} `json:"reason"`
		} `json:"leg"`
	} `json:"flight"`
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
	fmt.Printf("INFO:I am listening to partition %s:%s...\n", entityPath, partitionid)
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

	var cycle CycleEvent
	var etd EtdEvent

	json.Unmarshal(event.Data, &cycle)
	json.Unmarshal(event.Data, &etd)

	if cycle.Event == "CYCLE" {
		//keyStr := fmt.Sprintf("%s-%d %s", member.City, member.Counter, member.DtTm)
		//fmt.Printf("%s:%s\n", entityPath, keyStr)
		fmt.Printf("%s:%s %s %s %s %s\n", entityPath, cycle.Event, cycle.Flight.TrackingID, cycle.Flight.Key.FltNum, cycle.Flight.Key.DepSta, cycle.Flight.Leg.Stations.Arr)
	} else {
		fmt.Printf("%s:%s %s %s %s %s\n", entityPath, etd.Event, etd.Flight.TrackingID, etd.Flight.Key.FltNum, etd.Flight.Key.DepSta, etd.Flight.Leg.Stations.Arr)
	}
	return nil
}
