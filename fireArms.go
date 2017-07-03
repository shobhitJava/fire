package main

import (
	
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"

)

var logger = shim.NewLogger("FireArms")

//ALL_ELEMENENTS Key to refer the list of application
const ALL_ELEMENENTS = "ALL_APP"




//FireArms Chaincode default interface
type FireArms struct {
}





//Append a new fireArms appid to the master list
func updateMasterRecords(stub shim.ChaincodeStubInterface, appId string) error {
	var recordList []string
	recBytes, _ := stub.GetState(ALL_ELEMENENTS)

	err := json.Unmarshal(recBytes, &recordList)
	if err != nil {
		return errors.New("Failed to unmarshal updateMasterReords ")
	}
	recordList = append(recordList, appId)
	bytesToStore, _ := json.Marshal(recordList)
	logger.Info("After addition" + string(bytesToStore))
	stub.PutState(ALL_ELEMENENTS, bytesToStore)
	return nil
}



// Creating a new fireArm Application
func createApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Info("createfireArm application called")
	var id string
	var data map[string]string
	valAsbytes, err := stub.GetState("id")
	if err != nil {
		jsonResp:= "{\"Error\":\"Failed to get state for id\"}"
		return nil, errors.New(jsonResp)
	}
	fmt.Println((string)(valAsbytes))
	id=(string)(valAsbytes)
	
	fmt.Println("new id is"+id)
	uniqueId,_:=strconv.Atoi(id)
	appId:="AppId"+id
	newid:=uniqueId+1
	stub.PutState("id",[]byte(strconv.Itoa(newid)))
	payload := args[0]
	json.Unmarshal([]byte(payload), &data)
	email:=data["appemail"]	
	fmt.Println("emiail is: "+email)
	stub.PutState(email,[]byte(appId))	
	fmt.Println("new Payload is " + payload)
	
		stub.PutState(appId, []byte(payload))

		updateMasterRecords(stub, appId)
		logger.Info("Created the FireArms")
	
	return nil, nil
}



//Validate a input string as number or not
func validateNumber(str string) float64 {
	if netCharge, err := strconv.ParseFloat(str, 64); err == nil {
		return netCharge
	}
	return float64(-1.0)
}

//Update the existing record with the mofied key value pair
func updateRecord(existingRecord map[string]string, fieldsToUpdate map[string]string) (string, error) {
	for key, value := range fieldsToUpdate {

		existingRecord[key] = value
	}
	outputMapBytes, _ := json.Marshal(existingRecord)
	logger.Info("updateRecord: Final json after update " + string(outputMapBytes))
	return string(outputMapBytes), nil
}


func probe() []byte {
	ts := time.Now().Format(time.UnixDate)
	output := "{\"status\":\"Success\",\"ts\" : \"" + ts + "\" }"
	return []byte(output)
}


// Init initializes the smart contracts
func (t *FireArms) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Info("Init called")
	//Place an empty arry
	stub.PutState(ALL_ELEMENENTS, []byte("[]"))
	stub.PutState("id",[]byte("1"))
	return nil, nil
}

// Invoke entry point
func (t *FireArms) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Info("Invoke called")

	if function == "createApplication" {
		createApplication(stub, args)
	} 

	return nil, nil
}

// Query the rcords form the  smart contracts
func (t *FireArms) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Info("Query called")
	if function == "getAppById" {
		return getAppById(stub, args[0])
	}else if function == "getAppByEmailId" {
		return getAppByEmailId(stub, args[0])
	}else if function == "getAllApp" {
		return getAllApp(stub)
	}   
	
	return nil, nil
}

//Get a single Application
func getAppById(stub shim.ChaincodeStubInterface, args string) ([]byte, error) {
	logger.Info("getAppById called with AppId: " + args )

	var outputRecord map[string]string
	appid := args//AppId
	recBytes, _ := stub.GetState(appid)
	json.Unmarshal(recBytes, &outputRecord)
	outputBytes, _ := json.Marshal(outputRecord)
	logger.Info("Returning records from getAppId " + string(outputBytes))
	return outputBytes, nil
}

//Get a single Application on the basis of email id
func getAppByEmailId(stub shim.ChaincodeStubInterface, args string) ([]byte, error) {
	logger.Info("getAppById called with AppId: " + args )
	
	id,err:=stub.GetState(args)
	if err != nil {
		jsonResp:= "{\"Error\":\"Failed to get state for id\"}"
		return nil, errors.New(jsonResp)
	}
	if (string)(id) == "" {
		jsonResp:= "{\"Error\":\"Failed to get state for this emaili" +args+"\"}"
		return nil, errors.New(jsonResp)
	}
	fmt.Println("The email id is"+(string)(id))
	
	var outputRecord map[string]string
	appid := (string)(id)//AppId
	recBytes, _ := stub.GetState(appid)
	json.Unmarshal(recBytes, &outputRecord)
	outputBytes, _ := json.Marshal(outputRecord)
	logger.Info("Returning records from getAppId " + string(outputBytes))
	return outputBytes, nil
}
//Get all the Application based on the status 
//to do
func getAllApp(stub shim.ChaincodeStubInterface) ([]byte, error) {
	logger.Info("getAllApp called" )
	
	return nil, nil
}

//Main method
func main() {
	logger.SetLevel(shim.LogInfo)
	
	err := shim.Start(new(FireArms))
	if err != nil {
		fmt.Printf("Error starting FireArms: %s", err)
	}
}
