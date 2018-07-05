package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/examples/chaincode/go/mylib"
	"os"
)

var logger = shim.NewLogger("patentchain")

func main() {
	logger.SetLevel(shim.LogDebug)
    logLevel, _ := shim.LogLevel(os.Getenv("SHIM_LOGGING_LEVEL"))
    shim.SetLoggingLevel(logLevel)
	err := shim.Start(new(ChainCode))
	if err != nil {
		logger.Error("Error starting patentchain chaincode: %s", err)
	}
}

//=================================================================================================================================
//   Init : Trigger at the time of chaincode deployment
//=================================================================================================================================
func (t *ChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

//=================================================================================================================================
//   Invoke : Perform transactions
//=================================================================================================================================
func (t *ChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if len(args)<2{
		logger.Error(mylib.ArgumentErrorMessage)
		shim.Error(mylib.ArgumentErrorMessage)
	}

	if args[0]=="storePatentDetails" {
		return t.storePatentDetails(stub, args)
    }else if args[0]=="updatePatentDetails" {
		return t.updatePatentDetails(stub, args)
	}else if args[0]=="updatePatentDocument" {
		return t.updatePatentDocument(stub, args)
	}else if args[0]=="validateAppRefNumber" {
		return t.validateAppRefNumber(stub, args)
	}else if args[0]=="submitPatentDetails" {
		return t.submitPatentDetails(stub, args)
	}else if args[0]=="deletePatentInfo" {
		return t.deletePatentInfo(stub, args[1])
	}else if args[0]=="rejectPatentApplication" {
		return t.rejectPatentApplication(stub, args)
	}else if args[0]=="requestAdditionalInformation" {
		return t.requestAdditionalInformation(stub, args)
	}else if args[0]=="updateUnreadNotifications" {
		return t.updateUnreadNotifications(stub, args)
	}else if function=="query"{
		switch args[0] {
		case "getPatentDetails":
			return t.getPatentDetails(stub, args)
		case "getPatentsByApplicant":
			return t.getPatentsByApplicant(stub, args[1])
		case "getPatentInfoByKey":
			return t.getPatentInfoByKey(stub, args[1])
		case "getDraftPatent":
            return t.getDraftPatent(stub, args)	
        case "getSubmittedPatent":
            return t.getSubmittedPatent(stub, args)	
        case "getRejectedPatent":
            return t.getRejectedPatent(stub, args)	
        case "getSearchReportAttachedPatent":
            return t.getSearchReportAttachedPatent(stub, args)
		case "getPatentHistory":
			return t.getHistoryForPatent(stub, args)
		case "getNotifications":
			return t.getNotifications(stub, args)
		case "getPatentSearchsubmit":
            return t.getPatentSearchsubmit(stub, args)
        case "getPatentSearchapplicant":
            return t.getPatentSearchapplicant(stub, args)
        case "getPatentSearchtitle":
            return t.getPatentSearchtitle(stub, args)   
        case "getPatentSearch":
            return t.getPatentSearch(stub, args)
        case "getPatentSearchcreate":
            return t.getPatentSearchcreate(stub, args)  
        case "getPatentappSearch":
            return t.getPatentappSearch(stub, args)
        case "getPatentpriSearch":
            return t.getPatentpriSearch(stub, args)
        case "getPatentListForPublic":
			return t.getPatentListForPublic(stub, args)    	
		default:
			logger.Error("Invalid function")
			return shim.Error("Invalid function")
		}
	}else{
		logger.Error("Unknown function call")
		return shim.Error("Unknown function call")
	}

	return shim.Success([]byte("Error in chaincode calling"))
}

