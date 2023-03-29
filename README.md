
# Wallet

Implementation of Wallet type that holds bitcoins and allows the following operations: to deposit and withdraw money and get the wallet's balance.

## Usage


Firstly, run this command in order to load module and set up dependencies:

```go
go get "github.com/superhacker2002/wallet"
```

Then import module in your program:
```
import "github.com/superhacker2002/wallet"
```

And use wallet:
```
package main

import "github.com/superhacker2002/wallet"

func main() {
	wallet := wallet.NewWallet(10.0)
	wallet.Withdraw(5.0)
}
```

## Running the tests

Run unit tests using 
```
go test -v
```
