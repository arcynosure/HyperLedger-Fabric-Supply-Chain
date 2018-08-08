package blockchain

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
)

// Creating a Vegetable Asset  - Adding a new record to the ledger

func (setup *FabricSetup) CreateVeg(key, value string) (string, error) {

	// Prepare arguments

	var args []string
	args = append(args, "invoke")
	args = append(args, "create")
	args = append(args, key)
	args = append(args, value)

	eventID := "eventCreateVeg"

	// Add data that will be visible in the proposal, like a description of the invoke request
	
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in Create Vegetable invoke")

	// Register a notification handler on the client

	notifier := make(chan *chclient.CCEvent)
	rce, err := setup.client.RegisterChaincodeEvent(notifier, setup.ChainCodeID, eventID)
	if err != nil {
		return "", fmt.Errorf("failed to register chaincode event: %v", err)
	}

	// Create a request (proposal) and send it

	response, err := setup.client.Execute(chclient.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}, TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to create Vegetable: %v", err)
	}

	// Wait for the result of the submission

	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %s\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}

	// Unregister the notification handler previously created on the client
	err = setup.client.UnregisterChaincodeEvent(rce)

	return response.TransactionID.ID, nil
}

// Changing Ownership of a Vegetable - Updating a record

func (setup *FabricSetup) ChangeVegOwner(key, value string) (string, error) {

	// Prepare arguments

	var args []string
	args = append(args, "invoke")
	args = append(args, "changeOwner")
	args = append(args, key)
	args = append(args, value)

	eventID := "eventChangeVegOwner"

	// Add data that will be visible in the proposal, like a description of the invoke request
	
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in Create Veg invoke")

	// Register a notification handler on the client

	notifier := make(chan *chclient.CCEvent)
	rce, err := setup.client.RegisterChaincodeEvent(notifier, setup.ChainCodeID, eventID)
	if err != nil {
		return "", fmt.Errorf("failed to register chaincode evet: %v", err)
	}

	// Create a request (proposal) and send it
	response, err := setup.client.Execute(chclient.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3])}, TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to create Veg: %v", err)
	}

	// Wait for the result of the submission
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %s\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}

	// Unregister the notification handler previously created on the client
	err = setup.client.UnregisterChaincodeEvent(rce)

	return response.TransactionID.ID, nil
}

// Updating Records of a Vegetable - Updating all fields

func (setup *FabricSetup) UpdateVegRecord(key, value string) (string, error) {

	// Prepare arguments

	var args []string
	args = append(args, "invoke")
	args = append(args, "updateRecord")
	args = append(args, key)
	args = append(args, value)

	eventID := "eventUpdateRecords"

	// Add data that will be visible in the proposal, like a description of the invoke request
	
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in Update Veg Records invoke")

	// Register a notification handler on the client

	notifier := make(chan *chclient.CCEvent)
	rce, err := setup.client.RegisterChaincodeEvent(notifier, setup.ChainCodeID, eventID)
	if err != nil {
		return "", fmt.Errorf("failed to register chaincode evet: %v", err)
	}

	// Create a request (proposal) and send it

	response, err := setup.client.Execute(chclient.Request{ChaincodeID: setup.ChainCodeID, Fcn: args[0], Args: [][]byte{[]byte(args[0]), []byte(args[1]), []byte(args[2]), []byte(args[3])}, TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to create Veg: %v", err)
	}

	// Wait for the result of the submission
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %s\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}

	// Unregister the notification handler previously created on the client
	err = setup.client.UnregisterChaincodeEvent(rce)

	return response.TransactionID.ID, nil
}
