package entity

type Inventory struct {
	ID                string `` //
	CalculationMethod string `` //
}

type InventoryRequest struct {
}

func NewInventory(req InventoryRequest) (*Inventory, error) {

	var obj Inventory

	return &obj, nil
}
