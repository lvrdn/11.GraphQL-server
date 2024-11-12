package main

import "context"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Storage interface {
	//*Catalog must have *Catalog.ID value obtained on a previous step
	FillCatalogFieldsWithData(*Catalog) error

	//*Item must have *Item.ID value obtained on a previous step
	FillItemFieldsWithData(item *Item) error

	//*Seller must have *Seller.ID value obtained on a previous step
	FillSellerFieldsWithData(*Seller) error

	//get user cart with user id
	GetUserCart(userID int) ([]*CartItem, error)

	//action must be `add`(means add to cart) or `remove`(means remove from cart)
	EditUserCart(userID, itemID, quantity int, action string) ([]*CartItem, error)

	QuantityInCart(userID, itemID int) (int, error)
}

type Resolver struct {
	StorageData   Storage
	UserIdFromCtx func(cxt context.Context) (int, error)
}
