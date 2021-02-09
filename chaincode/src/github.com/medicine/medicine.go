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


type manufacturer struct {
	Name     	      		 string          			`json:"Name"`
	EmailId           		 string          			`json:"EmailId"`
	ShipmentRecord		   []string                     `json:"ShipmentRecord"`
}

type wholesaler struct {
	Name     	    string        `json:"Name"`
	EmailId         string        `json:"EmailId"`
	ShipmentRecord []string       `json:"ShipmentRecord"`
}

type pharmacist struct {
	Name     	    string        `json:"Name"`
	EmailId         string        `json:"EmailId"`
	ShipmentRecord []string       `json:"ShipmentRecord"`
}

type consumer struct {
	Name     	    string        `json:"Name"`
	EmailId         string        `json:"EmailId"`
	ShipmentRecord []string       `json:"ShipmentRecord"`
}


type product struct {
	ProductId      	    string        `json:"ProductId"`
	Name          		string        `json:"Name"`
	Owner			    string		  `json:"Owner"`
	Manufacturer  	    string 		  `json:"Manufacturer"`
	MedicineCategory	string		  `json:"MedicineCategory"` 
	Batch 	            string		  `json:"Batch"`
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
	if function == "enrollManufacturerUnit" { // create a new Milk-Farm with initialisation
		return t.enrollManufacturerUnit(stub, args)
	} else if function == "readManufacturerUnit" { // read Milk-Farm 
		return t.readManufacturerUnit(stub, args)
	} else if function == "enrollWholesalerUnit" { // read Milk-Farm 
		return t.enrollWholesalerUnit(stub, args)
	} else if function == "readWholesalerUnit" { // read Milk-Farm 
		return t.readWholesalerUnit(stub, args)
	} else if function == "enrollPharmacistUnit" { // read Milk-Farm 
		return t.enrollPharmacistUnit(stub, args)
	} else if function == "readPharmacistUnit" { // read Milk-Farm 
		return t.readPharmacistUnit(stub, args)
	} else if function == "enrollConsumer" { // read Milk-Farm 
		return t.enrollConsumer(stub, args)
	} else if function == "readConsumer" { // read Milk-Farm 
		return t.readConsumer(stub, args)
	} else if function == "createProduct" { // read Milk-Farm 
		return t.createProduct(stub, args)
	} else if function == "readProduct" { // read Milk-Farm 
		return t.readProduct(stub, args)
	} else if function == "changeOwnershipProduct" { // read Milk-Farm 
		return t.changeOwnershipProduct(stub, args)
	} else if function == "getShipmentRecord" { // read Milk-Farm 
		return t.getShipmentRecord(stub, args)
	} else if function == "getHistoryForAsset" { // read Milk-Farm 
		return t.getHistoryForAsset(stub, args)
	}     	
	
	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
//               Initialise Milk-Farm
// ============================================================

func (t *SimpleChaincode) enrollManufacturerUnit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// ==== Input sanitation ====
	fmt.Println("- start enroll Manufacturer")
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

	Manufacturer := &manufacturer{Name, EmailId,ShipmentRecord}

	fmt.Println("Manufacturer is ", Manufacturer)	
	ManufacturerJSONasBytes, err := json.Marshal(Manufacturer)

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("After marshal ship is ", ManufacturerJSONasBytes)

	// === Save Manufacturer to state ===
	err = stub.PutState(EmailId, ManufacturerJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}


	// === Set Event Manufacturer  ===
	
	err = stub.SetEvent("sampleEvent", []byte(EmailId))
	if err != nil {
		fmt.Println("SetEvent Error", err)
	}
	// ==== Manufacturer saved. Return success ====
	fmt.Println("- end enroll Manufacturer")
	return shim.Success(nil)
}


// ===============================================
//           Reading Milk-Farm
// ===============================================

func (t *SimpleChaincode) readManufacturerUnit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var ManufacturerEmailId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting Milk-Farm EmailId to query")
	}

	ManufacturerEmailId = args[0]
	valAsbytes, err := stub.GetState(ManufacturerEmailId) //get the milk-Farm status from chaincode state.
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + ManufacturerEmailId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Shipment does not exist: " + ManufacturerEmailId + "\"}"
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

	 valUnit:= manufacturer{}

	json.Unmarshal(valAsbytes, &valUnit)

	fmt.Println("json response: ",valAsbytes)
	fmt.Println("Milk-Farm unmarshaled: ",valUnit)
	
	ShipmentRecordAsbytes,_ := json.Marshal(valUnit.ShipmentRecord)

	return shim.Success(ShipmentRecordAsbytes)
}

// =====================================================================
//               Initialise Wholesaler Unit
// =====================================================================

func (t *SimpleChaincode) enrollWholesalerUnit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	wholesaler := &wholesaler{Name,EmailId,ShipmentRecord}

	fmt.Println("wholesaler is ", wholesaler)
	wholesalerJSONasBytes, err := json.Marshal(wholesaler)

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("After marshal ship is ", wholesalerJSONasBytes)

	// === Save wholesaler to state ===
	err = stub.PutState(EmailId, wholesalerJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Set Event wholesaler  ===
	
	err = stub.SetEvent("sampleEvent", []byte(EmailId))
	if err != nil {
		fmt.Println("SetEvent Error", err)
	}
	// ==== wholesaler saved. Return success ====
	fmt.Println("- end enroll wholesaler")
	return shim.Success(nil)
}

// ================================================================
//           Reading wholesaler
// ================================================================

func (t *SimpleChaincode) readWholesalerUnit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var wholesalerEmailId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting Milk-Farm EmailId to query")
	}

	wholesalerEmailId = args[0]
	valAsbytes, err := stub.GetState(wholesalerEmailId) //get the milk-Farm status from chaincode state.
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + wholesalerEmailId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Shipment does not exist: " + wholesalerEmailId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// =====================================================================
//               Initialise pharmacist Unit
// =====================================================================

func (t *SimpleChaincode) enrollPharmacistUnit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// ==== Input sanitation ====
	fmt.Println("- start enroll pharmacist-Unit")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument email of enrollment must be a non-empty string")
	}
	
	EmailId := args[0]	
	Name := args[1]

	var ShipmentRecord1 []string
	ShipmentRecord1 = append(ShipmentRecord1,"") 
	ShipmentRecord := ShipmentRecord1


	fmt.Println("Name: ", Name ,"EmailId: ",EmailId,"ShipmentRecord: ",ShipmentRecord)

	pharmacistUnit := &pharmacist{Name,EmailId,ShipmentRecord}

	fmt.Println("pharmacist Unit is ", pharmacistUnit)
	pharmacistUnitJSONasBytes, err := json.Marshal(pharmacistUnit)

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("After marshal ship is ", pharmacistUnitJSONasBytes)

	// === Save pharmacist Unit to state ===
	err = stub.PutState(EmailId, pharmacistUnitJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Set Event pharmacist Unit  ===
	
	err = stub.SetEvent("sampleEvent", []byte(EmailId))
	if err != nil {
		fmt.Println("SetEvent Error", err)
	}
	// ==== pharmacist Unit saved. Return success ====
	fmt.Println("- end enroll pharmacist Unit")
	return shim.Success(nil)
}

// ================================================================
//           Reading pharmacist Unit
// ================================================================

func (t *SimpleChaincode) readPharmacistUnit(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var pharmacistUnitEmailId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting Milk-Farm EmailId to query")
	}

	pharmacistUnitEmailId = args[0]
	valAsbytes, err := stub.GetState(pharmacistUnitEmailId) //get the milk-Farm status from chaincode state.
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + pharmacistUnitEmailId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Shipment does not exist: " + pharmacistUnitEmailId + "\"}"
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

	var ShipmentRecord1 []string
	ShipmentRecord1 = append(ShipmentRecord1,"") 
	ShipmentRecord := ShipmentRecord1


	fmt.Println("Name: ", Name ,"EmailId: ",EmailId,"ShipmentRecord: ",ShipmentRecord)

	consumer := &consumer{Name,EmailId,ShipmentRecord}

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

// ================================================================
//           Reading consumer
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
	MedicineCategory := args[4]
	Batch := args[5]
	
	fmt.Println("ProductId: ", ProductId ,"Name: ",Name,"owner: ",Owner,"manufacturer: ",Manufacturer,"MedicineCategory: ",MedicineCategory,"batch :",Batch)

//	========= Append Packaging Unit Shipments ===============	

	manufacturerAsBytes, _ := stub.GetState(Manufacturer)

	manufacturer := manufacturer{}

	json.Unmarshal(manufacturerAsBytes,&manufacturer)

	manufacturer.ShipmentRecord = append(manufacturer.ShipmentRecord,ProductId)

	manufacturerAsBytes, err = json.Marshal(manufacturer)
	if err != nil {
		return shim.Error(err.Error())
	}
	// === Save Packaging Unit to state ===

	err = stub.PutState(Manufacturer, manufacturerAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

// ============ Create Batch  ====================================	
	product := &product{ProductId, Name,Owner,Manufacturer,MedicineCategory,Batch}

	fmt.Println("product Unit is ", product)
	productJSONasBytes, err := json.Marshal(product)

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("After marshal ship is ", productJSONasBytes)

	// === Save product to state ===
	err = stub.PutState(ProductId, productJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Set Event product Unit  ===
	
	err = stub.SetEvent("sampleEvent", []byte(ProductId))
	if err != nil {
		fmt.Println("SetEvent Error", err)
	}
	// ==== product Unit saved. Return success ====
	fmt.Println("- end enroll product Unit")
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
	valUnit := manufacturer{}

	json.Unmarshal(valAsBytes,&valUnit)
	valUnit.ShipmentRecord = append(valUnit.ShipmentRecord,ProductId)

// =========== Delete the Shipment from the existing Owner == 

	currentOwner := Product.Owner

	fmt.Println("BatchId is ",ProductId,"currentOwner is ", currentOwner)

	valAsBytes, _ = stub.GetState(currentOwner)
	manufacturer := manufacturer{}

	json.Unmarshal(valAsBytes, &manufacturer)
	ShipmentRecord := manufacturer.ShipmentRecord

	fmt.Println("Shipment Records: ",ShipmentRecord)

	for k, v := range ShipmentRecord {
			if ProductId == v {
				i = k
		}
	}

	fmt.Println("index: ",i)

	newShipmentRecord := append(ShipmentRecord[:i], ShipmentRecord[i+1:]...)

	manufacturer.ShipmentRecord = newShipmentRecord

	valAsBytes, err = json.Marshal(manufacturer)
	if err != nil {
	return shim.Error(err.Error())
	}

	err = stub.PutState(currentOwner, valAsBytes)
	if err != nil {
	return shim.Error(err.Error())
	}
  
  // =========== Change the owner =================================

	fmt.Println("ProductId is ",ProductId,"NewOwner is ", NewOwner)

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

	