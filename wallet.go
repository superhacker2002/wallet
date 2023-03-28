package wallet

import (
	"fmt"
	"sync"
)

type Bitcoin float64

type Wallet struct {
	balance Bitcoin
	mutex   sync.Mutex
}

func (w *Wallet) Deposit(amount Bitcoin) error {
	if amount <= 0 {
		return fmt.Errorf("impossible to deposit negative or zero amount of bitcoins: %f", amount)
	}
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.balance += amount
	return nil
}

func (w *Wallet) Withdraw(amount Bitcoin) error {
	if amount <= 0 {
		return fmt.Errorf("impossible to withdraw negative or zero amount of bitcoins: %f", amount)
	}
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.balance < amount {
		return fmt.Errorf("insufficient funds on the wallet balance: %f", w.balance)
	}
	w.balance -= amount
	return nil
}

func (w *Wallet) Balance() Bitcoin {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.balance
}
