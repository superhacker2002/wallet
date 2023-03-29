package wallet

import (
	"math/rand"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

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
		wantError error
	}{
		{
			name:      "Deposit positive amount",
			amount:    100.0,
			expected:  100.0,
			wantError: nil,
		},
		{
			name:      "Deposit negative amount",
			amount:    -10.0,
			expected:  0.0,
			wantError: ErrNonPositiveAmount,
		},
		{
			name:      "Deposit zero amount",
			amount:    0.0,
			expected:  0.0,
			wantError: ErrNonPositiveAmount,
		},
		{
			name:      "Deposit large amount",
			amount:    1e10,
			expected:  1e10,
			wantError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet := NewWallet(0.0)
			err := wallet.Deposit(test.amount)
			assert.ErrorIs(t, err, test.wantError, "wallet.Withdraw(%v)", test.amount)
			assert.Equal(t, test.expected, wallet.Balance(), "wallet.Withdraw(%v)", test.amount)
		})
	}
}

func TestWallet_Balance(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		initAmount Bitcoin
	}{
		{
			name:       "Zero balance",
			initAmount: 0.0,
		},
		{
			name:       "Positive balance",
			initAmount: 100.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet := NewWallet(test.initAmount)
			assert.Equal(t, test.initAmount, wallet.Balance())
		})
	}
}

func TestWallet_ParallelWithdraw(t *testing.T) {
	const goroutinesCount = 100
	wallet := NewWallet(1000.0)
	g := new(errgroup.Group)
	for i := 0; i < goroutinesCount; i++ {
		g.Go(func() error {
			return wallet.Withdraw(10.0)
		})
	}
	err := g.Wait()
	assert.NoError(t, err, "one of withdraw returns error")
	assert.Equal(t, Bitcoin(0.0), wallet.Balance())
}

func TestWallet_ParallelDeposit(t *testing.T) {
	const goroutinesCount = 100
	wallet := NewWallet(0.0)
	g := new(errgroup.Group)
	for i := 0; i < goroutinesCount; i++ {
		g.Go(func() error {
			return wallet.Deposit(10.0)
		})
	}
	err := g.Wait()
	assert.NoError(t, err, "one of deposit returns error")
	assert.Equal(t, Bitcoin(1000.0), wallet.Balance())
}

func TestWallet_ParallelRandomMethods(t *testing.T) {
	const goroutinesCount = 100
	balance := Bitcoin(1000.0)
	wallet := NewWallet(balance)
	g := new(errgroup.Group)

	seed := time.Now().UnixNano()
	rand.Seed(seed)
	t.Log(seed)

	for i := 0; i < goroutinesCount; i++ {
		if function := rand.Intn(1); function == 0 {
			balance += 10.0
			g.Go(func() error {
				return wallet.Deposit(10.0)
			})
		} else {
			balance -= 10.0
			g.Go(func() error {
				return wallet.Withdraw(10.0)
			})
		}
	}
	err := g.Wait()
	assert.NoError(t, err, "one of deposit or withdraw returns error")
	assert.Equal(t, balance, wallet.Balance())
}
