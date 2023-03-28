package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWallet_Withdraw(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		initAmount Bitcoin
		amount     Bitcoin
		expected   Bitcoin
		wantError  error
	}{
		{
			name:       "Withdraw equal amount",
			initAmount: 100.0,
			amount:     100.0,
			expected:   0.0,
			wantError:  nil,
		},
		{
			name:       "Withdraw lower amount",
			initAmount: 100.0,
			amount:     10.0,
			expected:   90.0,
			wantError:  nil,
		},
		{
			name:       "Withdraw zero amount",
			initAmount: 100.0,
			amount:     0.0,
			expected:   100.0,
			wantError:  ErrNonPositiveAmount,
		},
		{
			name:       "Withdraw from zero balance",
			initAmount: 0.0,
			amount:     100.0,
			expected:   0.0,
			wantError:  ErrInsufficientBalance,
		},
		{
			name:       "Withdraw negative amount",
			initAmount: 100.0,
			amount:     -10.0,
			expected:   100.0,
			wantError:  ErrNonPositiveAmount,
		},
		{
			name:       "Withdraw large amount",
			initAmount: 1e11,
			amount:     1e10,
			expected:   1e11 - 1e10,
			wantError:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet := NewWallet(test.initAmount)
			err := wallet.Withdraw(test.amount)
			assert.ErrorIs(t, err, test.wantError, "wallet.Withdraw(%v)", test.amount)
			assert.Equal(t, test.expected, wallet.Balance(), "wallet.Withdraw(%v)", test.amount)
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
