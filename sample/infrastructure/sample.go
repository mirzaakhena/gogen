package infrastructure

import (
	"database/sql"
	"fmt"
	"log"
)

type Inventory struct {
	Quantity int
	Price    float64
}

type LastInventoryProvider interface {
	GetLastInventory(code string) (*Inventory, error)
}

func ReduceStock(code string, qtyRequested int, p LastInventoryProvider) (*Inventory, error) {

	if p == nil {
		return nil, fmt.Errorf("LastInventoryProvider must not nil")
	}

	if qtyRequested <= 0 {
		return nil, fmt.Errorf("quantity must greater than 0")
	}

	lastInventory, err := p.GetLastInventory(code)
	if err != nil {
		return nil, err
	}

	qtyReminder := lastInventory.Quantity - qtyRequested
	if qtyReminder < 0 {
		return nil, fmt.Errorf("quantity is not enough")
	}

	return &Inventory{
		Quantity: qtyReminder,
		Price:    lastInventory.Price,
	}, nil

}

func AddNewStock(code string, qty int, price float64, p LastInventoryProvider) (*Inventory, error) {

	if p == nil {
		return nil, fmt.Errorf("LastInventoryProvider must not nil")
	}

	if qty <= 0 {
		return nil, fmt.Errorf("quantity must greater than 0")
	}

	if price <= 0 {
		return nil, fmt.Errorf("price must greater than 0")
	}

	lastInventory, err := p.GetLastInventory(code)
	if err != nil {
		return nil, err
	}

	newQty := lastInventory.Quantity + qty
	lastTotalPrice := lastInventory.Price * float64(lastInventory.Quantity)
	newTotalPrice := price * float64(qty)

	pricePerQty := (lastTotalPrice + newTotalPrice) / float64(newQty)

	return &Inventory{
		Quantity: newQty,
		Price:    pricePerQty,
	}, nil

}

type lastInventoryProviderFromDB struct{}

func (lastInventoryProviderFromDB) GetLastInventory(code string) (*Inventory, error) {

	dataSourceName := "root:pass1@tcp(127.0.0.1:3306)/inventory_db"
	db, err := sql.Open("mysql", dataSourceName)

	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT quantity, price FROM inventories WHERE code = ? ORDER BY date DESC LIMIT 1"
	dbPrepare, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	resultStmt, err := dbPrepare.Query(code)
	if err != nil {
		return nil, err
	}
	defer resultStmt.Close()

	var inv Inventory

	if resultStmt.Next() {
		err = resultStmt.Scan(&inv.Quantity, &inv.Price)
		if err != nil {
			return nil, err
		}
	}

	return &inv, nil

}

func main() {

	dataProvider := lastInventoryProviderFromDB{}

	// try ReduceStock
	{
		newInv, err := ReduceStock("A44", 2, &dataProvider)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("%v\n", newInv)
	}

	// try AddNewStock
	{
		newInv, err := AddNewStock("A44", 2, 2000, &dataProvider)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("%v\n", newInv)
	}

}
