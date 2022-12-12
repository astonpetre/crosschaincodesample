package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// WalletContract contract for managing CRUD for Wallet
type WalletContract struct {
	contractapi.Contract
}

type Wallet struct {
	WalletID          string  `json:"walletID"`
	Owner             string  `json:"owner"`
	WalletType        string  `json:"walletType"`
	Balance           float64 `json:"balance"`
	TransactionAmount float64 `json:"transactionAmount"`
	TransactionState  string  `json:"transactionState"`
	TransactionDate   string  `json:"transactionDate"`
	LastModified      string  `json:"lastModified"`
	ActionBy          string  `json:"actionBy"`
	OrgMSP            string  `json:"orgMSP"`
	ActionPerformed   string  `json:"actionPerformed"`
}

//Read Wallet
func (s *WalletContract) ReadWallet(ctx contractapi.TransactionContextInterface, id string) (*Wallet, error) {
	walletJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if walletJSON == nil {
		return nil, fmt.Errorf("the transaction %s does not exist", id)
	}

	var wallet Wallet
	err = json.Unmarshal(walletJSON, &wallet)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// Debit Wallet
func (s *WalletContract) Distribute(ctx contractapi.TransactionContextInterface, senderWallet string, receiverWallet string, amount string, txnID string, txnDate string) error {

	/////Log Transaction
	sender, err := s.ReadWallet(ctx, senderWallet)
	if err != nil {
		return fmt.Errorf("failed to read sender wallet: %s", err)
	}

	receiver, err := s.ReadWallet(ctx, receiverWallet)
	if err != nil {
		return fmt.Errorf("failed to read receiver wallet: %s", err)
	}

	txnAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return fmt.Errorf("string to float conversion error: %s", err)
	}

	//txnID := "afasfafadsfasdfadf"
	senderPreBalance := fmt.Sprintf("%v", sender.Balance)
	receiverPreBalance := fmt.Sprintf("%v", receiver.Balance)
	senderPostBalanceCalc := sender.Balance - txnAmount
	senderPostBalance := fmt.Sprintf("%v", senderPostBalanceCalc)
	receiverPostBalanceCalc := receiver.Balance + txnAmount
	receiverPostBalance := fmt.Sprintf("%v", receiverPostBalanceCalc)
	//transactionDate := "20-12-2022"
	actionPerformed := "DISTRIBUTE"

	//Query the external transactions chaincode
	params := []string{"CreateTransaction", txnID, senderWallet, receiverWallet, amount, senderPreBalance, receiverPreBalance, senderPostBalance, receiverPostBalance, txnDate, actionPerformed}
	invokeArgs := make([][]byte, 11)

	for i, arg := range params {
		invokeArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode("transactions", invokeArgs, "mychannel")

	if response.Status != shim.OK {
		return fmt.Errorf("could not invoke chaincode: %s", response.Payload)
	}

	return nil
}
