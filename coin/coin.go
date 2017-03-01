/*
Copyright Hydursio Labs Inc. 2016 All Rights Reserved.
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
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	max_bankId    = "max_bankId"
	max_companyId = "max_companyId"
)

type Yeasycoin struct {
}

func (coin *Yeasycoin) Init(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()
	if function != "init" {
		// return nil, ErrInvalidFunction
		return shim.Error(ErrInvalidFunction.Error())
	}

	if len(args) != 2 {
		// return nil, ErrInvalidParams
		return shim.Error(ErrInvalidParams.Error())
	}

	bankName := args[0]
	bankNumber, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}

	// init center bank
	cbank := &Bank{
		Id:          0,
		Name:        bankName,
		TotalNumber: bankNumber,
		RestNumber:  bankNumber,
	}
	bankBytes, err := proto.Marshal(cbank)
	if err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}

	// put bank into blockchain
	if err := stub.PutState("bank_0", bankBytes); err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}

	// put something else into blockchain
	if err := stub.PutState(max_bankId, []byte("0")); err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}
	if err := stub.PutState(max_companyId, []byte("0")); err != nil {
		// return nil, err
		return shim.Error(err.Error())
	}

	// return bankBytes, nil
	return shim.Success(bankBytes)
}

const (
	invoke_createBank      = "createBank"
	invoke_createCompany   = "createCompany"
	invoke_issueCoin       = "issueCoin"
	invoke_issueCoinToBank = "issueCoinToBank"
	invoke_issueCoinToCp   = "issueCoinToCp"
	invoke_transfer        = "transfer"
)

func (coin *Yeasycoin) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()

	switch function {
	case invoke_createBank:
		return coin.createBank(stub, args)
	case invoke_createCompany:
		return coin.createCompany(stub, args)
	case invoke_issueCoin:
		return coin.issueCoin(stub, args)
	case invoke_issueCoinToBank:
		return coin.issueCoinToBank(stub, args)
	case invoke_issueCoinToCp:
		return coin.issueCoinToCp(stub, args)
	case invoke_transfer:
		return coin.transfer(stub, args)
	default:
		// return nil, ErrInvalidFunction
		return shim.Error(ErrInvalidParams.Error())
	}
}
