
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

type batch struct {
	ObjectType      	string        `json:"docType"`         //docType is used to distinguish the various types of objects in state database
	BatchId 			string        `json:"BatchId"` 
	FarmarId        	string        `json:"FarmarId"`
	AnimalHealthRecord  string        `json:"AnimalHealthRecord"`
	FatPercentage       string        `json:"FatPercentage"`
	DateTimeCollection  time.Time     `json:"DateTimeCollection"`
	TransitId		  	string        `json:"TransitId"`
	DealerId     		string 		  `json:"DealerId"`
	IngredientList      string      `json:"IngredientList"`
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
	} else if function == "setDealerInput" { // insert new temperature 
		return t.setDealerInput(stub, args)
	} else if function == "readBatch" { // read the shipment
		return t.readBatch(stub, args)
	} else if function == "getHistoryForShipment" { // get history of shipments
		return t.getHistoryForShipment(stub, args)
	}  else if function == "setDairyFarmInput" { //set fuel reading of shipment
		return t.setDairyFarmInput(stub, args)
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
	
	
	BatchId := args[0]
	FarmarId := args[1]
	AnimalHealthRecord := args[2]
	FatPercentage := args[3]
	DateTimeCollection := time.Now()
	TransitId := ""
	DealerId := ""
	IngredientList := ""
	

	objectType := "ship"
	batch:= &batch{objectType,BatchId,FarmarId,AnimalHealthRecord,FatPercentage,DateTimeCollection,TransitId,DealerId,IngredientList}

	
	fmt.Println("batch is ", batch)
	batchJSONasBytes, err := json.Marshal(batch)

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("After marshal ship is ", batchJSONasBytes)

	// === Save Batch to state ===
	err = stub.PutState(BatchId, batchJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent("sampleEvent", []byte(BatchId))
	if err != nil {
		fmt.Println("SetEvent Error", err)
	}
	// ==== Batch saved. Return success ====
	fmt.Println("- end init Batch")
	return shim.Success(nil)
}


// ===============================================
//           Reading Shipment
// ===============================================

func (t *SimpleChaincode) readBatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var BatchId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting Batch id to query")
	}

	BatchId = args[0]
	valAsbytes, err := stub.GetState(BatchId) //get the Batch status from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + BatchId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Batch does not exist: " + BatchId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ===============================================
//          Set Dealer Input
// ===============================================

func (t *SimpleChaincode) setDealerInput(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	
	BatchIdAsBytes, _ := stub.GetState(args[0])
	batch := batch{}

	json.Unmarshal(BatchIdAsBytes, &batch)

	batch.DealerId = args[1]
	BatchIdAsBytes, _ = json.Marshal(batch)
	stub.PutState(args[0], BatchIdAsBytes)

	return shim.Success(nil)
}

// ===============================================
//          Set Dairy Farm Input
// ===============================================
func (t *SimpleChaincode) setDairyFarmInput(stub shim.ChaincodeStubInterface, args []string) pb.Response {


	BatchAsBytes, _ := stub.GetState(args[0])
	batch := batch{}

	json.Unmarshal(BatchAsBytes, &batch)


	batch.IngredientList = args[1]
	BatchAsBytes, _ = json.Marshal(batch)
	stub.PutState(args[0], BatchAsBytes)

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