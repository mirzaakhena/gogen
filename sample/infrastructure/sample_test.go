package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type lastInventoryProviderFromHardcode struct{}

func (lastInventoryProviderFromHardcode) GetLastInventory(code string) (*Inventory, error) {
	return &Inventory{
		Quantity: 4,
		Price:    2000,
	}, nil
}

func TestAddNewStock_SamePrice(t *testing.T) {

	dataProvider := lastInventoryProviderFromHardcode{}
	newInv, _ := AddNewStock("", 2, 2000, &dataProvider)

	assert.Equal(t, newInv.Price, float64(2000))
	assert.Equal(t, newInv.Quantity, 6)

}

func TestAddNewStock_DifferentPrice(t *testing.T) {

	dataProvider := lastInventoryProviderFromHardcode{}
	newInv, _ := AddNewStock("", 2, 2600, &dataProvider)

	assert.Equal(t, newInv.Price, float64(2200))
	assert.Equal(t, newInv.Quantity, 6)

}

func TestReduceStock_Available(t *testing.T) {

	dataProvider := lastInventoryProviderFromHardcode{}
	newInv, _ := ReduceStock("", 2, &dataProvider)

	assert.Equal(t, newInv.Price, float64(2000))
	assert.Equal(t, newInv.Quantity, 2)

}

func TestReduceStock_NotEnough(t *testing.T) {

	dataProvider := lastInventoryProviderFromHardcode{}
	_, err := ReduceStock("", 8, &dataProvider)

	assert.EqualError(t, err, "quantity is not enough")

}
