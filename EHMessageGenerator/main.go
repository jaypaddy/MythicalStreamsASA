package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go"
	uuid "github.com/nu7hatch/gouuid"
)

//cycle to Register
type cycle struct {
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

//etd
type etd struct {
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

var citycode = []string{"DFW", "ORD", "MIA", "LAX", "PHI", "AUS", "SFO", "DBX", "MAA", "KOC", "LHR", "MAA", "BLR", "PAR", "LHR", "XYZ"}
var fltnum = []string{"1800", "1820", "1850", "1870", "1900", "1920", "1950", "1970", "2000", "2020", "2050", "2070"}

//https://notes.shichao.io/gopl/ch9/

var (
	mu          sync.Mutex // guards memberCount
	memberCount int
)

func main() {

	memberCount = 1

	//Read cycle.json
	cycleJSONByte, err := ioutil.ReadFile("cycle.json")
	if err != nil {
		fmt.Printf("Error reading cycle.json.\n File reading error:\n%s", err)
		return
	}
	//Read etd.json
	etdJSONByte, err := ioutil.ReadFile("etd.json")
	if err != nil {
		fmt.Printf("Error reading etd.json.\n File reading error:\n%s", err)
		return
	}

	//Unmarshal cycle
	var f interface{}
	//b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
	err = json.Unmarshal(etdJSONByte, &f)

	//Gather EH connection details
	connStr := os.Getenv("MSG_EVENTHUBNS")
	if connStr == "" {
		fmt.Println("Missing namespace : MSG_EVENTHUBNS")
		return
	}
	cycleEH := os.Getenv("MSG_CYCLEEH")
	if cycleEH == "" {
		fmt.Println("Missing topic : MSG_CYCLEEH")
		return
	}
	etdEH := os.Getenv("MSG_ETDEH")
	if etdEH == "" {
		fmt.Println("Missing topic : MSG_ETDEH")
		return
	}

	for _, cityCODE := range citycode {
		go GenMessage(cycleJSONByte, "CYCLE", connStr, cycleEH, cityCODE)
		go GenMessage(etdJSONByte, "ETD", connStr, etdEH, cityCODE)

	}

	//Wait for a signal to quit
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

}

//GenMessage - get the city code for filtering
func GenMessage(JSONByte []byte, eventType string, connStr string, ehName string, city string) {

	connStr = fmt.Sprintf("%s;EntityPath=%s", connStr, ehName)
	hub, err := eventhub.NewHubFromConnectionString(connStr)
	if err != nil {
		fmt.Printf("ERR:Error connecting to .\n")
	}

	for x := 0; x < 10000000000; x++ {
		//var printStr string
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		mu.Lock()
		memberCount++
		mu.Unlock()

		eventJSON, ticker := BuildMessage(JSONByte, city, eventType)
		event := eventhub.NewEvent([]byte(eventJSON))
		event.PartitionKey = &city

		err = hub.Send(ctx, event)
		if err != nil {
			fmt.Println(err)
			return
		}
		//rand.Seed(time.Now().UnixNano())
		//k := rand.Intn(10)
		//k := 1 //Second Wait between Member Generation

		dispInfo := fmt.Sprintf("%s", ticker)
		fmt.Printf("%s\n", dispInfo)
		time.Sleep(time.Duration(10) * time.Second)
	}

}

//BuildMessage - get the city code for filtering
func BuildMessage(jsonTemplate []byte, depCity string, eventType string) (string, string) {

	eventJSON := string(jsonTemplate)
	/*
		"trackingID": "<TRACKINGID>",
			"sourceTimeStamp": "<SRCTIMESTAMP>",
			"fltHubTimeStamp": "<FLTHUBTIMESTAMP>"
			"fltNum": "<FLTNUM>",
			"depSta": "<DEP>",
			"arr": "<ARR>",
	*/
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	dtTm := time.Now().Format(time.RFC3339)
	trackingID := uuid.String()
	eventJSON = strings.ReplaceAll(eventJSON, "<TRACKINGID>", trackingID)
	eventJSON = strings.ReplaceAll(eventJSON, "<SRCTIMESTAMP>", dtTm)
	eventJSON = strings.ReplaceAll(eventJSON, "<FLTHUBTIMESTAMP>", dtTm)
	fltNumDetail := fmt.Sprintf("AA%s", GetFltNum())
	eventJSON = strings.ReplaceAll(eventJSON, "<FLTNUM>", fltNumDetail)
	eventJSON = strings.ReplaceAll(eventJSON, "<DEP>", depCity)
	arrCity := GetCityCode()
	eventJSON = strings.ReplaceAll(eventJSON, "<ARR>", arrCity)

	retStr := fmt.Sprintf("%s:%s %s %s %s", eventType, trackingID, fltNumDetail, depCity, arrCity)

	return eventJSON, retStr

}

//GetCityCode - get the city code for filtering
func GetCityCode() string {
	l := len(citycode)
	rand.Seed(time.Now().UnixNano())
	k := rand.Intn(l - 1)
	return citycode[k]
}

//GetFltNum - get the Flt Num for filtering
func GetFltNum() string {
	l := len(fltnum)
	rand.Seed(time.Now().UnixNano())
	k := rand.Intn(l - 1)
	return fltnum[k]
}
