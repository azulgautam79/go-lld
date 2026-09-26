package main

// package strategy

import (
	"fmt"
	"strings"
)

// Strategy Interface
type PaymentStrategy interface {
	Pay(amount float64)
	Validate() bool
}

// Concrete Strategies
type CreditCardPayment struct {
	cardNumber string
	cvv        string
}

func NewCreditCardPayment(cardNumber, cvv string) *CreditCardPayment {
	return &CreditCardPayment{
		cardNumber: cardNumber,
		cvv:        cvv,
	}
}

func (c *CreditCardPayment) Validate() bool {
	return len(c.cardNumber) == 16 && len(c.cvv) == 3
}

func (c *CreditCardPayment) Pay(amount float64) {
	fmt.Printf("Paid $%.2f using Credit Card\n", amount)
}

// -----------------------------------------
type PaypalPayment struct {
	email string
}

func NewPayPalPayment(email string) *PaypalPayment {
	return &PaypalPayment{
		email: email,
	}
}

func (p *PaypalPayment) Validate() bool {
	return strings.Contains(p.email, "@")
}

func (p *PaypalPayment) Pay(amount float64) {
	fmt.Printf("Paid $%.2f using Paypal: %s\n", amount, p.email)
}

// -----------------------------------------
type UPIPayment struct {
	upiID string
}

func NewUPIPayment(upiID string) *UPIPayment {
	return &UPIPayment{
		upiID: upiID,
	}
}

func (u *UPIPayment) Validate() bool {
	return strings.Contains(u.upiID, "@")
}

func (u *UPIPayment) Pay(amount float64) {
	fmt.Printf("Paid $%.2f using UPI: %s\n", amount, u.upiID)
}

// ========== Item ==========

type Item struct {
	name  string
	price float64
}

func NewItem(name string, price float64) Item {
	return Item{
		name:  name,
		price: price,
	}
}

// ========== Context ==========
type ShoppingCart struct {
	items           []Item
	paymentStrategy PaymentStrategy
}

func (c *ShoppingCart) AddItem(item Item) {
	c.items = append(c.items, item)
}

func (c *ShoppingCart) CalculateTotal() float64 {
	var total float64

	for _, item := range c.items {
		total += item.price
	}

	return total
}

// Strategy can be set/changed at runtime
func (c *ShoppingCart) SetPaymentStrategy(strategy PaymentStrategy) {
	c.paymentStrategy = strategy
}

func (c *ShoppingCart) Checkout() {
	total := c.CalculateTotal()

	if c.paymentStrategy == nil {
		fmt.Println("Payment strategy not set")
		return
	}

	if c.paymentStrategy.Validate() {
		c.paymentStrategy.Pay(total)
	} else {
		fmt.Println("Payment validation failed")
	}
}

// ========== Usage ==========
func main() {

	cart := &ShoppingCart{}

	cart.AddItem(NewItem("Laptop", 999.99))
	cart.AddItem(NewItem("Mouse", 29.99))

	// Pay with Credit Card
	cart.SetPaymentStrategy(
		NewCreditCardPayment("1234567890123456", "123"),
	)
	cart.Checkout()

	// Same cart, different payment method
	cart.SetPaymentStrategy(
		NewPayPalPayment("user@gmail.com"),
	)
	cart.Checkout()

	// UPI
	cart.SetPaymentStrategy(
		NewUPIPayment("user@okbank"),
	)
	cart.Checkout()
}
