package mylib


var PatentStatusLookUp = map[string]PatentStatus{
      "Draft" : DRAFT,
      "Submitted" : SUBMITTED,
      "Rejected" : REJECTED,
      "SearchReportAttached" : SEARCHREPORTATTACHED,
}



var requestStatusLookUp = map[string]RequestStatus{
      "Incomplete" : INCOMPLETE,
      "Complete" : COMPLETE,
}

var RoleLookUp = map[string]Role{
      "Applicant" : APPLICANT,
      "FormalityOfficers" : FORMALITYOFFICERS,
      "NPO" : NPO,
      "Public" : PUBLIC,
}

var documentTypeLookUp = map[string]DocumentType{
      "Supporting Document" : SupportingDocument,
      "Request Letter" : RequestLetter,
      "Search Reports" : SearchReports,
      "Rejection Report" : RejectionReport,
      "Opinion of Patentability" : OpinionofPatentability,
}


var poCodeLookUp = map[string]PoCode{
      "IP5" : IP5,
      "EPO" : EPO,
      "NPO" : npo,
}

var countryCodeLookUp = map[string]CountryCode{
      "Japan" :  JP,
      "United States of America" :  US,
      "Albania" :  AL,
      "Austria" :  AT,
      "Belgium" :  BE,
      "Bulgaria" :  BG,
      "Switzerland" :  CH,
      "Cyprus" :  CY,
      "Czech Republic" :  CZ,
      "Germany" :  DE,
      "Denmark" :  DK,
      "Estonia" :  EE,
      "Spain" :  ES,
      "Finland" :  FI,
      "France" :  FR,
      "United Kingdom" :  GB,
      "Greece" :  GR,
      "Croatia" :  HR,
      "Hungary" :  HU,
      "Ireland" :  IE ,
      "Iceland" :  IS,
      "Italy" :  IT,
      "Liechtenstein" :  LI,
      "Lithuania" :  LT,
      "Luxembourg" :  LU,
      "Latvia" :  LV,
      "Monaco" :  MC,
      "Former Yugoslav Republic of Macedonia" :  MK,
      "Malta" :  MT,
      "Netherlands" :  NL,
      "Norway" :  NO,
      "Poland" :  PL,
      "Portugal" :  PT,
      "Romania" :  RO,
      "Sweden" :  SE,
      "Slovenia" :  SI,
      "Slovakia" :  SK,
      "San Marino" :  SM,
      "Turkey" :  TR,
}


//Need to generate a lookup for countrycode constant (enum) as well