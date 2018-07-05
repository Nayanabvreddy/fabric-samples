package mylib

//==============================================================================================================================
//	 Patent status constants
//==============================================================================================================================
type PatentStatus string
const (
	DRAFT  PatentStatus = "Draft"
	SUBMITTED  PatentStatus = "Submitted"
	REJECTED     PatentStatus = "Rejected"
	SEARCHREPORTATTACHED  PatentStatus = "SearchReportAttached"
)

//==============================================================================================================================
//	 Examiner Request Status constants
//==============================================================================================================================
type RequestStatus string
const (
	INCOMPLETE RequestStatus = "Incomplete"
    COMPLETE  RequestStatus = "Complete"
)


//==============================================================================================================================
//	 Patent Role constants
//==============================================================================================================================
type Role string
const (
	APPLICANT Role= "Applicant"
    FORMALITYOFFICERS  Role= "FormalityOfficers"
    NPO       Role = "NPO"
    PUBLIC    Role = "Public"
)

//==============================================================================================================================
//	 Document Type constants
//==============================================================================================================================
type DocumentType string
const (
	SupportingDocument DocumentType="Supporting Document"
	RequestLetter DocumentType="Request Letter"
	SearchReports DocumentType= "Search Reports"
	RejectionReport  DocumentType= "Rejection Report"
	OpinionofPatentability DocumentType="Opinion of Patentability"
)

//==============================================================================================================================
//	 PO code constants
//==============================================================================================================================
type PoCode string
const (
	IP5 PoCode ="IP5"
	EPO PoCode ="EPO"
	npo PoCode ="NPO"
)

//==============================================================================================================================
//  Country codes
//==============================================================================================================================

type CountryCode string
const (
	JP CountryCode= "Japan"
	US CountryCode= "United States of America"
	AL CountryCode= "Albania"
	AT CountryCode= "Austria"
	BE CountryCode= "Belgium"
	BG CountryCode= "Bulgaria"
	CH CountryCode= "Switzerland"
	CY CountryCode= "Cyprus"
	CZ CountryCode= "Czech Republic"
	DE CountryCode= "Germany"
	DK CountryCode= "Denmark"
	EE CountryCode= "Estonia"
	ES CountryCode= "Spain"
	FI CountryCode= "Finland"
	FR CountryCode= "France"
	GB CountryCode= "United Kingdom"
	GR CountryCode= "Greece"
	HR CountryCode= "Croatia"
	HU CountryCode= "Hungary"
	IE  CountryCode= "Ireland"
	IS CountryCode= "Iceland"
	IT CountryCode= "Italy"
	LI CountryCode= "Liechtenstein"
	LT CountryCode= "Lithuania"
	LU CountryCode= "Luxembourg"
	LV CountryCode= "Latvia"
	MC CountryCode= "Monaco"
	MK CountryCode= "Former Yugoslav Republic of Macedonia"
	MT CountryCode= "Malta"
	NL CountryCode= "Netherlands"
	NO CountryCode= "Norway"
	PL CountryCode= "Poland"
	PT CountryCode= "Portugal"
	RO CountryCode= "Romania"
	SE CountryCode= "Sweden"
	SI CountryCode= "Slovenia"
	SK CountryCode= "Slovakia"
	SM CountryCode= "San Marino"
	TR CountryCode= "Turkey"
)

//==============================================================================================================================
//  Config chain code
//==============================================================================================================================

const (
	ConfigChainCode = "conf2"
	ConfigChannel= "mychannel"
	ConfigQuery= "query"
	ConfigInvoke= "invoke"
)


//==============================================================================================================================
//	 Common chaincode error messages
//==============================================================================================================================
const (
	ArgumentErrorMessage = "Insufficient number of arguments"
	PutErrorMessage  = "Error in inserting data into ledger"
	GetStateErrorMessage  = "Unable to retrieve key value"
	DeleteStateErrorMessage  = "Unable to delete key value"
	MarshalErrorMessage  = "Unable to marshal json data"
	UnmarshalErrorMessage  = "Unable to unmarshal data"
	ConversionErrorMessage  = "Unable to parse data"
	SearchReportDescription = "Search Reaport for this Application has been uploaded"
	RejectReportDescription = "Patent application has been rejected Please find attached Document"
	OpinionOfPatentabilityDescription = "Opinion Of Patentability report for this Application has been uploaded"
	RequestLetterDescription = "Some document are missig in your Patent Application Please provide ASAP"

)

