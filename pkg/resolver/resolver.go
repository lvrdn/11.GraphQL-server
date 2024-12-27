package resolver

import (
	"context"
	"shop/pkg/generated"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	StorageData   Storage
	UserIdFromCtx func(cxt context.Context) (int, error)
}

type Storage interface {
	//*Catalog must have *Catalog.ID value obtained on a previous step
	FillCatalogFieldsWithData(*generated.Catalog) error

	//*Item must have *Item.ID value obtained on a previous step
	FillItemFieldsWithData(item *generated.Item) error

	//*Seller must have *Seller.ID value obtained on a previous step
	FillSellerFieldsWithData(*generated.Seller) error

	//get user cart with user id
	GetUserCart(userID int) ([]*generated.CartItem, error)

	//action must be `add`(means add to cart) or `remove`(means remove from cart)
	EditUserCart(userID, itemID, quantity int, action string) ([]*generated.CartItem, error)

	QuantityInCart(userID, itemID int) (int, error)
}
