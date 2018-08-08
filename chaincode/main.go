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


type HeroesServiceChaincode struct {
}

type Vegetable struct {
	Name   string `json:"name"`
	Id  string `json:"id"`
	Quality string `json:"quality"`
	Owner  string `json:"owner"`
}

func (t *HeroesServiceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### SupplyChain Init ###########")

	function, _ := stub.GetFunctionAndParameters()

	if function != "init" {
		return shim.Error("Unknown function call")
	}

	return shim.Success(nil)
}

// Invoke

func (t *HeroesServiceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### VegOwner Invoke ###########")

	// Get the function and arguments from the request
	function, args := stub.GetFunctionAndParameters()

	// Check whether it is an invoke request
	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// In order to manage multiple type of request, we will check the first argument.
	// Here we have one possible argument: query (every query request will read in the ledger without modification)
	// In order to manage multiple type of request, we will check the first argument.
	// Here we have one possible argument: query (every query request will read in the ledger without modification)
	if args[0] == "query" {
		return t.query(stub, args)
	}

	// The update argument will manage all update in the ledger
	if args[0] == "invoke" {
		return t.invoke(stub, args)
	}

	// Querying Single Record by Passing Veg ID => Key as parameter
	if args[0] == "queryone" {
		return t.queryone(stub, args)
	}

	// Getting History of a Record by passing Veg ID => Key as parameter.
	if args[0] == "gethistory" {
		return t.gethistory(stub, args)
	}

	// Adding a new transaction to the ledger
	if args[0] == "create" {
		return t.createVeg(stub, args)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown action, check the first argument")
}

// query
// Every readonly functions in the ledger will be here
func (t *HeroesServiceChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### VegetableOwnership query ###########")

	// Check whether the number of arguments is sufficient
	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Like the Invoke function, we manage multiple type of query requests with the second argument.
	
	if args[1] == "all" {

		// GetState by passing lower and upper limits
		resultsIterator, err := stub.GetStateByRange("", "")
		if err != nil {
			return shim.Error(err.Error())
		}
		defer resultsIterator.Close()

		// buffer is a JSON array containing QueryResults
		var buffer bytes.Buffer
		buffer.WriteString("[")

		bArrayMemberAlreadyWritten := false
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return shim.Error(err.Error())
			}
			// Add a comma before array members, suppress it for the first array member
			if bArrayMemberAlreadyWritten == true {
				buffer.WriteString(",")
			}
			buffer.WriteString("{\"Key\":")
			buffer.WriteString("\"")
			buffer.WriteString(queryResponse.Key)
			buffer.WriteString("\"")

			buffer.WriteString(", \"Record\":")
			// Record is a JSON object, so we write as-is
			buffer.WriteString(string(queryResponse.Value))
			buffer.WriteString("}")
			bArrayMemberAlreadyWritten = true
		}
		buffer.WriteString("]")

		fmt.Printf("- queryAllVeg:\n%s\n", buffer.String())

		return shim.Success(buffer.Bytes())
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown query action, check the second argument.")
}

// invoke
// Every functions that read and write in the ledger will be here
func (t *HeroesServiceChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("########### VegetableOwnership invoke ###########")

	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Changing Ownership of a Vegetable item by Accepting Key and Value
	if args[1] == "changeOwner" && len(args) == 4 {

		vegAsBytes, _ := stub.GetState(args[2])
		veg := Veg{}

		json.Unmarshal(vegAsBytes, &veg)
		veg.Owner = args[3]

		vegAsBytes, _ = json.Marshal(veg)
		stub.PutState(args[2], vegAsBytes)

		// Notify listeners that an event "eventInvoke" have been executed (check line 19 in the file invoke.go)
		err := stub.SetEvent("eventChangeVegOwner", []byte{})
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)
	}

	/*
		 Updating all fields of record
	*/
	if args[1] == "updateRecord" && len(args) == 4 {
		fmt.Println("Update All")
		var newVeg Veg
		json.Unmarshal([]byte(args[3]), &newVeg)
		var veg = Veg{ Name: newVeg.Name, Id: newVeg.Id, Quality: newVeg.Quality, Owner: newVeg.Owner}
		vegAsBytes, _ := json.Marshal(veg)

		// Updating Record

		stub.PutState(args[2], vegAsBytes)

		// Notify listeners that an event "eventInvoke" have been executed (check line 19 in the file invoke.go)
		
		err := stub.SetEvent("eventUpdateRecords", []byte{})
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)
	}

	// If the arguments given don’t match any function, we return an error

	return shim.Error("Unknown invoke action, check the second argument.")
}

//  Retrieves a single record from the ledger by accepting Key value

func (t *HeroesServiceChaincode) queryone(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// GetState retrieves the data from ledger using the Key

	vegAsBytes, _ := stub.GetState(args[1])

	// Transaction Response

	return shim.Success(vegAsBytes)

}

// Adds a new transaction to the ledger

func (s *HeroesServiceChaincode) createveg(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var newVeg Veg
	json.Unmarshal([]byte(args[2]), &newVeg)
	var veg = Veg{ Name: newVeg.Name, Id: newVeg.Id, Quality: newVeg.Quality, Owner: newVeg.Owner}
	vegAsBytes, _ := json.Marshal(veg)
	stub.PutState(args[1], vegAsBytes)

	// Notify listeners that an event "eventInvoke" have been executed (check line 19 in the file invoke.go)
	err := stub.SetEvent("eventCreateVeg", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Get History of a transaction by passing Key

func (s *HeroesServiceChaincode) gethistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	vegKey := args[1]
	fmt.Printf("##### start History of Record: %s\n", vegKey)

	resultsIterator, err := stub.GetHistoryForKey(vegKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
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
		//as-is (as the Value itself a JSON marble)
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

	fmt.Printf("- getHistoryForVegetable returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func main() {

	// Start the chaincode and make it ready for futures requests

	err := shim.Start(new(HeroesServiceChaincode))
	if err != nil {
		fmt.Printf("Error starting Heroes Service chaincode: %s", err)
	}
}
