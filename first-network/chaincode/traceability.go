package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"encoding/json"
	"fmt"
	"bytes"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)


//==============================================================================================================================
//	 Patent Chain Structure Definitions
//==============================================================================================================================

type SmartContract struct {
}

type Product struct{
	UpcCode				string          		`json:"upc_code"`
	MID				    string					`json:"manufacturer_id"`
	Size                string      		    `json:"size"`
	Quantity            string      		    `json:"quantity"`
	ActionDetails       []ActionDetails   		`json:"acion_details"`
	ProductName			string					`json:"product_name"`
}
type ActionDetails struct{
	//For Phase 3
	UpcCode				    string      `json:"upc_code"`
	UserID				     string		`json:"user_id"`
	//ends
	RecievedDate			string		`json:"recieved_date"`
	RecievedTime			string		`json:"recieved_time"`
	ShippedDate				string		`json:"shipped_date"`
	ShippedTime				string		`json:"shipped_time"`
	QuantityRecieved		string		`json:"quantity_recieved"`
	QuantityShipped		    string		`json:"quantity_shipped"`
	ItemCost                string      `json:"item_cost`
	Status 					string		`json:status`
	TxID                    string      `json:"tx_id"`
	ItemSize				string		`json:"item_size"`
	Shipped_By				string		`json:"shipped_by"`
	Product_Name			string		`json:"product_name"`
}

type ShippingDetails struct{
	ActionDetails       []ActionDetails   		`json:"acion_details"`
}


type ItemList struct{
	Items 				[]Item  	`json:"items"`
}

//For Phase III
type Item struct{
	ProductName     	 string      `json:"productname"`
	MID					 string		 `json:"manufacturer_id"`
	UpcCode				 string      `json:"upc_code"`
	Size				 string 	 `json:"size"`
	Quantity			 string		 `json:"quantity"`
	TxID                 string      `json:"tx_id"`
}

/*   For Phase I and II
 type Item struct{
 	ProductName     	 string      `json:"productname"`
 	MID					 string		 `json:"manufacturer_id"`
 	UpcCode				 string      `json:"upc_code"`
 }
*/

type Name struct{
	ItemName        string     `json:"item_name"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}


/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	_, args := APIstub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	if args[0] == "addProduct" {
		return s.addProduct(APIstub, args)
	}else if args[0] == "addProductItem" {
		return s.addProductItem(APIstub, args)
	}else if args[0] == "actionToProduct" {
		return s.actionToProduct(APIstub, args)
	} else if args[0] == "getProductDetails" {
		return s.getProductDetails(APIstub, args)
	}else if args[0] == "getProductList" {
		return s.getProductList(APIstub, args)
	}else if args[0] == "getProductsByMID" {
		return s.getProductsByMID(APIstub, args)
	}else if args[0] == "changeStatus" {
		return s.changeStatus(APIstub, args)
	}else if args[0] == "addUser" {
		return s.addUser(APIstub, args)
	}else if args[0] == "getUsersByRole" {
		return s.getUsersByRole(APIstub, args)
	}else if args[0] == "changeStatusOrShippingDetails" {
		return s.changeStatusOrShippingDetails(APIstub, args)
	}else if args[0] == "getStatusOrShippingDetails" {
		return s.getStatusOrShippingDetails(APIstub, args)
	}else if args[0] == "getItemsByMID" {
		return s.getItemsByMID(APIstub, args)
	}else if args[0] == "addSubUsers" {
		return s.addSubUsers(APIstub, args)
	}else if args[0] == "getSubUsersById" {
		return s.getSubUsersById(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}
//For Multiple manufacturers for Phase III
func (s *SmartContract) addProductItem(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	productListAsBytes, _ := APIstub.GetState("productList")
	var itemList ItemList

	if len(productListAsBytes) != 0 {
		err := json.Unmarshal(productListAsBytes, &itemList)
		if (err != nil) {
			return shim.Error("Unmashalling Error")
		}
	}

	for i := 0; i < len(itemList.Items); i++ {
		if itemList.Items[i].UpcCode == args[1]{
			return shim.Error("UPC Code already exists")
		}
	}

	var product Item;
	product.UpcCode = args[1];
	product.ProductName = args[2];
	product.Quantity = args[3]
	product.Size = args[4]
	product.MID = args[6]
	product.TxID = APIstub.GetTxID()
	
	itemList.Items = append(itemList.Items, product)

	itemListAsBytes,_ := json.Marshal(itemList)
	APIstub.PutState("productList", itemListAsBytes)

	return shim.Success([]byte("Product Added successfully"))
}

func (s *SmartContract) getItemsByMID(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {

	mid := args[1]
	fmt.Println("list-00-->>>",mid)

	listAsBytes,_ := APIstub.GetState("productList")
	var list ItemList
	if len(listAsBytes) != 0 {
		err := json.Unmarshal(listAsBytes, &list)
		if (err != nil) {
			return shim.Error("Unmashalling Error")
		}
	}
	fmt.Println("list--->>>",list)
	var reqList ItemList
	for i,_ := range list.Items{
		fmt.Println("list11--->>>",list.Items[i].MID)
		if list.Items[i].MID == mid{
			reqList.Items = append(reqList.Items, list.Items[i])
		}
	}
	itemListAsBytes,_ := json.Marshal(reqList)
	return shim.Success(itemListAsBytes)
}

//For Phase 1 and Phase II
func (s *SmartContract) addProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	productListAsBytes, _ := APIstub.GetState(args[1])

	if(productListAsBytes != nil){
		return shim.Error("UPC Code already exists")
	}
	var product Product;
	var actionDetails ActionDetails
	product.UpcCode = args[1];
	product.ProductName = args[2];
	product.Quantity = args[3]
	product.Size = args[4]
	product.MID = args[6]

	err := json.Unmarshal([]byte(args[5]), &actionDetails)
	if (err != nil) {
		return shim.Error("Unmashalling Error")
	}
    
	actionDetails.TxID = APIstub.GetTxID()
	product.ActionDetails = append(product.ActionDetails, actionDetails)

	productAsBytes, _ := json.Marshal(product)
	APIstub.PutState(args[1], productAsBytes)

	productNameAsBytes,_ := APIstub.GetState("productList")

	var item Item
	item.ProductName = args[2];
	item.MID = args[6];
	item.UpcCode=args[1];

	var itemList ItemList

	if len(productNameAsBytes) != 0 {
		err = json.Unmarshal(productNameAsBytes, &itemList)
		if (err != nil) {
			return shim.Error("Unmashalling Error")
		}
	}


	itemList.Items = append(itemList.Items, item)

	itemListAsBytes,_ := json.Marshal(itemList)
	APIstub.PutState("productList", itemListAsBytes)

	return shim.Success([]byte("Product Added successfully"))
}

func (s *SmartContract) getProductDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	productAsBytes, _:= APIstub.GetState(args[1])
	return shim.Success(productAsBytes)
}

func (s *SmartContract) getProductList(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	productNameListAsBytes, _:= APIstub.GetState("productList")
	return shim.Success(productNameListAsBytes)
}

func (s *SmartContract) actionToProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	productAsBytes, _ := APIstub.GetState(args[1])

	var product Product
	err := json.Unmarshal(productAsBytes, &product)
	if (err != nil) {
		return shim.Error("Unmashalling Error")
	}
	var actiondetails ActionDetails
	err = json.Unmarshal([]byte(args[2]), &actiondetails)
	if (err != nil) {
		return shim.Error("Unmashalling Error")
	}

	actiondetails.TxID = APIstub.GetTxID()
	product.ActionDetails = append(product.ActionDetails, actiondetails)
	productNewAsBytes, _ := json.Marshal(product)
	APIstub.PutState(args[1], productNewAsBytes)

	return shim.Success(productNewAsBytes)
}

func (s *SmartContract) getProductsByMID(stub shim.ChaincodeStubInterface, args[] string) sc.Response {

	mid := args[1]

	queryString := fmt.Sprintf("{\"selector\":{\"manufacturer_id\":{\"$eq\": \"%s\"}}}",  mid)
	fmt.Println("queryString",queryString)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	fmt.Println("queryResults",queryResults)
	if err != nil {
		return shim.Error("Unable to fetch records by for this mid")
	}

	return shim.Success(queryResults)
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

func (s *SmartContract) changeStatus(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {
	productAsBytes, _ := APIstub.GetState(args[1])    //UPC Code args[1]

	var product Product
	err := json.Unmarshal(productAsBytes, &product)
	if (err != nil) {
		return shim.Error("Unmashalling Error")
	}

	var actiondetails ActionDetails
	actiondetails.RecievedDate = args[2];
	actiondetails.RecievedTime = args[3];
	actiondetails.ShippedDate = args[4];
	actiondetails.ShippedTime = args[5];
	actiondetails.QuantityRecieved = args[6];
	actiondetails.QuantityShipped = args[7];
	actiondetails.ItemCost = args[8];
	actiondetails.Status = args[9];
	actiondetails.TxID = APIstub.GetTxID()
	actiondetails.ItemSize = args[10];
	actiondetails.Shipped_By = args[11];

	fmt.Println("actiondetails",actiondetails)
	product.ActionDetails = append(product.ActionDetails, actiondetails)
	productNewAsBytes, _ := json.Marshal(product)
	APIstub.PutState(args[1], productNewAsBytes)

	return shim.Success(productNewAsBytes)
}



func (s *SmartContract) changeStatusOrShippingDetails(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {
	shippingDetailsAsBytes, _ := APIstub.GetState("shippingDetails")
	var shippingDetails ShippingDetails
	if len(shippingDetailsAsBytes) != 0 {
		err := json.Unmarshal(shippingDetailsAsBytes, &shippingDetails)
		if (err != nil) {
			return shim.Error("Unmashalling Error")
		}
	}
	var actiondetails ActionDetails
	value,_ := strconv.Atoi(args[6])
	
	if args[9] == "DistributorReceipt" && value < 500 {
		return shim.Error("Please accept atleast 500 quantity ")
	}

	actiondetails.UpcCode = args[1];
	actiondetails.RecievedDate = args[2];
	actiondetails.RecievedTime = args[3];
	actiondetails.ShippedDate = args[4];
	actiondetails.ShippedTime = args[5];
	actiondetails.QuantityRecieved = args[6];
	actiondetails.QuantityShipped = args[7];
	actiondetails.ItemCost = args[8];
	actiondetails.Status = args[9];
	actiondetails.TxID = APIstub.GetTxID()
	actiondetails.ItemSize = args[10];
	actiondetails.UserID = args[11];
	actiondetails.Shipped_By = args[12];
	actiondetails.Product_Name = args[13];

	fmt.Println("actiondetails",actiondetails)
	shippingDetails.ActionDetails = append(shippingDetails.ActionDetails, actiondetails)
	shippingDetailsAsBytesNewAsBytes, _ := json.Marshal(shippingDetails)
	APIstub.PutState("shippingDetails", shippingDetailsAsBytesNewAsBytes)

	return shim.Success(shippingDetailsAsBytesNewAsBytes)
}

func (s *SmartContract) getStatusOrShippingDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	shippingDetailsAsBytes, _:= APIstub.GetState("shippingDetails")
	return shim.Success(shippingDetailsAsBytes)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

