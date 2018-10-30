
package main

import (
    "fmt"
    "github.com/hyperledger/fabric/common/flogging"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
    "encoding/json"
    "reflect"
    //"strconv"
)

// logger
var chaincodeLogger = flogging.MustGetLogger("BulkchainChaincode")


//审核状态
const (
	Check_State_Checking = "Checking"
	Check_State_CheckedResolved = "Resolved"
	Check_State_CheckedRejected = "Rejected"
    Check_State_Finished = "Finished"
)

type Member struct{
	MemberId string `json:MemberId`
	MemberName string `json:MemberName`
}

type Client struct{
	ClientId string `json:ClientId`
	ClientName string `json:ClientName`
}

//重量单位 以克为最小单位，UnitValue即克数
type Unit struct{
	UnitName string `json:UnitName`
	UnitValue int `json:UnitValue`
}

//单位

var Unit_Gram Unit= Unit{"克",1}
var Unit_Kg Unit = Unit{"千克",1000}
var Unit_Ton Unit= Unit{"吨",1000000}

//数量
type Weight struct{
	UnitType Unit `json:UnitType`
	Value int `json:Value`
}

// 品种信息
type VarietyInfo struct{
	VarietyCode string `json:VarietyCode`
	VarietyName string `json:VarietyName`
	TradeUnit Weight `json:TradeUnit`
	MinDeliveryHands int `json:MinDeliveryHands`
}

// 上市品种信息
var Variety_WH VarietyInfo = VarietyInfo{"WH","优质强筋小麦",Weight{Unit_Ton,20},5}
var Variety_CF VarietyInfo = VarietyInfo{"CF","棉花",Weight{Unit_Ton,5},8}
var Variety_SR VarietyInfo = VarietyInfo{"SR","白糖",Weight{Unit_Ton,10},1}
var Variety_OI VarietyInfo = VarietyInfo{"OI","菜籽油",Weight{Unit_Ton,10},1}
var Variety_RI VarietyInfo = VarietyInfo{"RI","早籼稻",Weight{Unit_Ton,20},1}
var Variety_PM VarietyInfo = VarietyInfo{"PM","普麦",Weight{Unit_Ton,50},1}
var Variety_RS VarietyInfo = VarietyInfo{"RS","油菜籽",Weight{Unit_Ton,10},1}
var Variety_JR VarietyInfo = VarietyInfo{"JR","粳稻",Weight{Unit_Ton,20},50}


var varieties map[string]VarietyInfo = make(map[string]VarietyInfo)
var requestsIndexNames map[string]string = make(map[string]string)
var warehouseReceiptIndexNames map[string]string = make(map[string]string)


type Goods struct{
	VarietyCode string `json:VarietyCode`
	Quantity int `json:Quantity`
	Quality string `json:Quality`
	Brand string `json:Brand`
	GoodsPackage string `json:GoodsPackage`
	GoodsSpecification string `json:GoodsSpecifications`
	ProductionPlace string `json:ProductionPlace`
	ProductionDate string `json:ProductionDate`
	ValidDate string `json:ValidDate`  //eg "2018-08-10"
}


//WarehouseReceipt Type
const (
	WarehouseReceipt_Type_Standard = "Standard"
	WarehouseReceipt_Type_NonStandard = "NonStandard"
)
//WarehouseReceipt State
const (
	//warehouse receipt state
	State_INBOUND = "Inbound"
	State_FLOWABLE = "Flowable"
	State_PLEDGED = "Pledged"
	State_OUTBOUNDREADY = "OutboundReady"
	State_OUTBOUNDED = "Outbounded"
	State_OUTBOUNDING = "Outbounding"
	State_REGISTERING = "Registering"
	State_PLEDGING = "Pledging"
	State_UNPLEDGING = "Unpledging"
	State_UNREGISTERING = "Unregistering"
	State_UNREGISTERED = "Unregistered"
	State_DELIVERYING = "Deliverying"
	//match state
	MatchState_Matching = "Matching"
	MatchState_Matched = "Matched"
	MatchState_Unmatched = "Unmatched"
	//confirm state
	ConfirmState_Confirming = "Confirming"
	ConfirmState_ConfirmResolved = "ConfirmResolved"
	ConfirmState_ConfirmRejected = "ConfirmRejected"
)

type Transaction struct {
	TxId string `json:TxId`
	TxType string `json:TxType`
	Content []byte `json:Info`
}

type Place struct {
	Address string `json:Address`
	Location string `json:Location`
}

type Period struct{
	StartDate string `json:StartDate`
	EndDate string `json:EndDate`
}
//交易类型
const (
	TxType_All = "*"
	VarietyType_All = "*"
	TxType_InboundRequest = "InboundRequest"
	TxType_RegisterRequest = "RegisterRequest"
	TxType_PledgeRequest = "PledgeRequest"
	TxType_UnpledgeRequest = "UnpledgeRequest"
	TxType_DeliveryRequest = "DeliveryRequest"
	TxType_UnregisterRequest = "UnregisterRequest"
	TxType_OutboundRequest = "OutboundRequest"

	DeliveryType_Seller = "Seller"
	DeliveryType_Buyer = "Buyer"

	PledgeType_Inside = "Inside"
	PledgeType_Outside = "Outside"
)



type WarehouseReceipt struct{
	WarehouseReceiptSeriesId string `json:WarehouseReceiptSeriesId`
	Quantity int `json:Quantity`
	BeginId int `json:BeginId`
	EndId int `json:ENdId`
	Type string `json:Type`
	State string `json:State`
	VarietyCode string `json:VarietyCode`
	Quality string `json:Quality` 
	GoodsQuantity int `json:GoodsQuantity`
	Signer string `json:Signer`
	SignPlace string `json:SignPlace`
	SignDate string `json:SignDate`
	GoodsHolder string `json:GoodsHolder`
	GoodsHolderId string `json:GoodsHolderId`
	StoragePeriod Period `json:StorgaePeriod`
	StoragePlace Place `json:StoragePlace`
	GoodsSaver string `json:GoodsSaver`
	GoodsSaverId string `json:GoodsSaverId`
	MemberName string `json:MemberName`
	MemberId string `json:MemberId`
	WarehouseReceiptHolder string `json:WarehouseReceiptHolder`
	WarehouseReceiptHolderId string `json:WarehouseReceiptHolderId`
	TransactionHistory []string `json:TransactionHistory`
	History []WarehouseReceipt `json:History`
}

type InboundRequest struct{
	//BasicInfo
	TransactionId string `json:TransactionId`
	TxType string `json:TxType`
	MemberId string `json:MemberId`
	MemberName string `json:MemberName`
	DateRequest string `json:DateRequest`
	MemberContact string `json:MemberContact`
	MemberContactPhoneNumber string `json:MemberContactPhoneNumber`
	ClientId string `json:ClientId`
	ClientName string `json:ClientName`
	ClientContact string `json:ClientContact`
	ClientContactPhoneNumber string `json:ClientContactPhoneNumber`
	GoodsListRequested []Goods `json:GoodsListRequested`
	ModeOfTransport string `json:ModeOfTransport`
	TargetWarehouseId string `json:TargetWarehouseId`
	TargetWarehouseName string `json:TargetWarehouseName`
	DateInPlan string `json:DateInPlan`
	//response by the warehouse 
	DateCheck string `json:DateCheck`
	CheckState string `json:CheckState`
	GoodsListPermitted []Goods `json:GoodsListPermitted`
	DatePermitted string `json:DatePermitted`
	Description string `json:Description`
	//response
	GoodsListIndeed []Goods `json:GoodsListIndeed`
	DateIndeed string `json:DateIndeed`
	WarehouseReceipts []WarehouseReceipt `json:WarehouseReceipts`
	DateCreate string `json:DateCreate`
}


type RegisterRequest struct {
	//BasicInfo
	TransactionId string `json:TransactionId`
	TxType string `json:TxType`
	MemberId string `json:MemberId`
	MemberName string `json:MemberName`
	DateRequest string `json:DateRequest`
	MemberContact string `json:MemberContact`
	MemberContactPhoneNumber string `json:MemberContactPhoneNumber`
	ClientId string `json:ClientId`
	ClientName string `json:ClientName`
	ClientContact string `json:ClientContact`
	ClientContactPhoneNumber string `json:ClientContactPhoneNumber`
	RegisteringSeriesId string `json:RegisteringSeriesId`
	RegisteringQuantity int `json:RegisteringQuantity`
	//response by the exchange
	DateCheck string `json:DateCheck`
	CheckState string `json:CheckState`
	Description string `json:Description`
	//supplementary info
	RegisteringWarehouseReceipt WarehouseReceipt  `json:RegisteringWarehouseReceipt`
	RegisteredWarehouseReceipt WarehouseReceipt `json:RegisteredWarehouseReceipt`
}


type DeliveryRequest struct{
	//BasicInfo
	TransactionId string `json:TransactionId`
	TxType string `json:TxType`
	MemberId string `json:MemberId`
	MemberName string `json:MemberName`
	DateRequest string `json:DateRequest`
	MemberContact string `json:MemberContact`
	MemberContactPhoneNumber string `json:MemberContactPhoneNumber`
	ClientId string `json:ClientId`
	ClientName string `json:ClientName`
	ClientContact string `json:ClientContact`
	ClientContactPhoneNumber string `json:ClientContactPhoneNumber`
	DeliveryType string `json:DeliveryType`
	//seller
	DeliveryWarehouseReceiptSeriesId string `json:DeliveryWarehouseReceiptSeriesId`
	//buyer
	DeliveryVarietyCode string `json:DeliveryVarietyCode`
	DeliveryQuantity int `json:DeliveryQuantity`
	//supplementary info
	TradeUnit int `json:TradeUnit`
	//response by the exchange
	DateCheck string `json:DateCheck`
	CheckState string `json:CheckState`
	Description string `json:Description`
	//match info
	MatchState string `json:MatchState`
	DateMatch string `json:DateMatch`
	MatchClientId string `json:MatchClientId`
	MatchClientName string `json:MatchClientName`
	MatchMemberId string `json:MatchMemberId`
	MatchMemberName string `json:MatchMemberName`
	//Match info 
	MatchTxId string `json:MatchTxId`
	//buyer
	MatchQuantity int `json:MatchQuantity`
	MatchWarehouseReceipt WarehouseReceipt `json:MatchWarehouseReceipts`
	//confirm match result
	ConfirmState string `json:ConfirmState`
}


type PledgeRequest struct {
	//BasicInfo
	TransactionId string `json:TransactionId`
	TxType string `json:TxType`
	MemberId string `json:MemberId`
	MemberName string `json:MemberName`
	DateRequest string `json:DateRequest`
	MemberContact string `json:MemberContact`
	MemberContactPhoneNumber string `json:MemberContactPhoneNumber`
	ClientId string `json:ClientId`
	ClientName string `json:ClientName`
	ClientContact string `json:ClientContact`
	ClientContactPhoneNumber string `json:ClientContactPhoneNumber`
	PledgeType string `json:PledgeType`
	TargetBankId string `json:TargetBankId`
	TargetBankName string `json:TargetBankName`
	PledgingWarehouseReceiptSeriesId string `json:PledgingWarehouseReceiptSeriesId`
	PledgingQuantity int `json:PledgingQuantity`
	AmountOfMoneyRequest int `json:AmountOfMoneyRequest`
	DateDDL string `json:DateDDL`
	//supplementary info
	PledgingWarehouseReceipt WarehouseReceipt `json:PledgingWarehouseReceipt`
	//response by the exchange or bank
	DateCheck string `json:DateCheck`
	CheckState string `json:CheckState`
	AmountOfMoneyLended int `json:AmountOfMoneyIndeed`
	AmountOfMoneyReturning int `json:AmountOfMoneyReturning`
	Description string `json:Description`
	//response by the client
	ConfirmState string `json:ConfirmState`
	//unpledge request from the client
	UnpledgeRequestDate string `json:UnpledgeRequestDate`
	//response from the exchange or the bank
	DateCheckUnpledge string `json:DateCheckUnpledge`
	CheckStateUnpledge string `json:CheckStateUnpledge`
	AmountOfMoneyReturned int `json:AmountOfMoneyReturned`
	DescriptionUnpledge string `json:DescriptionUnpledge`
}


type UnregisterRequest struct {
	//basic info
	TransactionId string `json:TransactionId`
	TxType string `json:TxType`
	MemberId string `json:MemberId`
	MemberName string `json:MemberName`
	DateRequest string `json:DateRequest`
	MemberContact string `json:MemberContact`
	MemberContactPhoneNumber string `json:MemberContactPhoneNumber`
	ClientId string `json:ClientId`
	ClientName string `json:ClientName`
	ClientContact string `json:ClientContact`
	ClientContactPhoneNumber string `json:ClientContactPhoneNumber`
	UnregisteringSeriesId string `json:UnregisteringSeriesId`
	UnregisteringQuantity int `json:UnregisteringQuantity`
	//response by the exchange
	DateCheck string `json:DateCheck`
	CheckState string `json:CheckState`
	Description string `json:Description`
	//supplementary info
	UnregisteringWarehouseReceipt WarehouseReceipt `json:UnregisteringWarehouseReceipt`
	UnregisteredWarehouseReceipt WarehouseReceipt `json:UnregisteredWarehouseReceipt`
}


type OutboundRequest struct {
	//basic info
	TransactionId string `json:TransactionId`
	TxType string `json:TxType`
	MemberId string `json:MemberId`
	MemberName string `json:MemberName`
	DateRequest string `json:DateRequest`
	MemberContact string `json:MemberContact`
	MemberContactPhoneNumber string `json:MemberContactPhoneNumber`
	ClientId string `json:ClientId`
	ClientName string `json:ClientName`
	ClientContact string `json:ClientContact`
	ClientContactPhoneNumber string `json:ClientContactPhoneNumber`
	OutboundingSeriesId string `json:OutboundingSeriesId`
	DateInPlan string `json:DateInPlan`
	//response by the warehouse
	DateCheck string `json:DateCheck`
	CheckState string `json:CheckState`
	DatePermitted string `json:DatePermitted`
	Description string `json:Description`
	//response by the client
	GoodsIndeed Goods `json:GoodsListIndeed`
	DateIndeed string `json:DateIndeed`
} 

type MatchingInfo struct {
	VarietyCode string `json:VarietyCode`
	BuyerInfo map[string]int `json:BuyerInfo`
	SellerInfo map[string]int `json:SellerInfo`
	BuyerMatchState map[string]bool `json:BuyerMatchState`
	SellerMatchState map[string]bool `json:SellerMatchState`
	BuyerMatchResult map[string]string `json:BuyerMatchResult`
	SellerMatchResult map[string]string `json:SellerMatchResult`
}


//chaincode
type mychaincode struct{

}

const (
	Role_Member = "Member"
	Role_Client = "Client"
	Role_Warehouse = "Warehouse"
	Role_Bank = "Bank"
	Role_Exchange = "Exchange"
)

const (
	Index_AllMembers = "Role~MemberId"
	Index_AllClients = "Role~ClientId"
	Index_AllWarehouses = "Role~WarehouseId"
	Index_AllBanks = "Role~Bank"

	Index_MemberId_ClientId = "MemberId~ClientId"

	Index_AllRequests = "TxType~TxId"
	Index_Member_TxId = "MemberId~TxId"
	Index_Client_TxId = "ClientId~TxId"
	Index_Warehouse_TxId = "WarehouseId~TxId" 
	Index_Exchange_TxId = "ExchangeId~TxId"
	Index_Bank_TxId = "BankId~TxId"

	Index_AllWarehouseReceipts = "All~WrId"
	Index_Member_WrId = "MemberId~WrId"
	Index_Client_WrId = "ClientId~WrId"
	Index_Bank_WrId = "BankId~WrId"
	Index_Exchange_WrId = "Exchange~WrId"
	Index_Warehouse_WrId = "Warehouse~WrId"
	Key_AllWarehouseReceipts = "AllWarehouseReceipts"
	Key_Exchange = "TheExchange"
	//for matching 
	Index_DeliveryType_VarietyCode_TxId = "DeliveryType-VarietyCode~TxId"
)


//prefix
const (
	Request_Prefix = "Request_"
	WarehouseReceipt_Prefix = "WarehouseReceipt_"
)

// chaincode response结构
type chaincodeRet struct {
    Code int // 0 success otherwise 1
    Des  string //description
}

// response message format
func getRetByte(code int,des string) []byte {
    var r chaincodeRet
    r.Code = code
    r.Des = des

    b,err := json.Marshal(r)

    if err!=nil {
        fmt.Println("marshal Ret failed")
        return nil
    }
    return b
}

// response message format
func getRetString(code int,des string) string {
    var r chaincodeRet
    r.Code = code
    r.Des = des

    b,err := json.Marshal(r)

    if err!=nil {
        fmt.Println("marshal Ret failed")
        return ""
    }
	chaincodeLogger.Infof("%s",string(b[:]))
    return string(b[:])
}

//check if the user type is valid
func isUserTypeValid(userType string) bool {
	types := [5]string{Role_Bank,Role_Warehouse,Role_Client,Role_Member,Role_Exchange}
	for _,v := range types {
		if v == userType {
			return true
		}
	}
	return false
}

func isCheckStateValid(state string) bool {
	states := [4]string{Check_State_Checking,Check_State_CheckedRejected,Check_State_CheckedResolved,Check_State_Finished}
	for _,v := range states {
		if v == state {
			return true
		}
	}
	return false
}


func isPledgeTypeValid(ptype string) bool {
	types := [2]string{PledgeType_Outside,PledgeType_Inside}
	for _,v := range types {
		if v == ptype {
			return true
		}
	}
	return false
}


//args:Transaction
func (a *mychaincode) putRequest(stub shim.ChaincodeStubInterface,tx Transaction) ([]byte, bool) {
	bTx,err := json.Marshal(tx)
	if err!=nil {
		return nil, false
	}

	err = stub.PutState(Request_Prefix + tx.TxId, bTx)
	if err!=nil {
		return nil, false
	}
	return bTx, true
}


//args: 0 -TransactionId
func (a *mychaincode) getRequest(stub shim.ChaincodeStubInterface,TxId string) (Transaction, bool) {
	var tx Transaction
	key := Request_Prefix + TxId
	bTx,err := stub.GetState(key)
	if bTx == nil {
		return tx, false
	}
	err = json.Unmarshal(bTx,&tx)
	if err != nil {
		return tx, false
	}
	return tx, true
}


//args: WarehouseReceiptSeriesId
func (a *mychaincode) getWarehouseReceipt(stub shim.ChaincodeStubInterface, WrId string) (WarehouseReceipt,bool) {
	var wr WarehouseReceipt
	key := WarehouseReceipt_Prefix + WrId
	bWr,err := stub.GetState(key)
	if bWr == nil {
		return wr,false
	}
	err = json.Unmarshal(bWr,&wr)
	if err != nil {
		return wr,false
	}
	return wr,true
}


//args:WarehouseReceipt
func (a *mychaincode) putWarehouseReceipt(stub shim.ChaincodeStubInterface, wr WarehouseReceipt) ([]byte,bool) {
	bWr,err := json.Marshal(wr)
	if err != nil {
		return nil,false
	}
	err = stub.PutState(WarehouseReceipt_Prefix+wr.WarehouseReceiptSeriesId,bWr)
	if err != nil {
		return nil,false
	}
	return bWr,true
}

//submit inbound request
//args: 0 - {InboundRequest object} 
func (a *mychaincode) sendInboundRequest(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args)!=1 {
		res := getRetString(1,"BulkchainChaincode Invoke sendInboundRequest args!=1")
		return shim.Error(res)
	}

	var request InboundRequest
	err := json.Unmarshal([]byte(args[0]), &request)
	if err!=nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendInboundRequest unmarshal request failed")
		return shim.Error(res)
	}
	// TODO 根据票号 查找是否票号已存在
	// TODO 如stat中已有同号票据 返回error message

	request.TransactionId = stub.GetTxID()
	request.CheckState = Check_State_Checking
	request.TxType = TxType_InboundRequest

	//put search tables
	MemberIdTxIdKey, err := stub.CreateCompositeKey(Index_Member_TxId, []string{request.MemberId, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendInboundRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(MemberIdTxIdKey, []byte(TxType_InboundRequest))

	ClientIdTxIdKey, err := stub.CreateCompositeKey(Index_Client_TxId, []string{request.ClientId, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendInboundRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(ClientIdTxIdKey, []byte(TxType_InboundRequest))

	WarehouseIdTxIdKey, err := stub.CreateCompositeKey(Index_Warehouse_TxId, []string{request.TargetWarehouseId, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendInboundRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(WarehouseIdTxIdKey, []byte(TxType_InboundRequest))

	AllRequestsKey, err := stub.CreateCompositeKey(Index_AllRequests, []string{TxType_InboundRequest, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendInboundRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(AllRequestsKey, []byte{0x00})


	// 更改票据信息和状态并保存票据:票据状态设为新发布
	
    // 保存票据
    bReq,err := json.Marshal(request)
    if err != nil {
    	res := getRetString(1,"BulkchainChaincode Invoke sendInboundRequest marshal request error")
    	return shim.Error(res)
    }
    var tx Transaction
    tx.TxId = request.TransactionId
    tx.TxType = TxType_InboundRequest
    tx.Content = bReq
	_, bTx := a.putRequest(stub, tx)
	if !bTx {
		res := getRetString(1,"BulkchainChaincode Invoke sendInboundRequest put InboundRequest failed")
		return shim.Error(res)
	}

	res := getRetByte(0,"invoke sendInboundRequest success")
	return shim.Success(res)
}


//check InboundRequest by warehouse
//args: 0 - InboundRequest
func (a *mychaincode) checkInboundRequest(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	
	//TODO: check the creator  

	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode Invoke checkInboundRequest args!=1")
		return shim.Error(res)
	}
	var checkResult InboundRequest
	err := json.Unmarshal([]byte(args[0]), &checkResult)
	if err!=nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkInboundRequest unmarshal failed")
		return shim.Error(res)
	}


	tx,existReq := a.getRequest(stub,checkResult.TransactionId)
	if !existReq {
		res := getRetString(1,"BulkchainChaincode Invoke checkInboundRequest request not exist")
		return shim.Error(res)
	}
	var request InboundRequest
	err = json.Unmarshal(tx.Content,&request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkInboundRequest unmarshal request error")
		return shim.Error(res)
	}


	//validate
	if !isCheckStateValid(request.CheckState) {
		res := getRetString(1,"BulkchainChaincode Invoke checkInboundRequest wrong CheckState")
		return shim.Error(res)
	}
	if request.CheckState == Check_State_CheckedRejected || request.CheckState == Check_State_CheckedResolved {
		res := getRetString(1,"BulkchainChaincode Invoke checkInboundRequest request should not be checked multi times")
		return shim.Error(res)
	}


	request.DateCheck = checkResult.DateCheck
	request.CheckState = checkResult.CheckState
	request.GoodsListPermitted = checkResult.GoodsListPermitted
	request.DatePermitted = checkResult.DatePermitted
	request.Description = checkResult.Description

	//check the validity of CheckState
	if !isCheckStateValid(request.CheckState) {
		res := getRetString(1,"BulkchainChaincode Invoke checkInboundRequest CheckState is invalid")
		return shim.Error(res)
	}

	//TODO:check the goods list

	//TODO:check the Date

	//check the reality of data
	if !reflect.DeepEqual(request,checkResult) {
		res := getRetString(1,"BulkchainChaincode Invoke checkInboundRequest request info was modified unexpectedly")
		return shim.Error(res)
	} 	

	//put the request 
	bReq,err := json.Marshal(request)
    if err != nil {
    	res := getRetString(1,"BulkchainChaincode Invoke checkInboundRequest marshal request error")
    	return shim.Error(res)
    }
    tx.Content = bReq
	_, bTx := a.putRequest(stub, tx)
	if !bTx {
		res := getRetString(1,"BulkchainChaincode Invoke checkInboundRequest put InboundRequest failed")
		return shim.Error(res)
	}

	res := getRetByte(0,"invoke checkInboundRequest success")
	return shim.Success(res)
}

//invoke 
//RegisterInbound
//args: 0 - InboundRequest
//args: 1 - Signer
//args: 2 - SignPlace
//args: 3 - SignDate
//args: 4 - StartDate
//args: 5 - EndDate
//args: 6 - Address
//args: 7 - Location
func (a *mychaincode) registerInbound(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 8 {
		res := getRetString(1,"BulkchainChaincode registerInbound args != 8")
		return shim.Error(res)
	}

	var registerResult InboundRequest
	err := json.Unmarshal([]byte(args[0]), &registerResult)
	if err!=nil {
		res := getRetString(1,"BulkchainChaincode Invoke registerInbound unmarshal failed")
		return shim.Error(res)
	}

	//check existence
	tx,existReq := a.getRequest(stub,registerResult.TransactionId)
	if !existReq {
		res := getRetString(1,"BulkchainChaincode Invoke registerInbound request not exist")
		return shim.Error(res)
	}
	//get the request
	var request InboundRequest
	err = json.Unmarshal(tx.Content,&request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke registerInbound unmarshal request error")
		return shim.Error(res)
	}
	//check checkState
	if request.CheckState == Check_State_CheckedRejected {
		res := getRetString(1,"BulkchainChaincode Invoke registerInbound request has been rejected")
		return shim.Error(res)
	}


	//TODO:valid the register result
	if registerResult.DateIndeed == "" || registerResult.DateCreate == "" {
		res := getRetString(1,"BulkchainChaincode Invoke registerInbound DateIndeed or DateCreate should not be empty")
		return shim.Error(res)
	} 
	if registerResult.GoodsListIndeed == nil {
		res := getRetString(1,"BulkchainChaincode Invoke registerInbound GoodsListIndeed should not be null")
		return shim.Error(res)
	}
	if request.GoodsListIndeed != nil && len(request.GoodsListIndeed) > 0 {
		res := getRetString(1,"BulkchainChaincode Invoke registerInbound request has been registered")
		return shim.Error(res)
	}


	request.GoodsListIndeed = registerResult.GoodsListIndeed
	request.DateIndeed = registerResult.DateIndeed
	request.WarehouseReceipts = registerResult.WarehouseReceipts
	request.DateCreate = registerResult.DateCreate

	if !reflect.DeepEqual(request,registerResult) {
		res := getRetString(1,"BulkchainChaincode Invoke registerInbound request info was modified unexpectedly")
		return shim.Error(res)
	}


	


	//generate warehouse receipts
	bReq,err := json.Marshal(request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke registerInbound marshal request error")
		return shim.Error(res)
	}
	tx.Content = bReq
	var wrs []WarehouseReceipt
	signer := args[1]
	signPlace := args[2]
	signDate := args[3]
	storagePeriod := Period{args[4],args[5]}
	storagePlace := Place{args[6],args[7]}

	for _,v := range request.GoodsListIndeed {
		var wr WarehouseReceipt
		wr.WarehouseReceiptSeriesId = stub.GetTxID()
		if _,ok := varieties[v.VarietyCode];!ok {
			res := getRetString(1,"BulkchainChaincode Invoke registerInbound wrong variety type")
			return shim.Error(res)
		}
		wr.Quantity = v.Quantity / varieties[v.VarietyCode].MinDeliveryHands
		wr.BeginId = 0
		wr.EndId = wr.Quantity - 1
		wr.Type = WarehouseReceipt_Type_NonStandard
		wr.State = State_INBOUND
		wr.VarietyCode = v.VarietyCode
		wr.Quality = v.Quality
		wr.GoodsQuantity = varieties[v.VarietyCode].MinDeliveryHands
		wr.Signer = signer
		wr.SignPlace = signPlace
		wr.SignDate = signDate
		wr.GoodsHolder = request.TargetWarehouseName
		wr.GoodsHolderId = request.TargetWarehouseId
		wr.StoragePeriod = storagePeriod
		wr.StoragePlace = storagePlace
		wr.GoodsSaver = request.ClientName
		wr.GoodsSaverId = request.ClientId
		wr.MemberName = request.MemberName
		wr.MemberId = request.MemberId
		wr.WarehouseReceiptHolder = request.ClientName
		wr.WarehouseReceiptHolderId = request.ClientId
		wr.TransactionHistory = append(wr.TransactionHistory,tx.TxId)

		wrs =  append(wrs,wr)
		_,bWr := a.putWarehouseReceipt(stub,wr)
		if !bWr {
			res := getRetString(1,"BulkchainChaincode Invoke registerInbound put WarehouseReceipt failed")
			return shim.Error(res)
		}

		//put search tables
		MemberIdWrIdKey,err := stub.CreateCompositeKey(Index_Member_WrId, []string{WarehouseReceipt_Prefix + request.MemberId, wr.WarehouseReceiptSeriesId})
		if err != nil {
			res := getRetString(1,"BulkchainChaincode Invoke registerInbound put search table Member_WrId failed")
			return shim.Error(res)
		}
		stub.PutState(MemberIdWrIdKey, []byte(wr.VarietyCode))

		ClientIdWrIdKey,err := stub.CreateCompositeKey(Index_Client_WrId,[]string{WarehouseReceipt_Prefix + request.ClientId, wr.WarehouseReceiptSeriesId})
		if err != nil {
			res := getRetString(1,"BulkchainChaincode Invoke registerInbound put search table Client_WrId failed")
			return shim.Error(res)
		}
		stub.PutState(ClientIdWrIdKey, []byte(wr.VarietyCode))

		WarehouseIdWrIdKey, err := stub.CreateCompositeKey(Index_Warehouse_WrId, []string{WarehouseReceipt_Prefix + request.TargetWarehouseId, wr.WarehouseReceiptSeriesId})
		if err != nil {
			res := getRetString(1,"BulkchainChaincode Invoke registerInbound put search table Warehouse_WrId failed")
			return shim.Error(res)
		}
		stub.PutState(WarehouseIdWrIdKey, []byte(wr.VarietyCode))

		AllWrIdKey,err := stub.CreateCompositeKey(Index_AllWarehouseReceipts, []string{Key_AllWarehouseReceipts, wr.WarehouseReceiptSeriesId})
		if err != nil {
			res := getRetString(1,"BulkchainChaincode Invoke registerInbound put search table AllWarehouseReceipts failed")
			return shim.Error(res)
		}
		stub.PutState(AllWrIdKey, []byte(wr.VarietyCode))
	}
		
	request.WarehouseReceipts = wrs
    request.CheckState = Check_State_Finished
	bReq,err = json.Marshal(request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke registerInbound marshal request error")
		return shim.Error(res)
	}
	tx.Content = bReq
	_, bTx := a.putRequest(stub, tx)
	if !bTx {
		res := getRetString(1,"BulkchainChaincode Invoke registerInbound put InboundRequest failed")
		return shim.Error(res)
	}

	res :=  getRetByte(0,"BulkchainChaincode Invoke registerInbound success")
	return shim.Success(res)

}


//send Register Request to the exchange
//args: 0 - RegisterRequest
func (a *mychaincode) sendRegisterRequest(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode Invoke sendRegisterRequest args != 1")
		return shim.Error(res)
	}

	var request RegisterRequest
	err := json.Unmarshal([]byte(args[0]),&request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendRegisterRequest unmarshal RegisterRequest failed")
		return shim.Error(res)
	}
	request.TransactionId = stub.GetTxID()

	//TODO: validate
	//check the registering quantity
	/*
	if request.RegisteringQuantity <= 0 {
		res := getRetString(1,"BulkchainChaincode Invoke sendRegisterRequest RegisteringQuantity <= 0")
		return shim.Error(res)
	}
	*/
	wr,bWr := a.getWarehouseReceipt(stub,request.RegisteringSeriesId)
	if !bWr {
		res := getRetString(1,"BulkchainChaincode Invoke sendRegisterRequest the registering warehousereceipt does not exist")
		return shim.Error(res)
	}
	//check the quantity
	/*
	if request.RegisteringQuantity > wr.Quantity {
		res := getRetString(1,"BulkchainChaincode Invoke sendRegisterRequest RegisteringQuantity > maxQuantity")
		return shim.Error(res)
	}
	*/
	//check the holder of wr
	if wr.WarehouseReceiptHolderId != request.ClientId || wr.MemberId != request.MemberId {
		res := getRetString(1,"BulkchainChaincode Invoke sendRegisterRequest the client is not the holder of the WarehouseReceipt")
		return shim.Error(res)
	}
	//check the wr type and state
	if wr.Type != WarehouseReceipt_Type_NonStandard || wr.State != State_INBOUND {
		res := getRetString(1,"BulkchainChaincode Invoke sendRegisterRequest the operation is illegal because of the state and type of WarehouseReceipt ")
		return shim.Error(res)
	} 


	//put search tables
	MemberIdTxIdKey, err := stub.CreateCompositeKey(Index_Member_TxId, []string{request.MemberId, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendRegisterRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(MemberIdTxIdKey, []byte(TxType_RegisterRequest))

	ClientIdTxIdKey, err := stub.CreateCompositeKey(Index_Client_TxId, []string{request.ClientId, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendRegisterRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(ClientIdTxIdKey, []byte(TxType_RegisterRequest))

	ExchangeTxIdKey, err := stub.CreateCompositeKey(Index_Exchange_TxId,[]string{Key_Exchange, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendRegisterRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(ExchangeTxIdKey, []byte(TxType_RegisterRequest))

	AllRequestsKey, err := stub.CreateCompositeKey(Index_AllRequests, []string{TxType_RegisterRequest, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendRegisterRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(AllRequestsKey, []byte{0x00})

	//set wr state and record history
	wr.State = State_REGISTERING
	wr.TransactionHistory = append(wr.TransactionHistory,stub.GetTxID())
	_,pWr := a.putWarehouseReceipt(stub,wr)
	if !pWr {
		res := getRetString(1,"BulkchainChaincode Invoke sendRegisterRequest put WarehouseReceipt failed")
		return shim.Error(res)
	}

	// put request
	request.TxType = TxType_RegisterRequest
	request.RegisteringWarehouseReceipt = wr
	request.CheckState = Check_State_Checking
	if request.RegisteringQuantity <= 0 {
		request.RegisteringQuantity = wr.Quantity
	} else {
		//replace this line when needs that function of registeringQuantity
		request.RegisteringQuantity = wr.Quantity
	}

	var tx Transaction
	bReq,err := json.Marshal(request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendRegisterRequest marshal the request failed")
		return shim.Error(res)
	}
	tx.TxId = request.TransactionId
	tx.TxType = request.TxType
	tx.Content = bReq

	_,bTx := a.putRequest(stub,tx)
	if !bTx {
		res := getRetString(1,"BulkchainChaincode Invoke sendRegisterRequest put request failed")
		return shim.Error(res)
	} 

	res := getRetByte(0,"BulkchainChaincode Invoke sendRegisterRequest success")
	return shim.Success(res)


}

//the exchange check the register request
//args: 0 - RegisterRequest
func (a *mychaincode) checkRegisterRequest(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode Invoke checkRegisterRequest args != 1")
		return shim.Error(res)
	}

	var checkResult RegisterRequest
	err := json.Unmarshal([]byte(args[0]),&checkResult)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkRegisterRequest unmarshal request failed")
		return shim.Error(res)
	}

	tx,hasReq := a.getRequest(stub,checkResult.TransactionId)
	if !hasReq {
		res := getRetString(1,"BulkchainChaincode Invoke checkRegisterRequest request does not exist")
		return shim.Error(res)
	}

	var request RegisterRequest
	err = json.Unmarshal(tx.Content,&request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkRegisterRequest unmarshal request failed")
		return shim.Error(res)
	}

	//validate
	if !isCheckStateValid(checkResult.CheckState) {
		res := getRetString(1,"BulkchainChaincode Invoke checkRegisterRequest wrong CheckState")
		return shim.Error(res)
	}
	if request.CheckState != Check_State_Checking && (request.CheckState == Check_State_CheckedResolved || request.CheckState == Check_State_CheckedRejected) {
		res := getRetString(1,"BulkchainChaincode Invoke checkRegisterRequest request should not be checked multi times")
		return shim.Error(res)
	} 

	request.DateCheck = checkResult.DateCheck
	request.CheckState = checkResult.CheckState
	request.Description = checkResult.Description

	if request.CheckState != Check_State_CheckedRejected && request.CheckState != Check_State_CheckedResolved {
		res := getRetString(1,"BulkchainChaincode Invoke checkRegisterRequest CheckState is invalid")
		return shim.Error(res)
	}

	if !reflect.DeepEqual(request,checkResult) {
		res := getRetString(1,"BulkchainChaincode Invoke checkRegisterRequest request info was modified unexpectedly")
		return shim.Error(res)
	}

	//block#issue1
	//register the warehouse receipt
	//if the check state is ok 
	wr,bWr := a.getWarehouseReceipt(stub,request.RegisteringSeriesId)
	if !bWr {
		res := getRetString(1,"BulkchainChaincode Invoke checkRegisterRequest the registering warehousereceipt does not exist")
		return shim.Error(res)
	}
	if request.CheckState == Check_State_CheckedResolved {
		wr.State = State_FLOWABLE
		wr.Type = WarehouseReceipt_Type_Standard
	} else {
		wr.State = State_INBOUND
	}
	_,pWr := a.putWarehouseReceipt(stub,wr)
	if !pWr {
		res := getRetString(1,"BulkchainChaincode Invoke checkRegisterRequest put WarehouseReceipt failed")
		return shim.Error(res)
	}

	//put the request
	request.RegisteredWarehouseReceipt = wr
	bReq,err := json.Marshal(request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkRegisterRequest marshal request failed")
		return shim.Error(res)
	}
	tx.Content = bReq
	_,bTx := a.putRequest(stub,tx)
	if !bTx {
		res := getRetString(1,"BulkchainChaincode Invoke RegisterRequest put request failed")
		return shim.Error(res)
	}

	res := getRetByte(0,"BulkchainChaincode Invoke checkRegisterRequest success")
	return shim.Success(res)

}


//args: 0 - DeliveryRequest
func (a *mychaincode) sendDeliveryRequest(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode Invoke sendDeliveryRequest args != 1")
		return shim.Error(res)
	} 

	var request DeliveryRequest
	err := json.Unmarshal([]byte(args[0]),&request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendDeliveryRequest unmarshal request failed")
		return shim.Error(res)
	}
	request.TransactionId = stub.GetTxID()
	
	//validate
	//check the deliveryType
	if request.DeliveryType != DeliveryType_Buyer && request.DeliveryType != DeliveryType_Seller {
		res := getRetString(1,"BulkchainChaincode Invoke sendDeliveryRequest wrong DeliveryType")
		return shim.Error(res)
	}
	//check the warehouse receipt
	wr,bWr := a.getWarehouseReceipt(stub,request.DeliveryWarehouseReceiptSeriesId)
	if request.DeliveryType == DeliveryType_Buyer {
		if _,ok := varieties[request.DeliveryVarietyCode];!ok {
			res := getRetString(1,"BulkchainChaincode Invoke sendDeliveryRequest wrong variety code")
			return shim.Error(res)
		}
		if request.DeliveryQuantity <= 0 {
			res := getRetString(1,"BulkchainChaincode Invoke sendDeliveryRequest DeliveryQuantity invalid")
			return shim.Error(res)
		}

		request.TradeUnit = varieties[request.DeliveryVarietyCode].MinDeliveryHands
	} else {
		if !bWr {
			res := getRetString(1,"BulkchainChaincode Invoke sendDeliveryRequest WarehouseReceipt does not exist")
			return shim.Error(res)
		}
		request.DeliveryVarietyCode = wr.VarietyCode
		request.DeliveryQuantity = wr.Quantity
		request.TradeUnit = varieties[wr.VarietyCode].MinDeliveryHands 

		//check the holder
		if wr.WarehouseReceiptHolderId != request.ClientId || wr.MemberId != request.MemberId {
			res := getRetString(1,"BulkchainChaincode Invoke sendDeliveryRequest the client is not the holder of the WarehouseReceipt")
			return shim.Error(res)
		}
		///check the type and state
		if wr.Type != WarehouseReceipt_Type_Standard || wr.State != State_FLOWABLE {
			res := getRetString(1,"BulkchainChaincode Invoke sendDeliveryRequest the operation is illegal because of the state and type of WarehouseReceipt ")
			return shim.Error(res)
		}
	}
	


	//set wr state and record history
	if bWr {
		wr.State = State_DELIVERYING
		wr.TransactionHistory = append(wr.TransactionHistory,stub.GetTxID())
		_,pWr := a.putWarehouseReceipt(stub,wr)
		if !pWr {
			res := getRetString(1,"BulkchainChaincode Invoke sendDeliveryRequest put WarehouseReceipt failed")
			return shim.Error(res)
		}
	}

	//put search table
	MemberIdTxIdKey, err := stub.CreateCompositeKey(Index_Member_TxId, []string{request.MemberId, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendDeliveryRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(MemberIdTxIdKey, []byte(TxType_DeliveryRequest))

	ClientIdTxIdKey, err := stub.CreateCompositeKey(Index_Client_TxId, []string{request.ClientId, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendDeliveryRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(ClientIdTxIdKey, []byte(TxType_DeliveryRequest))

	ExchangeTxIdKey, err := stub.CreateCompositeKey(Index_Exchange_TxId,[]string{Key_Exchange, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendDeliveryRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(ExchangeTxIdKey, []byte(TxType_DeliveryRequest))

	AllRequestsKey, err := stub.CreateCompositeKey(Index_AllRequests, []string{TxType_RegisterRequest, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendDeliveryRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(AllRequestsKey, []byte{0x00})

	//put request
	
	request.TxType = TxType_DeliveryRequest
	request.CheckState = Check_State_Checking
	var tx Transaction
	bReq,err := json.Marshal(request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendDeliveryRequest marshal request failed")
		return shim.Error(res)
	}
	tx.TxId = request.TransactionId
	tx.TxType = request.TxType
	tx.Content = bReq
	_,bTx := a.putRequest(stub,tx)
	if !bTx {
		res := getRetString(1,"BulkchainChaincode Invoke sendDeliveryRequest put request failed")
		return shim.Error(res)
	} 

	res := getRetByte(0,"BulkchainChaincode Invoke sendDeliveryRequest success")
	return shim.Success(res)
}


//args: 0 - DeliveryRequest
func (a *mychaincode) checkDeliveryRequest(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode Invoke checkDeliveryRequest args != 1")
		return shim.Error(res)
	}

	var checkResult DeliveryRequest
	err := json.Unmarshal([]byte(args[0]),&checkResult)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkDeliveryRequest unmarshal failed")
		return shim.Error(res)
	}

	tx,hasReq := a.getRequest(stub,checkResult.TransactionId)
	if !hasReq {
		res := getRetString(1,"BulkchainChaincode Invoke checkDeliveryRequest request does not exist")
		return shim.Error(res)
	}


	var request DeliveryRequest
	err = json.Unmarshal(tx.Content,&request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkDeliveryRequest unmarshal request failed")
		return shim.Error(res)
	}

	//validate
	if !isCheckStateValid(checkResult.CheckState) {
		res := getRetString(1,"BulkchainChaincode Invoke checkDeliveryRequest wrong CheckState")
		return shim.Error(res)
	}
	if request.CheckState != Check_State_Checking && (request.CheckState == Check_State_CheckedRejected || request.CheckState == Check_State_CheckedResolved) {
		res := getRetString(1,"BulkchainChaincode Invoke checkDeliveryRequest request should not be checked multi times")
		return shim.Error(res)	
	}

	request.DateCheck = checkResult.DateCheck
	request.CheckState = checkResult.CheckState
	request.Description = checkResult.Description

	if !reflect.DeepEqual(request,checkResult) {
		res := getRetString(1,"BulkchainChaincode Invoke checkDeliveryRequest request was modified unexpectedly")
		return shim.Error(res)
	}

	//for matching 
	if request.CheckState == Check_State_CheckedResolved {
		request.MatchState = MatchState_Matching
		DeliveryTypeVarietyCodeTxIdKey,err := stub.CreateCompositeKey(Index_DeliveryType_VarietyCode_TxId,[]string{request.DeliveryType,request.DeliveryVarietyCode,request.TransactionId})
		if err != nil {
			res := getRetString(1,"BulkchainChaincode Invoke checkDeliveryRequest CreateCompositeKey failed")
			return shim.Error(res)
		}
		bQuantity,err := json.Marshal(request.DeliveryQuantity)
		if err != nil {
			res := getRetString(1,"BulkchainChaincode Invoke checkDeliveryRequest marshal DeliveryQuantity failed")
			return shim.Error(res)
		}
		stub.PutState(DeliveryTypeVarietyCodeTxIdKey,bQuantity)
	}
	//put the request
	bReq,err := json.Marshal(request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkDeliveryRequest marshal request failed")
		return shim.Error(res)
	}
	tx.Content = bReq
	_,bTx := a.putRequest(stub,tx)
	if !bTx {
		res := getRetString(1,"BulkchainChaincode Invoke checkDeliveryRequest put request failed")
		return shim.Error(res)
	}	
	res := getRetByte(0,"BulkchainChaincode Invoke checkDeliveryRequest success")
	return shim.Success(res)
}


//args: 0 - DateMatch
func (a *mychaincode) matchDeliveryRequest(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest args != 1")
		return shim.Error(res)
	}

	for varietyCode,_ := range varieties {
		//get iterator
		buyerTxIdIterator,err := stub.GetStateByPartialCompositeKey(Index_DeliveryType_VarietyCode_TxId,[]string{DeliveryType_Buyer,varietyCode})	
		if err != nil {
			res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest GetStateByPartialCompositeKey of buyers failed")
			return shim.Error(res)
		}
		defer buyerTxIdIterator.Close()
		sellerTxIdIterator,err := stub.GetStateByPartialCompositeKey(Index_DeliveryType_VarietyCode_TxId,[]string{DeliveryType_Seller,varietyCode})
		if err != nil {
			res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest GetStateByPartialCompositeKey of sellers failed")
			return shim.Error(res)
		}
		defer sellerTxIdIterator.Close()
		
		//get txid and quantity lists
		//for buyer requests
		buyerReqs := make(map[string]int)
		for buyerTxIdIterator.HasNext() {
			kv,_ := buyerTxIdIterator.Next()
			_,compositeKeyParts,err := stub.SplitCompositeKey(kv.Key)
			if err != nil {
				res := getRetString(1,"BulkchainChaincode matchDeliveryRequest SplitCompositeKey error")
				return shim.Error(res)
			}
			txId := compositeKeyParts[2]
			
			bQuantity := kv.Value
			var quantity int
			err = json.Unmarshal(bQuantity,&quantity)
			if err != nil {
				res := getRetString(1,"BulkchainChaincode matchDeliveryRequest unmarshal quantity failed")
				return shim.Error(res)
			}
		
			buyerReqs[txId] = quantity
		}

		//for sellers requests
		sellerReqs := make(map[string]int)
		for sellerTxIdIterator.HasNext() {
			kv,_ := sellerTxIdIterator.Next()
			_,compositeKeyParts,err := stub.SplitCompositeKey(kv.Key)
			if err != nil {
				res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest SplitCompositeKey failed")
				return shim.Error(res)
			}
			txId := compositeKeyParts[2]

			bQuantity := kv.Value
			var quantity int
			err = json.Unmarshal(bQuantity,&quantity)
			if err != nil {
				res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest unmarshal quantity failed")
				return shim.Error(res)
			}

			sellerReqs[txId] = quantity
		}

		var matchRate = 0.1
		//matching 
		matchSellerToBuyerResults := make(map[string]string)
		matchSellerState := make(map[string]bool)
		matchBuyerState := make(map[string]bool)
		for sellerTxId,sellerQuantity := range sellerReqs {
			matchSellerState[sellerTxId] = false
			for buyerTxId,buyerQuantity := range buyerReqs {
				if _,ok := matchBuyerState[buyerTxId];!ok {
					matchBuyerState[buyerTxId] = false
				}

				if float64(buyerQuantity) <= float64(sellerQuantity) * (1.0 + matchRate) && float64(buyerQuantity) >= float64(sellerQuantity) * (1.0 - matchRate) && matchBuyerState[buyerTxId] == false {
					matchSellerToBuyerResults[sellerTxId] = buyerTxId
					matchSellerState[sellerTxId] = true
					matchBuyerState[buyerTxId] = true
				} 
			}
			
		}

		//put matching result to the delivery request
		for sellerTxId,buyerTxId := range matchSellerToBuyerResults {
			//get seller request
			bSellerTx,hasTx := a.getRequest(stub,sellerTxId)
			if !hasTx {
				res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest getRequest of sellerTxId failed")
				return shim.Error(res)
			}
			var sellerRequest DeliveryRequest
 			err = json.Unmarshal(bSellerTx.Content,&sellerRequest) 
 			if err != nil {
 				res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest unmarshal sellerRequest failed")
 				return shim.Error(res)
 			}
 			//get buyer request
 			bBuyerTx,hasTx := a.getRequest(stub,buyerTxId)
 			if !hasTx {
 				res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest getRequest of buyerTxId failed")
 				return shim.Error(res)
 			}
 			var buyerRequest DeliveryRequest
 			err = json.Unmarshal(bBuyerTx.Content,&buyerRequest)
 			if err != nil {
 				res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest unmarshal buyerRequest failed")
 				return shim.Error(res)
 			}

 			//get seller wr
 			wr,hasWr := a.getWarehouseReceipt(stub,sellerRequest.DeliveryWarehouseReceiptSeriesId)
 			if !hasWr {
 				res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest getWarehouseReceipt failed")
 				return shim.Error(res)
 			}

 			//update seller request
 			sellerRequest.MatchState = MatchState_Matched
 			sellerRequest.DateMatch = args[0]
 			sellerRequest.MatchClientId = buyerRequest.ClientId
 			sellerRequest.MatchClientName = buyerRequest.ClientName
 			sellerRequest.MatchMemberId = buyerRequest.MemberId
 			sellerRequest.MatchMemberName = buyerRequest.MemberName
 			sellerRequest.MatchTxId = buyerTxId
 			//update buyer request
 			buyerRequest.MatchState = MatchState_Matched
 			buyerRequest.DateMatch = args[0]
 			buyerRequest.MatchClientId = sellerRequest.ClientId
 			buyerRequest.MatchClientName = sellerRequest.ClientName
 			buyerRequest.MatchMemberId = sellerRequest.MemberId
 			buyerRequest.MatchMemberName = sellerRequest.MemberName
 			buyerRequest.MatchTxId = sellerTxId
 			buyerRequest.MatchQuantity = wr.Quantity
 			buyerRequest.MatchWarehouseReceipt = wr
 			//put seller request
 			bSellerRequest,err := json.Marshal(sellerRequest)
 			if err != nil {
 				res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest marshal sellerRequest failed")
 				return shim.Error(res)
 			}
			bSellerTx.Content = bSellerRequest
			_,ok := a.putRequest(stub,bSellerTx)
			if !ok {
				res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest matchDeliveryRequest put sellerRequest failed")
				return shim.Error(res)
			}
			//put buyer request
			bBuyerRequest,err := json.Marshal(buyerRequest)
			if err != nil {
				res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest marshal buyerRequest failed")
				return shim.Error(res)
			}
			bBuyerTx.Content = bBuyerRequest
			_,ok = a.putRequest(stub,bBuyerTx)
			if !ok {
				res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest put buyerRequest failed")
				return shim.Error(res)
			} 			

		}
		//recored the unmatched
		for sellerTxId,isMatched := range matchSellerState {
			if !isMatched {
				//get the seller request
				bSellerTx,hasTx := a.getRequest(stub,sellerTxId)
				if !hasTx {
					res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest getRequest of sellerTxId failed")
					return shim.Error(res)
				}
				var sellerRequest DeliveryRequest
	 			err = json.Unmarshal(bSellerTx.Content,&sellerRequest) 
	 			if err != nil {
	 				res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest unmarshal sellerRequest failed")
	 				return shim.Error(res)
	 			}

	 			sellerRequest.MatchState = MatchState_Unmatched
	 			sellerRequest.DateMatch = args[0]

	 			//put seller request
	 			bSellerRequest,err := json.Marshal(sellerRequest)
	 			if err != nil {
	 				res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest marshal sellerRequest failed")
	 				return shim.Error(res)
	 			}
				bSellerTx.Content = bSellerRequest
				_,ok := a.putRequest(stub,bSellerTx)
				if !ok {
					res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest matchDeliveryRequest put sellerRequest failed")
					return shim.Error(res)
				}
			}
		}

		for buyerTxId,isMatched := range matchBuyerState {
			if !isMatched {
				//get buyer request
	 			bBuyerTx,hasTx := a.getRequest(stub,buyerTxId)
	 			if !hasTx {
	 				res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest getRequest of buyerTxId failed")
	 				return shim.Error(res)
	 			}
	 			var buyerRequest DeliveryRequest
	 			err = json.Unmarshal(bBuyerTx.Content,&buyerRequest)
	 			if err != nil {
	 				res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest unmarshal buyerRequest failed")
	 				return shim.Error(res)
	 			}

	 			buyerRequest.MatchState = MatchState_Unmatched
	 			buyerRequest.DateMatch = args[0]

	 			//put buyer request
				bBuyerRequest,err := json.Marshal(buyerRequest)
				if err != nil {
					res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest marshal buyerRequest failed")
					return shim.Error(res)
				}
				bBuyerTx.Content = bBuyerRequest
				_,ok := a.putRequest(stub,bBuyerTx)
				if !ok {
					res := getRetString(1,"BulkchainChaincode Invoke matchDeliveryRequest put buyerRequest failed")
					return shim.Error(res)
				} 	

			}
		}
		//put matching result in search tables
		

		//delete matching tables
		
	}

	//return matching results
	return shim.Success(nil)	
}

//send PledgeRequest
//args: 0 - PledgeReuqest
func (a *mychaincode) sendPledgeRequest(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode Invoke sendPledgeRequest args != 1")
		return shim.Error(res)
	}

	var request PledgeRequest
	err := json.Unmarshal([]byte(args[0]),&request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendPledgeRequest unmarshal request failed")
		return shim.Error(res)
	}
	request.TransactionId = stub.GetTxID()

	//get warehouse receipt
	wr,bWr := a.getWarehouseReceipt(stub,request.PledgingWarehouseReceiptSeriesId)
	if !bWr {
		res := getRetString(1,"BulkchainChaincode Invoke sendPledgeRequest WarehouseReceipt does not existed")
		return shim.Error(res)
	}

	//validate
	//check the holder of wr
	if wr.WarehouseReceiptHolderId != request.ClientId || wr.MemberId != request.MemberId {
		res := getRetString(1,"BulkchainChaincode Invoke sendPledgeRequest the client is not the holder of warehousereceipt")
		return shim.Error(res)
	}
	//check the wr type and state
	if wr.Type != WarehouseReceipt_Type_Standard || wr.State != State_FLOWABLE {
		res := getRetString(1,"BulkchainChaincode Invoke sendPledgeRequest the operation is illegal because of the state and type of WarehouseReceipt")
		return shim.Error(res)
	}

	//check the pledge type
	if !isPledgeTypeValid(request.PledgeType) {
		res := getRetString(1,"BulkchainChaincode Invoke sendPledgeRequest PledgeType is invalid")
		return shim.Error(res)
	}
	//check bankname and bankid
	if request.PledgeType == PledgeType_Outside {
		if request.TargetBankName == "" || request.TargetBankId == "" {
			res := getRetString(1,"BulkchainChaincode Invoke sendPledgeRequest TargetBankId and TargetBankName should not be empty")
			return shim.Error(res)
		} 
	}

	
	//put search table
	MemberIdTxIdKey, err := stub.CreateCompositeKey(Index_Member_TxId, []string{request.MemberId, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendPledgeRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(MemberIdTxIdKey, []byte(TxType_PledgeRequest))

	ClientIdTxIdKey, err := stub.CreateCompositeKey(Index_Client_TxId, []string{request.ClientId, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendPledgeRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(ClientIdTxIdKey, []byte(TxType_PledgeRequest))


	if request.PledgeType == PledgeType_Inside {
		ExchangeTxIdKey, err := stub.CreateCompositeKey(Index_Exchange_TxId,[]string{Key_Exchange, request.TransactionId})
		if err != nil {
			res := getRetString(1,"BulkchainChaincode Invoke sendPledgeRequest put search table failed")
			return shim.Error(res)
		}
		stub.PutState(ExchangeTxIdKey, []byte(TxType_PledgeRequest))
	} else {
		BankIdTxIdKey, err := stub.CreateCompositeKey(Index_Bank_TxId,[]string{request.TargetBankId,request.TransactionId})
		if err != nil {
			 res := getRetString(1,"BulkchainChaincode Invoke sendPledgeRequest put search table failed")
			 return shim.Error(res)
		}
		stub.PutState(BankIdTxIdKey, []byte(TxType_PledgeRequest))
	}
	

	AllRequestsKey, err := stub.CreateCompositeKey(Index_AllRequests, []string{TxType_PledgeRequest, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendPledgeRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(AllRequestsKey, []byte{0x00})


	//set wr state and record history
	wr.State = State_PLEDGING
	wr.TransactionHistory = append(wr.TransactionHistory,request.TransactionId)
	_,pWr := a.putWarehouseReceipt(stub,wr)
	if !pWr {
		res := getRetString(1,"BulkchainChaincode Invoke sendPledgeRequest putWarehouseReceipt failed")
		return shim.Error(res)
	}

	//put request
	request.TxType = TxType_PledgeRequest
	request.PledgingWarehouseReceipt = wr
	request.CheckState = Check_State_Checking
	var tx Transaction
	bReq,err := json.Marshal(request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendPledgeRequest marshal request failed")
		return shim.Error(res)
	}
	tx.TxId = request.TransactionId
	tx.Content = bReq
	tx.TxType = TxType_PledgeRequest

	_,bTx := a.putRequest(stub,tx)
	if !bTx {
		res := getRetString(1,"BulkchainChaincode Invoke sendPledgeRequest putRequest failed")
		return shim.Error(res)
	}

	res := getRetByte(0,"BulkchainChaincode Invoke sendPledgeRequest success")
	return shim.Success(res)
}


//check PledgeRequest
//args: 0 - PledgeRequest
func (a *mychaincode) checkPledgeRequest(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode Invoke checkPledgeRequest args != 1")
		return shim.Error(res)
	}

	var checkResult PledgeRequest
	err := json.Unmarshal([]byte(args[0]),&checkResult)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkPledgeRequest unmarshal request failed")
		return shim.Error(res)
	}

	tx,hasReq := a.getRequest(stub,checkResult.TransactionId)
	if !hasReq {
		res := getRetString(1,"BulkchainChaincode Invoke checkPledgeRequest request does not exist")
		return shim.Error(res)
	}

	var request PledgeRequest
	err = json.Unmarshal(tx.Content,&request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkPledgeRequest unmarshal request failed")
		return shim.Error(res)
	}

	//validate
	if !isCheckStateValid(checkResult.CheckState) {
		res := getRetString(1,"BulkchainChaincode Invoke checkPledgeRequest wrong CheckState")
		return shim.Error(res)
	}
	if request.CheckState != Check_State_Checking && (request.CheckState == Check_State_CheckedResolved || request.CheckState == Check_State_CheckedRejected) {
		res := getRetString(1,"BulkchainChaincode Invoke checkPledgeRequest request should not be checked multi times")
		return shim.Error(res)
	} 

	request.DateCheck = checkResult.DateCheck
	request.CheckState = checkResult.CheckState
	request.Description = checkResult.Description
	request.AmountOfMoneyLended = checkResult.AmountOfMoneyLended
	request.AmountOfMoneyReturning = checkResult.AmountOfMoneyReturning

	if checkResult.CheckState != Check_State_CheckedRejected && checkResult.CheckState != Check_State_CheckedResolved {
		res := getRetString(1,"BulkchainChaincode Invoke checkPledgeRequest CheckState is invalid")
		return shim.Error(res)
	}

	if !reflect.DeepEqual(request,checkResult) {
		res := getRetString(1,"BulkchainChaincode Invoke checkPledgeRequest request info was modified unexpectedly")
		return shim.Error(res)
	}


	//the check wr state  
	wr,bWr := a.getWarehouseReceipt(stub,request.PledgingWarehouseReceiptSeriesId)
	if !bWr {
		res := getRetString(1,"BulkchainChaincode Invoke checkPledgeRequest the pledging warehousereceipt does not exist")
		return shim.Error(res)
	}
	if request.CheckState == Check_State_CheckedResolved {
		wr.State = State_PLEDGING
	} else {
		wr.State = State_FLOWABLE
	}
	_,pWr := a.putWarehouseReceipt(stub,wr)
	if !pWr {
		res := getRetString(1,"BulkchainChaincode Invoke checkPledgeRequest put WarehouseReceipt failed")
		return shim.Error(res)
	}

	//put the request
	request.PledgingWarehouseReceipt = wr
	request.ConfirmState = ConfirmState_Confirming
	bReq,err := json.Marshal(request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkPledgeRequest marshal request failed")
		return shim.Error(res)
	}
	tx.Content = bReq
	_,bTx := a.putRequest(stub,tx)
	if !bTx {
		res := getRetString(1,"BulkchainChaincode Invoke checkPledgeRequest put request failed")
		return shim.Error(res)
	}

	res := getRetByte(0,"BulkchainChaincode Invoke checkPledgeRequest success")
	return shim.Success(res)
}


//confirm the PledgeRequest
//args: 0 - PledgeRequest
func (a *mychaincode) confirmPledgeRequest(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode Invoke confirmPledgeRequest args != 1")
		return shim.Error(res)
	}

	var checkResult PledgeRequest
	err := json.Unmarshal([]byte(args[0]),&checkResult)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke confirmPledgeRequest unmarshal request failed")
		return shim.Error(res)
	}

	tx,hasReq := a.getRequest(stub,checkResult.TransactionId)
	if !hasReq {
		res := getRetString(1,"BulkchainChaincode Invoke confirmPledgeRequest request does not exist")
		return shim.Error(res)
	}

	var request PledgeRequest
	err = json.Unmarshal(tx.Content,&request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke confirmPledgeRequest unmarshal request failed")
		return shim.Error(res)
	}

	//validate
	if checkResult.ConfirmState != ConfirmState_ConfirmRejected && checkResult.ConfirmState != ConfirmState_ConfirmResolved {
		res := getRetString(1,"BulkchainChaincode Invoke confirmPledgeRequest ConfirmState is invalid")
		return shim.Error(res)
	}
	if request.CheckState != Check_State_CheckedResolved || request.ConfirmState != ConfirmState_Confirming {
		res := getRetString(1,"BulkchainChaincode Invoke confirmPledgeRequest request should be resolved")
		return shim.Error(res)
	} 

	request.ConfirmState = checkResult.ConfirmState

	if !reflect.DeepEqual(request,checkResult) {
		res := getRetString(1,"BulkchainChaincode Invoke confirmPledgeRequest request info was modified unexpectedly")
		return shim.Error(res)
	}

	//check the wr 
	wr,bWr := a.getWarehouseReceipt(stub,request.PledgingWarehouseReceiptSeriesId)
	if !bWr {
		res := getRetString(1,"BulkchainChaincode Invoke confirmPledgeRequest the registering warehousereceipt does not exist")
		return shim.Error(res)
	}
	if request.ConfirmState == ConfirmState_ConfirmRejected {
		wr.State = State_FLOWABLE
	} else {
		wr.State = State_PLEDGED
	}
	_,pWr := a.putWarehouseReceipt(stub,wr)
	if !pWr {
		res := getRetString(1,"BulkchainChaincode Invoke confirmPledgeRequest put WarehouseReceipt failed")
		return shim.Error(res)
	}

	//put the request
	request.PledgingWarehouseReceipt = wr
	bReq,err := json.Marshal(request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke confirmPledgeRequest marshal request failed")
		return shim.Error(res)
	}
	tx.Content = bReq
	_,bTx := a.putRequest(stub,tx)
	if !bTx {
		res := getRetString(1,"BulkchainChaincode Invoke confirmPledgeRequest put request failed")
		return shim.Error(res)
	}

	res := getRetByte(0,"BulkchainChaincode Invoke checkRegisterRequest success")
	return shim.Success(res)
}



//sendUnpledgeRequest
//args: 0 - PledgeRequest
func (a *mychaincode) sendUnpledgeRequest(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnpledgeRequest args != 1")
		return shim.Error(res)
	}

	var secRequest PledgeRequest
	err := json.Unmarshal([]byte(args[0]),&secRequest)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnpledgeRequest unmarshal RegisterRequest failed")
		return shim.Error(res)
	}

	tx,hasReq := a.getRequest(stub,secRequest.TransactionId)
	if !hasReq {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnpledgeRequest request does not exist")
		return shim.Error(res)
	}

	var request PledgeRequest
	err = json.Unmarshal(tx.Content,&request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnpledgeRequest unmarshal request failed")
		return shim.Error(res)
	}

	//validate
	if request.ConfirmState != ConfirmState_ConfirmResolved {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnpledgeRequest ConfirmState should be resolved")
		return shim.Error(res)
	}
	request.UnpledgeRequestDate = secRequest.UnpledgeRequestDate
	if !reflect.DeepEqual(request,secRequest) {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnpledgeRequest request info was modified unexpectedly")
		return shim.Error(res)
	}

	wr,bWr := a.getWarehouseReceipt(stub,request.PledgingWarehouseReceiptSeriesId)
	if !bWr {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnpledgeRequest the pledging WarehouseReceipt does not eixst")
		return shim.Error(res)
	}
	//check type and state
	if wr.State != State_PLEDGED || wr.Type != WarehouseReceipt_Type_Standard {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnpledgeRequest operation is illegal because of the state and the type of warehousereceipt")
		return shim.Error(res)
	}
	//set and put wr
	wr.State = State_UNPLEDGING
	_,pWr := a.putWarehouseReceipt(stub,wr)
	if !pWr {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnpledgeRequest putWarehouseReceipt failed")
		return shim.Error(res)
	}

	//put the request
	request.CheckStateUnpledge = Check_State_Checking
	bReq,err := json.Marshal(request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnpledgeRequest marshal request failed")
		return shim.Error(res)
	}
	tx.Content = bReq

	_,bTx := a.putRequest(stub,tx)
	if !bTx {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnpledgeRequest putRequest failed")
		return shim.Error(res)
	}

	res := getRetByte(0,"BulkchainChaincode Invoke sendUnpledgeRequest success")
	return shim.Success(res)	

}

//check UnpledgeRequest
//args: 0 - PledgeRequest
func (a *mychaincode) checkUnpledgeRequest(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnpledgeRequest args != 1")
		return shim.Error(res)
	}

	var checkResult PledgeRequest
	err := json.Unmarshal([]byte(args[0]),&checkResult)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnpledgeRequest unmarshal request failed")
		return shim.Error(res)
	}

	tx,hasReq := a.getRequest(stub,checkResult.TransactionId)
	if !hasReq {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnpledgeRequest request does not exist")
		return shim.Error(res)
	}

	var request PledgeRequest
	err = json.Unmarshal(tx.Content,&request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnpledgeRequest unmarshal request failed")
		return shim.Error(res)
	}

	//validate
	if !isCheckStateValid(checkResult.CheckStateUnpledge) {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnpledgeRequest wrong CheckState")
		return shim.Error(res)
	}
	if request.CheckStateUnpledge != Check_State_Checking && (request.CheckStateUnpledge == Check_State_CheckedResolved || request.CheckStateUnpledge == Check_State_CheckedRejected) {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnpledgeRequest request should not be checked multi times")
		return shim.Error(res)
	} 

	request.DateCheckUnpledge = checkResult.DateCheckUnpledge
	request.CheckStateUnpledge = checkResult.CheckStateUnpledge
	request.AmountOfMoneyReturned = checkResult.AmountOfMoneyReturned
	request.DescriptionUnpledge = checkResult.DescriptionUnpledge

	if request.CheckStateUnpledge != Check_State_CheckedRejected && request.CheckStateUnpledge != Check_State_CheckedResolved {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnpledgeRequest CheckState is invalid")
		return shim.Error(res)
	}

	if !reflect.DeepEqual(request,checkResult) {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnpledgeRequest request info was modified unexpectedly")
		return shim.Error(res)
	}

	//check the wr
	wr,bWr := a.getWarehouseReceipt(stub,request.PledgingWarehouseReceiptSeriesId)
	if !bWr {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnpledgeRequest the Unpledging warehousereceipt does not exist")
		return shim.Error(res)
	}
	//checo wr type and state
	if request.CheckStateUnpledge == Check_State_CheckedResolved {
		wr.State = State_FLOWABLE
	} else {
		wr.State = State_PLEDGED
	}
	_,pWr := a.putWarehouseReceipt(stub,wr)
	if !pWr {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnpledgeRequest put WarehouseReceipt failed")
		return shim.Error(res)
	}

	//put the request
	request.PledgingWarehouseReceipt = wr
	bReq,err := json.Marshal(request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnpledgeRequest marshal request failed")
		return shim.Error(res)
	}
	tx.Content = bReq
	_,bTx := a.putRequest(stub,tx)
	if !bTx {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnpledgeRequest put request failed")
		return shim.Error(res)
	}

	res := getRetByte(0,"BulkchainChaincode Invoke checkUnpledgeRequest success")
	return shim.Success(res)
}


//send unregister request to the exchange
//args: 0 - UnregisterRequest
func (a *mychaincode) sendUnregisterRequest(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnregisterRequest args != 1")
		return shim.Error(res)
	}

	var request UnregisterRequest
	err := json.Unmarshal([]byte(args[0]),&request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnregisterRequest unmarshal request failed")
		return shim.Error(res)
	}
	request.TransactionId = stub.GetTxID()


	wr,bWr := a.getWarehouseReceipt(stub,request.UnregisteringSeriesId)
	if !bWr {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnregisterRequest the unregistering warehousereceipt does not exist")
		return shim.Error(res)
	}


	//validate
	//check the holder
	if wr.WarehouseReceiptHolderId != request.ClientId || wr.MemberId != request.MemberId {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnregisterRequest the client is not the holder of the warehousereceipt")
		return shim.Error(res)
	}
	//check the wr type and state 
	if wr.Type != WarehouseReceipt_Type_Standard || wr.State != State_FLOWABLE {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnregisterRequest the operation is illegal  because of the type and the state of WarehouseReceipt")
		return shim.Error(res)
	}

	//put search tables
	MemberIdTxIdKey, err := stub.CreateCompositeKey(Index_Member_TxId, []string{request.MemberId, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnregisterRequest put search table memberid-txid failed")
		return shim.Error(res)
	}
	stub.PutState(MemberIdTxIdKey, []byte(TxType_RegisterRequest))

	ClientIdTxIdKey, err := stub.CreateCompositeKey(Index_Client_TxId, []string{request.ClientId, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnregisterRequest put search table clientid-txid failed")
		return shim.Error(res)
	}
	stub.PutState(ClientIdTxIdKey, []byte(TxType_RegisterRequest))

	ExchangeTxIdKey, err := stub.CreateCompositeKey(Index_Exchange_TxId,[]string{Key_Exchange, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnregisterRequest put search table exchange-txid failed")
		return shim.Error(res)
	}
	stub.PutState(ExchangeTxIdKey, []byte(TxType_RegisterRequest))

	AllRequestsKey, err := stub.CreateCompositeKey(Index_AllRequests, []string{TxType_RegisterRequest, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnregisterRequest put search table allrequests failed")
		return shim.Error(res)
	}
	stub.PutState(AllRequestsKey, []byte{0x00})

	//set wr state and record history
	wr.State = State_UNREGISTERING
	wr.TransactionHistory = append(wr.TransactionHistory,stub.GetTxID())
	_,pWr := a.putWarehouseReceipt(stub,wr)
	if !pWr {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnregisterRequest putWarehouseReceipt failed")
		return shim.Error(res)
	}

	//put request
	request.TxType = TxType_UnregisterRequest
	request.UnregisteringWarehouseReceipt = wr
	request.CheckState = Check_State_Checking
	var tx Transaction
	bReq,err := json.Marshal(request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnregisterRequest marshal request failed")
		return shim.Error(res)
	}
	tx.TxId = request.TransactionId
	tx.TxType = request.TxType
	tx.Content = bReq
	_,bTx := a.putRequest(stub,tx)
	if !bTx {
		res := getRetString(1,"BulkchainChaincode Invoke sendUnregisterRequest putRequest failed")
		return shim.Error(res)
	}

	res := getRetByte(0,"BulkchainChaincode Invoke sendUnregisterRequest success")
	return shim.Success(res)
}

//the exchange check the unregister request
//args: 0 - UnregisterRequest
func (a *mychaincode) checkUnregisterRequest(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnregisterRequest args != 1")
		return shim.Error(res)
	}

	var checkResult UnregisterRequest
	err := json.Unmarshal([]byte(args[0]),&checkResult)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnregisterRequest unmarshal checkResult failed")
		return shim.Error(res)
	}

	tx,hasReq := a. getRequest(stub,checkResult.TransactionId)
	if !hasReq {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnregisterRequest getRequest failed")
		return shim.Error(res)
	}

	var request UnregisterRequest
	err = json.Unmarshal(tx.Content,&request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnregisterRequest unmarshal request failed")
		return shim.Error(res)
	}

	//validate
	if !isCheckStateValid(checkResult.CheckState) {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnregisterRequest wrong checkState")
		return shim.Error(res)
	}
	if request.CheckState != Check_State_Checking && (request.CheckState == Check_State_CheckedResolved || request.CheckState == Check_State_CheckedRejected) {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnregisterRequest request should not be checked multi times")
		return shim.Error(res)
	}

	request.DateCheck = checkResult.DateCheck
	request.CheckState = checkResult.CheckState
	request.Description = checkResult.Description

	if request.CheckState != Check_State_CheckedResolved && request.CheckState != Check_State_CheckedRejected {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnregisterRequest checkState is invalid")
		return shim.Error(res)
	}

	if !reflect.DeepEqual(request,checkResult) {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnregisterRequest request info was modified unexpectedly")
		return shim.Error(res)
	}

	//block#issue1
	//register the warehouse receipt
	//if the check state is ok 
	wr,bWr := a.getWarehouseReceipt(stub,request.UnregisteringSeriesId)
	if !bWr {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnregisterRequest the unregistering warehousereceipt does not exist")
		return shim.Error(res)
	}
	if request.CheckState == Check_State_CheckedResolved {
		wr.State = State_UNREGISTERED
		wr.Type = WarehouseReceipt_Type_NonStandard
	} else {
		wr.State = State_FLOWABLE
	}
	_,pWr := a.putWarehouseReceipt(stub,wr)
	if !pWr {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnregisterRequest put warehousereceipt failed")
		return shim.Error(res)
	}

	//put the request
	request.UnregisteredWarehouseReceipt = wr
	bReq,err := json.Marshal(request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnregisterRequest marshal request failed")
		return shim.Error(res)
	}
	tx.Content = bReq
	_,bTx := a.putRequest(stub,tx)
	if !bTx {
		res := getRetString(1,"BulkchainChaincode Invoke checkUnregisterRequest put request failed")
		return shim.Error(res)
	}

	res := getRetByte(0,"BulkchainChaincode Invoke checkUnregisterRequest success")
	return shim.Success(res)
}


//submit outbound request
//args:0 - OutboundRequest
func (a *mychaincode) sendOutboundRequest(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode Invoke sendOutboundRequest args != 1")
		return shim.Error(res)
	}

	var request OutboundRequest
	err := json.Unmarshal([]byte(args[0]),&request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendOutboundRequest unmarshal request failed")
		return shim.Error(res)
	}
	request.TransactionId = stub.GetTxID()
	
	//validate
	wr,bWr := a.getWarehouseReceipt(stub,request.OutboundingSeriesId)
	if !bWr {
		res := getRetString(1,"BulkchainChaincode Invoke sendOutboundRequest getWarehouseReceipt failed")
		return shim.Error(res)
	}
	//check the holder of wr
	if wr.WarehouseReceiptHolderId != request.ClientId || wr.MemberId != request.MemberId {
		res := getRetString(1,"BulkchainChaincode Invoke sendOutboundRequest the client is not the holder of the WarehouseReceipt")
		return shim.Error(res)
	}
	//check the wr type and state
	if wr.Type != WarehouseReceipt_Type_NonStandard || (wr.State != State_INBOUND && wr.State != State_UNREGISTERED) {
		res := getRetString(1,"BulkchainChaincode Invoke sendOutboundRequest the operation is illegal because of the state and type of WarehouseReceipt")
		return shim.Error(res)
	}

	//put search tables
	MemberIdTxIdKey, err := stub.CreateCompositeKey(Index_Member_TxId, []string{request.MemberId, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendOutboundRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(MemberIdTxIdKey, []byte(TxType_OutboundRequest))

	ClientIdTxIdKey, err := stub.CreateCompositeKey(Index_Client_TxId, []string{request.ClientId, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendOutboundRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(ClientIdTxIdKey, []byte(TxType_RegisterRequest))

	WarehouseIdTxIdKey, err := stub.CreateCompositeKey(Index_Warehouse_TxId,[]string{wr.GoodsHolderId, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendOutboundRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(WarehouseIdTxIdKey, []byte(TxType_RegisterRequest))

	AllRequestsKey, err := stub.CreateCompositeKey(Index_AllRequests, []string{TxType_OutboundRequest, request.TransactionId})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendOutboundRequest put search table failed")
		return shim.Error(res)
	}
	stub.PutState(AllRequestsKey, []byte{0x00})


	//set wr and record history
	wr.State = State_OUTBOUNDING
	wr.TransactionHistory = append(wr.TransactionHistory,request.TransactionId)
	_,pWr := a.putWarehouseReceipt(stub,wr)
	if !pWr {
		res := getRetString(1,"BulkchainChaincode Invoke sendOutboundRequest putWarehouseReceipt failed")
		return shim.Error(res)
	}

	//put request
	request.CheckState = Check_State_Checking
	request.TxType = TxType_OutboundRequest
	var tx Transaction
	bReq,err := json.Marshal(request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke sendOutboundRequest marshal request failed")
		return shim.Error(res)
	}
	tx.TxId = request.TransactionId
	tx.TxType = request.TxType
	tx.Content = bReq

	_,bTx := a.putRequest(stub,tx)
	if !bTx {
		res := getRetString(1,"BulkchainChaincode Invoke sendOutboundRequest putRequest failed")
		return shim.Error(res)
	}

	res := getRetByte(0,"BulkchainChaincode Invoke sendOutboundRequest success")
	return shim.Success(res)


}


//the warehouse check the register request
//args: 0 - outboundRequest
func (a *mychaincode) checkOutboundRequest(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode Invoke checkOutboundRequest args != 1")
		return shim.Error(res)
	}

	var checkResult OutboundRequest
	err := json.Unmarshal([]byte(args[0]),&checkResult)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkOutboundRequest unmarshal checkResult failed")
		return shim.Error(res)
	}

	tx,hasReq := a.getRequest(stub,checkResult.TransactionId)
	if !hasReq {
		res := getRetString(1,"BulkchainChaincode Invoke checkOutboundRequest request does not exist")
		return shim.Error(res)
	}

	var request OutboundRequest
	err = json.Unmarshal(tx.Content,&request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkOutboundRequest unmarshal request failed")
		return shim.Error(res)
	}

	//validate
	if !isCheckStateValid(checkResult.CheckState) {
		res := getRetString(1,"BulkchainChaincode Invoke checkOutboundRequest wrong checkState")
		return shim.Error(res)
	}
	if request.CheckState != Check_State_Checking && (request.CheckState == Check_State_CheckedResolved || request.CheckState == Check_State_CheckedRejected) {
		res := getRetString(1,"BulkchainChaincode Invoke checkOutboundRequest request should not be checked multi times")
		return shim.Error(res)
	}

	request.DateCheck = checkResult.DateCheck
	request.CheckState = checkResult.CheckState
	request.DatePermitted = checkResult.DatePermitted
	request.Description = checkResult.Description

	if request.CheckState != Check_State_CheckedRejected && request.CheckState != Check_State_CheckedResolved {
		res := getRetString(1,"BulkchainChaincode Invoke checkOutboundRequest checkState is invalid")
		return shim.Error(res)
	}

	if !reflect.DeepEqual(request,checkResult) {
		res :=getRetString(1,"BulkchainChaincode Invoke checkOutboundRequest request info was modified unexpectedly")
		return shim.Error(res)
	}

	//check the warehouse receipt
	wr,bWr := a.getWarehouseReceipt(stub,request.OutboundingSeriesId)
	if !bWr {
		res := getRetString(1,"BulkchainChaincode Invoke checkOutboundRequest the outbounding warehousereceipt does not exist")
		return shim.Error(res)
	}
	if request.CheckState == Check_State_CheckedResolved {
		wr.State = State_OUTBOUNDREADY
	} else {
		wr.State = State_UNREGISTERED
	}
	_,pWr := a.putWarehouseReceipt(stub,wr)
	if !pWr {
		res := getRetString(1,"BulkchainChaincode Invoke checkOutboundRequest put WarehouseReceipt failed")
		return shim.Error(res)
	}

	//put the request
	bReq,err := json.Marshal(request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke checkOutboundRequest marshal request failed")
		return shim.Error(res)
	}
	tx.Content = bReq
	_,bTx := a.putRequest(stub,tx)
	if !bTx {
		res := getRetString(1,"BulkchainChaincode Invoke checkOutboundRequest putRequest failed")
		return shim.Error(res)
	}

	res := getRetByte(0,"BulkchainChaincode Invoke checkOutboundRequest success")
	return shim.Success(res)

}

//invoke
//registerOutbound
//args: 0 - OutboundRequest
func (a *mychaincode) registerOutbound(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode Invoke registerOutbound args != 1")
		return shim.Error(res)
	}

	var registerResult OutboundRequest
	err := json.Unmarshal([]byte(args[0]),&registerResult)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke registerOutbound unmarshal request failed")
		return shim.Error(res)
	}

	//validate
	tx,existReq := a.getRequest(stub,registerResult.TransactionId)
	if !existReq {
		res := getRetString(1,"BulkchainChaincode Invoke registerOutbound the request does not exist")
		return shim.Error(res)
	}
	var request OutboundRequest
	err = json.Unmarshal(tx.Content,&request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke registerOutbound unmarshal request failed")
		return shim.Error(res)
	}

	if request.CheckState != Check_State_CheckedResolved {
		res := getRetString(1,"BulkchainChaincode Invoke registerOutbound CheckState is not resolved or has been outbounded")
		return shim.Error(res)
	}

	if registerResult.DateIndeed == ""  {
		res :=getRetString(1,"BulkchainChaincode Invoke registerOutbound DateIndeed or GoodsIndeed should not be empty")
		return shim.Error(res)
	}
	/*
	if request.GoodsIndeed != nil {
		res := getRetString(1,"BulkchainChaincode Invoke registerOutbound goods has been outbounded")
		return shim.Error(res)
	}
	*/


	//check the warehouse receipt
	wr,bWr := a.getWarehouseReceipt(stub,request.OutboundingSeriesId)
	if !bWr {
		res := getRetString(1,"BulkchainChaincode Invoke registerOutbound the WarehouseReceipt does not exist")
		return shim.Error(res)
	}
	wr.State = State_OUTBOUNDED
	_,pWr := a.putWarehouseReceipt(stub,wr)
	if !pWr {
		res := getRetString(1,"BulkchainChaincode Invoke registerOutbound put WarehouseReceipt failed")
		return shim.Error(res)
	}


	request.DateIndeed = registerResult.DateIndeed
	request.GoodsIndeed = registerResult.GoodsIndeed

	if !reflect.DeepEqual(request,registerResult) {
		res := getRetString(1,"BulkchainChaincode Invoke registerOutbound the request info was been modified unexpectedly")
		return shim.Error(res)
	}

	//put request
	request.CheckState = Check_State_Finished
	bReq,err := json.Marshal(request)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke registerOutbound marshal request failed")
		return shim.Error(res)
	}
	tx.Content = bReq
	_,bTx := a.putRequest(stub,tx)
	if !bTx {
		res := getRetString(1,"BulkchainChaincode Invoke registerOutbound putRequest failed")
		return shim.Error(res)
	}

	res := getRetByte(0,"BulkchainChaincode Invoke registerOutbound success")
	return shim.Success(res)

}



//query
//args: 0 - warehouseReceiptSeriesId
func (a *mychaincode) queryWarehouseReceiptTransactionHistory(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode Invoke queryWarehouseReceiptTransactionHistory args != 1")
		return shim.Error(res)
	}

	seriesId := args[0]
	wr,bWr := a.getWarehouseReceipt(stub,seriesId)
	if !bWr {
		shim.Success(nil)
	}
	//get each transaction
	var results []string
	for _,txId := range wr.TransactionHistory {
		tx,hasReq := a.getRequest(stub,txId)
		if !hasReq {
			res := getRetString(1,"BulkchainChaincode Invoke queryWarehouseReceiptTransactionHistory the request was lost or does not exist")
			return shim.Error(res)
		}
		results = append(results,string(tx.Content))
	}

	bResults,err := json.Marshal(results)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode Invoke queryWarehouseReceiptTransactionHistory marshal results failed")
		return shim.Error(res)
	}
	return shim.Success(bResults)
}


/*
func (a *mychaincode) queryMatchingInfo(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	
}
*/

//query
//queryRequestById
//args: 0 - TxId
func (a *mychaincode) queryRequestById(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode queryRequestById args != 1")
		return shim.Error(res)
	}

	TxId := args[0]
	tx,bTx := a.getRequest(stub,TxId)
	if !bTx {
		return shim.Success(nil)
	}
	return shim.Success(tx.Content)

}

//query
//queryMyRequests
//args: 0 - UserId
//args: 1 - UserType
//args: 2 - TxType
func (a *mychaincode) queryMyRequests(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 3 {
		res := getRetString(1,"BulkchainChaincode queryMyRequests args != 3")
		return shim.Error(res)
	}

	userType := args[1]
	TxType := args[2]
	if !isUserTypeValid(userType){
		res := getRetString(1,"BulkchainChaincode queryMyRequests wrong userType")
		return shim.Error(res)
	} 
	requestsIterator,err := stub.GetStateByPartialCompositeKey(requestsIndexNames[userType], []string{args[0]})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode queryMyRequests get requests list error")
		return shim.Error(res)
	}
	defer requestsIterator.Close()

	var requestsList []string

	for requestsIterator.HasNext() {
		kv,_ := requestsIterator.Next()
		_,compositeKeyParts,err := stub.SplitCompositeKey(kv.Key)
		if err != nil {
			res := getRetString(1,"BulkchainChaincode queryMyRequests SplitCompositeKey error")
			return shim.Error(res)
		}
		tx,bReq := a.getRequest(stub,compositeKeyParts[1])
		if !bReq {
			res := getRetString(1,"BulkchainChaincode queryMyRequests fatal error, transaction lost unexpectedly")
			return shim.Error(res)
		}
		if TxType == TxType_All || TxType == tx.TxType {
			requestsList = append(requestsList,string(tx.Content))
		}	
	}

	b, err := json.Marshal(requestsList)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode queryMyWaitBill marshal requestsList error")
		return shim.Error(res)
	}
	return shim.Success(b)
}

//query
//queryWarehouseReceiptById
//args: 0 - SeriesId
func (a *mychaincode) queryWarehouseReceiptById(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 1 {
		res := getRetString(1,"BulkchainChaincode queryRequestById args != 1")
		return shim.Error(res)
	}
	seriesId := args[0]
	wr,bWr := a.getWarehouseReceipt(stub,seriesId)
	if !bWr {
		return shim.Success(nil)
	}
	b,err := json.Marshal(wr)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode queryRequestById marshal WarehouseReceipt error")
		return shim.Error(res)
	}
	return shim.Success(b)
}

//query
//queryMyWarehouseReceipts
//args: 0 - UserId
//args: 1 - UserType
//args: 2 - VarietyType
func (a *mychaincode) queryMyWarehouseReceipts(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		res := getRetString(1,"BulkchainChaincode queryMyWarehouseReceipts args != 3")
		return shim.Error(res)
	}

	userType := args[1]
	varietyType := args[2]
	if !isUserTypeValid(userType){
		res := getRetString(1,"BulkchainChaincode queryMyWarehouseReceipts wrong userType")
		return shim.Error(res)
	}
	
	wrsIterator,err := stub.GetStateByPartialCompositeKey(warehouseReceiptIndexNames[userType], []string{WarehouseReceipt_Prefix + args[0]})
	if err != nil {
		res := getRetString(1,"BulkchainChaincode queryMyWarehouseReceipts get WarehouseReceipts list error")
		return shim.Error(res)
	}
	defer wrsIterator.Close()

	var wrsList []WarehouseReceipt

	for wrsIterator.HasNext() {
		kv,_ := wrsIterator.Next()
		_,compositeKeyParts,err := stub.SplitCompositeKey(kv.Key)
		if err != nil {
			res := getRetString(1,"BulkchainChaincode queryMyWarehouseReceipts SplitCompositeKey error")
			return shim.Error(res)
		}
		wr,bWr := a.getWarehouseReceipt(stub,compositeKeyParts[1])
		if !bWr {
			res := getRetString(1,"BulkchainChaincode queryMyWarehouseReceipts fatal error, WarehouseReceipt lost unexpectedly")
			return shim.Error(res)
		}
		if varietyType == VarietyType_All || varietyType == wr.VarietyCode {
			wrsList = append(wrsList,wr)
		}	
	}

	b, err := json.Marshal(wrsList)
	if err != nil {
		res := getRetString(1,"BulkchainChaincode queryMyWarehouseReceipts marshal wrsList error")
		return shim.Error(res)
	}
	return shim.Success(b)
}




// chaincode Init 接口
func (a *mychaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	//varieties info 
	varieties["WH"] = Variety_WH
	varieties["CF"] = Variety_CF
	varieties["SR"] = Variety_SR
	varieties["OI"] = Variety_OI
	varieties["RI"] = Variety_RI
	varieties["PM"] = Variety_PM
	varieties["RS"] = Variety_RS
	varieties["JR"] = Variety_JR

	//index names table
	requestsIndexNames[Role_Member] = Index_Member_TxId
	requestsIndexNames[Role_Client] = Index_Client_TxId
	requestsIndexNames[Role_Warehouse] = Index_Warehouse_TxId
	requestsIndexNames[Role_Bank] = Index_Bank_TxId
	requestsIndexNames[Role_Exchange] = Index_Exchange_TxId

	warehouseReceiptIndexNames[Role_Member] = Index_Member_WrId
	warehouseReceiptIndexNames[Role_Client] = Index_Client_WrId
	warehouseReceiptIndexNames[Role_Warehouse] = Index_Warehouse_WrId
	warehouseReceiptIndexNames[Role_Bank] = Index_Bank_WrId
	warehouseReceiptIndexNames[Role_Exchange] = Index_Exchange_WrId

	return shim.Success(nil)
}





// chaincode Invoke 接口
func (a *mychaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    function,args := stub.GetFunctionAndParameters()
	chaincodeLogger.Info("%s%s","BulkchainChaincode function=",function)
	chaincodeLogger.Info("%s%s","BulkchainChaincode args=",args)

    switch function {
    	//invoke
    	case "sendInboundRequest":
    		return a.sendInboundRequest(stub,args)
    	case "checkInboundRequest":
    		return a.checkInboundRequest(stub,args)
    	case "registerInbound":
    		return a.registerInbound(stub,args)
    	case "sendRegisterRequest":
    		return a.sendRegisterRequest(stub,args)
    	case "checkRegisterRequest":
    		return a.checkRegisterRequest(stub,args)
    	case "sendDeliveryRequest":
    		return a.sendDeliveryRequest(stub,args)
    	case "checkDeliveryRequest":
    		return a.checkDeliveryRequest(stub,args)
    	case "matchDeliveryRequest":
    		return a.matchDeliveryRequest(stub,args)
    	case "sendPledgeRequest":
    		return a.sendPledgeRequest(stub,args)
    	case "checkPledgeRequest":
    		return a.checkPledgeRequest(stub,args)
    	case "confirmPledgeRequest":
    		return a.confirmPledgeRequest(stub,args) 
    	case "sendUnpledgeRequest":
    		return a.sendUnpledgeRequest(stub,args)
    	case "sendUnregisterRequest":
    		return a.sendUnregisterRequest(stub,args)
    	case "checkUnregisterRequest":
    		return a.checkUnregisterRequest(stub,args)
    	case "sendOutboundRequest":
    		return a.sendOutboundRequest(stub,args)
    	case "checkOutboundRequest":
    		return a.checkOutboundRequest(stub,args)
    	case "registerOutbound":
    		return a.registerOutbound(stub,args)
    	//query
    	case "queryMyRequests":
    		return a.queryMyRequests(stub,args)
    	case "queryMyWarehouseReceipts":
    		return a.queryMyWarehouseReceipts(stub,args)
    	case "queryRequestById":
    		return a.queryRequestById(stub,args)
    	case "queryWarehouseReceiptById":
    		return a.queryWarehouseReceiptById(stub,args)
    	case "queryWarehouseReceiptTransactionHistory":
    		return a.queryWarehouseReceiptTransactionHistory(stub,args)
    }



    res := getRetString(1,"BulkchainChaincode Unkown method!")
	chaincodeLogger.Info("%s",res)
	chaincodeLogger.Infof("%s",res)
    return shim.Error(res)
}


func main() {
	
    if err := shim.Start(new(mychaincode)); err != nil {
        fmt.Printf("Error starting mychaincode: %s", err)
    }
}
