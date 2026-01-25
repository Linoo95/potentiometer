package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sdoque/mbaigo/components"
	"github.com/sdoque/mbaigo/forms"
	"github.com/sdoque/mbaigo/usecases"
	"github.com/tarm/serial"
)

type Traits struct {
	Port string `json:"port"`
	Baud int    `json:"baud"`
}

// /dev/ttyACM0

// UnitAsset type models the unit asset (interface) of the system
type UnitAsset struct {
	Name        string              `json:"name"`
	Owner       *components.System  `json:"-"`
	Details     map[string][]string `json:"details"`
	ServicesMap components.Services `json:"-"`
	CervicesMap components.Cervices `json:"-"`
	//
	Traits
	RudderPosition int
	lastADC        int
	uartPort       *serial.Port
}

// GetName returns the name of the Resource.
func (ua *UnitAsset) GetName() string {
	return ua.Name
}

// GetServices returns the services of the Resource.
func (ua *UnitAsset) GetServices() components.Services {
	return ua.ServicesMap
}

// GetCervices returns the list of consumed services by the Resource.
func (ua *UnitAsset) GetCervices() components.Cervices {
	return ua.CervicesMap
}

// GetDetails returns the details of the Resource.
func (ua *UnitAsset) GetDetails() map[string][]string {
	return ua.Details
}

// GetTraits returns the traits of the Resource.
func (ua *UnitAsset) GetTraits() any {
	return ua.Traits
}

// ensure UnitAsset implements components.UnitAsset
var _ components.UnitAsset = (*UnitAsset)(nil)

//-------------------------------------Instantiate a unit asset template

// initTemplate initializes a UnitAsset with default values.
func initTemplate() components.UnitAsset {
	// Define the services that expose the capabilities of the unit asset(s)
	position := components.Service{
		Definition:  "rudder position",
		SubPath:     "position",
		Details:     map[string][]string{"Forms": {"SignalA_v1a"}, "Unit": {"Percent"}},
		RegPeriod:   30,
		Description: "return rudder position based on potentiometer ADC value",
	}

	// var uat components.UnitAsset // this is an interface, which we then initialize
	uat := &UnitAsset{
		Name:    "Potentiometer_1",
		Details: map[string][]string{"Model": {"Potentiometer"}, "Location": {"Microcontroller"}},
		ServicesMap: components.Services{
			position.SubPath: &position, // Inline assignment of the rotation service
		},
	}
	return uat
}

// newResource creates the Resource resource with its pointers and channels based on the configuration using the tConfig structs
func newResource(configuredAsset usecases.ConfigurableAsset, sys *components.System) (components.UnitAsset, func()) {
	ua := &UnitAsset{
		Name:        configuredAsset.Name,
		Owner:       sys,
		Details:     configuredAsset.Details,
		ServicesMap: usecases.MakeServiceMap(configuredAsset.Services),
	}
	traits, err := UnmarshalTraits(configuredAsset.Traits)
	if err != nil {
		log.Println("Warning: could not unmarshal traits:", err)
	} else if len(traits) > 0 {
		ua.Traits = traits[0]
	}
	// open UART connection
	log.Printf("Potentiometer %s using UART %s @ %d baud", ua.Name, ua.Port, ua.Baud)
	configure := &serial.Config{
		Name: ua.Traits.Port,
		Baud: ua.Traits.Baud,
	}

	port, err := serial.OpenPort(configure)

	if err != nil {
		log.Fatal("failed to open UART:", err)
	}

	ua.uartPort = port

	go ua.readUART()

	cleanup := func() {
		log.Println("closing potentiometer", ua.Name)
		if ua.uartPort != nil {
			_ = ua.uartPort.Close()
		}
	}

	return ua, cleanup
}
func (ua *UnitAsset) readUART() {
	reader := bufio.NewReader(ua.uartPort)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("UART read error:", err)
			continue
		}
		log.Println("Raw UART line:", fmt.Sprintf("%q", line))

		line = strings.TrimSpace(line)
		adc, err := strconv.Atoi(line)
		if err != nil {
			log.Println("Failed to parse ADC:", line, err)
			continue
		}
		ua.lastADC = adc
		ua.RudderPosition = adcToPercent(adc)
	}
}

// UnmarshalTraits unmarshals a slice of json.RawMessage into a slice of Traits.
func UnmarshalTraits(rawTraits []json.RawMessage) ([]Traits, error) {
	var traitsList []Traits
	for _, raw := range rawTraits {
		var t Traits
		if err := json.Unmarshal(raw, &t); err != nil {
			return nil, fmt.Errorf("failed to unmarshal trait: %w", err)
		}
		traitsList = append(traitsList, t)
	}
	return traitsList, nil
}

func (ua *UnitAsset) position(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		form := ua.getPositiom()
		usecases.HTTPProcessGetRequest(w, r, &form)
	default:
		http.Error(w, "not supported", http.StatusNotFound)
	}
}

func adcToPercent(adc int) int {
	return (adc * 100) / 4095
}
func (ua *UnitAsset) getPositiom() (f forms.SignalA_v1a) {
	f.NewForm()

	f.Value = float64(ua.RudderPosition)
	f.Unit = "Percent"
	f.Timestamp = time.Now()

	return f
}
func (ua *UnitAsset) Serving(w http.ResponseWriter, r *http.Request, servicePath string) {
	switch servicePath {
	case "position":
		ua.position(w, r)
	default:
		http.Error(
			w,
			"Invalid service request",
			http.StatusBadRequest,
		)
	}
}
