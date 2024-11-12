package main

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"
	"strconv"
)

// Parent is the resolver for the parent field.
func (r *catalogResolver) Parent(ctx context.Context, obj *Catalog) (*Catalog, error) {
	if obj.Parent == nil {
		return nil, fmt.Errorf("this catalog dont have parent catalog")
	}
	err := r.StorageData.FillCatalogFieldsWithData(obj.Parent)
	if err != nil {
		return nil, err
	}

	return obj.Parent, nil
}

// Childs is the resolver for the childs field.
func (r *catalogResolver) Childs(ctx context.Context, obj *Catalog) ([]*Catalog, error) {
	if obj.Childs == nil {
		return nil, fmt.Errorf("this catalog dont have childs catalog")
	}
	for i := 0; i < len(obj.Childs); i++ {
		err := r.StorageData.FillCatalogFieldsWithData(obj.Childs[i])
		if err != nil {
			return nil, err
		}
	}
	return obj.Childs, nil
}

// Items is the resolver for the items field.
func (r *catalogResolver) Items(ctx context.Context, obj *Catalog, limit *int, offset *int) ([]*Item, error) {

	if obj.Items == nil {
		return nil, fmt.Errorf("this catalog dont have items")
	}

	if *offset+*limit < len(obj.Items) {
		obj.Items = obj.Items[*offset : *offset+*limit]
	} else if *offset < len(obj.Items) {
		obj.Items = obj.Items[*offset:]
	} else {
		obj.Items = obj.Items[:0]
	}

	for i := 0; i < len(obj.Items); i++ {

		err := r.StorageData.FillItemFieldsWithData(obj.Items[i])
		if err != nil {
			return nil, err
		}
	}
	return obj.Items, nil
}

// Parent is the resolver for the parent field.
func (r *itemResolver) Parent(ctx context.Context, obj *Item) (*Catalog, error) {
	if obj.Parent == nil {
		return nil, fmt.Errorf("this item dont have parent catalog")
	}
	err := r.StorageData.FillCatalogFieldsWithData(obj.Parent)
	if err != nil {
		return nil, err
	}

	return obj.Parent, nil
}

// Seller is the resolver for the seller field.
func (r *itemResolver) Seller(ctx context.Context, obj *Item) (*Seller, error) {
	err := r.StorageData.FillSellerFieldsWithData(obj.Seller)
	if err != nil {
		return nil, err
	}

	return obj.Seller, nil
}

// InCart is the resolver for the inCart field.
func (r *itemResolver) InCart(ctx context.Context, obj *Item) (int, error) {
	userID, err := r.UserIdFromCtx(ctx)
	if err != nil {
		return 0, err
	}

	quantity, err := r.StorageData.QuantityInCart(userID, *obj.ID)
	if err != nil {
		return 0, err
	}

	return quantity, nil
}

// AddToCart is the resolver for the AddToCart field.
func (r *mutationResolver) AddToCart(ctx context.Context, in *CartInput) ([]*CartItem, error) {
	userID, err := r.UserIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	itemID := in.ItemID
	quantity := in.Quantity

	cartItems, err := r.StorageData.EditUserCart(userID, itemID, quantity, "add")
	if err != nil {
		return nil, err
	}

	return cartItems, nil
}

// RemoveFromCart is the resolver for the RemoveFromCart field.
func (r *mutationResolver) RemoveFromCart(ctx context.Context, in CartInput) ([]*CartItem, error) {
	userID, err := r.UserIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	itemID := in.ItemID
	quantity := in.Quantity

	cartItems, err := r.StorageData.EditUserCart(userID, itemID, quantity, "remove")
	if err != nil {
		return nil, err
	}

	return cartItems, nil
}

// Catalog is the resolver for the Catalog field.
func (r *queryResolver) Catalog(ctx context.Context, id *string) (*Catalog, error) {
	idINT, err := strconv.Atoi(*id)
	if err != nil {
		return nil, fmt.Errorf("id must be number")
	}
	catalog := &Catalog{ID: &idINT}
	err = r.StorageData.FillCatalogFieldsWithData(catalog)
	if err != nil {
		return nil, err
	}
	return catalog, nil
}

// Seller is the resolver for the Seller field.
func (r *queryResolver) Seller(ctx context.Context, id *string) (*Seller, error) {
	idINT, err := strconv.Atoi(*id)
	if err != nil {
		return nil, fmt.Errorf("id must be number")
	}
	seller := &Seller{ID: &idINT}
	err = r.StorageData.FillSellerFieldsWithData(seller)
	if err != nil {
		return nil, err
	}
	return seller, nil
}

// MyCart is the resolver for the MyCart field.
func (r *queryResolver) MyCart(ctx context.Context) ([]*CartItem, error) {
	userID, err := r.UserIdFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	cartItems, err := r.StorageData.GetUserCart(userID)
	if err != nil {
		return nil, err
	}

	return cartItems, nil
}

// Items is the resolver for the items field.
func (r *sellerResolver) Items(ctx context.Context, obj *Seller, limit *int, offset *int) ([]*Item, error) {

	if obj.Items == nil {
		return nil, fmt.Errorf("this seller dont have items")
	}

	if *offset+*limit < len(obj.Items) {
		obj.Items = obj.Items[*offset : *offset+*limit]
	} else if *offset < len(obj.Items) {
		obj.Items = obj.Items[*offset:]
	} else {
		obj.Items = obj.Items[:0]
	}

	for i := 0; i < len(obj.Items); i++ {

		err := r.StorageData.FillItemFieldsWithData(obj.Items[i])
		if err != nil {
			return nil, err
		}
	}
	return obj.Items, nil
}

// Catalog returns CatalogResolver implementation.
func (r *Resolver) Catalog() CatalogResolver { return &catalogResolver{r} }

// Item returns ItemResolver implementation.
func (r *Resolver) Item() ItemResolver { return &itemResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Seller returns SellerResolver implementation.
func (r *Resolver) Seller() SellerResolver { return &sellerResolver{r} }

type catalogResolver struct{ *Resolver }
type itemResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type sellerResolver struct{ *Resolver }