/*
Copyright Hydrusio Labs Inc. 2016 All Rights Reserved.
Written by mint.zhao.chiu@gmail.com. github.com: https://www.github.com/mintzhao

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package coin

import (
	"fmt"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func (coin *Yeasycoin) createBank(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		// return nil, ErrInvalidParams
		return shim.Error(ErrInvalidParams.Error())
	}

	bankName := args[0]

	// get max bank id now
	maxBankIdBytes, err := stub.GetState(max_bankId)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	maxBankId, err := strconv.ParseInt(string(maxBankIdBytes), 10, 64)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	bankId := maxBankId + 1

	bank := &Bank{
		Id:   bankId,
		Name: bankName,
	}

	// put bank into blockchain
	bankBytes, err := proto.Marshal(bank)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	if err := stub.PutState(fmt.Sprintf("bank_%v", bankId), bankBytes); err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	if err := stub.PutState(max_bankId, []byte(strconv.FormatInt(bankId, 10))); err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}

	// return bankBytes, nil
	return shim.Success(bankBytes)
}

func (coin *Yeasycoin) createCompany(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		// return nil, ErrInvalidParams
		return shim.Error(ErrInvalidParams.Error())
	}

	companyName := args[0]

	// get max company id
	maxCompanyIdBytes, err := stub.GetState(max_companyId)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	maxCompanyId, err := strconv.ParseInt(string(maxCompanyIdBytes), 10, 64)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	companyId := maxCompanyId + 1

	company := &Company{
		Id:   companyId,
		Name: companyName,
	}

	companyBytes, err := proto.Marshal(company)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	if err := stub.PutState(fmt.Sprintf("company_%v", companyId), companyBytes); err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	if err := stub.PutState(max_companyId, []byte(strconv.FormatInt(companyId, 10))); err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}

	// return companyBytes, nil
	return shim.Success(companyBytes)
}

func (coin *Yeasycoin) issueCoin(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		// return nil, ErrInvalidParams
		return shim.Error(ErrInvalidParams.Error())

	}

	coinNumber, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	// consider: if using time.Now() at chaincode, maybe the tx can not be the same at other vps
	timestamp := args[1]

	tx := &Transaction{
		FromType:  Transaction_FROM_CENTERBANK,
		FromId:    0,
		ToType:    Transaction_TO_CENTERBANK,
		ToId:      0,
		Timestamp: timestamp,
		Number:    coinNumber,
	}
	if txHash, err := HashTx(tx); err != nil {
		// return nil, err
		return shim.Error(err.Error())
	} else {
		tx.Id = txHash
	}
	txBytes, err := proto.Marshal(tx)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	if err := stub.PutState(fmt.Sprintf("tx_%s", tx.Id), txBytes); err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}

	// get centerbank
	cbankBytes, err := stub.GetState("bank_0")
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	bank, err := ParseBank(cbankBytes)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	bank.TotalNumber += coinNumber
	bank.RestNumber += coinNumber

	cbankBytes, err = proto.Marshal(bank)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	if err := stub.PutState("bank_0", cbankBytes); err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}

	// return txBytes, nil
	return shim.Success(txBytes)
}

func (coin *Yeasycoin) issueCoinToBank(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		// return nil, ErrInvalidParams
		return shim.Error(ErrInvalidParams.Error())
	}

	bankId, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	transferCount, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	timestamp := args[2]

	// transfer form centerbank to bank
	tx := &Transaction{
		FromType:  Transaction_FROM_CENTERBANK,
		FromId:    0,
		ToType:    Transaction_TO_BANK,
		ToId:      bankId,
		Timestamp: timestamp,
		Number:    transferCount,
	}
	if txHash, err := HashTx(tx); err != nil {
		// return nil, err
		return shim.Error(err.Error())
	} else {
		tx.Id = txHash
	}
	txBytes, err := proto.Marshal(tx)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	if err := stub.PutState(fmt.Sprintf("tx_%s", tx.Id), txBytes); err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}

	// get center bank & bank
	cbankBytes, err := stub.GetState("bank_0")
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	cbank, err := ParseBank(cbankBytes)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	cbank.RestNumber -= transferCount

	cbankBytes, err = proto.Marshal(cbank)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	if err := stub.PutState("bank_0", cbankBytes); err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}

	bankKey := fmt.Sprintf("bank_%v", bankId)
	bankBytes, err := stub.GetState(bankKey)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	bank, err := ParseBank(bankBytes)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	bank.TotalNumber += transferCount
	bank.RestNumber += transferCount

	bankBytes, err = proto.Marshal(bank)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	if err := stub.PutState(bankKey, bankBytes); err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}

	// return txBytes, nil
	return shim.Success(txBytes)
}

func (coin *Yeasycoin) issueCoinToCp(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	return shim.Success(nil)
}

func (coin *Yeasycoin) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	return shim.Success(nil)
}
