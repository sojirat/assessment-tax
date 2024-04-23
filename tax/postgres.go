package tax

// func (p *Postgres) ProductById(id string) (product.DB, error) {
// 	var productDb product.DB
// 	rows, err := p.Db.Query("SELECT product_id,name,category FROM product WHERE product_id = $1", id)
// 	if err != nil {
// 		return product.DB{}, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var p product.DB
// 		err = rows.Scan(&p.ProductID, &p.Name, &p.Category)
// 		if err != nil {
// 			log.Fatal(err)
// 			return productDb, err
// 		}
// 		productDb = p
// 	}

// 	return productDb, nil
// }
