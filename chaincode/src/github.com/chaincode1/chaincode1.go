
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode poc simple Chaincode implementation
type SimpleChaincode struct {
}

type ship struct {
	ObjectType      string        `json:"docType"`         //docType is used to distinguish the various types of objects in state database
	ShipmentOrderNo string        `json:"ShipmentOrderNo"` 
	Status          string        `json:"Status"`
	FirstName       string        `json:"FirstName"`
	LastName        string        `json:"LastName"`
	EmailId         string        `json:"EmailId"`
	AccountBalance  int           `json:"AccountBalance"`
	Temperature     []temperature `json:"Temperature"`
	Fuel            []fuel        `json:"Fuel"`
}

type temperature struct {
	TempReading   int       `json:"TempReading"`
	TimeOfReading time.Time `json:"timeOfReading"`
}

type fuel struct {
	FuelReading   int       `json:"fuelReading"`
	TimeOfReading time.Time `json:"timeOfReading"`
}

// ===================================================================================
//                                  Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting the chaincode: %s", err)
	}
}

// ===================================================================================
//                       Init initializes chaincode
// ===================================================================================

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Entering Contract init successfully")
	return shim.Success(nil)
}

//===============================================================
//              Invoke - Our entry point for Invocations
// ==============================================================

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "initShipment" { // create a new shipment with initialisation
		return t.initShipment(stub, args)
	} else if function == "setTemp" { // insert new temperature 
		return t.setTemp(stub, args)
	} else if function == "readShipment" { // read the shipment
		return t.readShipment(stub, args)
	} else if function == "getHistoryForShipment" { // get history of shipments
		return t.getHistoryForShipment(stub, args)
	} else if function == "shipmentDelivered" { //set the shipment is delivered
		return t.shipmentDelivered(stub, args)
	} else if function == "setFuel" { //set fuel reading of shipment
		return t.setFuel(stub, args)
	}
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
//               Initialise Shipment
// ============================================================

func (t *SimpleChaincode) initShipment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// ==== Input sanitation ====
	fmt.Println("- start init Shipment")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument Shipment must be a non-empty string")
	}
	
	ShipmentOrderNo := args[0]
	Status := args[1]
	FirstName := args[2]
	LastName := args[3]
	EmailId := args[4]
	AccountBalance, err := strconv.Atoi(args[5])
	
	tempReading, err := strconv.Atoi(args[6])
	fmt.Println("Temperature is %d", tempReading)
	if err != nil {
		return shim.Error("5th argument must be a numeric string")
	}

	fuelReading, err := strconv.Atoi(args[7])
	fmt.Println("fuel is %d", fuelReading)
	if err != nil {
		return shim.Error("6th argument must be a numeric string")
	}
	fmt.Println("Fuel is %d", fuelReading)

	temp1 := temperature{tempReading, time.Now()}
	temps := []temperature{}

	fuel1 := fuel{fuelReading, time.Now()}
	fuels := []fuel{}

	objectType := "ship"
	ship := &ship{objectType, ShipmentOrderNo, Status, FirstName, LastName, EmailId, AccountBalance, temps, fuels}

	ship.Temperature = append(ship.Temperature, temp1)
	ship.Fuel = append(ship.Fuel, fuel1)
	fmt.Println("Ship is ", ship)
	shipJSONasBytes, err := json.Marshal(ship)

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("After marshal ship is ", shipJSONasBytes)

	// === Save Shipment to state ===
	err = stub.PutState(ShipmentOrderNo, shipJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent("sampleEvent", []byte(ShipmentOrderNo))
	if err != nil {
		fmt.Println("SetEvent Error", err)
	}
	// ==== Shipment saved. Return success ====
	fmt.Println("- end init Shipment")
	return shim.Success(nil)
}

// ===============================================
//           Reading Shipment
// ===============================================

func (t *SimpleChaincode) readShipment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var ShipmentOrder, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting Shipment id to query")
	}

	ShipmentOrder = args[0]
	valAsbytes, err := stub.GetState(ShipmentOrder) //get the shipment status from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + ShipmentOrder + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Shipment does not exist: " + ShipmentOrder + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ===============================================
//          Set Temperature Reading
// ===============================================

func (t *SimpleChaincode) setTemp(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var MinimumTemperature = 2
	var MaximumTemperature = 10


	ShipmentOrderNo := args[0]

	tempAsBytes, _ := stub.GetState(args[0])
	Ship := ship{}

	json.Unmarshal(tempAsBytes, &Ship)

	if Ship.Status != "IN-TRANSIT" {
		return shim.Error("The shipment is not IN-TRANSIT")
	}

	tempReading, err := strconv.Atoi(args[1])

	temp1 := temperature{tempReading, time.Now()}

	if err != nil {
		return shim.Error("2nd argument must be a numeric string")
	}

	if tempReading <= MinimumTemperature {

		err = stub.SetEvent("MinTempEvent", []byte(ShipmentOrderNo))
		if err != nil {
			fmt.Println("SetEvent Error", err)
		}
		fmt.Println("Less than Minimum Temperature ", tempReading)
	}

	if tempReading >= MaximumTemperature {

		err = stub.SetEvent("MaxTempEvent", []byte(ShipmentOrderNo))
		if err != nil {
			fmt.Println("SetEvent Error", err)
		}
		fmt.Println("Temeperature is Greater than maximum :", tempReading)
	}

	fmt.Println("temp1 is ", temp1)

	if len(Ship.Temperature) >= 4 {

		var x = len(Ship.Temperature)
		fmt.Println("Length of temperature ", x)

		var i = 0
		var deviation = 0

		for i = 0; i < 4; i++ {
			//fmt.Println("Length of temperature ", x, Ship.Temperature[x].TempReading)
			reading := Ship.Temperature[x-1].TempReading
			fmt.Println("Reading is ",reading)
			if reading <= MinimumTemperature || reading >= MaximumTemperature {
				deviation++
				fmt.Println("Deviation ois ", deviation)
			}
			x--;
		}
		if deviation == 4 {
			err = stub.SetEvent("TempDeviationEvent", []byte(ShipmentOrderNo))
			if err != nil {
				fmt.Println("SetEvent Error", err)

			}

			fmt.Println("Temperature deviated continuously")
			Ship.AccountBalance += 100
		}
	}

	Ship.Temperature = append(Ship.Temperature, temp1)
	tempAsBytes, _ = json.Marshal(Ship)
	stub.PutState(args[0], tempAsBytes)

	return shim.Success(nil)
}

// ===============================================
//          Set Temperature Reading
// ===============================================
func (t *SimpleChaincode) setFuel(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var MinimumFuel = 20

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ShipmentOrderNo := args[0]

	fuelAsBytes, _ := stub.GetState(args[0])
	Ship := ship{}

	json.Unmarshal(fuelAsBytes, &Ship)
	//car.Owner = args[1]

	if Ship.Status != "IN-TRANSIT" {
		return shim.Error("The shipment is not IN-TRANSIT")
	}

	fuelReading, err := strconv.Atoi(args[1])

	fuel1 := fuel{fuelReading, time.Now()}

	if fuelReading <= MinimumFuel {

		err = stub.SetEvent("MinFuelEvent", []byte(ShipmentOrderNo))
		if err != nil {
			fmt.Println("SetEvent Error", err)
		}
		fmt.Println("Less than the Minimum Fuel ", fuelReading)
	}
	Ship.Fuel = append(Ship.Fuel, fuel1)
	fuelAsBytes, _ = json.Marshal(Ship)
	stub.PutState(args[0], fuelAsBytes)

	return shim.Success(nil)
}

// ===============================================
//          Shipment Delivered
// ===============================================

func (t *SimpleChaincode) shipmentDelivered(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	shipAsBytes, _ := stub.GetState(args[0])
	Ship := ship{}

	json.Unmarshal(shipAsBytes, &Ship)
	Ship.Status = "DELIVERED"

	shipAsBytes, _ = json.Marshal(Ship)
	stub.PutState(args[0], shipAsBytes)

	return shim.Success(nil)
}

// ===============================================
//                History for Shipment
// ===============================================

func (t *SimpleChaincode) getHistoryForShipment(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	shipmentOrderNo := args[0]

	fmt.Printf("- start getHistoryForShipment: %s\n", shipmentOrderNo)

	resultsIterator, err := stub.GetHistoryForKey(shipmentOrderNo)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the Shipment
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON shipment)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForShipment returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}