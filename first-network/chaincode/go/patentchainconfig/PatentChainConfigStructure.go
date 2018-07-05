package main

import (
	"github.com/hyperledger/fabric/examples/chaincode/go/mylib"
	//"../mylib"
)

//==============================================================================================================================
//	 Patent Chain Structure Definitions
//==============================================================================================================================

type ChainCode struct {
}



type PO struct {
	POName                    string                                `json:"po_name"`
    POAddress                 string                                `json:"po_address"`
    Code                      mylib.PoCode                                `json:"po_code"`
    CountryCode               mylib.CountryCode                           `json:"po_country_code"`
}

//To hold the data when init called
type adminData struct {
    POInfo                      []PO                                `json:"po_info"`
    Role                        []Roles                              `json:"role"`
}

type Administrators struct {
    UserId                    string                              `json:"user_id"`
    Password                  string                              `json:"password"`
}

type Roles struct {
	Role                      string                                `json:"role_name"`
}