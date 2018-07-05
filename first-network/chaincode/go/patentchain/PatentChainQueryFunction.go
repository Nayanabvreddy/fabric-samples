package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"fmt"
	"github.com/hyperledger/fabric/examples/chaincode/go/mylib"
	"encoding/json"
	"bytes"
	"time"
	"strconv"
)


//=================================================================================================================================
//   getPatentInfoByStatus : Return patent info by status
//==================================================================================================================================
func(t *ChainCode) getPatentDetails(stub shim.ChaincodeStubInterface, args[] string) pb.Response {

if len(args) != 4 {
		logger.Error("No. of arguements not sufficient")
        return shim.Error("No. of arguements not sufficient")
    }

    //role := args[1]

    key := args[2]

    //pub_modulus := args[3]

    patent, err := t.retrievePatent(stub, key);
    if err != nil {
    	logger.Error(err)
        return shim.Error("Error retrieving profile " + err.Error())
    }

/*    if actorLookUp[role] != TEAM_LEAD &&
        actorLookUp[role] != KYC_ANALYST && 
        actorLookUp[role] != COMPLIANCE_OFFICER &&
        actorLookUp[role] != AUDITOR &&
        actorLookUp[role] != ADMIN &&
        actorLookUp[role] != CUSTOMER {
        return shim.Error("not authenticated for this action")
    }

    if actorLookUp[role] == CUSTOMER && profile.Customer_modulus != pub_modulus {
        return shim.Error("not authenticated for this action")
    }*/

    // encrypt aes key before sending
/*	patent.AESkey, err = encrypt_using_public_key(patent.AESkey, pub_modulus)
    if err != nil {
		fmt.Println(err)
        return shim.Error("Error encrypting record")
    }*/

    bytes, err := json.Marshal(patent)

    if err != nil {
    	logger.Error(err)
        return shim.Error("Error converting profile record")
    } 

    return shim.Success(bytes)
}

func(t *ChainCode) getDraftPatent(stub shim.ChaincodeStubInterface, args[] string) pb.Response {

    if len(args) != 3 {
		logger.Error("No. of arguements not sufficient")
        return shim.Error("No. of arguements not sufficient")
    }

    record, err := t.getPatentListByStatus(stub, args, mylib.DRAFT)
    if err != nil {
    	logger.Error("error retreiving records " + err.Error())
        return shim.Error("error retreiving records " + err.Error())
    }

    return shim.Success(record)
  
}

func(t *ChainCode) getSubmittedPatent(stub shim.ChaincodeStubInterface, args[] string) pb.Response {

    if len(args) != 3 {
    	logger.Error("No. of arguements not sufficient")
        return shim.Error("No. of arguements not sufficient")
    }

    record, err := t.getPatentListByStatus(stub, args, mylib.SUBMITTED)
    if err != nil {
    	logger.Error("error retreiving records " + err.Error())
        return shim.Error("error retreiving records " + err.Error())
    }

	return shim.Success(record)
  
}
func(t *ChainCode) getRejectedPatent(stub shim.ChaincodeStubInterface, args[] string) pb.Response {

    if len(args) != 3 {
    	logger.Error("No. of arguements not sufficient")
        return shim.Error("No. of arguements not sufficient")
    }

    record, err := t.getPatentListByStatus(stub, args, mylib.REJECTED)
    if err != nil {
    	logger.Error("error retreiving records " + err.Error())
        return shim.Error("error retreiving records " + err.Error())
    }

    return shim.Success(record)
  
}

func(t *ChainCode) getSearchReportAttachedPatent(stub shim.ChaincodeStubInterface, args[] string) pb.Response {

    if len(args) != 3 {
    	logger.Error("No. of arguements not sufficient")
        return shim.Error("No. of arguements not sufficient")
    }

    record, err := t.getPatentListByStatus(stub, args, mylib.SEARCHREPORTATTACHED)
    if err != nil {
    	logger.Error("error retreiving records " + err.Error())
        return shim.Error("error retreiving records " + err.Error())
    }

    return shim.Success(record)
  
}





//=================================================================================================================================
//   getPatentInfoByKey : Return patent info by key
//==================================================================================================================================
func (t *ChainCode) getPatentInfoByKey(stub shim.ChaincodeStubInterface, key string) pb.Response  {
	patentInfoAsBytes, err := stub.GetState(key)
	if(err != nil){
		logger.Error(mylib.GetStateErrorMessage)
		shim.Error(mylib.GetStateErrorMessage)
	}
	return shim.Success(patentInfoAsBytes)
}

//=================================================================================================================================
//   getPatentsByApplicant : Return all patents by applicant
//==================================================================================================================================
func (t *ChainCode) getPatentsByApplicant(stub shim.ChaincodeStubInterface, applicantEmail string) pb.Response {
   queryString := fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_email\": \"%s\"}}}]}}",  "Draft", applicantEmail)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		logger.Error("Unable to fetch records by applicant")
		return shim.Error("Unable to fetch records by applicant")
	}
	return shim.Success(queryResults)
}

func(t *ChainCode) getPatentSearch(stub shim.ChaincodeStubInterface, args[] string) pb.Response {	
    record, err := t.getSearchPatentList(stub, args)
    if err != nil {
        return shim.Error("error retreiving records " + err.Error())
    }

    return shim.Success(record)
  
} 

func(t *ChainCode) getPatentappSearch(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
	  record, err := t.getSearchPatentappList(stub, args)
	  if err != nil {
		  return shim.Error("error retreiving records " + err.Error())
	  }
  
	  return shim.Success(record)
	
  }

  func(t *ChainCode) getPatentpriSearch(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
      record, err := t.getSearchPatentpriList(stub, args)
	  if err != nil {
		  return shim.Error("error retreiving records " + err.Error())
	  }
  
	  return shim.Success(record)
	
  }

func (t *ChainCode) getPatentSearchcreate(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
	createdAt := args[1]

	queryString := fmt.Sprintf("{\"selector\":{\"created_at\":{\"$gt\":\"%s\"}}}", createdAt)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error("Unable to fetch records by applicant")
	}
	fmt.Println("function getPatentSearch end")
	return shim.Success(queryResults)
}

func (t *ChainCode) getPatentSearchsubmit(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
	submittedTo := args[1]
	fmt.Println("Submitted query for",submittedTo)

	queryString := fmt.Sprintf("{\"selector\":{\"submitted_to\":{\"npo_name\":\"%s\"}}}",  submittedTo)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error("Unable to fetch records by applicant")
	}
	fmt.Println("function getPatentSearch end")
	return shim.Success(queryResults)
}

func (t *ChainCode) getPatentSearchapplicant(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
   
	applicant := args[1]
    queryString := fmt.Sprintf("{\"selector\":{\"Applicants\":{\"$elemMatch\":{\"ext_name\": \"%s\"}}}}",  applicant)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error("Unable to fetch records by applicant")
	}
	fmt.Println("function getPatentSearch end")
	return shim.Success(queryResults)
}

func (t *ChainCode) getPatentSearchtitle(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
 
	patenttitle := args[1]
  
    queryString := fmt.Sprintf("{\"selector\":{\"patent_title\":{\"$eq\": \"%s\"}}}",  patenttitle)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error("Unable to fetch records by patent title")
	}
	fmt.Println("function getPatentSearch end")
	return shim.Success(queryResults)
}


func (t *ChainCode) getHistoryForPatent(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		logger.Error("No. of arguements not sufficient")
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	patentRefNo := args[1]

	resultsIterator, err := stub.GetHistoryForKey(patentRefNo)
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

	return shim.Success(buffer.Bytes())
}

func(t *ChainCode) getNotifications(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
	if len(args) != 3 {
		logger.Error("No. of arguements not sufficient")
		return shim.Error("No. of arguements not sufficient")
	}

	var notifications Notifications
	notificationsAsBytes, err := stub.GetState("notifications")
	if(err != nil){
		shim.Error(mylib.GetStateErrorMessage)
	}

	if(len(notificationsAsBytes) != 0) {
		//Unmarshalling package json string to native go structure
		err = json.Unmarshal(notificationsAsBytes, &notifications)
		if (err != nil) {
			return shim.Error(mylib.UnmarshalErrorMessage)
		}
	}else{
		notifications.Count=0
		notificationBytes, err := json.Marshal(notifications)
		if(err != nil){
			shim.Error(mylib.MarshalErrorMessage)
		}
		notificationsAsBytes=notificationBytes
	}

	return shim.Success(notificationsAsBytes)

}

