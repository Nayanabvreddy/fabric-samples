package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	bytes "bytes"

	"encoding/json"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"crypto/rand"
	"errors"
	"github.com/hyperledger/fabric/examples/chaincode/go/mylib"
	"crypto/rsa"
	"math/big"
	"crypto/elliptic"
	"crypto/ecdsa"
)

func(t *ChainCode) getPatentListByStatus(stub shim.ChaincodeStubInterface, args[] string, holder mylib.PatentStatus)([] byte, error) {

    role := args[1]
    var queryString,submittedBy string

    if mylib.RoleLookUp[role] == mylib.APPLICANT {
		queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_email\": \"%s\"}}}]}}",  holder, args[2])
    }

	if ( (mylib.RoleLookUp[role] == mylib.FORMALITYOFFICERS) || (mylib.RoleLookUp[role] == mylib.NPO) ) {
		if holder !=mylib.DRAFT {
			queryString = fmt.Sprintf("{\"selector\":{\"current_status\":{\"$eq\": \"%s\"}}}", holder)
		}
	}

	var patentHolder PatentStatusHolder
	var patentInfo PatentInfo

	resultsIterator, error := stub.GetQueryResult(queryString)
	if error != nil {
		return nil, error
	}
	defer resultsIterator.Close()
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
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


    // append results of all holders to these bytes
    patentCount, err := t.getPatentCountByStatus(stub, args[1],args[2])
    if err != nil {
    	logger.Error(err)
        return nil, errors.New("error retreiving holder strings " +  err.Error())
    }

	patentHolder.PatentCount = patentCount
    bytes, err := json.Marshal(patentHolder)
    if err != nil {
    	logger.Error(err)
        return nil, errors.New("Error converting status holder record")
    }

    return bytes, nil
}

func(t *ChainCode) getSearchPatentList(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {

	role := args[1]
	email := args[2]
	status := args[3]
	createAt := args[4]
	 submitto := args[5]
	applicant := args[6]
	title := args[7]
	oldDate := args[8]

	fmt.Println("Submitted query email ",email)
	fmt.Println("Submitted query status ",status)
	fmt.Println("Submitted query created at ",createAt)
 fmt.Println("Submitted query for ",submitto)
	fmt.Println("created query by ",applicant)
	fmt.Println("Submitted query title ",title)
    var queryString,submittedBy string

   
	if ( (mylib.RoleLookUp[role] == mylib.FORMALITYOFFICERS) || (mylib.RoleLookUp[role] == mylib.NPO) )  {
		if mylib.PatentStatusLookUp[status] !=mylib.DRAFT {
			
			if ((createAt != "undefined" || createAt != "") && (submitto == "undefined" || submitto == "") && (applicant == "undefined" || applicant == "") && (title == "udefined" || title == "")) {
				
				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"created_at\":{\"$gt\":\"%s\"}}]}}", mylib.PatentStatusLookUp[status], createAt)
			
				}else if ((createAt == "undefined" || createAt == "") && (submitto == "undefined" || submitto == "" ) && (applicant == "undefined" || applicant == "") && (title == "undefined" || title == "")) {
				   fmt.Println("query in ")
					queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}}]}}", mylib.PatentStatusLookUp[status])									   
			    } else if ((createAt == "undefined" || createAt == "") && (submitto != "undefined" || submitto != "" ) && (applicant == "undefined" || applicant == "") && (title == "undefined" || title == "") ) {

				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"submitted_to\":{\"npo_name\": \"%s\"}}]}}", mylib.PatentStatusLookUp[status], submitto)
				
				} else if ((createAt == "undefined" || createAt == "") && (submitto == "undefined" || submitto == "" ) && (applicant != "undefined" || applicant != "") && (title == "undefined" || title == "") ){

				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_name\":{\"$regex\": \"%s+\"}}}}]}}", mylib.PatentStatusLookUp[status], applicant)
			   
				} else if ((createAt == "undefined" || createAt == "") && (submitto == "undefined" || submitto == "" ) && (applicant == "undefined" || applicant == "") && (title != "undefined" || title != "") ) {

				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"patent_title\":{\"$regex\": \"%s+\"}}]}}", mylib.PatentStatusLookUp[status], title)
			   
			    } else if ((createAt != "undefined" || createAt != "") && (submitto != "undefined" || submitto != "" ) && (applicant == "undefined" || applicant == "") && (title == "undefined" || title == "") ) {

				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"created_at\":{\"$gt\":\"%s\"}},{\"submitted_to\":{\"npo_name\": \"%s\"}}]}}", mylib.PatentStatusLookUp[status], createAt,submitto )		
				 
				} else if ((createAt != "undefined" || createAt != "") && (submitto == "undefined" || submitto == "" ) && (applicant != "undefined" || applicant != "") && (title == "undefined" || title == "") ) {
				   
				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"created_at\":{\"$gt\":\"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_name\":{\"$regex\": \"%s+\"}}}}]}}", mylib.PatentStatusLookUp[status], createAt, applicant )
						   
				} else if ((createAt != "undefined" || createAt != "") && (submitto == "undefined" || submitto == "" ) && (applicant == "undefined" || applicant == "") && (title != "undefined" || title != "") ) {
				   
				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"created_at\":{\"$gt\":\"%s\"}},{\"patent_title\":{\"$regex\": \"%s+\"}}]}}", mylib.PatentStatusLookUp[status], createAt, title )
							   
			    } else if ((createAt == "undefined" || createAt == "") && (submitto != "undefined" || submitto != "" ) && (applicant != "undefined" || applicant != "") && (title == "undefined" || title == "") ) {
				   
				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"submitted_to\":{\"npo_name\": \"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_name\":{\"$regex\": \"%s+\"}}}}]}}", mylib.PatentStatusLookUp[status], submitto, applicant )
								   
				} else if ((createAt == "undefined" || createAt == "") && (submitto != "undefined" || submitto != "" ) && (applicant == "undefined" || applicant == "") && (title != "undefined" || title != "") ) {
				   
					queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"submitted_to\":{\"npo_name\": \"%s\"}},{\"patent_title\":{\"$regex\": \"%s+\"}}]}}", mylib.PatentStatusLookUp[status], submitto, title )
									   
			    } else if ((createAt == "undefined" || createAt == "") && (submitto != "undefined" || submitto != "" ) && (applicant != "undefined" || applicant != "") && (title != "undefined" || title != "") ) {
				   
					queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_name\":{\"$regex\": \"%s+\"}}}},{\"patent_title\":{\"$regex\": \"%s+\"}}]}}", mylib.PatentStatusLookUp[status], applicant, title )
									   
			    } else if ((createAt != "undefined" || createAt != "") && (submitto != "undefined" || submitto != "" ) && (applicant != "undefined" || applicant != "") && (title == "undefined" || title == "") ) {
				   
					queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"created_at\":{\"$gt\":\"%s\"}},{\"submitted_to\":{\"npo_name\": \"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_name\":{\"$regex\": \"%s+\"}}}}]}}", mylib.PatentStatusLookUp[status], createAt, submitto, applicant )
									   
				} else if ((createAt != "undefined" || createAt != "") && (submitto == "undefined" || submitto == "" ) && (applicant != "undefined" || applicant != "") && (title != "undefined" || title != "") ) {
				   
					queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"created_at\":{\"$gt\":\"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_name\":{\"$regex\": \"%s+\"}}}},{\"patent_title\":{\"$regex\": \"%s+\"}}]}}", mylib.PatentStatusLookUp[status], createAt, applicant, title )
									   
			    } else if ((createAt == "undefined" || createAt == "") && (submitto != "undefined" || submitto != "" ) && (applicant != "undefined" || applicant != "") && (title != "undefined" || title != "") ) {
				   
					queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"submitted_to\":{\"npo_name\": \"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_name\":{\"$regex\": \"%s+\"}}}},{\"patent_title\":{\"$regex\": \"%s+\"}}]}}", mylib.PatentStatusLookUp[status], submitto, applicant, title )
									   
			    } else if ((createAt != "undefined" || createAt != "")  && (submitto != "undefined" || submitto != "" ) && (applicant == "undefined" || applicant == "") && (title != "undefined" || title != "") ) {
				   
					queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"created_at\":{\"$gt\":\"%s\"}},{\"submitted_to\":{\"npo_name\": \"%s\"}},{\"patent_title\":{\"$regex\": \"%s+\"}}]}}", mylib.PatentStatusLookUp[status], createAt , submitto, title )
									   
			    }else{

			    	queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"created_at\":{\"$gt\":\"%s\"}},{\"submitted_to\":{\"npo_name\": \"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_name\":{\"$regex\": \"%s+\"}},{\"patent_title\":{\"$regex\": \"%s+\"}}]}}", mylib.PatentStatusLookUp[status], createAt, submitto, applicant, title)
			    }
		}
	} 

	if (mylib.RoleLookUp[role] == mylib.PUBLIC) {
			 if ((createAt != "undefined" || createAt != "") && (submitto == "undefined" || submitto == "") && (applicant == "undefined" || applicant == "") && (title == "udefined" || title == "")) {
				
				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"created_at\":{\"$gt\":\"%s\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, createAt, oldDate)
			
				}else if ((createAt == "undefined" || createAt == "") && (submitto == "undefined" || submitto == "" ) && (applicant == "undefined" || applicant == "") && (title == "undefined" || title == "")) {
				   fmt.Println("query in Public")
					queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, oldDate)									   
			    } else if ((createAt == "undefined" || createAt == "") && (submitto != "undefined" || submitto != "" ) && (applicant == "undefined" || applicant == "") && (title == "undefined" || title == "") ) {

				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"submitted_to\":{\"npo_name\": \"%s\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, submitto, oldDate)
				
				} else if ((createAt == "undefined" || createAt == "") && (submitto == "undefined" || submitto == "" ) && (applicant != "undefined" || applicant != "") && (title == "undefined" || title == "") ){

				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_name\":{\"$regex\": \"%s+\"}}}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, applicant, oldDate)
			   
				} else if ((createAt == "undefined" || createAt == "") && (submitto == "undefined" || submitto == "" ) && (applicant == "undefined" || applicant == "") && (title != "undefined" || title != "") ) {

				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"patent_title\":{\"$regex\": \"%s+\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, title, oldDate)
			   
			    } else if ((createAt != "undefined" || createAt != "") && (submitto != "undefined" || submitto != "" ) && (applicant == "undefined" || applicant == "") && (title == "undefined" || title == "") ) {

				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"created_at\":{\"$gt\":\"%s\"}},{\"submitted_to\":{\"npo_name\": \"%s\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, createAt,submitto, oldDate)		
				 
				} else if ((createAt != "undefined" || createAt != "") && (submitto == "undefined" || submitto == "" ) && (applicant != "undefined" || applicant != "") && (title == "undefined" || title == "") ) {
				   
				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"created_at\":{\"$gt\":\"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_name\":{\"$regex\": \"%s+\"}}}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, createAt, applicant, oldDate)
						   
				} else if ((createAt != "undefined" || createAt != "") && (submitto == "undefined" || submitto == "" ) && (applicant == "undefined" || applicant == "") && (title != "undefined" || title != "") ) {
				   
				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"created_at\":{\"$gt\":\"%s\"}},{\"patent_title\":{\"$regex\": \"%s+\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, createAt, title, oldDate)
							   
			    } else if ((createAt == "undefined" || createAt == "") && (submitto != "undefined" || submitto != "" ) && (applicant != "undefined" || applicant != "") && (title == "undefined" || title == "") ) {
				   
				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"submitted_to\":{\"npo_name\": \"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_name\":{\"$regex\": \"%s+\"}}}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, submitto, applicant, oldDate)
								   
				} else if ((createAt == "undefined" || createAt == "") && (submitto != "undefined" || submitto != "" ) && (applicant == "undefined" || applicant == "") && (title != "undefined" || title != "") ) {
				   
					queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"submitted_to\":{\"npo_name\": \"%s\"}},{\"patent_title\":{\"$regex\": \"%s+\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, submitto, title, oldDate)
									   
			    } else if ((createAt == "undefined" || createAt == "") && (submitto != "undefined" || submitto != "" ) && (applicant != "undefined" || applicant != "") && (title != "undefined" || title != "") ) {
				   
					queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_name\":{\"$regex\": \"%s+\"}}}},{\"patent_title\":{\"$regex\": \"%s+\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, applicant, title, oldDate)
									   
			    } else if ((createAt != "undefined" || createAt != "") && (submitto != "undefined" || submitto != "" ) && (applicant != "undefined" || applicant != "") && (title == "undefined" || title == "") ) {
				   
					queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"created_at\":{\"$gt\":\"%s\"}},{\"submitted_to\":{\"npo_name\": \"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_name\":{\"$regex\": \"%s+\"}}}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, createAt, submitto, applicant, oldDate)
									   
				} else if ((createAt != "undefined" || createAt != "") && (submitto == "undefined" || submitto == "" ) && (applicant != "undefined" || applicant != "") && (title != "undefined" || title != "") ) {
				   
					queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"created_at\":{\"$gt\":\"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_name\":{\"$regex\": \"%s+\"}}}},{\"patent_title\":{\"$regex\": \"%s+\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, createAt, applicant, title, oldDate)
									   
			    } else if ((createAt == "undefined" || createAt == "") && (submitto != "undefined" || submitto != "" ) && (applicant != "undefined" || applicant != "") && (title != "undefined" || title != "") ) {
				   
					queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"submitted_to\":{\"npo_name\": \"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_name\":{\"$regex\": \"%s+\"}}}},{\"patent_title\":{\"$regex\": \"%s+\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, submitto, applicant, title, oldDate)
									   
			    } else if ((createAt != "undefined" || createAt != "")  && (submitto != "undefined" || submitto != "" ) && (applicant == "undefined" || applicant == "") && (title != "undefined" || title != "") ) {
				   
					queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"created_at\":{\"$gt\":\"%s\"}},{\"submitted_to\":{\"npo_name\": \"%s\"}},{\"patent_title\":{\"$regex\": \"%s+\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, createAt , submitto, title, oldDate)
									   
			    } else{

			    	queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"created_at\":{\"$gt\":\"%s\"}},{\"submitted_to\":{\"npo_name\": \"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_name\":{\"$regex\": \"%s+\"}},{\"patent_title\":{\"$regex\": \"%s+\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, createAt, submitto, applicant, title, oldDate)
			    }
	} 

	var patentHolder PatentStatusHolder
	var patentInfo PatentInfo
fmt.Println("query string is ", queryString)
	resultsIterator, error := stub.GetQueryResult(queryString)
	if error != nil {
		return nil, error
	}
	defer resultsIterator.Close()
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
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
				fmt.Println("The action detail with the ID ",i)
				actionStatus := actionDetails[i].Status
				if actionStatus == mylib.SUBMITTED {
					userFirstName := actionDetails[i].ActionBy.FirstName
					userLastName := actionDetails[i].ActionBy.LastName
					submittedBy = userFirstName+" "+userLastName
					fmt.Println("The action status",actionStatus)
				}
			}
			basicInfo.SubmittedTo = patentInfo.SubmittedTo.POName
			basicInfo.SubmittedBy = submittedBy
		}

		patentHolder.PatentList = append(patentHolder.PatentList, basicInfo)

	}

    // append results of all holders to these bytes
        if mylib.RoleLookUp[role] == mylib.PUBLIC {
    patentCount, err := t.getPatentCountByStatusForPublic(stub, args[1],args[2],oldDate)
			if err != nil {
				return nil, errors.New("error retreiving holder strings " +  err.Error())
			}

			patentHolder.PatentCount = patentCount
     } else{
        patentCount, err := t.getPatentCountByStatus(stub, args[1],args[2])
			if err != nil {
				return nil, errors.New("error retreiving holder strings " +  err.Error())
			}

			patentHolder.PatentCount = patentCount
    }

    bytes, err := json.Marshal(patentHolder)
    if err != nil {
        return nil, errors.New("Error converting status holder record")
    }

    return bytes, nil
}

func(t *ChainCode) getSearchPatentappList(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {

	role := args[1]
	email := args[2]
	status := args[3]
	apprefer := args[4]
	oldDate := args[5]

	fmt.Println("Submitted query role ",role)  
	fmt.Println("Submitted query email ",email)
	fmt.Println("Submitted query status ",status)
	fmt.Println("Submitted query created at ",apprefer)

    var queryString,submittedBy string

	if ( (mylib.RoleLookUp[role] == mylib.FORMALITYOFFICERS) || (mylib.RoleLookUp[role] == mylib.NPO) ) {
		if mylib.PatentStatusLookUp[status] !=mylib.DRAFT {
			 if (apprefer != "undefined" || apprefer != "") {
				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"application_ref_number\":{\"$eq\":\"%s\"}}]}}", mylib.PatentStatusLookUp[status], apprefer)
				} 
		}
	} 

	if (mylib.RoleLookUp[role] == mylib.PUBLIC) {		
			 if (apprefer != "undefined" || apprefer != "") {
				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"application_ref_number\":{\"$eq\":\"%s\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, apprefer, oldDate)
				} 		
	} 

	var patentHolder PatentStatusHolder
	var patentInfo PatentInfo

	resultsIterator, error := stub.GetQueryResult(queryString)
	if error != nil {
		return nil, error
	}
	defer resultsIterator.Close()
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
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
				fmt.Println("The action detail with the ID ",i)
				actionStatus := actionDetails[i].Status
				if actionStatus == mylib.SUBMITTED {
					userFirstName := actionDetails[i].ActionBy.FirstName
					userLastName := actionDetails[i].ActionBy.LastName
					submittedBy = userFirstName+" "+userLastName
					fmt.Println("The action status",actionStatus)
				}
			}
			basicInfo.SubmittedTo = patentInfo.SubmittedTo.POName
			basicInfo.SubmittedBy = submittedBy
		}
		patentHolder.PatentList = append(patentHolder.PatentList, basicInfo)
	}

	// append results of all holders to these bytes
	if mylib.RoleLookUp[role] == mylib.PUBLIC {
		patentCount, err := t.getPatentCountByStatusForPublic(stub, args[1],args[2],oldDate)
		if err != nil {
			return nil, errors.New("error retreiving holder strings " +  err.Error())
		}

		patentHolder.PatentCount = patentCount
	} else{
		patentCount, err := t.getPatentCountByStatus(stub, args[1],args[2])
		if err != nil {
			return nil, errors.New("error retreiving holder strings " +  err.Error())
		}

		patentHolder.PatentCount = patentCount
	}
    bytes, err := json.Marshal(patentHolder)
    if err != nil {
        return nil, errors.New("Error converting status holder record")
    }

    return bytes, nil
}

func(t *ChainCode) getSearchPatentpriList(stub shim.ChaincodeStubInterface, args[] string)([] byte, error) {

	role := args[1]
	email := args[2]
	status := args[3]
	prinumber := args[4]
	oldDate := args[5]

	fmt.Println("Submitted query email ",email)
	fmt.Println("Submitted query status ",status)
	fmt.Println("Submitted query created at ",prinumber)
    fmt.Println("Submitted query date ",oldDate)

    var queryString,submittedBy string

	if ( (mylib.RoleLookUp[role] == mylib.FORMALITYOFFICERS) || (mylib.RoleLookUp[role] == mylib.NPO) )  {
		if mylib.PatentStatusLookUp[status] !=mylib.DRAFT {
			   if (prinumber != "undefined" || prinumber != "") {
			   fmt.Println("not public role")				
				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"priority_number\":{\"$regex\": \"^%s\"}}]}}", mylib.PatentStatusLookUp[status], prinumber)			
				} 
		}
	} 
    
    if (mylib.RoleLookUp[role] == mylib.PUBLIC) {
		if mylib.PatentStatusLookUp[status] !=mylib.DRAFT {
			 if (prinumber != "undefined" || prinumber != "") {				
fmt.Println("public role")
				queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$ne\": \"%s\"}},{\"priority_number\":{\"$regex\": \"^%s\"}},{\"submitted_at\":{\"$lte\": \"%s\"}}]}}", mylib.DRAFT, prinumber, oldDate)			
				} 
		}
	} 

	fmt.Println("queryString for public", queryString)

	var patentHolder PatentStatusHolder
	var patentInfo PatentInfo

	resultsIterator, error := stub.GetQueryResult(queryString)
	if error != nil {
		return nil, error
	}
	defer resultsIterator.Close()
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
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
				fmt.Println("The action detail with the ID ",i)
				actionStatus := actionDetails[i].Status
				if actionStatus == mylib.SUBMITTED {
					userFirstName := actionDetails[i].ActionBy.FirstName
					userLastName := actionDetails[i].ActionBy.LastName
					submittedBy = userFirstName+" "+userLastName
					fmt.Println("The action status",actionStatus)
				}
			}
			basicInfo.SubmittedTo = patentInfo.SubmittedTo.POName
			basicInfo.SubmittedBy = submittedBy
		}

		patentHolder.PatentList = append(patentHolder.PatentList, basicInfo)

	}

	// append results of all holders to these bytes
	if mylib.RoleLookUp[role] == mylib.PUBLIC {
		patentCount, err := t.getPatentCountByStatusForPublic(stub, args[1],args[2],oldDate)
		if err != nil {
			return nil, errors.New("error retreiving holder strings " +  err.Error())
		}

		patentHolder.PatentCount = patentCount
	} else{
		patentCount, err := t.getPatentCountByStatus(stub, args[1],args[2])
		if err != nil {
			return nil, errors.New("error retreiving holder strings " +  err.Error())
		}

		patentHolder.PatentCount = patentCount
	}
	fmt.Println("query public patentCount ", patentHolder.PatentCount)
    bytes, err := json.Marshal(patentHolder)
    if err != nil {
        return nil, errors.New("Error converting status holder record")
    }

    return bytes, nil
}


func(t *ChainCode) getPatentCountByStatus(stub shim.ChaincodeStubInterface,role string ,applicantEmail string)(PatentStatusCount, error) {

    var patentCount PatentStatusCount
    var queryString string
	userRole := role
    for _, value := range mylib.PatentStatusLookUp {
		count:=0

        //queryString := fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_email\": \"%s\"}}}]}}", value , applicantEmail)
		if mylib.RoleLookUp[userRole] == mylib.APPLICANT {
			queryString = fmt.Sprintf("{\"selector\":{\"$and\":[{\"current_status\":{\"$eq\": \"%s\"}},{\"Applicants\":{\"$elemMatch\":{\"ext_email\": \"%s\"}}}]}}", value , applicantEmail)
		}

		if ( (mylib.RoleLookUp[role] == mylib.FORMALITYOFFICERS) || (mylib.RoleLookUp[role] == mylib.NPO) ) {
				queryString = fmt.Sprintf("{\"selector\":{\"current_status\":{\"$eq\": \"%s\"}}}", value)
		}

		resultsIterator, error := stub.GetQueryResult(queryString)
		if error != nil {
			logger.Error(error)
			return patentCount, errors.New("Error creating status holders")
		}
		defer resultsIterator.Close()
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return patentCount, err
			}
			if queryResponse != nil {
				count= count+1
			}
		}
		if mylib.RoleLookUp[userRole] == mylib.APPLICANT {
			switch value {

			case mylib.DRAFT:
				patentCount.Draft_count = count

			case mylib.SUBMITTED:
				patentCount.Submitted_count = count

			case mylib.REJECTED:
				patentCount.Rejected_count = count

			case mylib.SEARCHREPORTATTACHED:
				patentCount.SearchReportAttached_count = count


			}
		}

		if ( mylib.RoleLookUp[userRole] == mylib.FORMALITYOFFICERS ) || ( mylib.RoleLookUp[userRole] == mylib.NPO ) || (mylib.RoleLookUp[role] == mylib.PUBLIC)  {
			switch value {

			case mylib.SUBMITTED:
				patentCount.Submitted_count = count

			case mylib.REJECTED:
				patentCount.Rejected_count = count

			case mylib.SEARCHREPORTATTACHED:
				patentCount.SearchReportAttached_count = count


			}
		}
    }

    return patentCount, nil
}


//==============================================================================================================================
// Invoke another chaincode
//==============================================================================================================================
func(t *ChainCode) invokeConfigChaincode(stub shim.ChaincodeStubInterface, invokeArgs [][]byte)(string, error) {
	response := stub.InvokeChaincode(mylib.ConfigChainCode, invokeArgs, mylib.ConfigChannel)
	if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to invoke chaincode. Got error: %s")
		logger.Error(errStr)
		return "null",errors.New(errStr)
	}
	return string(response.Payload),nil
}

// ToChaincodeArgs converts string args to []byte args
func ToChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

//==============================================================================================================================
// save_changes - Writes to the ledger the KycProfile struct passed in a JSON format. Uses the shim file's
//                method 'PutState'.
//==============================================================================================================================
func(t *ChainCode) saveChanges(stub shim.ChaincodeStubInterface, patentInfo PatentInfo)(bool, error) {

    bytes, err := json.Marshal(patentInfo)

    if err != nil {
        logger.Error(err)
        return false, errors.New("Error converting patent record")
    }

    err = stub.PutState(patentInfo.ApplicationRefNumber, bytes)
    if err != nil {
        logger.Error(err)
        return false, errors.New("Error storing patent record")
    }

    return true, nil
}

func(t *ChainCode) retrievePatent(stub shim.ChaincodeStubInterface, key string)(PatentInfo, error) {

	var patent PatentInfo

	bytes, err := stub.GetState(key);

	if err != nil {
		logger.Error(err)
		return patent, errors.New("retrieve_profile: Error retrieving profile with ID = " + key + " there was an error")
	}
	if bytes == nil {
		return patent, errors.New("retrieve_profile: Error retrieving profile with ID = " + key + " ; bytes are nil")
	}

	err = json.Unmarshal(bytes, & patent);

	if err != nil {
		logger.Error(err)
		return patent, errors.New("retrieve_profile: Corrupt record" + string(bytes))
	}

	return patent, nil
}


func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return buffer.Bytes(), nil
}

//==============================================================================================================================
// encrypt_data_using_aes_key -
//==============================================================================================================================

func encrypt_data_using_aes_key(data string, aes_key string)(string, error){

    key := []byte(aes_key)
    plaintext := []byte(data)

    block, err := aes.NewCipher(key)
    if err != nil {
    	logger.Error(err)
        return "null", errors.New("Error generating cipher")
    }

    ciphertext := make([]byte, len(plaintext))

    iv := []byte{'\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f'}

    stream := cipher.NewCTR(block, iv)
    stream.XORKeyStream(ciphertext, plaintext)

    base := base64.StdEncoding.EncodeToString(ciphertext)
    return string(base), nil

}

func generate_random_aes_key()(string, error) {

   size := 16
   rb := make([]byte,size)
   _, err := rand.Read(rb)

   if err != nil {
   	  logger.Error(err)
      return "null", errors.New("Error generating key")
   }

   key := base64.URLEncoding.EncodeToString(rb)
   return key, nil

}


//==============================================================================================================================
// decrypt_data_using_aes_key -
//==============================================================================================================================

func decrypt_data_using_aes_key(data string, aes_key string)(string, error){

  key := []byte(aes_key)

  ciphertext, err := base64.StdEncoding.DecodeString(data)
  if err != nil {
  	logger.Error(err)
    return "null", errors.New("Error decoding cipher")
  }

  block, err := aes.NewCipher(key)
  if err != nil {
  	 logger.Error(err)
     return "null", errors.New("Error generating cipher")
  }

  decrypt := make([]byte, len(ciphertext))

  iv := []byte{'\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f','\x0f'}

  stream := cipher.NewCTR(block, iv)
  stream.XORKeyStream(decrypt , ciphertext)

  return string(decrypt), nil

}

func encrypt_using_public_key(data string, modulus string)(string, error){

  //generate public key from modulus and e as 3 (hardcoded)
  decN, err := base64.StdEncoding.DecodeString(modulus)
  if err != nil {
  	  logger.Error(err)
      return "null", errors.New("Error generating cipher")
  }
  n := big.NewInt(0)
  n.SetBytes(decN)

  pKey := &rsa.PublicKey{N: n, E: 3}

    //encrpyt the aes key now
  secretMessage := []byte(data)
  rng := rand.Reader

  ciphertext, err := rsa.EncryptPKCS1v15(rng, pKey, secretMessage)
  if err != nil {
  	  logger.Error(err)
      return "null", errors.New("Error generating cipher")
  }
  testKeyGeneration(elliptic.P256(), "p256")
  base := base64.StdEncoding.EncodeToString(ciphertext)
  return string(base), nil

}

func testKeyGeneration(c elliptic.Curve, tag string) {

	priv, err := ecdsa.GenerateKey(c, rand.Reader)

	if err != nil {
		logger.Error(err)
		return 
	}
	fmt.Println("public key fetched")
	fmt.Println(priv.PublicKey)
	fmt.Println(priv.PublicKey.X)
	fmt.Println(priv.PublicKey.Y)
	if !c.IsOnCurve(priv.PublicKey.X, priv.PublicKey.Y) {

		fmt.Println("%s: public key invalid: %s", tag, err)

	}

/*	pubkeyCurve := elliptic.P256()
	var h hash.Hash
	h = md5.New()
	signhash := h.Sum(nil)*/


}


