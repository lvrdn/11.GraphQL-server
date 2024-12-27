package storage

import (
	"encoding/json"
	"os"
)

func AddInitialData(filename string, st *Storage) error {

	dataFromFile, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})

	err = json.Unmarshal(dataFromFile, &data)
	if err != nil {
		return err
	}

	fillCatsAndItems(data["catalog"], 0, st.Catalog, st.Item)
	fillSellers(data["sellers"], st.Seller, st.Item)

	return nil
}

func fillCatsAndItems(x interface{}, parentID int, catalogSt *CatalogStorage, itemSt *ItemStorage) interface{} {

	slice, ok := x.([]interface{})
	if ok {
		ids := make([]int, 0)
		for _, elem := range slice {
			id := fillCatsAndItems(elem, parentID, catalogSt, itemSt).(int)
			ids = append(ids, id)
		}
		return ids
	}

	data, ok := x.(map[string]interface{})
	if ok {
		if _, ok := data["seller_id"]; ok {

			id := int(data["id"].(float64))
			inStock := int(data["in_stock"].(float64))
			sellerID := int(data["seller_id"].(float64))
			name := data["name"].(string)

			itemSt.Items = append(itemSt.Items, &Item{
				ID:       id,
				Name:     name,
				InStock:  inStock,
				SellerID: sellerID,
				ParentID: parentID,
			})

			return id
		}

		id := int(data["id"].(float64))
		name := data["name"].(string)

		if childs, ok := data["childs"]; ok {

			ids := fillCatsAndItems(childs, id, catalogSt, itemSt).([]int)
			catalogSt.Catalogs = append(catalogSt.Catalogs, &Catalog{
				ID:       id,
				Name:     name,
				ChildsID: ids,
			})
		}

		if items, ok := data["items"]; ok {
			ids := fillCatsAndItems(items, id, catalogSt, itemSt).([]int)
			catalogSt.Catalogs = append(catalogSt.Catalogs, &Catalog{
				ID:      id,
				Name:    name,
				ItemsID: ids,
			})
		}

		return id

	}
	return nil
}

func fillSellers(x interface{}, sellerSt *SellerStorage, itemSt *ItemStorage) {

	for _, elem := range x.([]interface{}) {
		data := elem.(map[string]interface{})

		id := int(data["id"].(float64))
		name := data["name"].(string)
		deals := int(data["deals"].(float64))

		ids := make([]int, 0)
		for _, item := range itemSt.Items {
			if item.SellerID == id {
				ids = append(ids, item.ID)
			}
		}

		sellerSt.Sellers = append(sellerSt.Sellers, &Seller{
			ID:      id,
			Name:    name,
			Deals:   deals,
			ItemsID: ids,
		})
	}
}
