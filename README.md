# Leaning Chaincode by TDD

Hyperledger の ChaincodeをTDD（テスト駆動開発）をしながら学ぶための手順です。


## 1. Go言語の環境設定

### 1.1. Cloud9のWorkspaceを作成

[Cloud9](https://c9.io)のアカウントを作成し、新しいWorkspaceを作成します。

* Workspace name : learning-chiancode (Anything OK)
* Description(ワークスペースの説明) : hyperledgerのchaincodeの学習環境 (Anything OK)
* Hosted workspace(他のユーザから環境が見えるか) : Public(無料アカウントでは１つだけPrivateが作れる。あとで変更できます。)
* Clone from Git or Mercurial URL : Blank（特に設定不要）
* Choose a template : Blank(他でもいいですが、今回の目的では標準的なUbuntu環境で十分です)

![new workspace](static/cloud9_new_workspace.png)


Workspaceを作成すると下記のような画面が表示されます。
ディレクトリ構造や、ファイルエディタ、Terminalでの処理を行えます。

![workspace](static/c9_workspace.png)

### 1.2. Go言語のVersion確認

Cloud9のWorkspaceには標準でGo言語がインストールされています。

```bash
$ go version
go version go1.7.3 darwin/amd64
```

### 1.3. GOPATH環境設定

$GOPATH は設定済みです。

```bash
$ echo $GOPATH
/home/ubuntu/workspace
```

変更したい場合は、下記などのようにすることで変更可能です。
（特に変更する必要はありません。）

```bash
$ vi ~/.bashrc
# Add the below line at the end of file.
export GOPATH="/home/ubuntu/workspace/new/gopath"
```

```bash
$ source ~/.bashrc
```

### 1.4. 動作確認
go run サブコマンドを使うことでソースコードをビルドすると同時に実行する

```bash
$ vi helloworld.go
package main

import (
  "fmt"
)

func main() {
  fmt.Println("Hello, World!")
}
```

```bash
$ go run helloworld.go
Hello, World!
```

## 2. Hyperledgerの環境準備

### 2.1. サンプルのチェーンコードをクローンする

```bash
# Create the parent directories on your GOPATH
$ mkdir -p $GOPATH/src/github.com/hyperledger
$ cd $GOPATH/src/github.com/hyperledger 

# Clone the appropriate release codebase into $GOPATH/src/github.com/hyperledger/fabric
# Note that the v0.6 release is a branch of the repository.  It is defined below after the -b argument
$ git clone -b v0.6 http://gerrit.hyperledger.org/r/fabric
```

### 2.2. 動作確認

Buildしてエラーがないか確認します。

```bash
$ mkdir -p ~/workspace/build_test
$ cd ~/workspace/build_test
$ wget https://raw.githubusercontent.com/IBM-Blockchain/example02/v2.0/chaincode/chaincode_example02.go
$ go build ./
```

Buildでエラーが発生しなければ、```build_test``` という実行ファイルが出来ているはずです。

## 3. Getting started with TDD

### 3.1. モックのソースコードの事前準備

mock stub のソース (```varunmockstub.go```) を下記のディレクトリに置きます。

$GOPATH/src/github.com/hyperledger/fabric/core/chaincode/shim/


```bash
$ cd $GOPATH/src/github.com/hyperledger/fabric/core/chaincode/shim/
$ wget https://raw.githubusercontent.com/tohosokawa/learningCCbyTDD/cloud9/varunmockstub.go
```

上記のソースコードのオリジナルは、[チュートリアルページ](https://www.ibm.com/developerworks/cloud/library/cl-ibm-blockchain-chaincode-testing-using-golang/index.html#artdownload)の下段に
```Code samples in this tutorial```として置かれているものです。

### 3.2. 開発コードの準備

Workディレクトリ(sample_tdd) を作成

```bash
$ mkdir -p ~/workspace/sample_tdd
$ cd ~/workspace/sample_tdd
```

sample_tddディレクトリを右クリックしたり、terminalから以下2つのファイルを作成します。

1. sample_chaincode_test.go : sample_chaincode.goのテストを記述するファイル。(テスト対象コード)_test.go という命名規則がある。
2. sample_chaincode.go : 今回のローンアプリのユースケースを記述する。

```bash
$ touch sample_chaincode.go
$ touch sample_chaincode_test.go
```

```sample_chaincode_test.go``` を以下のように編集します。

```go
package main
import (
    "fmt"
    "testing"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)
```

testingパッケージをimportしているのですが、これはGoパッケージの自動テストを実装するために行います。

[Goの自動テストについての参考URLはこちらから。](http://golang.jp/pkg/testing)

さらに先程コピーしたテスト用のスタブファイル（CustomMockStub）とchaincodeを実装するために、shimをimportしています。

エディタで編集して保存すると、自動的に変更される場合（インデントの変更や空白行の挿入等）があります。
これはCloud9の標準設定では、Go言語の場合はファイル保存時に ```gofmt -w "$file"``` というコマンドで整形が行われるためです。
(gofmtはGo言語標準の整形ツールです。)


## 4. CreateLoanApplicationの実装

### 4.1. 実装の要求

これからsample_chaincode.goに実装する CreateLoanApplication() の要求仕様は下記です。

1. CreateLoanApplicationは ```loan application ID```と```loan applicationを表すJSON```と```ChaincodeStubInterface```を引数にとる。
2. 生成された ```loan application を表すserialized JSON```と```errorオブジェクト```を返却する。
3. 入力値が不足していたり無効な場合は ```validation error``` を throw する


sample_chaincode_test.goを以下のように編集します。
（TDDでは実装の前に要求仕様からテストを書きます。）


```go
package main
import (
    "fmt"
    "testing"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

func TestCreateLoanApplication (t *testing.T) {
    fmt.Println("Entering TestCreateLoanApplication")
    attributes := make(map[string][]byte)
    //Create a custom MockStub that internally uses shim.MockStub
    stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
    if stub == nil {
        t.Fatalf("MockStub creation failed")
    }
}
```

Golang testing packageを実行するために、functionの名前は必ずTest* にします。
またこのtest function は```*testing.T```を引数にとります。

このテストコードは```SampleChaincode```をテストするためのコードです。

この状態でgo testを実行して稼働させると、sample_chaincode.goに何も入っていないためエラーになったというメッセージが出力されます。

```bash
$ cd ~/workspace/sample_tdd
$ go test
 can't load package: package .:
 sample_chaincode.go:1:1:1 expected 'package', found 'EOF'
```

何も実装していないため、当然テストは失敗しますが、これがTDDでのRedの段階です。

TDDでの開発サイクルは Red/Green/Refactor と呼ばれます。

* <font color=red>Red</font> : テストコードを書き、実行して**失敗する**
* <font color=green>Green</font> : テストに通る最低限のコードを書く
* <font color=blue>Refactor</font> : コードのリファクタリング（重複部分の関数化など）を行う。


それでは、以下のとおりにsample_chaincode.go を編集します。

### 4.2. SampleChaincodeの実装

```go
package main
```

この状態で実行してみます。

```bash
$ go test
# _/home/ubuntu/workspace/sample_tdd
./sample_chaincode_test.go:13: undefined: SampleChaincode
FAIL    _/home/ubuntu/workspace/sample_tdd [build failed]
```

SampleChaincodeの定義がないのでエラーとなりました。そこで、以下のようにSampleChaincodeを実装します。

```go
package main

type SampleChaincode struct {
}
```

ここで実行してみると以下のようなエラーになります。

```bash
$ go test
# _/home/ubuntu/workspace/sample_tdd
./sample_chaincode_test.go:13: cannot use new(SampleChaincode) (type *SampleChaincode) as type shim.Chaincode in argument to shim.NewCustomMockStub:
        *SampleChaincode does not implement shim.Chaincode (missing Init method)
FAIL    _/home/ubuntu/workspace/sample_tdd [build failed]
```

chaincodeのInit, Query,Invokeなどの関数を定義するshim.Chaincodeが実装されていないためです。

そこで以下の通りに必要なshimのimportを含めて実装します。


```sample_chaincode.go
package main
import (
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

type SampleChaincode struct {
}

func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    return nil, nil
}
 
func (t *SampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    return nil, nil
}
 
func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    return nil, nil
}
```

SampleChaincodeはsample_chaincode_test.goでshim.NewCustomMockStub()にChaincode型の値として渡しているため、
下記に定義されているようにInit, Invoke, Queryの３つの関数を実装する必要があります。

$GOPATH/src/github.com/hyperledger/fabric/chaincode/shim/interfaces.go

```go
type Chaincode interface {
	// Init is called during Deploy transaction after the container has been
	// established, allowing the chaincode to initialize its internal data
	Init(stub ChaincodeStubInterface, function string, args []string) ([]byte, error)

	// Invoke is called for every Invoke transactions. The chaincode may change
	// its state variables
	Invoke(stub ChaincodeStubInterface, function string, args []string) ([]byte, error)

	// Query is called for Query transactions. The chaincode may only read
	// (but not modify) its state variables and return the result
	Query(stub ChaincodeStubInterface, function string, args []string) ([]byte, error)
}
```

必要なメソッドを定義するとテストがとおるようになります。
これがTDDのGreen Stage です。

```bash
$ go test
Entering TestCreateLoanApplication
2017/05/24 05:10:26 MockStub( mockStub &{} )
PASS
ok      _/home/ubuntu/workspace/sample_tdd      0.039s
```

### 4.3. CreateLoanApplicationの実装

sample_chaincode.goに下記のCreateLoanApplication() を実装します。

```go
func CreateLoanApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering CreateLoanApplication")
    return nil, nil
}
```

CreateLoanApplication()にて、```fmt.Println```というメソッドを使用するため、
下記のようにimportに "fmt" を追加します。

```go
import (
    "fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
```

sample_chaincode_test.go に CreateLoanApplication()のテスト関数として
下記の TestCreateLoanApplicationValidation()を追加します。

CreateLoanApplication()に空の入力をした場合にエラーが返ることを確認するテストです。

```sample_chaincode_test.go
func TestCreateLoanApplicationValidation(t *testing.T) {
    fmt.Println("Entering TestCreateLoanApplicationValidation")
    attributes := make(map[string][]byte)
    stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
    if stub == nil {
        t.Fatalf("MockStub creation failed")
    }
 
    stub.MockTransactionStart("t123")
    // Start transactional context
    _, err := CreateLoanApplication(stub, []string{})
    if err == nil {
        t.Fatalf("Expected CreateLoanApplication to return validation error")
    }
    // End transactional context
    stub.MockTransactionEnd("t123")
}
```

**テストを記述する際にはstub.MockTransactionStart(ID)とstub.MockTransactionEnd(ID)に注意してください。**

帳票への書き込みが発生する場合には、必ずtransactionが開始している状態でなければいけません。
CreateLoanApplication()でも書き込みが発生するため、stub.MockTransactionStart(ID)で
transactionを開始し、必ず同じIDでstub.MockTransactionEnd(ID)を呼ぶことで完了しています。


この状態で実行してみると、予想どおりvalidationエラーになります。（Red Stage)

```bash
$ go test
Entering TestCreateLoanApplication
2017/05/24 07:43:49 MockStub( mockStub &{} )
Entering TestCreateLoanApplicationValidation
2017/05/24 07:43:49 MockStub( mockStub &{} )
Entering CreateLoanApplication
--- FAIL: TestCreateLoanApplicationValidation (0.00s)
        sample_chaincode_test.go:30: Expected CreateLoanApplication to return validation error
FAIL
exit status 1
FAIL    _/home/ubuntu/workspace/sample_tdd      0.033s
```

次に、テストを通すためにCreateLoanApplication()を下記のように修正します。
(２つ目の戻り値を nilから errorsにしています)

```go
func CreateLoanApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering CreateLoanApplication")
    return nil, errors.New("Expected atleast two arguments for loan application creation") // <-- return errors
}
```

戻り値にerrors.New()を使うため、importに "errors" を追加します。

```go
import (
    "errors"
    "fmt"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)
```

ここでテストを実行します。

```bash
$ go test
Entering TestCreateLoanApplication
2017/05/24 08:18:02 MockStub( mockStub &{} )
Entering TestCreateLoanApplicationValidation
2017/05/24 08:18:02 MockStub( mockStub &{} )
Entering CreateLoanApplication
PASS
ok      _/home/ubuntu/workspace/sample_tdd      0.072s
```

正常に終了することが確認できました。(Green Stage)

### 4.4. Refactor1

このままではCreateLoanApplication()がエラーのメッセージを返すだけなので、
他の仕様を満たしていない部分がエラーになるように
下記のテストを追加してリファクタリングを行います。


```sample_chaincode_test.go
var loanApplicationID = "la1"
var loanApplication = `{"id":"` + loanApplicationID + `","propertyId":"prop1","landId":"land1","permitId":"permit1","buyerId":"vojha24","personalInfo":{"firstname":"Varun","lastname":"Ojha","dob":"dob","email":"varun@gmail.com","mobile":"99999999"},"financialInfo":{"monthlySalary":16000,"otherExpenditure":0,"monthlyRent":4150,"monthlyLoanPayment":4000},"status":"Submitted","requestedAmount":40000,"fairMarketValue":58000,"approvedAmount":40000,"reviewedBy":"bond","lastModifiedDate":"21/09/2016 2:30pm"}`

func TestCreateLoanApplicationValidation2(t *testing.T) {
	fmt.Println("Entering TestCreateLoanApplicationValidation2")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	stub.MockTransactionStart("t123")
	_, err := CreateLoanApplication(stub, []string{loanApplicationID, loanApplication})
	if err != nil {
		t.Fatalf("Expected CreateLoanApplication to succeed")
	}
	stub.MockTransactionEnd("t123")

}
```

追加した1行目、2行目でloanの申し込みデータを生成しています。これで実行してみます。

```bash
$ go test
Entering TestCreateLoanApplication
2017/05/24 09:02:37 MockStub( mockStub &{} )
Entering TestCreateLoanApplicationValidation
2017/05/24 09:02:37 MockStub( mockStub &{} )
Entering CreateLoanApplication
Entering TestCreateLoanApplicationValidation2
2017/05/24 09:02:37 MockStub( mockStub &{} )
Entering CreateLoanApplication
--- FAIL: TestCreateLoanApplicationValidation2 (0.00s)
        sample_chaincode_test.go:49: Expected CreateLoanApplication to succeed
FAIL
exit status 1
FAIL    _/home/ubuntu/workspace/sample_tdd      0.029s
```

正しい値をCreateLoanApplication()に渡していますが、エラーが返るためにテストが予想どおり失敗します。
入力値の数をチェックするように、CreateLoanApplication()を以下のように書き換えます。

```sample_chaincode.go
func CreateLoanApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering CreateLoanApplication")
    if len(args) < 2 {  // <-- This check code is inserted.
        fmt.Println("Invalid number of args")
        return nil, errors.New("Expected atleast two arguments for loan application creation")
    }
    return nil, nil
}
```

これで実行すれば、テストがPASSされます。

```bash
$ go test
Entering TestCreateLoanApplication
2017/05/24 09:11:50 MockStub( mockStub &{} )
Entering TestCreateLoanApplicationValidation
2017/05/24 09:11:50 MockStub( mockStub &{} )
Entering CreateLoanApplication
Invalid number of args
Entering TestCreateLoanApplicationValidation2
2017/05/24 09:11:50 MockStub( mockStub &{} )
Entering CreateLoanApplication
PASS
ok      _/home/ubuntu/workspace/sample_tdd      0.025s
```

### 4.5. refactor2

このあとにloanの申し込みデータが生成され、Blockchainに書き込まれるかをテストします。
以下の関数を sample_chaincode_test.go に追加します。

コードの前半は前のテストと同様にstubのセットアップです。

14行目は10行目のinvokedで成城通り作成されたloan application objectを検索します。
stub.GetState(loanApplicationID) はkeyに対応したバイト配列値を検索します。
この場合はloan application IDをレジャーから検索します。
18行目では検索されたバイト配列をLoanApplicationに戻しています。

```sample_chaincode_test.go
func TestCreateLoanApplicationValidation3(t *testing.T) {
    ////////////////////////////////////
    // Start setup
	fmt.Println("Entering TestCreateLoanApplicationValidation3")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	stub.MockTransactionStart("t123")
	CreateLoanApplication(stub, []string{loanApplicationID, loanApplication})
	stub.MockTransactionEnd("t123")
    // End setup
    ////////////////////////////////////


    ////////////////////////////////////
    // 上のCreateLoanApplicationで生成Ledgerにデータが生成されていることの確認
	var la LoanApplication
	bytes, err := stub.GetState(loanApplicationID)  
	if err != nil {
		t.Fatalf("Could not fetch loan application with ID " + loanApplicationID)
	}
    // 上のCreateLoanApplicationでLedgerにデータが生成されていることの確認
    ////////////////////////////////////

	err = json.Unmarshal(bytes, &la)
	if err != nil {
		t.Fatalf("Could not unmarshal loan application with ID " + loanApplicationID)
	}
	var errors = []string{}
	var loanApplicationInput LoanApplication
	err = json.Unmarshal([]byte(loanApplication), &loanApplicationInput)
	if la.ID != loanApplicationInput.ID {
		errors = append(errors, "Loan Application ID does not match")
	}
	if la.PropertyId != loanApplicationInput.PropertyId {
		errors = append(errors, "Loan Application PropertyId does not match")
	}
	if la.PersonalInfo.Firstname != loanApplicationInput.PersonalInfo.Firstname {
		errors = append(errors, "Loan Application PersonalInfo.Firstname does not match")
	}
	//Can be extended for all fields
	if len(errors) > 0 {
		t.Fatalf("Mismatch between input and stored Loan Application")
		for j := 0; j < len(errors); j++ {
			fmt.Println(errors[j])
		}
	}
}
```

また、json.Unmarshal()などを使用するため、importに```encoding/json```を追加してください。

```sample_chaincode_test.go
import (
	"encoding/json" // <-追加
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"testing"
)
```

この状態で実行するとLoanApplicationという変数型が未定義というエラーになります。

```bash
$ go test
# _/home/ubuntu/workspace/sample_tdd
./sample_chaincode_test.go:68: undefined: LoanApplication
./sample_chaincode_test.go:78: undefined: LoanApplication
FAIL    _/home/ubuntu/workspace/sample_tdd [build failed]
```

そこで、下記のLoanApplicationなどの定義をsample_chaincode.goに追加します。

```sample_chaincode.go
type PersonalInfo struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	DOB       string `json:"DOB"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
}

type FinancialInfo struct {
	MonthlySalary      int `json:"monthlySalary"`
	MonthlyRent        int `json:"monthlyRent"`
	OtherExpenditure   int `json:"otherExpenditure"`
	MonthlyLoanPayment int `json:"monthlyLoanPayment"`
}

type LoanApplication struct {
	ID                     string        `json:"id"`
	PropertyId             string        `json:"propertyId"`
	LandId                 string        `json:"landId"`
	PermitId               string        `json:"permitId"`
	BuyerId                string        `json:"buyerId"`
	AppraisalApplicationId string        `json:"appraiserApplicationId"`
	SalesContractId        string        `json:"salesContractId"`
	PersonalInfo           PersonalInfo  `json:"personalInfo"`
	FinancialInfo          FinancialInfo `json:"financialInfo"`
	Status                 string        `json:"status"`
	RequestedAmount        int           `json:"requestedAmount"`
	FairMarketValue        int           `json:"fairMarketValue"`
	ApprovedAmount         int           `json:"approvedAmount"`
	ReviewerId             string        `json:"reviewerId"`
	LastModifiedDate       string        `json:"lastModifiedDate"`
}
```

これを実行するとinputする値がないのでエラーになるはずです。


```bash
$ go test
Entering TestCreateLoanApplication
2017/05/25 01:44:52 MockStub( mockStub &{} )
Entering TestCreateLoanApplicationValidation
2017/05/25 01:44:52 MockStub( mockStub &{} )
Entering CreateLoanApplication
Invalid number of args
Entering TestCreateLoanApplicationValidation2
2017/05/25 01:44:52 MockStub( mockStub &{} )
Entering CreateLoanApplication
Entering TestCreateLoanApplicationValidation3
2017/05/25 01:44:52 MockStub( mockStub &{} )
Entering CreateLoanApplication
2017/05/25 01:44:52 MockStub mockStub Getting la1 []
--- FAIL: TestCreateLoanApplicationValidation3 (0.00s)
        sample_chaincode_test.go:75: Could not unmarshal loan application with ID la1
FAIL
exit status 1
FAIL    _/home/ubuntu/workspace/sample_tdd      0.022s
```

そこで、レジャーにloan applicationを入力する記述にします。
sample_chaincode.goのCreateLaonApplicationに以下を記述。

```sample_chaincode.go
func CreateLoanApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering CreateLoanApplication")
	if len(args) < 2 {
		fmt.Println("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for loan application creation")
	}

    // 下記から追加
	var loanApplicationId = args[0]
	var loanApplicationInput = args[1]
	//TODO: Include schema validation here

	err := stub.PutState(loanApplicationId, []byte(loanApplicationInput))
	if err != nil {
		fmt.Println("Could not save loan application to ledger", err)
		return nil, err
	}

	fmt.Println("Successfully saved loan application")
	return []byte(loanApplicationInput), nil  <-- 戻り値を nil, nilから変更
}
```

これでloanApplicationIdとloanApplicationInputのJSON文字列を検索します。

PutStateによって、keyとvalueのペアを保管します。この場合はapplication IDがkeyであり、loan application JSON文字列がvalueになります。実行すると格納されるbyte配列が帰って正常に終了します。

テストを実行すると下記のようになります。

```bash
Entering TestCreateLoanApplication
2017/05/25 01:49:41 MockStub( mockStub &{} )
Entering TestCreateLoanApplicationValidation
2017/05/25 01:49:41 MockStub( mockStub &{} )
Entering CreateLoanApplication
Invalid number of args
Entering TestCreateLoanApplicationValidation2
2017/05/25 01:49:41 MockStub( mockStub &{} )
Entering CreateLoanApplication
2017/05/25 01:49:41 MockStub mockStub Putting la1 [xx xx xx xx ...]
2017/05/25 01:49:41 MockStub mockStub Key la1 is first element in list
Successfully saved loan application
Entering TestCreateLoanApplicationValidation3
2017/05/25 01:49:41 MockStub( mockStub &{} )
Entering CreateLoanApplication
2017/05/25 01:49:41 MockStub mockStub Putting la1 [xx xx xx xx ...]
2017/05/25 01:49:41 MockStub mockStub Key la1 is first element in list
Successfully saved loan application
2017/05/25 01:49:41 MockStub mockStub Getting la1 [xx xx xx xx ...]
PASS
ok      _/home/ubuntu/workspace/sample_tdd      0.028s
```


## 5. Invoke methodの実装

shim.Chaincode.Invokeの実装を行います。

### 5.1 要求

1. Invoke() は、入力関数名の引数をチェックし、適切なハンドラに実行を委任する必要があります。
2. Invoke() は、無効な入力関数名の場合にエラーを返す必要があります。
3. Invoke() は、チェーンコードの呼び出し側/呼び出し側のトランザクション証明書に基づいて、アクセス制御とアクセス管理を実装/委任する必要があります。 CreateLoanApplicationメソッドを呼び出すには、Bank_Adminだけを許可する必要があります。

### 5.2 Functionality outlined in Requirement 3

最初のテストでは、上記の要件3で概説した機能を検証します。

まずは、sample_chaincode_test.goに以下のようにusernameやroleを定義する。

```sample_chaincode_test.go
func TestInvokeValidation(t *testing.T) {
	fmt.Println("Entering TestInvokeValidation")

	attributes := make(map[string][]byte)
	attributes["username"] = []byte("vojha24")
	attributes["role"] = []byte("client")

	stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	_, err := stub.MockInvoke("t123", "CreateLoanApplication", []string{loanApplicationID, loanApplication})
	if err == nil {
		t.Fatalf("Expected unauthorized user error to be returned")
	}

}
```

caller/invokerの権限属性を上記のようにユーザーで定義ができます。
13行目にMockInvokeの呼び出し方が記述されており、transaction ID, function名、そしてinputの引数でコールします。

しかし当然ながらこのままではエラーとなります。実行結果は以下の通り。

```bash
--- FAIL: TestInvokeValidation (0.00s)
        sample_chaincode_test.go:112: Expected unauthorized user error to be returned
FAIL
exit status 1
FAIL    _/home/ubuntu/workspace/sample_tdd      0.025s
```

sample_chaincode.goのInvoke()関数を下記のように修正しエラー応答を記述しておきます。

```sample_chaincode.go
func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("Entering Invoke")
    return nil, errors.New("unauthorized user")
}
```

テストの実行結果は正常実行されました。

```bash
$ go test
Entering TestInvokeValidation
2017/05/25 02:35:17 MockStub( mockStub &{} )
Entering Invoke
PASS
ok      _/home/ubuntu/workspace/sample_tdd      0.019s
```

### 5.3. Bank_Admin roleでのテスト

次にBank_Adminのrole権限を以下のように記述します。sample_chaincode_test.goに以下のようなTestInvokeValidation2を記述します。

```sample_chaincode_test.go
func TestInvokeValidation2(t *testing.T) {
	fmt.Println("Entering TestInvokeValidation")

	attributes := make(map[string][]byte)
	attributes["username"] = []byte("vojha24")
	attributes["role"] = []byte("Bank_Admin")

	stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	_, err := stub.MockInvoke("t123", "CreateLoanApplication", []string{loanApplicationID, loanApplication})
	if err != nil {
		t.Fatalf("Expected CreateLoanApplication to be invoked")
	}

}
```

sample_chaincode.goも以下のようにusernameとroleを読み込むようにし、権限の確認を行います。

```sample_chaincode.go
func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("Entering Invoke")
     
    ubytes, _ := stub.ReadCertAttribute("username")
    rbytes, _ := stub.ReadCertAttribute("role")
 
    username := string(ubytes)
    role := string(rbytes)
 
    if role != "Bank_Admin" {
        return nil, errors.New("caller with " + username + " and role " + role + " does not have 
         access to invoke CreateLoanApplication")
    }
    return nil, nil
}
```

これで実行すると正常に終了し、role権限のチェックが実装されました。

```
Entering TestInvokeValidation
2017/05/25 02:42:43 MockStub( mockStub &{} )
Entering Invoke
Entering TestInvokeValidation
2017/05/25 02:42:43 MockStub( mockStub &{} )
Entering Invoke
PASS
ok      _/home/ubuntu/workspace/sample_tdd      0.025s
```

### 5.4. Check input function name

次は関数名のチェックです。
sample_chaincode_test.goに以下の関数を追加します。

```sample_chaincode_test.go
func TestInvokeFunctionValidation(t *testing.T) {
    fmt.Println("Entering TestInvokeFunctionValidation")
 
    attributes := make(map[string][]byte)
    attributes["username"] = []byte("vojha24")
    attributes["role"] = []byte("Bank_Admin")
 
    stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
    if stub == nil {
        t.Fatalf("MockStub creation failed")
    }
 
    _, err := stub.MockInvoke("t123", "InvalidFunctionName", []string{})
    if err == nil {
        t.Fatalf("Expected invalid function name error")
    }
 
}
```

実行するとエラーになることがわかります。

```bash
--- FAIL: TestInvokeFunctionValidation (0.00s)
        sample_chaincode_test.go:150: Expected invalid function name error
FAIL
exit status 1
FAIL    _/home/ubuntu/workspace/sample_tdd      0.028s
```

そこで、sample_chaincode.goのInvokeに若干変更を行います。

```sample_chaincode.go
func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Entering Invoke")

	ubytes, _ := stub.ReadCertAttribute("username")
	rbytes, _ := stub.ReadCertAttribute("role")

	username := string(ubytes)
	role := string(rbytes)

	if role != "Bank_Admin" {
		return nil, errors.New("caller with " + username + " and role " + role + " does not have access to invoke CreateLoanApplication")
	}
	return nil, errors.New("Invalid function name") // <-- change nil to errors
}
```

これで実行するとTestInvokeFunctionValidationは成功しますが、
TestInvokeValidation2は失敗します。

```bash
Entering Invoke
--- FAIL: TestInvokeValidation2 (0.00s)
        sample_chaincode_test.go:131: Expected CreateLoanApplication to be invoked
Entering TestInvokeFunctionValidation
2017/05/25 03:57:34 MockStub( mockStub &{} )
Entering Invoke
FAIL
exit status 1
FAIL    _/home/ubuntu/workspace/sample_tdd      0.023s
```

### 5.5. Check the correct function name

正しく動くように、TestInvokeFunctionValidation2を追加します。

```sample_chaincode_test.go
func TestInvokeFunctionValidation2(t *testing.T) {
    fmt.Println("Entering TestInvokeFunctionValidation2")
 
    attributes := make(map[string][]byte)
    attributes["username"] = []byte("vojha24")
    attributes["role"] = []byte("Bank_Admin")
 
    stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
    if stub == nil {
        t.Fatalf("MockStub creation failed")
    }
 
    _, err := stub.MockInvoke("t123", "CreateLoanApplication", []string{})
    if err != nil {
        t.Fatalf("Expected CreateLoanApplication function to be invoked")
    }
 
}
```

このあとに実行すると正しくエラーメッセージが出力されていることが確認できます。
```bash
Entering Invoke
--- FAIL: TestInvokeValidation2 (0.00s)
        sample_chaincode_test.go:131: Expected CreateLoanApplication to be invoked
Entering TestInvokeFunctionValidation
2017/05/25 04:13:06 MockStub( mockStub &{} )
Entering Invoke
Entering TestInvokeFunctionValidation2
2017/05/25 04:13:06 MockStub( mockStub &{} )
Entering Invoke
--- FAIL: TestInvokeFunctionValidation2 (0.00s)
        sample_chaincode_test.go:168: Expected CreateLoanApplication function to be invoked
FAIL
exit status 1
FAIL    _/home/ubuntu/workspace/sample_tdd      0.021s
```

これでInvokeに正しい関数名を返すようにします。
sample_chaincode.goのInvokeの箇所を以下のように編集。

```sample_chaincode.go
func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("Entering Invoke")
 
    ubytes, _ := stub.ReadCertAttribute("username")
    rbytes, _ := stub.ReadCertAttribute("role")
 
    username := string(ubytes)
    role := string(rbytes)
 
    if role != "Bank_Admin" {
        return nil, errors.New("caller with " + username + " and role " + role + " does not have access to invoke CreateLoanApplication")
    }
     
    if function == "CreateLoanApplication" {
        return CreateLoanApplication(stub, args)
    }
    return nil, errors.New("Invalid function name. Valid functions ['CreateLoanApplication']")
}
```
このとき実際にCreateLoanApplicationメソッドが呼び出されInvokeされたかどうかのテストを行います。
理想的にはスパイオブジェクトを用いてやるべきですが、簡素化するためinvokeメソッドからのアウトプットから確認します。
TestInvokeFunctionValidation2を以下のように書き換えます。

```sample_chaincode_test.go
func TestInvokeFunctionValidation2(t *testing.T) {
	fmt.Println("Entering TestInvokeFunctionValidation2")

	attributes := make(map[string][]byte)
	attributes["username"] = []byte("vojha24")
	attributes["role"] = []byte("Bank_Admin")

	stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	// Add and modify the following lines
	bytes, err := stub.MockInvoke("t123", "CreateLoanApplication", []string{loanApplicationID, loanApplication})
	if err != nil {
		t.Fatalf("Expected CreateLoanApplication function to be invoked")
	}
	//A spy could have been used here to ensure CreateLoanApplication method actually got invoked.
	var la LoanApplication
	err = json.Unmarshal(bytes, &la)
	if err != nil {
		t.Fatalf("Expected valid loan application JSON string to be returned from CreateLoanApplication method")
	}

}
```

これでストを実行すれば下記のように正常終了するはずです。

```bash
Successfully saved loan application
PASS
ok      _/home/ubuntu/workspace/sample_tdd      0.024s
```

## 6. 非決定的な関数のテスト

[ここ](https://www.ibm.com/developerworks/cloud/library/cl-ibm-blockchain-chaincode-development-using-golang/index.html#N1028D)
に書かれているように Chaincode は決定的でなければいけません。

例を示します。 
4つのピアすべてがピアを検証している、4ピアのHyperledger Fabricベースのブロックチェーンネットワークを利用してください。

これは、トランザクションをブロックチェーンに書き込む必要がある場合は、
4人のピアがすべて、その元帳のローカルコピーでトランザクションを独立して実行することを意味します。
簡単に言えば、4つのピアのそれぞれは、同じ入力を持つ独立したチェーンコード機能を独立して実行し、ローカルレジ係状態を更新します。
このようにして、4人のピアはすべて同じ元帳状態になります。

したがって、ピアによってチェーンコードが4回実行されると、
同じ結果が得られ、同じ元帳状態になる必要があります。
これは確定的チェーンコードと呼ばれます。


下記のコードは、CreateLoanApplication関数の非決定的なバージョンを示しています。
つまり、同じ入力でこの関数を複数回実行すると、結果が異なることになります。

```sample_chaincode.go
func NonDeterministicFunction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering NonDeterministicFunction")
	//Use random number generator to generate the ID
	var random = rand.New(rand.NewSource(time.Now().UnixNano()))
	var loanApplicationID = "la1" + strconv.Itoa(random.Intn(1000))
	var loanApplication = args[0]
	var la LoanApplication
	err := json.Unmarshal([]byte(loanApplication), &la)
	if err != nil {
		fmt.Println("Could not unmarshal loan application", err)
		return nil, err
	}
	la.ID = loanApplicationID
	laBytes, err := json.Marshal(&la)
	if err != nil {
		fmt.Println("Could not marshal loan application", err)
		return nil, err
	}
	err = stub.PutState(loanApplicationID, laBytes)
	if err != nil {
		fmt.Println("Could not save loan application to ledger", err)
		return nil, err
	}

	fmt.Println("Successfully saved loan application")
	return []byte(loanApplicationID), nil
}
```

import には下記を追加してください。

```sample_chaincode.go
import (
	"encoding/json" // <-- add
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"math/rand"     // <-- add
	"strconv"       // <-- add
	"time"          // <-- add
)
```

上記のメソッドは、ローンアプリケーションIDを入力の一部として渡した
元のCreateLoanApplicationメソッドとは異なり、
乱数ジェネレータを使用してIDを生成し、渡されたローンアプリケーションコンテンツに追加します。

4行目と5行目は、ローンアプリケーションIDの生成方法を示しています。
19行目は更新されたローン申請内容を元帳に保存します。

```sample_chaincode_test.go
func TestNonDeterministicFunction(t *testing.T) {
	fmt.Println("Entering TestNonDeterministicFunction")
	attributes := make(map[string][]byte)
	const peerSize = 4
	var stubs [peerSize]*shim.CustomMockStub
	var responses [peerSize][]byte
	var loanApplicationCustom = `{"propertyId":"prop1","landId":"land1","permitId":"permit1","buyerId":"vojha24","personalInfo":{"firstname":"Varun","lastname":"Ojha","dob":"dob","email":"varun@gmail.com","mobile":"99999999"},"financialInfo":{"monthlySalary":16000,"otherExpenditure":0,"monthlyRent":4150,"monthlyLoanPayment":4000},"status":"Submitted","requestedAmount":40000,"fairMarketValue":58000,"approvedAmount":40000,"reviewedBy":"bond","lastModifiedDate":"21/09/2016 2:30pm"}`
	//Simulate execution of the chaincode function by multiple peers on their local ledgers
	for j := 0; j < peerSize; j++ {
		stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
		if stub == nil {
			t.Fatalf("MockStub creation failed")
		}
		stub.MockTransactionStart("tx" + string(j))
		resp, err := NonDeterministicFunction(stub, []string{loanApplicationCustom})
		if err != nil {
			t.Fatalf("Could not execute NonDeterministicFunction ")
		}
		stub.MockTransactionEnd("tx" + string(j))
		stubs[j] = stub
		responses[j] = resp
	}

	for i := 0; i < peerSize; i++ {
		if i < (peerSize - 1) {
			la1Bytes, _ := stubs[i].GetState(string(responses[i]))
			la2Bytes, _ := stubs[i+1].GetState(string(responses[i+1]))
			la1 := string(la1Bytes)
			la2 := string(la2Bytes)
			if la1 != la2 {
				//TODO: Compare individual values to find mismatch
				t.Fatalf("Expected all loan applications to be identical. Non Deterministic chaincode error")
			}
		}
		//All loan applications retrieved from each of the peer's ledger's match. Function is deterministic

	}

}
```