package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"encoding/json"
	"strconv"
	"strings"
	//"../mylib"
	"github.com/hyperledger/fabric/examples/chaincode/go/mylib"
)

func main() {
	err := shim.Start(new(ChainCode))
	if err != nil {
		fmt.Printf("Error starting patentchain chaincode: %s", err)
	}
}

//=================================================================================================================================
//   Init : Trigger at the time of chaincode deployment
//=================================================================================================================================
func (t *ChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	//Generating first application reference number
	applicationReferenceNumber := fmt.Sprintf("00/0000/" + fmt.Sprintf("%06d", 0))
	err := stub.PutState("ApplicationRefNumber", []byte(applicationReferenceNumber))
	if err != nil {
		return shim.Error(err.Error())
	}

	//Saving PO data for admin
	var poData adminData
	po := args[0]
   	//Un-marshalling location json string to native go structure
	err = json.Unmarshal([]byte(po), &poData)
	if err != nil {
		return shim.Error(mylib.UnmarshalErrorMessage)
	}
    for in:= range poData.POInfo{
        //poData.POInfo[in].Code = poCodeLookUp[string(poData.POInfo[in].Code)]
        //poData.POInfo[in].CountryCode = countryCodeLookUp[string(poData.POInfo[in].CountryCode)]
		//Marshalling final poData
        toPoInfo, err := json.Marshal(poData.POInfo[in])
    	if err != nil {
        return shim.Error(mylib.MarshalErrorMessage)
    	}

		//saving poData
    	err = stub.PutState(poData.POInfo[in].POName, []byte(toPoInfo))
		if err != nil {
		return shim.Error(err.Error())
		}
    }

    //Saving roles
	val, err := json.Marshal(poData.Role)
	if(err != nil){
		shim.Error(mylib.MarshalErrorMessage)
	}
	err = stub.PutState("roles", []byte(val))
	if err != nil {
		return shim.Error(err.Error())
	}

	//saving administrator data


    //dataString='{"po_info": [{"po_name": "IP5- JPO","po_address": "123, Tokyo","po_code": "IP5","po_country_code": "JP"},{"po_name": "IP5- USPTO","po_address": "5667, Seattle","po_code": "EPO","po_country_code": "US"}]}';
        
	return shim.Success(nil)
}

//=================================================================================================================================
//   Invoke : Perform transactions
//=================================================================================================================================
func (t *ChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if len(args)<2{
		shim.Error(mylib.ArgumentErrorMessage)
	}
	if args[0]=="generateApplicationRefNumber" {
		return t.generateApplicationRefNumber(stub, args)
	}else if function=="query"{
		switch args[0] {
		case "getPOInfoByKey":
			return t.getPOInfoByKey(stub, args[1])
		case "getRoles":
			return t.getRoles(stub)

		default:
			return shim.Error("Invalid function")
		}
	}else{
		return shim.Error("Unknown function call")
	}

	return shim.Success([]byte("Error in chaincode calling"))
}

//=================================================================================================================================
//   getRoles : Return roles
//==================================================================================================================================
func (t *ChainCode) getRoles(stub shim.ChaincodeStubInterface) pb.Response  {
	rolesAsBytes, err := stub.GetState("roles")
	if(err != nil){
		shim.Error(mylib.GetStateErrorMessage)
	}
	return shim.Success(rolesAsBytes)
}


//=================================================================================================================================
//   getPOInfoByKey : Return po info by key
//==================================================================================================================================
func (t *ChainCode) getPOInfoByKey(stub shim.ChaincodeStubInterface, poname string) pb.Response  {
	poInfoAsBytes, err := stub.GetState(poname)
	if(err != nil){
		shim.Error(mylib.GetStateErrorMessage)
	}
	return shim.Success(poInfoAsBytes)
}



//=================================================================================================================================
//   generateApplicationRefNumber : Return Application Reference Number
//==================================================================================================================================
func (t *ChainCode) generateApplicationRefNumber(stub shim.ChaincodeStubInterface, args []string) pb.Response  {
	sequenceNumberAsBytes, err := stub.GetState("ApplicationRefNumber")
	if err != nil {
		shim.Error(mylib.GetStateErrorMessage)
	}

	applicationReferenceNumber:=string(sequenceNumberAsBytes[:])
	sequenceAsInteger :=strings.Split(applicationReferenceNumber, "/")[2]
	sequenceNumber, err :=  strconv.ParseInt(sequenceAsInteger,10,32)
	if err != nil {
		return shim.Error(mylib.ConversionErrorMessage)
	}

	countryCode:=args[1];

	/*poname := args[1]
	stateAsBytes , err := stub.GetState(poname)
	if err != nil {
		shim.Error(mylib.GetStateErrorMessage)
	}

	//Un-marshalling json string to native go structure
	err = json.Unmarshal(stateAsBytes, &poInfo)
	if err != nil {
		return shim.Error(mylib.UnmarshalErrorMessage)
	}*/

	year := args[2]
	applicationReferenceNumber = fmt.Sprintf(string(countryCode) + "/" + year + "/" + fmt.Sprintf("%06d", sequenceNumber + 1))
	err = stub.PutState("ApplicationRefNumber", []byte(applicationReferenceNumber))
	if err != nil {
		return shim.Error(mylib.PutErrorMessage)
	}

	patentInfoAsBytes := []byte(applicationReferenceNumber)
	return shim.Success(patentInfoAsBytes)
}