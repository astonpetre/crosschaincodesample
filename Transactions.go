/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TransactionsContract contract for managing CRUD for Transactions
type TransactionsContract struct {
	contractapi.Contract
}

type Transactions struct {
	TxnID               string  `json:"txnID"`
	Sender              string  `json:"sender"`
	Receiver            string  `json:"receiver"`
	TransactionAmount   float64 `json:"transactionsAmount"`
	SenderPreBalance    float64 `json:"senderPreBalance"`
	ReceiverPreBalance  float64 `json:"receiverPreBalance"`
	SenderPostBalance   float64 `json:"senderPostBalance"`
	ReceiverPostBalance float64 `json:"receiverPostBalance"`
	TransactionDate     string  `json:"transactionDate"`
	ActionBy            string  `json:"actionBy"`
	ActionOrg           string  `json:"actionOrg"`
	ActionPerformed     string  `json:"actionPerformed"`
}

func (s *TransactionsContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	transactions := []Transactions{
		{
			TxnID:               "M-df8J3WAbvmqewrwqeRUhmQ==",
			Sender:              "13425qewrqwerqw23423423dfd",
			Receiver:            "asfasfasdfasdf1425245234sd",
			TransactionAmount:   100,
			SenderPreBalance:    300,
			ReceiverPreBalance:  500,
			SenderPostBalance:   200,
			ReceiverPostBalance: 600,
			TransactionDate:     "2022.10.27 13:37:25",
			ActionBy:            "ADMIN",
			ActionOrg:           "Org1MSP",
			ActionPerformed:     "CreateWallet",
		},
		{
			TxnID:               "M-df8J3WAbvmqewrwqeRUhnZ==",
			Sender:              "13425qewrqwerqw23423423dfd",
			Receiver:            "asfasfasdfasdf1425245234sd",
			TransactionAmount:   100,
			SenderPreBalance:    200,
			ReceiverPreBalance:  600,
			SenderPostBalance:   100,
			ReceiverPostBalance: 700,
			TransactionDate:     "2022.10.27 13:37:55",
			ActionBy:            "ADMIN",
			ActionOrg:           "Org1MSP",
			ActionPerformed:     "CreateWallet",
		},
	}

	for _, transaction := range transactions {
		transactionJSON, err := json.Marshal(transaction)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(transaction.TxnID, transactionJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// TestExists returns true when asset with given ID exists in world state
func (s *TransactionsContract) TransactionExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	transactionJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return transactionJSON != nil, nil
}

// CreateTest creates a new instance of Test
// CreateAsset issues a new asset to the world state with given details.
func (s *TransactionsContract) CreateTransaction(
	ctx contractapi.TransactionContextInterface,
	txnID string,
	sender string,
	receiver string,
	transactionAmount string,
	senderPreBalance string,
	receiverPreBalance string,
	senderPostBalance string,
	receiverPostBalance string,
	transactionDate string,
	actionPerformed string,
) error {
	exists, err := s.TransactionExists(ctx, txnID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", txnID)
	}

	_transactionAmount, err := strconv.ParseFloat(transactionAmount, 64)
	if err != nil {
		return fmt.Errorf("txn amount Error %s", err)
	}
	_senderPreBalance, err := strconv.ParseFloat(senderPreBalance, 64)
	if err != nil {
		return fmt.Errorf("sender pre balance error %s", err)
	}
	_receiverPreBalance, err := strconv.ParseFloat(receiverPreBalance, 64)
	if err != nil {
		return fmt.Errorf("receiver pre balance error %s", err)
	}
	_senderPostBalance, err := strconv.ParseFloat(senderPostBalance, 64)
	if err != nil {
		return fmt.Errorf("sender post balance error %s", err)
	}
	_receiverPostBalance, err := strconv.ParseFloat(receiverPostBalance, 64)
	if err != nil {
		return fmt.Errorf("receiver post balance error %s", err)
	}

	transaction := Transactions{
		TxnID:               txnID,
		Sender:              sender,
		Receiver:            receiver,
		TransactionAmount:   _transactionAmount,
		SenderPreBalance:    _senderPreBalance,
		ReceiverPreBalance:  _receiverPreBalance,
		SenderPostBalance:   _senderPostBalance,
		ReceiverPostBalance: _receiverPostBalance,
		TransactionDate:     transactionDate,
		ActionBy:            "ADMIN",
		ActionOrg:           "Org1MSP",
		ActionPerformed:     actionPerformed,
	}

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(txnID, transactionJSON)
}

// ReadTest retrieves an instance of Test from the world state
func (s *TransactionsContract) ReadTransaction(ctx contractapi.TransactionContextInterface, id string) (*Transactions, error) {
	transactionJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if transactionJSON == nil {
		return nil, fmt.Errorf("the transaction %s does not exist", id)
	}

	var transaction Transactions
	err = json.Unmarshal(transactionJSON, &transaction)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// DeleteTest deletes an instance of Test from the world state
func (s *TransactionsContract) DeleteTransaction(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.TransactionExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// GetAllTests returns all assets found in world state
func (s *TransactionsContract) GetAllTransactions(ctx contractapi.TransactionContextInterface) ([]*Transactions, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var transactions []*Transactions
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var transaction Transactions
		err = json.Unmarshal(queryResponse.Value, &transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}
