package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
 import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)


//==============================================================================================================================
//	 Patent Chain Structure Definitions
//==============================================================================================================================

type SmartContract struct {
}

type ProductList struct{
	Products				[]Product		`json:"productsList"`
}

type Product struct{
	ProductID				string          `json:"product_id"`
	ProductDescription 		string			`json:"product_desc"`
	ProductName				string			`json:"product_name"`
	SellableSize			int			`json:"priceperpackcust"`
	DistributorSize 		int			`json:"priceperpackdistributor"`
}

type ManufactureList struct{
	Manufacturers  			[]Manufacturer	`json:"listManufacturers"`
}

type Manufacturer struct{
	MID						string			`json:"mid"`
	Products				[]Product		`json:"products"`
	UPCEFFECTIVEDATE		string			`json:"upc_date"`
}

type Distributor struct{
	DistributorID 			string 			`json:"distributor_id"`
	UPCEFFECTIVEDATE		string			`json:"upc_effective_date"`
}

type Retailer struct{
	RetailerID				string			`json:"retailer_id"`
	UPCEFFECTIVEDATE		string			`json:"upc_effective_date"`
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

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "addProduct" {
		return s.addProduct(APIstub, args)
	} else if function == "getProduct" {
		return s.getProduct(APIstub, args)
	}
	// } else if function == "queryAllProducts" {
	// 	return s.queryAllProducts(APIstub)
	// } 

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) addProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var product Product;
	product.ProductID = args[0];
	product.ProductDescription = args[1];
	product.ProductName = args[2];
	product.SellableSize = args[3];
	product.DistributorSize = args[4];
	fmt.println("product--->>",product);


	var productList ProductList;
	productlist.Products = append(productlist.Products,product)
	fmt.println("productList--->>",productList);

	productAsBytes, _ := json.Marshal(product)
	APIstub.PutState(args[0], productAsBytes)

	return shim.Success(carAsBytes)
}

func (s *SmartContract) getProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	
	fmt.println("product ID--->>",args[0]);

	productAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(productAsBytes)
}



// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

