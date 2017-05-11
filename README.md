Go 環境設定

# 1. Homebrewでのインストール

```
$ brew install go
==> Downloading https://homebrew.bintray.com/bottles/go-1.6.3.el_capitan.bottle.1.tar.gz
######################################################################## 100.0%
==> Pouring go-1.6.3.el_capitan.bottle.1.tar.gz
==> Caveats
As of go 1.2, a valid GOPATH is required to use the `go get` command:
  https://golang.org/doc/code.html#GOPATH

You may wish to add the GOROOT-based install location to your PATH:
  export PATH=$PATH:/usr/local/opt/go/libexec/bin
==> Summary
🍺  /usr/local/Cellar/go/1.6.3: 5,778 files, 328.4M
% which go
/usr/local/bin/go
```

# 2. Version確認

```
$ go version
go version go1.5.1 darwin/amd64
```
# 3. 稼働確認
go run サブコマンドを使うことでソースコードをビルドすると同時に実行する

```
$ vi helloworld.go
package main

import (
  "fmt"
)

func main() {
  fmt.Println("Hello, World!")
}
```

# 4. GOPATH環境設定

[こちらを参考にしました。](http://blog.amedama.jp/entry/2015/10/06/231038)
[こちらも参考にしました。](http://www.songmu.jp/riji/entry/prepare-golang-development-environment-on-mac.html)

$GOPATH だけを決めればOKです。$GOPATH はどこでも良いので、ここでは $HOME/go を $GOPATH に設定します。

```
$ export GOPATH=$HOME/go
$ export PATH=$PATH:$GOPATH/bin
```

必要に応じてシェル起動を定義

```
$ echo 'export GOPATH=$HOME/go >> ~/.bashrc
$ echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
```

$GOPATHで指定したディレクトリの作成

```
$ mkdir -p ${GOPATH}
```


# 5. サンプルのチェーンコードをクローンする

```
# Create the parent directories on your GOPATH
$ mkdir -p $GOPATH/src/github.com/hyperledger
$ cd $GOAPTH/src/github.com/hyperledger

# Clone the appropriate release codebase into $GOPATH/src/github.com/hyperledger/fabric
# Note that the v0.5 release is a branch of the repository.  It is defined below after the -b argument
$ git clone -b v0.6 http://gerrit.hyperledger.org/r/fabric
```
失敗したのでこちら

```
$ git clone -b v0.6.1-preview https://github.com/hyperledger/fabric.git
```

失敗したのでこちら

```
$ git clone -b v0.6 http://gerrit.hyperledger.org/r/fabric
```

Buildしてエラーがないか確認します。

```
$ go build
```

# 6. テストドリブン開発の事前準備

ソースをDownloadする。
varunmockstub.go を $GOPATH/src/github.com/Hyperledger/fabric/core/chaincode/shim/ にコピーする。

Workディレクトリを作成

```
$ cd worktemp
$ mkdir sample_tdd
$ cd sample_tdd
```

適当なエディタを開いて、以下2つのファイルを作成します。

1. sample_chaincode_test.go 
2. sample_chaincode.go

sample_chaincode_test.goを以下のように編集します。

```
package main
import (
    // "encoding/json" //コメントアウト
    "fmt"
    "testing"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)
```
testingパッケージをimportしているのですが、これはGoパッケージの自動テストを実装するために行います。

[Goの自動テストについての参考URLはこちらから。](http://golang.jp/pkg/testing)

さらに先程コピーしたテスト用のスタブファイル（CustomMockStub）とchaincodeを実装するために、shimをimportしています。

# 7. CreateLoanApplicationの実装

loan application IDやloan applicationをinputとして、ChaincodeStubInterfaceを使用します。

sample_chaincode_test.goを以下のように編集します。

```
package main
import (
//    "encoding/json" //コメントアウト
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

loan applicationが正しく作られたか、エラーの場合はエラーメッセージを出力するようにします。

Golang testing packageを実行するために、functionの名前は必ずTest* にします。

この状態でgo testを実行して稼働させてみます。

```
$ go test
 can't load package: package .:
 sample_chaincode.go:1:1:1 expected 'package', found 'EOF'
```

sample_chaincode.goに何も入っていないためエラーになったというメッセージです。
それでは、以下のとおりにsample_chaincode.go を編集します。

```
package main
```

この状態で実行してみます。

```
$ go test
./sample_chaincode_test.go:13: undefined: SampleChaincode
```

SampleChaincodeの定義がないのでエラーとなりました。そこで、以下のようにSampleChaincodeを実装します。

```
package main

type SampleChaincode struct {
}
```

ここで実行してみると以下のようなエラーになります。

```
$ go test
./sample_chaincode_test.go:13: cannot use new(SampleChaincode) (type *SampleChaincode) as type shim.Chaincode in argument to shim.NewCustomMockStub:
	*SampleChaincode does not implement shim.Chaincode (missing Init method)
```

chaincodeのInit, Query,Invokeなどの関数を定義するshim.Chaincodeが実装されていないためです。

そこで以下の通りに必要なshimのimportを含めて実装します。


```sample_chaincode.go
package main
import (
  //  "encoding/json"
  //  "fmt"
  //  "testing"
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

テストを実行すると正しくアプリケーションが作成されました！

```
$ go test
Entering TestCreateLoanApplication
2017/05/11 14:05:09 MockStub( mockStub &{} )
PASS
```

CreateLoanApplicationも実装します。以下をsample_chaincode.goに追加します。

```sample_chaincode.go
package main
import (
  //  "encoding/json"
    "fmt"
  //  "testing"
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

func CreateLoanApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering CreateLoanApplication")
    return nil, nil
}
```

sample_chaincode_test.goに以下の通りに、TestCreateLoanApplicationValidationの関数を追加します。ただし、
CreateLoanApplication method が呼び込まれるまえにトランザクションを開始する必要があります。
なぜなら、CreateLoanApplication method はloanアプリケーションをレジャーに保存するためです。

```sample_chaincode_test.go
func TestCreateLoanApplication (t *testing.T) {
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
```

この状態で実行してみると、予想どおりvalidationエラーになりました。

```
$ go test
Entering TestCreateLoanApplication
2017/05/11 14:34:37 MockStub( mockStub &{} )
Entering TestCreateLoanApplicationValidation
2017/05/11 14:34:37 MockStub( mockStub &{} )
Entering CreateLoanApplication
--- FAIL: TestCreateLoanApplicationValidation (0.00s)
	sample_chaincode_test.go:30: Expected CreateLoanApplication to return validation error
FAIL
exit status 1
```

以下のように、sample_chaincode.goのCreateLoanApplicationの返り値にStringのエラーメッセージを返すように修正します。
このとき、errorsをimportしておきます。

```
package main
import (
    "errors"
  //  "encoding/json"
    "fmt"
  //  "testing"
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

func CreateLoanApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering CreateLoanApplication")
    return nil, errors.New("Expected atleast two arguments for loan application creation")
}
```

ここでテストを実行します。

```
$ go test
Entering TestCreateLoanApplication
2017/05/11 15:10:31 MockStub( mockStub &{} )
Entering TestCreateLoanApplicationValidation
2017/05/11 15:10:31 MockStub( mockStub &{} )
Entering CreateLoanApplication
PASS
ok
```

正常に終了することが確認できました。

このままではエラーのメッセージを返すだけなので、別のテストを実行するために以下の関数を追加します。


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
1行目、2行目でloan applicationのデータを生成しています。これで実行してみます。


```
$ go test
Entering TestCreateLoanApplication
2017/05/11 20:10:14 MockStub( mockStub &{} )
Entering TestCreateLoanApplicationValidation
2017/05/11 20:10:14 MockStub( mockStub &{} )
Entering CreateLoanApplication
Entering TestCreateLoanApplicationValidation2
2017/05/11 20:10:14 MockStub( mockStub &{} )
Entering CreateLoanApplication
--- FAIL: TestCreateLoanApplicationValidation2 (0.00s)
	sample_chaincode_test.go:49: Expected CreateLoanApplication to succeed
FAIL
exit status 1
```

案の定エラーになるので、以下のように書き換えます。


```sample_chaincode.go
・・・・
func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    return nil, nil
}

func CreateLoanApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering CreateLoanApplication")
    if len(args) < 2 {
        fmt.Println("Invalid number of args")
        return nil, errors.New("Expected atleast two arguments for loan application creation")
    }
    return nil, nil
}
```

これで実行すれば、テストがPASSされます。

```
$ go test
Entering TestCreateLoanApplication
2017/05/11 20:23:28 MockStub( mockStub &{} )
Entering TestCreateLoanApplicationValidation
2017/05/11 20:23:28 MockStub( mockStub &{} )
Entering CreateLoanApplication
Invalid number of args
Entering TestCreateLoanApplicationValidation2
2017/05/11 20:23:28 MockStub( mockStub &{} )
Entering CreateLoanApplication
PASS
ok
```

このあとにloan applicationが生成され、Blockchainに書き込まれるかをテストします。
以下のようにsample_chaincode_test.goに記述します。

```
package main
import (
    "encoding/json" 
    "fmt"
    "testing"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

・・・・（途中省略）

func TestCreateLoanApplicationValidation3(t *testing.T) {
    fmt.Println("Entering TestCreateLoanApplicationValidation3")
    attributes := make(map[string][]byte)
    stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
    if stub == nil {
        t.Fatalf("MockStub creation failed")
    }
 
    stub.MockTransactionStart("t123")
    CreateLoanApplication(stub, []string{loanApplicationID, loanApplication})
    stub.MockTransactionEnd("t123")
 
    var la LoanApplication
    bytes, err := stub.GetState(loanApplicationID)
    if err != nil {
        t.Fatalf("Could not fetch loan application with ID " + loanApplicationID)
    }
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
1-12行目は前述と同様にstubのセットアップ。14行目は10行目のinvokedで成城通り作成されたloan application objectを検索します。
stub.GetState(loanApplicationID) はkeyに対応したバイト配列値を検索します。
この場合はloan application IDをレジャーから検索します。
18行目では検索されたバイト配列をLoanApplicationに戻しています。
以下のようなエラーになります。

```
$ go test
# _/Users/morizumiyuusuke/Documents/sample_tdd
./sample_chaincode_test.go:67: undefined: LoanApplication
./sample_chaincode_test.go:77: undefined: LoanApplication
FAIL
```

LoanApplicationをsample_chaincode.goに記述します。

```sample_chaincode.go
package main
import (
    "errors"
//    "encoding/json"
    "fmt"
  //  "testing"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

type SampleChaincode struct {
}

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

func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    return nil, nil
}
・・・
```
これを実行するとinputする値がないのでエラーになるはずです。


```
$ go test
Entering TestCreateLoanApplication
2017/05/11 21:26:41 MockStub( mockStub &{} )
Entering TestCreateLoanApplicationValidation
2017/05/11 21:26:41 MockStub( mockStub &{} )
Entering CreateLoanApplication
Invalid number of args
Entering TestCreateLoanApplicationValidation2
2017/05/11 21:26:41 MockStub( mockStub &{} )
Entering CreateLoanApplication
Entering TestCreateLoanApplicationValidation3
2017/05/11 21:26:41 MockStub( mockStub &{} )
Entering CreateLoanApplication
2017/05/11 21:26:41 MockStub mockStub Getting la1 []
--- FAIL: TestCreateLoanApplicationValidation3 (0.00s)
	sample_chaincode_test.go:74: Could not unmarshal loan application with ID la1
FAIL
exit status 1
FAIL
```

そこで、レジャーにloan applicationを入力する記述にします。
sample_chaincode.goのCreateLaonApplicationに以下を記述。

```
・・・
func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    return nil, nil
}

func CreateLoanApplication(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Println("Entering CreateLoanApplication")
 
    if len(args) < 2 {
        fmt.Println("Invalid number of args")
        return nil, errors.New("Expected atleast two arguments for loan application creation")
    }
 
    var loanApplicationId = args[0]
    var loanApplicationInput = args[1]
    //TODO: Include schema validation here
 
    err := stub.PutState(loanApplicationId, []byte(loanApplicationInput))
    if err != nil {
        fmt.Println("Could not save loan application to ledger", err)
        return nil, err
    }
 
    fmt.Println("Successfully saved loan application")
    return []byte(loanApplicationInput), nil
 
}
```

これでloanApplicationIdとloanApplicationInputのJSON文字列を検索します。

PutStateによって、keyとvalueのペアを保管します。この場合はapplication IDがkeyであり、loan application JSON文字列がvalueになります。実行すると格納されるbyte配列が帰って正常に終了します。

# 8. invoke methodの実装


invokeの主な役割としてはしかるべき権限を持った人かのチェックや、正しいfunctionの名前で実行されているかなどをチェックを行う。
まずは、sample_chaincode_test.goに以下のようにusernameやroleを定義する。
```
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

```
--- FAIL: TestInvokeValidation (0.00s)
	sample_chaincode_test.go:111: Expected unauthorized user error to be returned
FAIL
exit status 1
```

sample_chaincode.goにInvokeの関数を修正しエラー応答を記述しておきます。

```
func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("Entering Invoke")
    return nil, errors.New("unauthorized user")
}
```

テストの実行結果は正常実行されました。


```
$ go test
Entering TestInvokeValidation
2017/05/12 00:35:37 MockStub( mockStub &{} )
Entering Invoke
PASS
ok
```

次にBank_Adminのrole権限を以下のように記述します。sample_chaincode_test.goに以下のようなTestInvokeValidation2を記述します。


```
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

```
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

```
Entering Invoke
--- FAIL: TestInvokeFunctionValidation (0.00s)
	sample_chaincode_test.go:149: Expected invalid function name error
FAIL
exit status 1
```

そこで、sample_chaincode.goのInvokeに若干変更を行います。

```
func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("Entering Invoke")
 
    ubytes, _ := stub.ReadCertAttribute("username")
    rbytes, _ := stub.ReadCertAttribute("role")
 
    username := string(ubytes)
    role := string(rbytes)
 
    if role != "Bank_Admin" {
        return nil, errors.New("caller with " + username + " and role " + role + " does not have access to invoke CreateLoanApplication")
    }
 
    return nil, errors.New("Invalid function name")
}
```
これで実行するとInvalidな名前の関数が指定されたとメッセージが出力されます。
正しく動くように、TestInvokeFunctionValidation2を追加します。

```
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
```
Entering TestInvokeFunctionValidation2
2017/05/12 01:27:49 MockStub( mockStub &{} )
Entering Invoke
--- FAIL: TestInvokeFunctionValidation2 (0.00s)
	sample_chaincode_test.go:168: Expected CreateLoanApplication function to be invoked
FAIL
exit status 1
FAIL
```

これでInvokeに正しい関数名を返すようにします。
sample_chaincode.goのInvokeの箇所を以下のように編集。

```
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
理想的にはスパイオブジェクトを用いてやるべきですが、簡素化するためinvokeメソッドからのアウトプットから確認します。TestInvokeFunctionValidation2を以下のように書き換えます。

```
func TestInvokeFunctionValidation2(t *testing.T) {
    fmt.Println("Entering TestInvokeFunctionValidation2")
 
    attributes := make(map[string][]byte)
    attributes["username"] = []byte("vojha24")
    attributes["role"] = []byte("Bank_Admin")
 
    stub := shim.NewCustomMockStub("mockStub", new(SampleChaincode), attributes)
    if stub == nil {
        t.Fatalf("MockStub creation failed")
    }
 
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
