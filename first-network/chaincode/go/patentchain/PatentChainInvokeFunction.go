package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	"github.com/hyperledger/fabric/examples/chaincode/go/mylib"
	"fmt"
	"math/rand"
	"strconv"
)
const (
	aes_key = "M0yR_X9lTHOV90j8NCjcsA=="
)

//=================================================================================================================================
//   validateAppRefNumber : Validate application reference number or generate a new one
//=================================================================================================================================
func (t *ChainCode) validateAppRefNumber(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		logger.Error("Incorrect number of arguements")
		return shim.Error("Incorrect number of arguements")
	}
	var ApplicationRefNumber string

	if len(args[1]) == 0 {
		//Calling Configuration chaincode for ApplicationRefNumber
		var chaincodeResponse string
		invokeArgs := ToChaincodeArgs(mylib.ConfigInvoke, "generateApplicationRefNumber", args[2], args[3])
		chaincodeResponse, err := t.invokeConfigChaincode(stub,invokeArgs)
		if err != nil {
			logger.Error("Error saving patent")
			return shim.Error("Error saving patent")
		}
		ApplicationRefNumber = chaincodeResponse
	} else {
		//check if patent already exists or not
		record, err := stub.GetState(args[1])
		if record == nil {
			//Return same Application Reference Number
			ApplicationRefNumber = args[1]
		}else{
			if(err != nil){
				logger.Error(mylib.GetStateErrorMessage)
				shim.Error(mylib.GetStateErrorMessage)
			}
			logger.Error("Application Reference Number already exists")
			return shim.Error("Application Reference Number already exists")
		}
	}

	return shim.Success([]byte(ApplicationRefNumber))
}

//=================================================================================================================================
//   storePackageDetails : Store whole package details into the ledger
//=================================================================================================================================
func (t *ChainCode) storePatentDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response{

    if len(args) != 4 {
        logger.Error("Incorrect number of arguements")
        return shim.Error("Incorrect number of arguements")
    }

    //User 
    var userDetail User
    err := json.Unmarshal([]byte(args[1]), & userDetail)
	if err != nil {
        logger.Error("Error parsing user info")
		return shim.Error("Error parsing user info")
    }


    if userDetail.Role != mylib.APPLICANT {
    	logger.Error("profile not being initiated by a customer")
        return shim.Error("profile not being initiated by a customer")
    } 

    //PatentInfo
    var patentInfo PatentInfo

    err = json.Unmarshal([]byte(args[2]), & patentInfo)
	if err != nil {
        logger.Error(err)
		return shim.Error("Error parsing profile info")
    }

    //Ext User
	var applicantIndex=0;
	for in, _:= range patentInfo.Applicants{
		if patentInfo.Applicants[in].Email == userDetail.Email{
			applicantIndex++
		}
	}
	if applicantIndex == 0 {
		var extUser ExtUser
		extUser.Name =  fmt.Sprintf(userDetail.FirstName+ " " + userDetail.LastName)
		extUser.Email = userDetail.Email
		patentInfo.Applicants = append(patentInfo.Applicants,extUser)
	}

	if(patentInfo.PriorityNumber == "") {
		patentInfo.PriorityNumber	= patentInfo.ApplicationRefNumber
	}

    //Generate AES key
   /* aes_key, err := generate_random_aes_key()
    if err != nil {
        fmt.Printf("patent: Error saving patent: %s", err)
		return shim.Error("Error saving patent")
    }*/

    patentInfo.AESkey = aes_key
    patentInfo.CurrentStatus = mylib.DRAFT
    patentInfo.UpdatedAt=patentInfo.CreatedAt

    // Todo: Call config chaincode for submitted_to. For the moment, passing JSON



    //Action
    var action ActionDetail
    action.Status = patentInfo.CurrentStatus
	action.ActionBy = userDetail
    action.ActionDate = patentInfo.CreatedAt
 	
    patentInfo.ActionDetails = append(patentInfo.ActionDetails,action)

    //update Patent Data
    patentJson := args[3]
    //TODO : check if correct format or not
    var patentData Patent
    err = json.Unmarshal([]byte(patentJson), & patentData);
    if err != nil {
    	logger.Error(err)
        return shim.Error("Error parsing data info")
    }

    for in, _:= range patentData.Documents{
        patentData.Documents[in].Uploadedby = userDetail
    }

    toCrypt_data, err := json.Marshal(patentData)
    if err != nil {
    	logger.Error(err)
        return shim.Error("Error parsing data info")
    }

    encrypted_string, err := encrypt_data_using_aes_key(string(toCrypt_data),patentInfo.AESkey)
    if err != nil {
        logger.Error(err)
        return shim.Error("Error saving patent")
    }

    patentInfo.PatentData = encrypted_string
    _, err = t.saveChanges(stub, patentInfo)
    if err != nil {
        logger.Error(err)
        return shim.Error("Error saving changes")
    }
    /*
    // Add this profile to statusholder REQUESTED
    var basic_info KycProfileBasicInfo

    basic_info.Customer_id = profile.Customer_id
    basic_info.Customer_name = profile.Customer_name
    basic_info.Customer_domain = profile.Customer_domain
    basic_info.Customer_country = profile.Customer_country

    fmt.Println("appending to ledger");
    _, err = t.append_profile_to_status_holder(stub, REQUESTED, basic_info)
    if err != nil {
        fmt.Printf("profile: Error saving changes: %s", err);
        return shim.Error("Error saving changes")
    }*/

    return shim.Success([]byte("Patent Saved Successfully"))
}

//=================================================================================================================================
//   updatePatentDetails : Store whole package details into the ledger
//=================================================================================================================================
func (t *ChainCode) updatePatentDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	if len(args) != 4 {
		logger.Error("Incorrect number of arguements")
		return shim.Error("Incorrect number of arguements")
	}

	//User
	var userDetail User
	err := json.Unmarshal([]byte(args[1]), & userDetail)
	if err != nil {
		logger.Error(err)
		return shim.Error("Error parsing user info")
	}


	if userDetail.Role != mylib.APPLICANT {
		logger.Error("profile not being initiated by a customer")
		return shim.Error("profile not being initiated by a customer")
	}

	//PatentInfo
	var patentInfo PatentInfo

	err = json.Unmarshal([]byte(args[2]), & patentInfo)
	if err != nil {
		logger.Error(err)
		return shim.Error("Error parsing profile info")
	}

	if len(patentInfo.ApplicationRefNumber) <= 0 {
		logger.Error("patent not exists")
		return shim.Error("patent not exists")
	}

	//Calling Configuration chaincode for ApplicationRefNumber
	if(patentInfo.PriorityNumber == "") {
		patentInfo.PriorityNumber	= patentInfo.ApplicationRefNumber
	}

	//Ext User
	var applicantIndex=0;
	for in, _:= range patentInfo.Applicants{
		if patentInfo.Applicants[in].Email == userDetail.Email{
			applicantIndex++
		}
	}
	if applicantIndex == 0 {
		var extUser ExtUser
		extUser.Name =  fmt.Sprintf(userDetail.FirstName+ " " + userDetail.LastName)
		extUser.Email = userDetail.Email
		patentInfo.Applicants = append(patentInfo.Applicants,extUser)
	}



	//check if patent already exists or not
	record, err := t.retrievePatent(stub,patentInfo.ApplicationRefNumber)
	if err != nil {
		logger.Error(err)
		return shim.Error("ApplicationRefNumber not exists")
	}

	//Generate AES key
	/*aes_key, err := generate_random_aes_key()
	if err != nil {
		fmt.Printf("patent: Error saving patent: %s", err)
		return shim.Error("Error saving patent")
	}*/

	patentInfo.AESkey = aes_key
	patentInfo.CurrentStatus = record.CurrentStatus
	patentInfo.CreatedAt=record.CreatedAt
	patentInfo.UpdatedAt=patentInfo.CreatedAt

	// Todo: Call config chaincode for submitted_to. For the moment, passing JSON

	//Action
	var action ActionDetail
	action.Status = patentInfo.CurrentStatus
	action.ActionBy = userDetail
	action.ActionDate = patentInfo.UpdatedAt

	patentInfo.ActionDetails = append(record.ActionDetails,action)

	//update Patent Data
	patentJson := args[3]
	//TODO : check if correct format or not
	var patentData Patent
	err = json.Unmarshal([]byte(patentJson), & patentData);
	if err != nil {
		logger.Error(err)
		return shim.Error("Error parsing data info")
	}

	for in, _:= range patentData.Documents{
		patentData.Documents[in].Uploadedby = userDetail
	}

	//Document versioning
	decrypted_data, err := decrypt_data_using_aes_key(record.PatentData, record.AESkey)
	if err != nil {
		logger.Error(err)
		return shim.Error("Not able to decrpt")
	}

	var originalPatentData Patent

	err = json.Unmarshal([]byte(decrypted_data), & originalPatentData);
	if err != nil {
		logger.Error(err)
		return shim.Error("Not able to decrpt")
	}

	patentData.Documents = append(originalPatentData.Documents, patentData.Documents...)
	//End Document versioning

	toCrypt_data, err := json.Marshal(patentData)
	if err != nil {
		logger.Error(err)
		return shim.Error("Error parsing data info")
	}

	encrypted_string, err := encrypt_data_using_aes_key(string(toCrypt_data),patentInfo.AESkey)
	if err != nil {
		logger.Error(err)
		return shim.Error("Error saving patent")
	}

	patentInfo.PatentData = encrypted_string
	_, err = t.saveChanges(stub, patentInfo)
	if err != nil {
		logger.Error(err)
		return shim.Error("Error saving changes")
	}

	return shim.Success([]byte(patentInfo.ApplicationRefNumber))
}

//=================================================================================================================================
//   updatePatentDetails : Store whole package details into the ledger
//=================================================================================================================================
func (t *ChainCode) updatePatentDocument(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	if len(args) != 5 {
		logger.Error("Incorrect number of arguements")
		return shim.Error("Incorrect number of arguements")
	}
	key := args[3]
	//User
	var userDetail User
	err := json.Unmarshal([]byte(args[1]), & userDetail)
	if err != nil {
		logger.Error(err)
		return shim.Error("Error parsing user info")
	}


	if userDetail.Role != mylib.APPLICANT {
		logger.Error("profile not being initiated by a customer")
		return shim.Error("profile not being initiated by a customer")
	}

	//PatentInfo
	var patentInfo PatentInfo

	//check if patent already exists or not
	patentInfo, err = t.retrievePatent(stub,key)
	if err != nil {
		logger.Error(err)
		return shim.Error("ApplicationRefNumber not exists")
	}

	patentInfo.UpdatedAt = args[4]

	// Todo: Call config chaincode for submitted_to. For the moment, passing JSON
	//Action
	var action ActionDetail
	action.Status = patentInfo.CurrentStatus
	action.ActionBy = userDetail
	action.ActionDate = patentInfo.UpdatedAt

	patentInfo.ActionDetails = append(patentInfo.ActionDetails,action)

	//update Patent Data
	patentJson := args[2]
	//TODO : check if correct format or not
	var patentData Patent
	err = json.Unmarshal([]byte(patentJson), & patentData);
	if err != nil {
		logger.Error(err)
		return shim.Error("Error parsing data info")
	}

	for in, _:= range patentData.Documents{
		patentData.Documents[in].Uploadedby = userDetail
	}

	//Document versioning
	decrypted_data, err := decrypt_data_using_aes_key(patentInfo.PatentData, patentInfo.AESkey)
	if err != nil {
		logger.Error(err)
		return shim.Error("Not able to decrpt")
	}

	var originalPatentData Patent

	err = json.Unmarshal([]byte(decrypted_data), & originalPatentData);
	if err != nil {
		logger.Error(err)
		return shim.Error("Not able to decrpt")
	}

	originalPatentData.Documents = append(patentData.Documents,originalPatentData.Documents...);
	//End Document versioning

	toCrypt_data, err := json.Marshal(originalPatentData)
	if err != nil {
		logger.Error(err)
		return shim.Error("Error parsing data info")
	}

	encrypted_string, err := encrypt_data_using_aes_key(string(toCrypt_data),patentInfo.AESkey)
	if err != nil {
		logger.Error(err)
		return shim.Error("Error saving patent")
	}

	patentInfo.PatentData = encrypted_string
	_, err = t.saveChanges(stub, patentInfo)
	if err != nil {
		logger.Error(err)
		return shim.Error("Error saving changes")
	}

	return shim.Success([]byte(patentInfo.ApplicationRefNumber))
}



//=================================================================================================================================
//   declarePackage : Change status of package and request custom for authentication
//=================================================================================================================================
func (t *ChainCode) submitPatentDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	key := args[1]

	//User 
    var userDetail User
    err := json.Unmarshal([]byte(args[2]), & userDetail)
	if err != nil {
        logger.Error(err)
		return shim.Error("Error parsing user info")
    }

    if userDetail.Role != mylib.APPLICANT {
    	logger.Error("profile not being initiated by a applicant")
        return shim.Error("profile not being initiated by a applicant")
    } 

    //Patent Info
	var patentInfo PatentInfo
	patentInfoAsBytes, err := stub.GetState(key)
	if(err != nil){
		logger.Error(err)
		shim.Error(mylib.GetStateErrorMessage)
	}

	//Unmarshalling package json string to native go structure
	err = json.Unmarshal(patentInfoAsBytes, &patentInfo)
	if (err != nil) {
		return shim.Error(mylib.UnmarshalErrorMessage)
	}
	patentInfo.CurrentStatus = mylib.SUBMITTED
	patentInfo.UpdatedAt = args[3]
	patentInfo.SubmittedAt = patentInfo.UpdatedAt

	//Action
    var action ActionDetail
    action.Status = mylib.SUBMITTED
	action.ActionBy = userDetail
    action.ActionDate = patentInfo.UpdatedAt
 	
    patentInfo.ActionDetails = append(patentInfo.ActionDetails,action)

	//Marshalling final patent info
	val, err := json.Marshal(patentInfo)
	if(err != nil){
		logger.Error(err)
		shim.Error(mylib.MarshalErrorMessage)
	}

	//Updating patent corresponding to its application refrence number
	err = stub.PutState(key, []byte(val))
	if(err != nil){
		logger.Error(err)
		return shim.Error(mylib.PutErrorMessage)
	}

	return shim.Success([]byte("patent is submitted successfully"))
}

//=================================================================================================================================
//   deletePatentInfo : Delete patent info from world state
//=================================================================================================================================
func (t *ChainCode) deletePatentInfo(stub shim.ChaincodeStubInterface, key string) pb.Response {

	var patentInfo PatentInfo
	patentInfoAsBytes, err := stub.GetState(key)
	if err != nil {
		logger.Error(err)
		shim.Error(mylib.GetStateErrorMessage)
	}

	//Unmarshalling package json string to native go structure
	err = json.Unmarshal(patentInfoAsBytes, &patentInfo)
	if err != nil {
		logger.Error(err)
		return shim.Error(mylib.UnmarshalErrorMessage)
	}

	if patentInfo.CurrentStatus==mylib.DRAFT {
		err = stub.DelState(key)
		if err !=nil {
			logger.Error(err)
			return shim.Error(mylib.DeleteStateErrorMessage)
		}
	}else {
		logger.Error("Patent is not in draft state")
		return shim.Error("Patent is not in draft state")
	}

	return shim.Success([]byte("Patent is deleted successfully"))
}

//=================================================================================================================================
//   rejectPatentApplication : Reject patent Application
//=================================================================================================================================

func (t *ChainCode) rejectPatentApplication(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
	key := args[1]
	fmt.Println("The rejectPatentApplication is called")
	//User
	var userDetail User
	err := json.Unmarshal([]byte(args[2]), & userDetail)
	if err != nil {
		logger.Error(err)
		return shim.Error("Error parsing user info")
	}

	if userDetail.Role != mylib.FORMALITYOFFICERS {
		logger.Error("Only Users with Formality Officer Permission can Reject the Patent Application")
		return shim.Error("Only Users with Formality Officer Permission can Reject the Patent Application")
	}

	//Patent Info
	var patentInfo PatentInfo
	patentInfoAsBytes, err := stub.GetState(key)
	if(err != nil){
		logger.Error(err)
		shim.Error(mylib.GetStateErrorMessage)
	}


	//Unmarshalling package json string to native go structure
	err = json.Unmarshal(patentInfoAsBytes, &patentInfo)
	if (err != nil) {
		logger.Error(err)
		return shim.Error(mylib.UnmarshalErrorMessage)
	}
	patentInfo.CurrentStatus = mylib.REJECTED
	patentInfo.UpdatedAt = args[3]

	//Action
	var action ActionDetail
	action.Status = mylib.REJECTED
	action.ActionBy = userDetail
	action.ActionDate = patentInfo.UpdatedAt

	patentInfo.ActionDetails = append(patentInfo.ActionDetails,action)

	//Marshalling final patent info
	val, err := json.Marshal(patentInfo)
	if(err != nil){
		logger.Error(err)
		shim.Error(mylib.MarshalErrorMessage)
	}

	//Updating patent corresponding to its application refrence number
	err = stub.PutState(key, []byte(val))
	if(err != nil){
		logger.Error(err)
		return shim.Error(mylib.PutErrorMessage)
	}

	return shim.Success([]byte("Patent is rejected successfully"))
}


//=================================================================================================================================
//   requestAdditionalInformation : request additional information by formality officer
//=================================================================================================================================
func (t *ChainCode) requestAdditionalInformation(stub shim.ChaincodeStubInterface, args []string) pb.Response{

	if len(args) != 8 {
		logger.Error("Incorrect number of arguements")
		return shim.Error("Incorrect number of arguements")
	}

	key := args[1]

	//User
	var userDetail User
	err := json.Unmarshal([]byte(args[2]), & userDetail)
	if err != nil {
		logger.Error(err)
		return shim.Error("Error parsing user info")
	}

	//
	//if userDetail.Role != mylib.FORMALITYOFFICERS {
	//	return shim.Error("Only Users with Formality Officer role can request for additional information")
	//}

	var patentInfo PatentInfo
	patentInfoAsBytes, err := stub.GetState(key)
	if(err != nil){
		logger.Error(err)
		shim.Error(mylib.GetStateErrorMessage)
	}

	//Unmarshalling package json string to native go structure
	err = json.Unmarshal(patentInfoAsBytes, &patentInfo)
	if (err != nil) {
		logger.Error(err)
		return shim.Error(mylib.UnmarshalErrorMessage)
	}

	var basicInfo PatentBasicInfo
	basicInfo.ApplicationRefNumber = patentInfo.ApplicationRefNumber
	basicInfo.PriorityNumber = patentInfo.PriorityNumber
	basicInfo.PatentTitle = patentInfo.PatentTitle
	basicInfo.CurrentStatus = patentInfo.CurrentStatus
	basicInfo.CreatedAt = patentInfo.CreatedAt
	basicInfo.UpdatedAt = patentInfo.UpdatedAt
	var submittedBy string
	if userDetail.Role != mylib.APPLICANT{
		actionDetails := patentInfo.ActionDetails
		for i:=0;i<len(actionDetails) ;i++  {
			actionStatus := actionDetails[i].Status
			if actionStatus == mylib.SUBMITTED {
				submittedBy = actionDetails[i].ActionBy.Email
			}
		}
		basicInfo.SubmittedTo = patentInfo.SubmittedTo.POName
		basicInfo.SubmittedBy = submittedBy
	}


	patentInfo.AESkey = aes_key

	decrypted_data, err := decrypt_data_using_aes_key(patentInfo.PatentData, patentInfo.AESkey)
	if err != nil {
		logger.Error(err)
		return shim.Error("Not able to decrpt")
	}

	var patentData Patent

	err = json.Unmarshal([]byte(decrypted_data), & patentData);
	if err != nil {
		logger.Error(err)
		return shim.Error("Not able to decrpt")
	}

	requestNumber := "REQ/2018/"
	// count :=1
	// The Request Number Starts
	if len(patentData.FormalityOfficerRequests) == 0 {
		requestNumber = requestNumber + strconv.Itoa(1)
	} else if len(patentData.FormalityOfficerRequests) > 0 {
		requestNumber = requestNumber + strconv.Itoa(len(patentData.FormalityOfficerRequests)+1)
	}

	// The Request Number Ends


	requestData := args[3]

	var foRequest FormalityOfficerRequest
	foRequest.RequestNumber =  requestNumber

	err = json.Unmarshal([]byte(requestData), & foRequest)
	if err != nil {
		logger.Error(err)
		return shim.Error("Error in unmarshalling request data")
	}


	docData := args[4]
	var documentDetail Document

	err = json.Unmarshal([]byte(docData), & documentDetail)
	if err != nil {
		logger.Error(err)
		return shim.Error("Error in unmarshalling document data data")
	}

	//documentDetail.DocumentType = mylib.RequestLetter
	switch args[5] {
	case "Reject_Letter":
		documentDetail.DocumentType = mylib.RejectionReport
	case "Search_Report":
		documentDetail.DocumentType = mylib.SearchReports
	case "Opinion_On_Patentability":
		documentDetail.DocumentType = mylib.OpinionofPatentability
	case "Request_Letter":
		documentDetail.DocumentType = mylib.RequestLetter
	case "Supporting_Document":
		documentDetail.DocumentType = mylib.SupportingDocument
		}


	documentDetail.Uploadedby = userDetail

	foRequest.RequestedFrom = userDetail
	foRequest.Documents = append(foRequest.Documents, documentDetail)


	if(args[5] == "Supporting_Document") {
		foRequest.Status = mylib.COMPLETE
	}else{
		foRequest.Status = mylib.INCOMPLETE
	}

	if(args[7] != "-1" && args[7]!="") {
		index := args[7]
		i, err := strconv.Atoi(index)
		if err != nil {
			logger.Error(err)
			return shim.Error("Error saving patent")
		}
		patentData.FormalityOfficerRequests[i].Status=mylib.COMPLETE
	}


	if(args[5]!="Supporting_Document") {
		patentData.FormalityOfficerRequests = append(patentData.FormalityOfficerRequests, foRequest)
	}else if(args[5]=="Supporting_Document"){
		patentData.Documents=append(patentData.Documents,documentDetail)
	}

	toCrypt_data, err := json.Marshal(patentData)
	if err != nil {
		logger.Error(err)
		return shim.Error("Error parsing data info")
	}

	encrypted_string, err := encrypt_data_using_aes_key(string(toCrypt_data),patentInfo.AESkey)
	if err != nil {
		logger.Error(err)
		return shim.Error("Error saving patent")
	}

	patentInfo.PatentData = encrypted_string

	// Making the entry for capturing the data for the action performed by the user in ActionDetails array starts
	var action ActionDetail
	if args[5] == "Reject_Letter"{
		patentInfo.CurrentStatus = mylib.REJECTED
		action.Status = mylib.REJECTED
	}

	if args[5] == "Search_Report"{
		patentInfo.CurrentStatus = mylib.SEARCHREPORTATTACHED
		action.Status = mylib.SEARCHREPORTATTACHED
	}

	patentInfo.UpdatedAt = args[6]
	action.ActionBy = userDetail
	action.ActionDate = patentInfo.UpdatedAt

	// Making the entry for capturing the data for the action performed by the user in ActionDetails array ends

	patentInfo.ActionDetails = append(patentInfo.ActionDetails,action)



	var notification Notification
	notification.NotificationID=strconv.Itoa(rand.Int())
	notification.ApplicationRefNumber =  patentInfo.ApplicationRefNumber;
	notification.ActionDate =  args[6]
	notification.Documents = append(notification.Documents, documentDetail)
	notification.Status	= args[5]
	notification.ActionBy =	userDetail
	notification.PatentTitle = patentInfo.PatentTitle
	notification.PatentBasicInfo = basicInfo
	switch args[5] {
	case "Search_Report":
		notification.Description=mylib.SearchReportDescription
	case "Reject_Letter":
		notification.Description=mylib.RejectReportDescription
	case "Opinion_On_Patentability":
		notification.Description=mylib.OpinionOfPatentabilityDescription
	case "Request_Letter":
		notification.Description=mylib.RequestLetterDescription
	default:
	}

	//Marshalling final patent info
	val, err := json.Marshal(patentInfo)
	if(err != nil){
		logger.Error(err)
		shim.Error(mylib.MarshalErrorMessage)
	}

	//Updating patent corresponding to its application refrence number
	err = stub.PutState(key, []byte(val))
	if(err != nil){
		logger.Error(err)
		return shim.Error(mylib.PutErrorMessage)
	}
	
	//For noification Tarun
	var notifications Notifications


	 notificationInfoAsBytes, err := stub.GetState("notifications")
	 if(err != nil){
	 	logger.Error(err)
	 	shim.Error(mylib.GetStateErrorMessage)
	 	}
		if len(notificationInfoAsBytes) != 0 {
//			Unmarshalling package json string to native go structure
			err = json.Unmarshal(notificationInfoAsBytes, &notifications)
			if (err != nil) {
				logger.Error(err)
				return shim.Error(mylib.UnmarshalErrorMessage)
			}
		}
	//fmt.Println("Notificationnnnnn: ",notifications, " ",notification.ApplicationRefNumber)


	if(args[5] == "Supporting_Document") {
		notifications.FONotifications=append(notifications.FONotifications, notification)
	}else{
		notifications.ApplicantNotifications=append(notifications.ApplicantNotifications, notification)
	}
	
	//Marshalling notificaions
	val1, err := json.Marshal(notifications)
	if(err != nil){
		logger.Error(err)	
		shim.Error(mylib.MarshalErrorMessage)
	}

	//For notification
	err = stub.PutState("notifications",[]byte(val1))
	if(err != nil){
		logger.Error(err)		
		return shim.Error(mylib.PutErrorMessage)
	}

	return shim.Success([]byte("Additional request by formality officer successful"))
}


func(t *ChainCode) updateUnreadNotifications(stub shim.ChaincodeStubInterface, args[] string) pb.Response {

	if len(args) != 4 {
		logger.Error("No. of arguements not sufficient")
		return shim.Error("No. of arguements not sufficient")
	}

	fmt.Println("Notificationnnnnn:11111454545 startsss")

	var notifications Notifications
	notificationsAsBytes, err := stub.GetState("notifications")
	if(err != nil){
		logger.Error(err)
		shim.Error(mylib.GetStateErrorMessage)
	}

	//Unmarshalling package json string to native go structure
	err = json.Unmarshal(notificationsAsBytes, &notifications)
	if (err != nil) {
		logger.Error(err)
		return shim.Error(mylib.UnmarshalErrorMessage)
	}

	role := args[1]
	notificationID:=args[3]
	if(mylib.RoleLookUp[role] == mylib.APPLICANT){
		for i:=0;i<len(notifications.ApplicantNotifications) ;i++  {
			if(notificationID == notifications.ApplicantNotifications[i].NotificationID){
				notifications.ApplicantNotifications=append(notifications.ApplicantNotifications[:i], notifications.ApplicantNotifications[i+1:]...)
			}
		}


	}else if(mylib.RoleLookUp[role] == mylib.FORMALITYOFFICERS){
		for i:=0;i<len(notifications.FONotifications) ;i++  {
			if(notificationID == notifications.FONotifications[i].NotificationID){
				notifications.FONotifications=append(notifications.FONotifications[:i], notifications.FONotifications[i+1:]...)
			}
		}
	}

	//Marshalling notificaions
	val1, err := json.Marshal(notifications)
	if(err != nil){
		logger.Error(err)
		shim.Error(mylib.MarshalErrorMessage)
	}

	//For notification
	err = stub.PutState("notifications",[]byte(val1))
	if(err != nil){
		logger.Error(err)
		return shim.Error(mylib.PutErrorMessage)
	}

	return shim.Success(notificationsAsBytes)

}

//=================================================================================================================================
//   getPatentListForPublic : getPatentListForPublic  info from world state
//=================================================================================================================================
func (t *ChainCode) getPatentListForPublic(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
    
    fmt.Println("Called")
 	role := args[1]
 	oldDate := args[3]
  
    var queryString,submittedBy string

  //	queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}}]}}",  "DRAFT","1515918178")
	queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}",  "Draft",oldDate)
	var patentHolder PatentStatusHolder
	var patentInfo PatentInfo

	fmt.Println("The Query String",queryString)

	resultsIterator, error := stub.GetQueryResult(queryString)
	if error != nil {
		return shim.Error("Error in the resultsIterator ")
	}
	defer resultsIterator.Close()
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error("Error in the queryResponse ")
		}
		fmt.Println("We have a resultset")
		err = json.Unmarshal(queryResponse.Value, &patentInfo);
		var basicInfo PatentBasicInfo
		basicInfo.ApplicationRefNumber = patentInfo.ApplicationRefNumber
		basicInfo.PriorityNumber = patentInfo.PriorityNumber
		basicInfo.PatentTitle = patentInfo.PatentTitle
		basicInfo.CurrentStatus = patentInfo.CurrentStatus
		basicInfo.CreatedAt = patentInfo.CreatedAt
		basicInfo.UpdatedAt = patentInfo.UpdatedAt
		basicInfo.SubmittedAt = patentInfo.SubmittedAt

		if mylib.RoleLookUp[role] != mylib.APPLICANT{
			actionDetails := patentInfo.ActionDetails
			for i:=0;i<len(actionDetails) ;i++  {
				actionStatus := actionDetails[i].Status
				if actionStatus == mylib.SUBMITTED {
					userFirstName := actionDetails[i].ActionBy.FirstName
					userLastName := actionDetails[i].ActionBy.LastName
					submittedBy = userFirstName+" "+userLastName
					
				}
			}
			basicInfo.SubmittedTo = patentInfo.SubmittedTo.POName
			basicInfo.SubmittedBy = submittedBy
		}

		patentHolder.PatentList = append(patentHolder.PatentList, basicInfo)

	}
	fmt.Println("The Patent Holder PatentList Object :",patentHolder.PatentList)
    fmt.Println("Now calling the Count Function")
    patentCount, err := t.getPatentCountByStatusForPublic(stub, args[1],args[2],oldDate)
    if err != nil {
        return shim.Error("Error in the getPatentCountByStatusForPublic ")
    }

	patentHolder.PatentCount = patentCount
	
	fmt.Println("The Patent Holder PatentCount Object :",patentHolder.PatentCount)
    publicViewRecord, err := json.Marshal(patentHolder)
    if err != nil {
        return shim.Error("Error converting status holder record")
    }

    return shim.Success(publicViewRecord)
	
}


func(t *ChainCode) getPatentCountByStatusForPublic(stub shim.ChaincodeStubInterface,role string ,applicantEmail string,oldDate string)(PatentStatusCount, error) {

    var patentCount PatentStatusCount
    var queryStringSubmitted,queryStringRejected,queryStringSearchReport string
    searchReportCount :=0
    submittedCount:=0
    rejectedCount :=0

		queryStringSubmitted = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}",  "Submitted",oldDate)
		queryStringRejected = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}",  "Rejected",oldDate)
		queryStringSearchReport = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}",  "SearchReportAttached",oldDate)

		resultsIterator, error := stub.GetQueryResult(queryStringSubmitted)
		if error != nil {
			return patentCount, error
		}
		defer resultsIterator.Close()
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return patentCount, err
			}
			if queryResponse != nil {
				submittedCount= submittedCount+1
			}
		}
		patentCount.Submitted_count = submittedCount

		resultsIterator, error = stub.GetQueryResult(queryStringRejected)
		if error != nil {
			return patentCount, error
		}
		defer resultsIterator.Close()
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return patentCount, err
			}
			if queryResponse != nil {
				rejectedCount= rejectedCount+1
			}
		}
		patentCount.Rejected_count = rejectedCount
        
		resultsIterator, error = stub.GetQueryResult(queryStringSearchReport)
		if error != nil {
			return patentCount, error
		}
		defer resultsIterator.Close()
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return patentCount, err
			}
			if queryResponse != nil {
				searchReportCount= searchReportCount+1
			}
		}
		patentCount.SearchReportAttached_count = searchReportCount
	    return patentCount, nil
}
