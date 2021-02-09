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

type basiceDetails  struct {
	TypeOf				string 	`json:"TypeOf"`
	Name				string  `json:"Name"`
	EmailId				string	`json:"EmailId"`
	PhoneNo				string	`json:"PhoneNo"`
	Address				string 	`json:"Address"`
	City				string	`json:"City"`
	State				string 	`json:"State"`
	ZipCode				string 	`json:"ZipCode"`
}

type storageUnits struct {
	StorageUnitsName				string 	`json:"StorageUnitsName"`
	StorageUnitsAddress				string  `json:"StorageUnitsAddress"`
	StorageUnitsCity				string	`json:"StorageUnitsCity"`
	StorageUnitsState				string	`json:"StorageUnitsState"`
	StorageUnitsZipCode				string 	`json:"StorageUnitsZipCode"`
	StorageUnitsLatitude			string	`json:"StorageUnitsLatitude"`
	StorageUnitsLongitude			string 	`json:"StorageUnitsLongitude"`
}

type automationIntegration struct {
	SystemTool 				string	`json:"SystemTool"`
	AuthenticationType 		string	`json:"AuthenticationType"`
	HostAddress				string	`json:"HostAddress"`
	UserName				string 	`json:"UserName"`
	Password				string	`json:"Password"`
	Port					string 	`json:"Port"`

}

type productInfo  struct {

	ProductId			string	 	`json:"ProductId"`
	ProductName			string	 	`json:"ProductName"`
	ProductType			string		`json:"ProductType"`
	ProductCategory		string 		`json:"ProductCategory"`
	ProductDescription	string 		`json:"ProductDescription"`
	ProductImage		string 		`json:"ProductImage"`

}


type excelData struct {

	FileName				string 				`json:"FileName"`
	FileData				[]productInfo 		`json:"FileData"`

}


type product  struct {
	AutomationIntegration  				automationIntegration 		`json:"AutomationIntegration"`
	ExcelData							excelData				    `json:"ExcelData"`
}


type entityDetails struct {
	BasiceDetails 			   basiceDetails 				`json:"BasiceDetails"`
	StorageUnits 		   	   []storageUnits   			`json:"StorageUnits"`
	Products 	  			   []product					`json:"Products"`
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
	if function == "enrollEntity" { // create a new shipment with initialisation
		return t.enrollEntity(stub, args)
	}  else if function == "readEntity" { // create a new shipment with initialisation
		return t.readEntity(stub, args)
	}  else if function == "getHistoryForShipment" { // create a new shipment with initialisation
		return t.getHistoryForShipment(stub, args)
	}	
	
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
//               Initialise Shipment
// ============================================================

func (t *SimpleChaincode) enrollEntity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var EntityAsJSON entityDetails 

	
	// ==== Input sanitation ====
	fmt.Println("- start init Shipment")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument Shipment must be a non-empty string")
	}
	
	 
	json.Unmarshal(args[0], &EntityAsJSON)
	 
	EmailId := EntityAsJSON.BasiceDetails.EmailId 

	fmt.Println("Entity Details in JSON is  ", EntityAsJSON)
	entityDetailsJSONasBytes, err := json.Marshal(EntityAsJSON)

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("After marshal Entity is ", entityDetailsJSONasBytes)

	// === Save Entity to state ===
	err = stub.PutState(EmailId, entityDetailsJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent("sampleEvent",entityDetailsJSONasBytes)
	if err != nil {
		fmt.Println("SetEvent Error", err)
	}
	// ==== Entity saved. Return success ====
	fmt.Println("- end init Entity")
	return shim.Success(nil)
}


// ===============================================
//           Reading Entity 
// ===============================================

func (t *SimpleChaincode) readEntity(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var EmailId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting Shipment id to query")
	}

	EmailId = args[0]
	valAsbytes, err := stub.GetState(EmailId) //get the shipment status from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + EmailId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Shipment does not exist: " + EmailId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
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