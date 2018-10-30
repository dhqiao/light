package chain


type Endorsements struct {
	Endorser string
	Signature string
}
type Response struct {
	Message string
	Status int
}


type Extension struct {
	ChaincodeId ChaincodeId
	Response Response
	Results string
}


type ProposalResponsePayload struct {
	Extension Extension
	ProposalHash string
} 

type Action struct {
	Endorsements []Endorsements
	ProposalResponsePayload ProposalResponsePayload
}

type ChaincodeId struct {
	Name string
	Path string
	Version string
}

type ParamsInput struct {
	Args []string
} 

type ChaincodeSpec struct {
	ChaincodeId ChaincodeId
	Input ParamsInput
	Timeout int
	Type string
} 
type PayloadInput struct {
	ChaincodeSpec ChaincodeSpec
} 
type ChaincodeProposalPayload struct {
	Input PayloadInput
} 
type TransctionPayload struct {
	Action Action `json:"action"`
	ChaincodeProposalPayload ChaincodeProposalPayload `json:"chaincode_proposal_payload"`
} 