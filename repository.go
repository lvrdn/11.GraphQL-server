package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ItemSt struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	InStock  int    `json:"in_stock"`
	SellerID int    `json:"seller_id"`
	ParentID int
}

type SellerSt struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Deals   int    `json:"deals"`
	ItemsID []int
}

type CatalogSt struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID int
	ChildsID []int
	ItemsID  []int
	Childs   []CatalogSt `json:"childs"`
	Items    []ItemSt    `json:"items"`
}

type StorageData struct {
	Items     []ItemSt
	Sellers   []SellerSt
	Catalogs  []CatalogSt
	UserCarts []Cart
}

func NewStorageInit() *StorageData {
	return &StorageData{
		Items:     make([]ItemSt, 0, 15),
		Sellers:   make([]SellerSt, 0, 15),
		Catalogs:  make([]CatalogSt, 0, 15),
		UserCarts: make([]Cart, 0, 15),
	}

}

func (st *StorageData) AddInitialData(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	dataFromFile, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	file.Close()

	data := &struct {
		Catalog CatalogSt  `json:"catalog"`
		Sellers []SellerSt `json:"sellers"`
	}{}

	err = json.Unmarshal(dataFromFile, data)
	if err != nil {
		return err
	}

	st.AddCatalogsAndItems(data.Catalog, 0)
	st.AddSellers(data.Sellers)
	return nil
}

func (st *StorageData) AddSellers(dataSellers []SellerSt) {
	for _, seller := range dataSellers {
		for _, item := range st.Items {
			if seller.ID == item.SellerID {
				seller.ItemsID = append(seller.ItemsID, item.ID)
			}
		}
		st.Sellers = append(st.Sellers, seller)
	}
}

func (st *StorageData) AddCatalogsAndItems(catalog CatalogSt, parentID int) {
	catalog.ParentID = parentID
	if catalog.Childs != nil {
		for i := 0; i < len(catalog.Childs); i++ {
			catalog.ChildsID = append(catalog.ChildsID, catalog.Childs[i].ID)
			st.AddCatalogsAndItems(catalog.Childs[i], catalog.ID)
		}
		catalog.Childs = nil

	} else if catalog.Items != nil {

		for i := 0; i < len(catalog.Items); i++ {
			catalog.ItemsID = append(catalog.ItemsID, catalog.Items[i].ID)
			catalog.Items[i].ParentID = catalog.ID

		}
		st.Items = append(st.Items, catalog.Items...)
		catalog.Items = nil
	}

	st.Catalogs = append(st.Catalogs, catalog)

}

func (st *StorageData) FillCatalogFieldsWithData(catalog *Catalog) error {
	if catalog.ID == nil {
		return fmt.Errorf("catalog must have id to fill other fields")
	}

	var name string
	var childsIDs []int
	var itemsIDs []int
	var parentID int
	founded := false

	for _, catalogFromStorage := range st.Catalogs {
		if *catalog.ID == catalogFromStorage.ID {
			name = catalogFromStorage.Name
			childsIDs = catalogFromStorage.ChildsID
			itemsIDs = catalogFromStorage.ItemsID
			parentID = catalogFromStorage.ParentID
			founded = true
			break
		}
	}
	if !founded {
		return fmt.Errorf("no catalog with this id")
	}

	catalog.Name = &name

	if parentID != 0 {
		catalog.Parent = &Catalog{ID: &parentID}
	}

	if len(childsIDs) != 0 {
		catalog.Childs = make([]*Catalog, 0)
		for i := 0; i < len(childsIDs); i++ {
			catalog.Childs = append(catalog.Childs, &Catalog{ID: &childsIDs[i]})
		}
	}

	if len(itemsIDs) != 0 {
		catalog.Items = make([]*Item, 0)
		for i := 0; i < len(itemsIDs); i++ {
			catalog.Items = append(catalog.Items, &Item{ID: &itemsIDs[i]})
		}
	}

	return nil
}

func (st *StorageData) FillItemFieldsWithData(item *Item) error {
	if item.ID == nil {
		return fmt.Errorf("item must have id to fill other fields")
	}

	var name string
	var parentID, sellerID int
	var inStockText string
	founded := false

	for _, itemFromStorage := range st.Items {
		if *item.ID == itemFromStorage.ID {
			name = itemFromStorage.Name
			sellerID = itemFromStorage.SellerID
			parentID = itemFromStorage.ParentID
			switch {
			case itemFromStorage.InStock <= 1:
				inStockText = "мало"
			case itemFromStorage.InStock >= 2 && itemFromStorage.InStock <= 3:
				inStockText = "хватает"
			case itemFromStorage.InStock > 3:
				inStockText = "много"
			}

			founded = true
			break
		}
	}
	if !founded {
		return fmt.Errorf("no item with this id")
	}

	item.Name = &name
	item.Parent = &Catalog{ID: &parentID}
	item.Seller = &Seller{ID: &sellerID}
	item.InStockText = inStockText

	return nil
}

func (st *StorageData) FillSellerFieldsWithData(seller *Seller) error {
	if seller.ID == nil {
		return fmt.Errorf("seller must have id to fill other fields")
	}

	var name string
	var deals int
	var itemsIDs []int
	founded := false

	for _, sellerFromStorage := range st.Sellers {
		if *seller.ID == sellerFromStorage.ID {
			name = sellerFromStorage.Name
			deals = sellerFromStorage.Deals
			itemsIDs = sellerFromStorage.ItemsID
			founded = true
			break
		}
	}

	if !founded {
		return fmt.Errorf("no seller with this id")
	}

	seller.Name = &name
	seller.Deals = deals

	if len(itemsIDs) != 0 {
		seller.Items = make([]*Item, 0)
		for i := 0; i < len(itemsIDs); i++ {
			seller.Items = append(seller.Items, &Item{ID: &itemsIDs[i]})
		}
	}

	return nil
}

func (st *StorageData) EditItemStock(itemID, quantity int, actionFromCart string) error {
	switch actionFromCart {
	case "add":
		for i, item := range st.Items {
			if item.ID == itemID {
				if st.Items[i].InStock < quantity {
					return fmt.Errorf("not enough quantity")
				}
				st.Items[i].InStock -= quantity //если добавляется в корзину, то из stock надо вычесть
				break
			}
		}
		return nil
	case "remove":
		for i, item := range st.Items {
			if item.ID == itemID {
				st.Items[i].InStock += quantity //если удаляется из корзины, то в stock надо добавить
				break
			}
		}
		return nil
	default:
		return fmt.Errorf("wrong action value, must be add or remove")
	}

}

func (st *StorageData) EditUserCart(userID, itemID, quantity int, action string) ([]*CartItem, error) {

	err := st.EditItemStock(itemID, quantity, action)
	if err != nil {
		return nil, err
	}

	for i, cart := range st.UserCarts {
		if *cart.UserID != userID {
			continue
		}
		for j, cartItem := range st.UserCarts[i].CartItems {
			if *cartItem.Item.ID != itemID {
				continue
			}
			switch action {
			case "add":
				st.UserCarts[i].CartItems[j].Quantity += quantity
				err := st.FillItemFieldsWithData(st.UserCarts[i].CartItems[j].Item)
				if err != nil {
					return nil, err
				}
				return st.UserCarts[i].CartItems, nil
			case "remove":
				if st.UserCarts[i].CartItems[j].Quantity <= quantity {
					st.UserCarts[i].CartItems[j] = st.UserCarts[i].CartItems[len(st.UserCarts[i].CartItems)-1]
					st.UserCarts[i].CartItems = st.UserCarts[i].CartItems[:len(st.UserCarts[i].CartItems)-1]
					return st.UserCarts[i].CartItems, nil
				}
				st.UserCarts[i].CartItems[j].Quantity -= quantity
				err := st.FillItemFieldsWithData(st.UserCarts[i].CartItems[j].Item)
				if err != nil {
					return nil, err
				}
				return st.UserCarts[i].CartItems, nil
			default:
				return nil, fmt.Errorf("wrong action value, must be add or remove")
			}

		}

		newItemInCart := &Item{ID: &itemID}
		err := st.FillItemFieldsWithData(newItemInCart)
		if err != nil {
			return nil, err
		}

		st.UserCarts[i].CartItems = append(st.UserCarts[i].CartItems,
			&CartItem{
				Quantity: quantity,
				Item:     newItemInCart,
			},
		)

		return st.UserCarts[i].CartItems, nil
	}

	newItemInCart := &Item{ID: &itemID}
	err = st.FillItemFieldsWithData(newItemInCart)
	if err != nil {
		return nil, err
	}
	st.UserCarts = append(st.UserCarts,
		Cart{
			UserID: &userID,
			CartItems: []*CartItem{
				{
					Quantity: quantity,
					Item:     newItemInCart,
				},
			},
		},
	)

	return st.UserCarts[0].CartItems, nil
}

func (st *StorageData) GetUserCart(userID int) ([]*CartItem, error) {
	for i := range st.UserCarts {
		if *st.UserCarts[i].UserID != userID {
			continue
		}

		for j := range st.UserCarts[i].CartItems {
			err := st.FillItemFieldsWithData(st.UserCarts[i].CartItems[j].Item)
			if err != nil {
				return nil, err
			}
		}
		return st.UserCarts[i].CartItems, nil
	}
	return nil, fmt.Errorf("no user with this user id")
}

func (st *StorageData) QuantityInCart(userID, itemID int) (int, error) {
	for i, cart := range st.UserCarts {
		if *cart.UserID != userID {
			continue
		}
		for _, cartItem := range st.UserCarts[i].CartItems {
			if *cartItem.Item.ID == itemID {
				return cartItem.Quantity, nil
			}
		}
	}
	return 0, nil
}
