package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"testing"
)

func TestCreateLoanApplication(t *testing.T) {
	fmt.Println("Entering TestCreateLoanApplication")
	attributes := make(map[string][]byte)
	//Create a custom MockStub that internally uses shim.MockStub
	stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}
}

func TestCreateLoanApplicationValidation(t *testing.T) {
	fmt.Println("Entering TestCreateLoanApplicationValidation")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	stub.MockTransactionStart("t123")
	_, err := CreateLoanApplication(stub, []string{})
	if err == nil {
		t.Fatalf("Expected CreateLoanApplication to return validation error")
	}
	stub.MockTransactionEnd("t123")
}
