package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type MilestonePaymentContract struct {
	contractapi.Contract
}

type Milestone struct {
	MilestoneID string  `json:"milestoneID"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
}

type Order struct {
	OrderID    string              `json:"orderID"`
	Milestones map[string]Milestone `json:"milestones"`
}

func (s *MilestonePaymentContract) CreateMilestonePayment(ctx contractapi.TransactionContextInterface, orderID string, milestoneID string, amount float64) error {
	orderAsBytes, err := ctx.GetStub().GetState(orderID)
	var order Order
	if err != nil || orderAsBytes == nil {
		order = Order{
			OrderID:    orderID,
			Milestones: make(map[string]Milestone),
		}
	} else {
		err = json.Unmarshal(orderAsBytes, &order)
		if err != nil {
			return fmt.Errorf("Failed to unmarshal order: %s", err.Error())
		}
	}

	order.Milestones[milestoneID] = Milestone{
		MilestoneID: milestoneID,
		Amount:      amount,
		Status:      "pending",
	}

	orderAsBytes, err = json.Marshal(order)
	if err != nil {
		return fmt.Errorf("Failed to marshal order: %s", err.Error())
	}

	return ctx.GetStub().PutState(orderID, orderAsBytes)
}

func (s *MilestonePaymentContract) UpdateMilestoneStatus(ctx contractapi.TransactionContextInterface, orderID string, milestoneID string, status string) error {
	orderAsBytes, err := ctx.GetStub().GetState(orderID)
	if err != nil || orderAsBytes == nil {
		return fmt.Errorf("Order %s does not exist", orderID)
	}

	var order Order
	err = json.Unmarshal(orderAsBytes, &order)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal order: %s", err.Error())
	}

	milestone, exists := order.Milestones[milestoneID]
	if !exists {
		return fmt.Errorf("Milestone %s does not exist", milestoneID)
	}

	milestone.Status = status
	order.Milestones[milestoneID] = milestone

	orderAsBytes, err = json.Marshal(order)
	if err != nil {
		return fmt.Errorf("Failed to marshal order: %s", err.Error())
	}

	return ctx.GetStub().PutState(orderID, orderAsBytes)
}

func (s *MilestonePaymentContract) ReleaseMilestonePayment(ctx contractapi.TransactionContextInterface, orderID string, milestoneID string) error {
	orderAsBytes, err := ctx.GetStub().GetState(orderID)
	if err != nil || orderAsBytes == nil {
		return fmt.Errorf("Order %s does not exist", orderID)
	}

	var order Order
	err = json.Unmarshal(orderAsBytes, &order)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal order: %s", err.Error())
	}

	milestone, exists := order.Milestones[milestoneID]
	if !exists {
		return fmt.Errorf("Milestone %s does not exist", milestoneID)
	}

	if milestone.Status != "completed" {
		return fmt.Errorf("Milestone %s is not completed", milestoneID)
	}

	fmt.Printf("Releasing payment of %.2f for milestone %s of order %s\n", milestone.Amount, milestoneID, orderID)

	milestone.Status = "paid"
	order.Milestones[milestoneID] = milestone

	orderAsBytes, err = json.Marshal(order)
	if err != nil {
		return fmt.Errorf("Failed to marshal order: %s", err.Error())
	}

	return ctx.GetStub().PutState(orderID, orderAsBytes)
}

func (s *MilestonePaymentContract) QueryMilestonePayments(ctx contractapi.TransactionContextInterface, orderID string) (map[string]Milestone, error) {
	orderAsBytes, err := ctx.GetStub().GetState(orderID)
	if err != nil || orderAsBytes == nil {
		return nil, fmt.Errorf("Order %s does not exist", orderID)
	}

	var order Order
	err = json.Unmarshal(orderAsBytes, &order)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal order: %s", err.Error())
	}

	return order.Milestones, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(MilestonePaymentContract))
	if err != nil {
		fmt.Printf("Error creating milestone payment chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting milestone payment chaincode: %s", err.Error())
	}
}
