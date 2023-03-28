package wallet

import (
	"testing"
)

func TestWallet_Withdraw(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		initAmount Bitcoin
		amount     Bitcoin
		expected   Bitcoin
		wantError  bool
	}{
		{name: "Withdraw equal amount", initAmount: 100.0, amount: 100.0, expected: 0.0, wantError: false},
		{name: "Withdraw lower amount", initAmount: 100.0, amount: 10.0, expected: 90.0, wantError: false},
		{name: "Withdraw zero amount", initAmount: 100.0, amount: 0.0, expected: 100.0, wantError: true},
		{name: "Withdraw from zero balance", initAmount: 0.0, amount: 100.0, expected: 0.0, wantError: true},
		{name: "Withdraw negative amount", initAmount: 100.0, amount: -10.0, expected: 100.0, wantError: true},
		{name: "Withdraw large amount", initAmount: 1e11, amount: 1e10, expected: 1e11 - 1e10, wantError: false},
	}

	for _, test := range tests {
		test := test // Capture range variable.
		t.Run(test.name, func(t *testing.T) {
			wallet := &Wallet{balance: test.initAmount}
			err := wallet.Withdraw(test.amount)
			if test.wantError {
				if err == nil {
					t.Errorf("An error was expected, but did not occur")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error occured: %v", err)
				}
				if wallet.Balance() != test.expected {
					t.Errorf("Balance is incorrect: expected %f, but got %f", test.expected, wallet.Balance())
				}
			}
		})
	}
}

func TestWallet_Deposit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		amount    Bitcoin
		expected  Bitcoin
		wantError bool
	}{
		{name: "Deposit positive amount", amount: 100.0, expected: 100.0, wantError: false},
		{name: "Deposit negative amount", amount: -10.0, expected: 0.0, wantError: true},
		{name: "Deposit zero amount", amount: 0.0, expected: 0.0, wantError: true},
		{name: "Deposit large amount", amount: 1e10, expected: 1e10, wantError: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet := &Wallet{balance: 0.0}
			err := wallet.Deposit(test.amount)
			if test.wantError {
				if err == nil {
					t.Errorf("An error was expected, but did not occur")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error occured: %v", err)
				}
				if wallet.Balance() != test.expected {
					t.Errorf("Balance is incorrect: expected %f, but got %f", test.expected, wallet.Balance())
				}
			}
		})
	}

}
