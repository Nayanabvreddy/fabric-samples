package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"encoding/json"
	"fmt"
	_"bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type User struct{
	ID					 string		   		    `json:"id"`
	Name				 string 	 	 		`json:"Name"`
	Role 				 string 				`json:"Role"`
	TxID                 string     			`json:"tx_id"`
	SubUserList			[]SubUser				`json:"sub_users_list"`
}

type SubUser struct{
	ID					 string		   		    `json:"id"`
	Name				 string 	 	 		`json:"Name"`
	Role				 string					`json:"Role"`
}

type SubUserList struct{
	List				[]SubUser				`json:"sub_user_list"`
}

type List struct{
	Userlist				[]User	 		 `json:"user_list"`
}

//For Phase 1 and Phase II
func (s *SmartContract) addUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("addUserr******")
	userAsBytes,_ := APIstub.GetState("userList")
	var user User
	user.Name = args[2];
	user.Role = args[3];
	user.ID=args[1];
	user.TxID = APIstub.GetTxID()

	var list List

	if len(userAsBytes) != 0 {
		err := json.Unmarshal(userAsBytes, &list)
		if (err != nil) {
			return shim.Error("Unmashalling Error")
		}
	}


	list.Userlist = append(list.Userlist, user)
	fmt.Println("list.Userlist******",list.Userlist)

	userListAsBytes,_ := json.Marshal(list)
	APIstub.PutState("userList", userListAsBytes)

	return shim.Success([]byte("User Added successfully"))
}

func (s *SmartContract) getUsersByRole(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {
	role := args[1]
	userAsBytes,_ := APIstub.GetState("userList")
	var list List
	if len(userAsBytes) != 0 {
		err := json.Unmarshal(userAsBytes, &list)
		if (err != nil) {
			return shim.Error("Unmashalling Error")
		}
	}
	var reqList List
	for i,_ := range list.Userlist{
		if list.Userlist[i].Role == role{
			reqList.Userlist = append(reqList.Userlist, list.Userlist[i])
		}
	}
	fmt.Println("reqList.Userlist******",reqList.U
	serlist)

	userListAsBytes,_ := json.Marshal(reqList)

	return shim.Success(userListAsBytes)
}

func (s *SmartContract) addSubUsers(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {
	userId := args[1]
	subUserId := args[2]
	Name := args[3]
	userAsBytes,_ := APIstub.GetState("userList")
	var list List
	if len(userAsBytes) != 0 {
		err := json.Unmarshal(userAsBytes, &list)
		if (err != nil) {
			return shim.Error("Unmashalling Error")
		}
	}

	var subUser SubUser
	subUser.ID = subUserId
	subUser.Name = Name
	subUser.Role = args[4]


	for i,_ := range list.Userlist{
		if list.Userlist[i].ID == userId{
			//Check if user already added
			for j,_ := range list.Userlist[i].SubUserList{
				if(list.Userlist[i].SubUserList[j].ID == subUser.ID){
					return shim.Error("User already added")
				}
			}	
			list.Userlist[i].SubUserList = append(list.Userlist[i].SubUserList,subUser)
			fmt.Println("subuserlist-->>",list.Userlist[i])
		}
	}

	userListAsBytesNew,_ := json.Marshal(list)
	APIstub.PutState("userList", userListAsBytesNew)

	return shim.Success([]byte("Sub User added succesfully"))
}

func (s *SmartContract) getSubUsersById(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {
	userId := args[1]
	userAsBytes,_ := APIstub.GetState("userList")
	var list List
	if len(userAsBytes) != 0 {
		err := json.Unmarshal(userAsBytes, &list)
		if (err != nil) {
			return shim.Error("Unmashalling Error")
		}
	}
	var subUserlist SubUserList
	for i,_ := range list.Userlist{
		if list.Userlist[i].ID == userId{
			subUserlist.List = list.Userlist[i].SubUserList
		}
	}
	fmt.Println("subuserList******",subUserlist)

	userListAsBytes,_ := json.Marshal(subUserlist)

	return shim.Success(userListAsBytes)
}