# TSTR

Simplify system tests.

## How to use

### Dependencies

- `Go` installed
- `dep` dependency manager

### Install it

```shell
$ go get github.com/nmeji/tstr
```

### Project Specific

```shell
$ dep ensure -add github.com/nmeji/tstr
```

## Features

### Environment

If you only need to access a few env vars, you can just get them straight away.

```go
import "github.com/nmeji/tstr/env"
```

#### Example

```go
env := env.StrVal("ENV")
fmt.Print(env)  // dev, qa, preprod
```

```go
port, err := env.IntVal("PORT")
fmt.Print(port)  // 443, 80
```

```go
toggle, err := env.BoolVal("TOGGLE")
fmt.Print(toggle)  // true, false
```

non-nil `err` means it failed converting the type to int or bool

### Test Data

Simplifies unmarshalling of test data/fixtures. Import only one package to unmarshal any of the ff file formats: `json, toml, csv`

```go
import "github.com/nmeji/tstr/testdata"
```

#### Basic Usage

```go
data, err := testdata.New("testdata/1.csv")
/*
 non-nil err means there is problem reading the test fixture
*/
```

### JSON

This uses the same unmarshaling as what Go provides.

```json
{
    "request-id": "D12837981",
    "payload": "example"
}
```

```go
data, _ := testdata.New("testdata/1.json")
testInput := struct{
    ID string `json:"request-id"`,
    Payload string,
} {}
err := data.Unmarshal(&testInput)

/*
 {ID:D12837981 Payload:example}
*/
```

### CSV

What Go provides is unmarshalling csv files to [][]string. What's new here is that you can unmarshal to array/slice of structs.

```csv
color,hex,my_comments
RED,#FF0000,red ferrari
GREEN,#00FF00,green-minded
BLUE,#0000FF,blue sky
```

```go
data, _ := testdata.New("testdata/1.csv")

result := []struct {
    Color string
    Hex   string
    Comment string `csv:"my_comments"`
}{}
err := data.Unmarshal(&result)

/*
[
    {Color:RED Hex:#FF0000 Comment:red ferrari}
    {Color:GREEN Hex:#00FF00 Comment:green-minded}
    {Color:BLUE Hex:#0000FF Comment:blue sky}
]
*/
```

### TOML

This also provides support for TOML files.

```toml
[qa]
id ="180b3347-210d-443f-92dc-369f51752188"
customer_idp_id = "3d234aad-a6f4-41de-aaa4-20cf762ec6da"
user_email = "sxaykvtiwco@mailinator.com"
vertical = "mobile"

[preprod]
id ="295c6316-0297-4b7d-a8a0-53dc967fba20"
customer_idp_id = "b1c47b13-bb71-41eb-b7a1-9bb6fac0db8c"
user_email = "rgfvocyjaim@mailinator.com"
vertical = "mobile"
```

```go
data, _ := testdata.New("testdata/1.toml")

type env struct {
    ID        string
    UserEmail string `toml:"user_email"`
    Vertical  string
}
var result struct{ QA env }  // change the property name from `QA` to `Preprod` if you want [preprod]
err := data.Unmarshal(&result)

/*
 {ID:180b3347-210d-443f-92dc-369f51752188 UserEmail:sxaykvtiwco@mailinator.com Vertical:mobile}
*/
```

### Response Assertions

```go
ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `{
        "books": [
                {"id": 1, "price": 88.94},
                {"id": 2, "price": 97.10},
                {"id": 3, "price": 11.87}
            ],
        "author": "niko"
    }`)
}))
defer ts.Close()

checker, err := request.Get(ts.URL)
if err != nil {
    t.Errorf("expected nil error, got %v", err)
}
books := []struct {
    ID    int
    Price float64
}{
    {ID: 2, Price: 97.10},
    {ID: 3, Price: 11.87},
}
checker.
    ExpectStatus(200).
    ExpectBody.ToHaveInJson("$.books[1,2]", books).
    ExpectBody.ToHaveInJson("$.author", "niko").
    ExpectBody.ToHaveInJson("$.books[0].price", 88.94).
    MakeAssertion(t)
```

## TODO

- [ ] Request Helpers
- [ ] More Assertion Matchers
    - [ ] Greater/Less Than Value
    - [ ] Range Values
    - [ ] List of Values
    - [ ] Regex Matcher
    - [ ] Negative Asserts (should not match [matcher])
- [ ] Expect Header
- [ ] (Test Data) XML Unmarshaller
- [ ] SOAP Assertion
    - [ ] Support Xpath filter

## License

[MIT](LICENSE)