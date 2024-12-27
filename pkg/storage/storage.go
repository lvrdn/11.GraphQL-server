package storage

import (
	"fmt"
	"shop/pkg/generated"
	"sync"
)

type ItemStorage struct {
	Items []*Item
	mu    *sync.RWMutex
}

type Item struct {
	ID       int
	Name     string
	InStock  int
	SellerID int
	ParentID int
}

type SellerStorage struct {
	Sellers []*Seller
	mu      *sync.RWMutex
}

type Seller struct {
	ID      int
	Name    string
	Deals   int
	ItemsID []int
}

type CatalogStorage struct {
	Catalogs []*Catalog
	mu       *sync.RWMutex
}

type Catalog struct {
	ID       int
	Name     string
	ParentID int
	ChildsID []int
	ItemsID  []int
}

type CartStorage struct {
	Carts []*generated.Cart
	mu    *sync.RWMutex
}

type Storage struct {
	Item     *ItemStorage
	Seller   *SellerStorage
	Catalog  *CatalogStorage
	UserCart *CartStorage
}

func NewStorageInit() *Storage {
	return &Storage{
		Item: &ItemStorage{
			Items: make([]*Item, 0),
			mu:    &sync.RWMutex{},
		},
		Seller: &SellerStorage{
			Sellers: make([]*Seller, 0),
			mu:      &sync.RWMutex{},
		},
		Catalog: &CatalogStorage{
			Catalogs: make([]*Catalog, 0),
			mu:       &sync.RWMutex{},
		},
		UserCart: &CartStorage{
			Carts: make([]*generated.Cart, 0),
			mu:    &sync.RWMutex{},
		},
	}

}

func (st *Storage) FillCatalogFieldsWithData(catalog *generated.Catalog) error {
	if catalog.ID == nil {
		return fmt.Errorf("catalog must have id to fill other fields")
	}

	var name string
	var childsIDs []int
	var itemsIDs []int
	var parentID int
	founded := false

	st.Catalog.mu.RLock()
	for _, catalogFromStorage := range st.Catalog.Catalogs {
		if *catalog.ID == catalogFromStorage.ID {
			name = catalogFromStorage.Name
			childsIDs = catalogFromStorage.ChildsID
			itemsIDs = catalogFromStorage.ItemsID
			parentID = catalogFromStorage.ParentID
			founded = true
			break
		}
	}
	st.Catalog.mu.RUnlock()

	if !founded {
		return fmt.Errorf("no catalog with this id")
	}

	catalog.Name = &name

	if parentID != 0 {
		catalog.Parent = &generated.Catalog{ID: &parentID}
	}

	if len(childsIDs) != 0 {
		catalog.Childs = make([]*generated.Catalog, 0)
		for i := 0; i < len(childsIDs); i++ {
			catalog.Childs = append(catalog.Childs, &generated.Catalog{ID: &childsIDs[i]})
		}
	}

	if len(itemsIDs) != 0 {
		catalog.Items = make([]*generated.Item, 0)
		for i := 0; i < len(itemsIDs); i++ {
			catalog.Items = append(catalog.Items, &generated.Item{ID: &itemsIDs[i]})
		}
	}

	return nil
}

func (st *Storage) FillItemFieldsWithData(item *generated.Item) error {
	if item.ID == nil {
		return fmt.Errorf("item must have id to fill other fields")
	}

	var name string
	var parentID, sellerID int
	var inStockText string
	founded := false

	st.Item.mu.RLock()
	for _, itemFromStorage := range st.Item.Items {
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
	st.Item.mu.RUnlock()

	if !founded {
		return fmt.Errorf("no item with this id")
	}

	item.Name = &name
	item.Parent = &generated.Catalog{ID: &parentID}
	item.Seller = &generated.Seller{ID: &sellerID}
	item.InStockText = inStockText

	return nil
}

func (st *Storage) FillSellerFieldsWithData(seller *generated.Seller) error {
	if seller.ID == nil {
		return fmt.Errorf("seller must have id to fill other fields")
	}

	var name string
	var deals int
	var itemsIDs []int
	founded := false

	st.Seller.mu.RLock()
	for _, sellerFromStorage := range st.Seller.Sellers {
		if *seller.ID == sellerFromStorage.ID {
			name = sellerFromStorage.Name
			deals = sellerFromStorage.Deals
			itemsIDs = sellerFromStorage.ItemsID
			founded = true
			break
		}
	}
	st.Seller.mu.RUnlock()

	if !founded {
		return fmt.Errorf("no seller with this id")
	}

	seller.Name = &name
	seller.Deals = deals

	if len(itemsIDs) != 0 {
		seller.Items = make([]*generated.Item, 0)
		for i := 0; i < len(itemsIDs); i++ {
			seller.Items = append(seller.Items, &generated.Item{ID: &itemsIDs[i]})
		}
	}

	return nil
}

func (st *Storage) EditItemStock(itemID, quantity int, actionFromCart string) error {
	st.Item.mu.Lock()
	defer st.Item.mu.Unlock()

	switch actionFromCart {
	case "add":
		for i, item := range st.Item.Items {
			if item.ID == itemID {
				if st.Item.Items[i].InStock < quantity {
					return fmt.Errorf("not enough quantity")
				}
				st.Item.Items[i].InStock -= quantity //если добавляется в корзину, то из stock надо вычесть
				break
			}
		}
		return nil
	case "remove":
		for i, item := range st.Item.Items {
			if item.ID == itemID {
				st.Item.Items[i].InStock += quantity //если удаляется из корзины, то в stock надо добавить
				break
			}
		}
		return nil
	default:
		return fmt.Errorf("wrong action value, must be add or remove")
	}

}

func (st *Storage) EditUserCart(userID, itemID, quantity int, action string) ([]*generated.CartItem, error) {

	err := st.EditItemStock(itemID, quantity, action)
	if err != nil {
		return nil, err
	}

	st.UserCart.mu.Lock()
	defer st.UserCart.mu.Unlock()

	for i, cart := range st.UserCart.Carts {
		if *cart.UserID != userID {
			continue
		}
		for j, cartItem := range st.UserCart.Carts[i].CartItems {
			if *cartItem.Item.ID != itemID {
				continue
			}
			switch action {
			case "add":
				st.UserCart.Carts[i].CartItems[j].Quantity += quantity
				err := st.FillItemFieldsWithData(st.UserCart.Carts[i].CartItems[j].Item)
				if err != nil {
					return nil, err
				}
				return st.UserCart.Carts[i].CartItems, nil
			case "remove":
				if st.UserCart.Carts[i].CartItems[j].Quantity <= quantity {
					st.UserCart.Carts[i].CartItems[j] = st.UserCart.Carts[i].CartItems[len(st.UserCart.Carts[i].CartItems)-1]
					st.UserCart.Carts[i].CartItems = st.UserCart.Carts[i].CartItems[:len(st.UserCart.Carts[i].CartItems)-1]
					return st.UserCart.Carts[i].CartItems, nil
				}
				st.UserCart.Carts[i].CartItems[j].Quantity -= quantity
				err := st.FillItemFieldsWithData(st.UserCart.Carts[i].CartItems[j].Item)
				if err != nil {
					return nil, err
				}
				return st.UserCart.Carts[i].CartItems, nil
			default:
				return nil, fmt.Errorf("wrong action value, must be add or remove")
			}

		}

		newItemInCart := &generated.Item{ID: &itemID}
		err := st.FillItemFieldsWithData(newItemInCart)
		if err != nil {
			return nil, err
		}

		st.UserCart.Carts[i].CartItems = append(st.UserCart.Carts[i].CartItems,
			&generated.CartItem{
				Quantity: quantity,
				Item:     newItemInCart,
			},
		)

		return st.UserCart.Carts[i].CartItems, nil
	}

	newItemInCart := &generated.Item{ID: &itemID}
	err = st.FillItemFieldsWithData(newItemInCart)
	if err != nil {
		return nil, err
	}
	st.UserCart.Carts = append(st.UserCart.Carts,
		&generated.Cart{
			UserID: &userID,
			CartItems: []*generated.CartItem{
				{
					Quantity: quantity,
					Item:     newItemInCart,
				},
			},
		},
	)

	return st.UserCart.Carts[0].CartItems, nil
}

func (st *Storage) GetUserCart(userID int) ([]*generated.CartItem, error) {

	st.UserCart.mu.RLock()
	defer st.UserCart.mu.RUnlock()

	for i := range st.UserCart.Carts {
		if *st.UserCart.Carts[i].UserID != userID {
			continue
		}

		for j := range st.UserCart.Carts[i].CartItems {
			err := st.FillItemFieldsWithData(st.UserCart.Carts[i].CartItems[j].Item)
			if err != nil {
				return nil, err
			}
		}
		return st.UserCart.Carts[i].CartItems, nil
	}
	return nil, fmt.Errorf("no user with this user id")
}

func (st *Storage) QuantityInCart(userID, itemID int) (int, error) {

	st.UserCart.mu.RLock()
	defer st.UserCart.mu.RUnlock()

	for i, cart := range st.UserCart.Carts {
		if *cart.UserID != userID {
			continue
		}
		for _, cartItem := range st.UserCart.Carts[i].CartItems {
			if *cartItem.Item.ID == itemID {
				return cartItem.Quantity, nil
			}
		}
	}
	return 0, nil
}
