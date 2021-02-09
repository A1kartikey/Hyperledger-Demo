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


type milkFarm struct {
	Name     	      		 string          			`json:"Name"`
	EmailId           		 string          			`json:"EmailId"`
	ShipmentRecord		   []string                     `json:"ShipmentRecord"`
}

type milkFarmList struct {
	MilkFarmStruct			  		[]milkFarmStruct 				`json:"MilkFarmStruct"`
}

type milkFarmStruct struct {
	Name     	      		 string          			`json:"Name"`
	EmailId           		 string          			`json:"EmailId"`	
}

type chillingAndProccessingUnit struct {
	Name     	    string        `json:"Name"`
	EmailId         string        `json:"EmailId"`
	ShipmentRecord []string       `json:"ShipmentRecord"`
}

type packagingUnit struct {
	Name     	    string        `json:"Name"`
	EmailId         string        `json:"EmailId"`
	ShipmentRecord []string       `json:"ShipmentRecord"`
}

type freight struct {
	Name     	    string        `json:"Name"`
	EmailId         string        `json:"EmailId"`
	ShipmentRecord []string       `json:"ShipmentRecord"`
}

type market struct {
	Name     	    string        `json:"Name"`
	EmailId         string        `json:"EmailId"`
	ShipmentRecord []string       `json:"ShipmentRecord"`
}

type consumer struct {
	Name     	    string        `json:"Name"`
	EmailId         string        `json:"EmailId"`
	ShipmentRecord []string       `json:"ShipmentRecord"`
}

type consumerList struct {
	ConsumerStruct 			  []consumerStruct 				`json:"consumerStruct"`
}

type consumerStruct struct {
	Name     	      		 string          			`json:"Name"`
	EmailId           		 string          			`json:"EmailId"`	
}


type batch struct {
	BatchId     	    string        `json:"BatchId"`
	Owner         		string        `json:"Owner"`
	Category     	    string        `json:"Category"`
	Manufacturer        string        `json:"Manufacturer"`
	Name     	    	string        `json:"Name"`
	Ph         			string        `json:"Ph"`
	Apperance     	    string        `json:"Apperance"`
	MilkType         	string        `json:"MilkType"`
}

type product struct {
	ProductId      	    string        `json:"ProductId"`
	Name          		string        `json:"Name"`
	Owner			    string		  `json:"Owner"`
	Manufacturer  	    string 		  `json:"Manufacturer"`
	Cost 			    string		  `json:"Cost"`
	MilkCategory	    string		  `json:"MilkCategory"` 
	PackageCategory     string		  `json:"PackageCategory"`
	Batch 	            string		  `json:"Batch"`
	Logs 				[]string 	  `json:Logs`
	Timestamp			[]time.Time	  `json:"Timestamp"`
	
}

var MilkFarmList = milkFarmList{}
var ConsumerList = consumerList{}

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
	if function == "enrollMilkFarm" { // create a new Milk-Farm with initialisation
		return t.enrollMilkFarm(stub, args)
	} else if function == "readMilkFarm" { // read Milk-Farm 
		return t.readMilkFarm(stub, args)
	} else if function == "enrollChillingAndProccessingUnit" { // read Milk-Farm 
		return t.enrollChillingAndProccessingUnit(stub, args)
	} else if function == "readChillingAndProccessingUnit" { // read Milk-Farm 
		return t.readChillingAndProccessingUnit(stub, args)
	} else if function == "enrollPackagingUnit" { // read Milk-Farm 
		return t.enrollPackagingUnit(stub, args)
	} else if function == "readPackagingUnit" { // read Milk-Farm 
		return t.readPackagingUnit(stub, args)
	} else if function == "enrollFreight" { // read Milk-Farm 
		return t.enrollFreight(stub, args)
	} else if function == "readFreight" { // read Milk-Farm 
		return t.readFreight(stub, args)
	} else if function == "enrollMarket" { // read Milk-Farm 
		return t.enrollMarket(stub, args)
	} else if function == "readMarket" { // read Milk-Farm 
		return t.readMarket(stub, args)
	} else if function == "enrollConsumer" { // read Milk-Farm 
		return t.enrollConsumer(stub, args)
	} else if function == "readConsumer" { // read Milk-Farm 
		return t.readConsumer(stub, args)
	} else if function == "addBatch" { // read Milk-Farm 
		return t.addBatch(stub, args)
	} else if function == "readBatch" { // read Milk-Farm 
		return t.readBatch(stub, args)
	} else if function == "createProduct" { // read Milk-Farm 
		return t.createProduct(stub, args)
	} else if function == "readProduct" { // read Milk-Farm 
		return t.readProduct(stub, args)
	} else if function == "changeOwnershipProduct" { // read Milk-Farm 
		return t.changeOwnershipProduct(stub, args)
	} else if function == "changeOwnershipBatch" { // read Milk-Farm 
		return t.changeOwnershipBatch(stub, args)
	} else if function == "changeOwnershipBatch" { // read Milk-Farm 
		return t.changeOwnershipBatch(stub, args)
	} else if function == "getShipmentRecord" { // read Milk-Farm 
		return t.getShipmentRecord(stub, args)
	} else if function == "getHistoryForAsset" { // read Milk-Farm 
		return t.getHistoryForAsset(stub, args)
	} else if function == "readMilkFarmList" { // read Milk-Farm 
		return t.readMilkFarmList(stub)
	} else if function == "readConsumerList" { // read Milk-Farm 
		return t.readConsumerList(stub, args)
	}     	
	
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
//               Initialise Milk-Farm
// ============================================================

func (t *SimpleChaincode) enrollMilkFarm(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// ==== Input sanitation ====
	fmt.Println("- start enroll Milk-Farm")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument email of enrollment must be a non-empty string")
	}

	EmailId := args[0]	
	Name := args[1]
// ====================== Enroll Milk Farm ===========================

	var ShipmentRecord1 []string
	ShipmentRecord1 = append(ShipmentRecord1,"") 
	ShipmentRecord := ShipmentRecord1

	fmt.Println("Name: ", Name ,"EmailId: ",EmailId,"Shipment Record: ",ShipmentRecord)

	MilkFarm := &milkFarm{Name, EmailId,ShipmentRecord}
	
	milkFarmStruct1 := milkFarmStruct{Name,EmailId}
		
	MilkFarmList.MilkFarmStruct = append(MilkFarmList.MilkFarmStruct,milkFarmStruct1) 

	fmt.Println("MilkFarm is ", MilkFarm)	
	MilkFarmJSONasBytes, err := json.Marshal(MilkFarm)

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("After marshal ship is ", MilkFarmJSONasBytes)

	// === Save Milk-Farm to state ===
	err = stub.PutState(EmailId, MilkFarmJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}


	// === Set Event Milk-Farm  ===
	
	err = stub.SetEvent("sampleEvent", []byte(EmailId))
	if err != nil {
		fmt.Println("SetEvent Error", err)
	}
	// ==== Milk-Farm saved. Return success ====
	fmt.Println("- end enroll Milk-Farm")
	return shim.Success(nil)
}



// ===============================================
//           Reading Milk-Farm list
// ===============================================
func (t *SimpleChaincode) readMilkFarmList(stub shim.ChaincodeStubInterface) pb.Response {
	var err error

	MilkFarmListAsbytes, err:= json.Marshal(MilkFarmList)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("After marshal MilkFarmListAsbytes is ", MilkFarmListAsbytes)
	
	return shim.Success(MilkFarmListAsbytes)
}

// ===============================================
//           Reading Milk-Farm
// ===============================================

func (t *SimpleChaincode) readMilkFarm(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var milkFarmEmailId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting Milk-Farm EmailId to query")
	}

	milkFarmEmailId = args[0]
	valAsbytes, err := stub.GetState(milkFarmEmailId) //get the milk-Farm status from chaincode state.
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + milkFarmEmailId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Shipment does not exist: " + milkFarmEmailId + "\"}"
		return shim.Error(jsonResp)
	}

	fmt.Println("json response",valAsbytes)

	return shim.Success(valAsbytes)
}

// ===============================================
//           Reading Milk-Farm Shipments
// ===============================================

func (t *SimpleChaincode) getShipmentRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var EmailId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting Milk-Farm EmailId to query")
	}

	EmailId = args[0]
	valAsbytes, err := stub.GetState(EmailId) //get the milk-Farm status from chaincode state.
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + EmailId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Shipment does not exist: " + EmailId + "\"}"
		return shim.Error(jsonResp)
	}

	 valUnit:= milkFarm{}

	json.Unmarshal(valAsbytes, &valUnit)

	fmt.Println("json response: ",valAsbytes)
	fmt.Println("Milk-Farm unmarshaled: ",valUnit)
	
	ShipmentRecordAsbytes,_ := json.Marshal(valUnit.ShipmentRecord)

	return shim.Success(ShipmentRecordAsbytes)
}

// =====================================================================
//               Initialise Chilling And Proccessing Unit
// =====================================================================

func (t *SimpleChaincode) enrollChillingAndProccessingUnit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// ==== Input sanitation ====
	fmt.Println("- start enroll Chilling-And-Proccessing-Unit")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument email of enrollment must be a non-empty string")
	}

	EmailId := args[0]	
	Name := args[1]

	var ShipmentRecord1 []string
	ShipmentRecord1 = append(ShipmentRecord1,"") 
	ShipmentRecord := ShipmentRecord1


	fmt.Println("Name: ", Name ,"EmailId: ",EmailId,"ShipmentRecord: ",ShipmentRecord)

	chillingAndProccessingUnit := &chillingAndProccessingUnit{Name,EmailId,ShipmentRecord}

	fmt.Println("Chilling And Proccessing Unit is ", chillingAndProccessingUnit)
	chillingAndProccessingUnitJSONasBytes, err := json.Marshal(chillingAndProccessingUnit)

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("After marshal ship is ", chillingAndProccessingUnitJSONasBytes)

	// === Save Chilling And Proccessing Unit to state ===
	err = stub.PutState(EmailId, chillingAndProccessingUnitJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Set Event Chilling And Proccessing Unit  ===
	
	err = stub.SetEvent("sampleEvent", []byte(EmailId))
	if err != nil {
		fmt.Println("SetEvent Error", err)
	}
	// ==== Chilling And Proccessing Unit saved. Return success ====
	fmt.Println("- end enroll Chilling And Proccessing Unit")
	return shim.Success(nil)
}

// ================================================================
//           Reading Chilling And Proccessing Unit
// ================================================================

func (t *SimpleChaincode) readChillingAndProccessingUnit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var chillingAndProccessingUnitEmailId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting Milk-Farm EmailId to query")
	}

	chillingAndProccessingUnitEmailId = args[0]
	valAsbytes, err := stub.GetState(chillingAndProccessingUnitEmailId) //get the milk-Farm status from chaincode state.
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + chillingAndProccessingUnitEmailId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Shipment does not exist: " + chillingAndProccessingUnitEmailId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// =====================================================================
//               Initialise Packaging Unit
// =====================================================================

func (t *SimpleChaincode) enrollPackagingUnit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// ==== Input sanitation ====
	fmt.Println("- start enroll Packaging-Unit")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument email of enrollment must be a non-empty string")
	}
	
	EmailId := args[0]	
	Name := args[1]

	var ShipmentRecord1 []string
	ShipmentRecord1 = append(ShipmentRecord1,"") 
	ShipmentRecord := ShipmentRecord1


	fmt.Println("Name: ", Name ,"EmailId: ",EmailId,"ShipmentRecord: ",ShipmentRecord)

	packagingUnit := &packagingUnit{Name,EmailId,ShipmentRecord}

	fmt.Println("packaging Unit is ", packagingUnit)
	packagingUnitJSONasBytes, err := json.Marshal(packagingUnit)

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("After marshal ship is ", packagingUnitJSONasBytes)

	// === Save packaging Unit to state ===
	err = stub.PutState(EmailId, packagingUnitJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Set Event packaging Unit  ===
	
	err = stub.SetEvent("sampleEvent", []byte(EmailId))
	if err != nil {
		fmt.Println("SetEvent Error", err)
	}
	// ==== packaging Unit saved. Return success ====
	fmt.Println("- end enroll packaging Unit")
	return shim.Success(nil)
}

// ================================================================
//           Reading Packaging Unit
// ================================================================

func (t *SimpleChaincode) readPackagingUnit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var packagingUnitEmailId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting Milk-Farm EmailId to query")
	}

	packagingUnitEmailId = args[0]
	valAsbytes, err := stub.GetState(packagingUnitEmailId) //get the milk-Farm status from chaincode state.
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + packagingUnitEmailId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Shipment does not exist: " + packagingUnitEmailId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// =====================================================================
//               Initialise Freight 
// =====================================================================

func (t *SimpleChaincode) enrollFreight(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// ==== Input sanitation ====
	fmt.Println("- start enroll freight")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument email of enrollment must be a non-empty string")
	}

	EmailId := args[0]	
	Name := args[1]

	var ShipmentRecord1 []string
	ShipmentRecord1 = append(ShipmentRecord1,"") 
	ShipmentRecord := ShipmentRecord1


	fmt.Println("Name: ", Name ,"EmailId: ",EmailId,"ShipmentRecord: ",ShipmentRecord)

	freight := &freight{Name,EmailId,ShipmentRecord}

	fmt.Println("freight Unit is ", freight)
	freightJSONasBytes, err := json.Marshal(freight)

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("After marshal ship is ", freightJSONasBytes)

	// === Save freight to state ===
	err = stub.PutState(EmailId, freightJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Set Event freight Unit  ===
	
	err = stub.SetEvent("sampleEvent", []byte(EmailId))
	if err != nil {
		fmt.Println("SetEvent Error", err)
	}
	// ==== freight Unit saved. Return success ====
	fmt.Println("- end enroll freight Unit")
	return shim.Success(nil)
}

// ================================================================
//           Reading Freight
// ================================================================

func (t *SimpleChaincode) readFreight(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var freightEmailId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting Milk-Farm EmailId to query")
	}

	freightEmailId = args[0]
	valAsbytes, err := stub.GetState(freightEmailId) //get the milk-Farm status from chaincode state.
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + freightEmailId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"freight does not exist: " + freightEmailId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}


// =====================================================================
//               Initialise Market 
// =====================================================================

func (t *SimpleChaincode) enrollMarket(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// ==== Input sanitation ====
	fmt.Println("- start enroll market")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument email of enrollment must be a non-empty string")
	}

	EmailId := args[0]	
	Name := args[1]

	var ShipmentRecord1 []string
	ShipmentRecord1 = append(ShipmentRecord1,"") 
	ShipmentRecord := ShipmentRecord1


	fmt.Println("Name: ", Name ,"EmailId: ",EmailId,"ShipmentRecord: ",ShipmentRecord)

	market := &market{Name,EmailId,ShipmentRecord}

	fmt.Println("market Unit is ", market)
	marketJSONasBytes, err := json.Marshal(market)

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("After marshal ship is ", marketJSONasBytes)

	// === Save market to state ===
	err = stub.PutState(EmailId, marketJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Set Event market Unit  ===
	
	err = stub.SetEvent("sampleEvent", []byte(EmailId))
	if err != nil {
		fmt.Println("SetEvent Error", err)
	}
	// ==== market Unit saved. Return success ====
	fmt.Println("- end enroll market Unit")
	return shim.Success(nil)
}

// ================================================================
//           Reading Market
// ================================================================

func (t *SimpleChaincode) readMarket(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var marketEmailId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting Milk-Farm EmailId to query")
	}

	marketEmailId = args[0]
	valAsbytes, err := stub.GetState(marketEmailId) //get the milk-Farm status from chaincode state.
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + marketEmailId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"market does not exist: " + marketEmailId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// =====================================================================
//               Initialise consumer 
// =====================================================================

func (t *SimpleChaincode) enrollConsumer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// ==== Input sanitation ====
	fmt.Println("- start enroll consumer")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument email of enrollment must be a non-empty string")
	}

	EmailId := args[0]	
	Name := args[1]

// ====================== Enroll consumer ===========================

	var ShipmentRecord1 []string
	ShipmentRecord1 = append(ShipmentRecord1,"") 
	ShipmentRecord := ShipmentRecord1


	fmt.Println("Name: ", Name ,"EmailId: ",EmailId,"ShipmentRecord: ",ShipmentRecord)

	consumer := &consumer{Name,EmailId,ShipmentRecord}


// ==================== Append  Consumer list  ================		
	consumerStruct1 := consumerStruct{Name,EmailId}
		
	ConsumerList.ConsumerStruct = append(ConsumerList.ConsumerStruct,consumerStruct1) 


	fmt.Println("consumer Unit is ", consumer)
	consumerJSONasBytes, err := json.Marshal(consumer)

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("After marshal ship is ", consumerJSONasBytes)

	// === Save consumer to state ===
	err = stub.PutState(EmailId, consumerJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}


	// === Set Event consumer Unit  ===
	
	err = stub.SetEvent("sampleEvent", []byte(EmailId))
	if err != nil {
		fmt.Println("SetEvent Error", err)
	}
	// ==== consumer Unit saved. Return success ====
	fmt.Println("- end enroll consumer Unit")
	return shim.Success(nil)
}

// ===============================================
//           Reading consumer list
// ===============================================
func (t *SimpleChaincode) readConsumerList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	ConsumerListAsbytes, err:= json.Marshal(ConsumerList)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("After marshal MilkFarmListAsbytes is ", ConsumerListAsbytes)
	
	return shim.Success(ConsumerListAsbytes)

}

// ================================================================
//           Reading Consumer
// ================================================================

func (t *SimpleChaincode) readConsumer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var consumerEmailId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting Milk-Farm EmailId to query")
	}

	consumerEmailId = args[0]
	valAsbytes, err := stub.GetState(consumerEmailId) //get the milk-Farm status from chaincode state.
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + consumerEmailId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"consumer does not exist: " + consumerEmailId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}


// =====================================================================
//               Initialise Batch 
// =====================================================================

func (t *SimpleChaincode) addBatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// ==== Input sanitation ====
	fmt.Println("- start add batch")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument batch Id must be a non-empty string")
	}

	BatchId := args[0]	
	Owner := args[1]
	Category := args[2]
	Manufacturer := args[3]
	Name := args[4]
	Ph := args[5]
	Apperance := args[6]
	MilkType := args[7]

	fmt.Println("batchId: ", BatchId ,"owner: ",Owner,"category: ",Category,"manufacturer: ",Manufacturer,"name: ",Name,"ph: ",Ph,"apperance: ",Apperance,"milkType :",MilkType)

//	========= Append Milk-Farm Shipments Records ===============	

	milkFarmAsBytes, _ := stub.GetState(Manufacturer)
	MilkFarm := milkFarm{}

	json.Unmarshal(milkFarmAsBytes, &MilkFarm)
	MilkFarm.ShipmentRecord = append(MilkFarm.ShipmentRecord,BatchId)

// ============ Create Batch  ====================================	
	batch := &batch{BatchId,Owner,Category,Manufacturer,Name,Ph,Apperance,MilkType}

	fmt.Println("batch is ", batch)
	batchJSONasBytes, err := json.Marshal(batch)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("After marshal ship is ", batchJSONasBytes)

	milkFarmAsBytes, err = json.Marshal(MilkFarm)
	if err != nil {
		return shim.Error(err.Error())
	}
	// === Save Milk-farm to state ===

	err = stub.PutState(Manufacturer, milkFarmAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	// === Save batch to state ===
	err = stub.PutState(BatchId, batchJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Set Event batch Unit  ===
	
	err = stub.SetEvent("sampleEvent", []byte(BatchId))
	if err != nil {
		fmt.Println("SetEvent Error", err)
	}
	// ==== batch saved. Return success ====
	fmt.Println("- end enroll batch Unit")
	return shim.Success(nil)
}

// ================================================================
//           Reading Batch
// ================================================================

func (t *SimpleChaincode) readBatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var batchId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting Milk-Farm EmailId to query")
	}

	batchId = args[0]
	valAsbytes, err := stub.GetState(batchId) //get the milk-Farm status from chaincode state.
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + batchId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"batch does not exist: " + batchId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// =====================================================================
//               Initialise Product 
// =====================================================================

func (t *SimpleChaincode) createProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// ==== Input sanitation ====
	fmt.Println("- start enroll consumer")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument email of enrollment must be a non-empty string")
	}

	ProductId  := args[0]	
	Name := args[1]
	Owner := args[2]
	Manufacturer := args[3]
	Cost := args[4]
	MilkCategory := args[5]
	PackageCategory := args[6]
	Batch := args[7]
	
	var Logs1 []string
	Logs1 = append(Logs1,"Product is Created")
	Logs :=Logs1 

	var Time1 []time.Time
	Time1 = append(Time1,time.Now())
	Time := Time1

	fmt.Println("ProductId: ", ProductId ,"Name: ",Name,"owner: ",Owner,"manufacturer: ",Manufacturer, "cost :",Cost,"milkCategory: ",MilkCategory,"packageCategory: ",PackageCategory,"batch :",Batch,"Logs: ",Logs,"Time : ",Time)

//	========= Append Packaging Unit Shipments ===============	

	packagingUnitAsBytes, _ := stub.GetState(Manufacturer)

	PackagingUnit := packagingUnit{}

	json.Unmarshal(packagingUnitAsBytes,&PackagingUnit)

	PackagingUnit.ShipmentRecord = append(PackagingUnit.ShipmentRecord,ProductId)

	packagingUnitAsBytes, err = json.Marshal(PackagingUnit)
	if err != nil {
		return shim.Error(err.Error())
	}
	// === Save Packaging Unit to state ===

	err = stub.PutState(Manufacturer, packagingUnitAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

// ============ Create Batch  ====================================	
	consumer := &product{ProductId, Name,Owner,Manufacturer,Cost,MilkCategory,PackageCategory,Batch,Logs,Time}

	fmt.Println("consumer Unit is ", consumer)
	consumerJSONasBytes, err := json.Marshal(consumer)

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("After marshal ship is ", consumerJSONasBytes)

	// === Save consumer to state ===
	err = stub.PutState(ProductId, consumerJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Set Event consumer Unit  ===
	
	err = stub.SetEvent("sampleEvent", []byte(ProductId))
	if err != nil {
		fmt.Println("SetEvent Error", err)
	}
	// ==== consumer Unit saved. Return success ====
	fmt.Println("- end enroll consumer Unit")
	return shim.Success(nil)
}

// ================================================================
//           Reading Product
// ================================================================

func (t *SimpleChaincode) readProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var productId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting Milk-Farm EmailId to query")
	}

	productId = args[0]
	valAsbytes, err := stub.GetState(productId) //get the milk-Farm status from chaincode state.
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + productId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"product does not exist: " + productId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// =====================================================================
//               Change Product Ownership 
// =====================================================================

func (t *SimpleChaincode) changeOwnershipProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var i int

	ProductId := args[0]

	productAsBytes, _ := stub.GetState(args[0])
	Product := product{}

	json.Unmarshal(productAsBytes, &Product)

	NewOwner:= args[1]

//	========= Append Next Stack Holder Shipments Records ===============

	valAsBytes, _ := stub.GetState(NewOwner)
	valUnit := packagingUnit{}

	json.Unmarshal(valAsBytes,&valUnit)
	valUnit.ShipmentRecord = append(valUnit.ShipmentRecord,ProductId)

// =========== Delete the Shipment from the existing Owner == 

	currentOwner := Product.Owner

	fmt.Println("BatchId is ",ProductId,"currentOwner is ", currentOwner)

	valAsBytes, _ = stub.GetState(currentOwner)
	MilkFarm := milkFarm{}

	json.Unmarshal(valAsBytes, &MilkFarm)
	ShipmentRecord := MilkFarm.ShipmentRecord

	fmt.Println("Shipment Records: ",ShipmentRecord)

	for k, v := range ShipmentRecord {
			if ProductId == v {
				i = k
		}
	}

	fmt.Println("index: ",i)

	newShipmentRecord := append(ShipmentRecord[:i], ShipmentRecord[i+1:]...)

	MilkFarm.ShipmentRecord = newShipmentRecord

	valAsBytes, err = json.Marshal(MilkFarm)
	if err != nil {
	return shim.Error(err.Error())
	}

	err = stub.PutState(currentOwner, valAsBytes)
	if err != nil {
	return shim.Error(err.Error())
	}
  
  // =========== Change the owner =================================

	fmt.Println("ProductId is ",ProductId,"NewOwner is ", NewOwner)

	Logs1 := "Product is moved from " + currentOwner + " to " + NewOwner 
	Product.Logs = append(Product.Logs,Logs1)
	fmt.Println("Product Logs: ",Product.Logs)

	Product.Timestamp = append(Product.Timestamp,time.Now())
	fmt.Println("Product Time Logs: ",Product.Timestamp)

	Product.Owner = NewOwner
	
	fmt.Println("Updated Product is: ",Product.Owner)

	valAsBytes, err = json.Marshal(valUnit)
	if err != nil {
		return shim.Error(err.Error())
	}
	// === Save valAsBytes to state ===

	err = stub.PutState(NewOwner, valAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save product to state ===
	productAsBytes, _ = json.Marshal(Product)
	stub.PutState(ProductId, productAsBytes)

	return shim.Success(nil)
}

// =====================================================================
//               Change Batch Ownership 
// =====================================================================

func (t *SimpleChaincode) changeOwnershipBatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var i int

	BatchId := args[0]

	batchAsBytes, _ := stub.GetState(args[0])
	Batch := batch{}

	json.Unmarshal(batchAsBytes, &Batch)

	NewOwner:= args[1]

//	========= Append Packaging Unit Shipments ===============

	packagingUnitAsBytes, _ := stub.GetState(NewOwner)
	PackagingUnit := packagingUnit{}
   

	json.Unmarshal(packagingUnitAsBytes,&PackagingUnit)
	PackagingUnit.ShipmentRecord = append(PackagingUnit.ShipmentRecord,BatchId)

// =========== Delete the Shipment from the existing Owner == 
	currentOwner := Batch.Owner

	fmt.Println("BatchId is ",BatchId,"currentOwner is ", currentOwner)

	milkFarmAsBytes, _ := stub.GetState(currentOwner)
	MilkFarm := milkFarm{}

	json.Unmarshal(milkFarmAsBytes, &MilkFarm)
	ShipmentRecord := MilkFarm.ShipmentRecord
	
	fmt.Println("Shipment Records: ",ShipmentRecord)
	
	for k, v := range ShipmentRecord {
			if BatchId == v {
		  i = k
	  }
  }

  fmt.Println("index: ",i)

  newShipmentRecord := append(ShipmentRecord[:i], ShipmentRecord[i+1:]...)

  MilkFarm.ShipmentRecord = newShipmentRecord

  milkFarmAsBytes, err = json.Marshal(MilkFarm)
  if err != nil {
	  return shim.Error(err.Error())
  }

  err = stub.PutState(currentOwner, milkFarmAsBytes)
  if err != nil {
	  return shim.Error(err.Error())
  }
  
  // =========== Change the owner =================================

	fmt.Println("BatchId is ",BatchId,"NewOwner is ", NewOwner)

	Batch.Owner = NewOwner

	fmt.Println("Updated Batch is: ",Batch)

	packagingUnitAsBytes, err = json.Marshal(PackagingUnit)
	if err != nil {
		return shim.Error(err.Error())
	}
	// === Save Packaging Unit to state ===

	err = stub.PutState(NewOwner, packagingUnitAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	// === Save batch to state ===

	batchAsBytes, _ = json.Marshal(Batch)
	stub.PutState(BatchId, batchAsBytes)

	return shim.Success(nil)	
}

// =====================================================================
//               Get History Ownership 
// =====================================================================

func (t *SimpleChaincode) getHistoryForAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := args[0]

	fmt.Printf("- start getHistoryForAsset: %s\n",key)

	resultsIterator, err := stub.GetHistoryForKey(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the Asset
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
		//as-is (as the Value itself a JSON Asset)
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

	fmt.Printf("- getHistoryForAsset returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

	