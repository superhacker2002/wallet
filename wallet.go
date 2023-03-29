package wallet

import (
	"errors"
	"fmt"
	"sync"
)

type Bitcoin float64

var ErrNonPositiveAmount = errors.New("negative or zero amount")
var ErrInsufficientBalance = errors.New("insufficient balance")

type Wallet struct {
	balance Bitcoin
	mutex   sync.RWMutex
}

func NewWallet(balance Bitcoin) *Wallet {
	return &Wallet{balance: balance}
}

func (w *Wallet) Deposit(amount Bitcoin) error {
	if amount <= 0 {
		return ErrNonPositiveAmount
	}
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.balance += amount
	return nil
}

func (w *Wallet) Withdraw(amount Bitcoin) error {
	if amount <= 0 {
		return ErrNonPositiveAmount
	}
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.balance < amount {
		return fmt.Errorf("%w: %f", ErrInsufficientBalance, w.balance)
	}
	w.balance -= amount
	return nil
}

func (w *Wallet) Balance() Bitcoin {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.balance
}
