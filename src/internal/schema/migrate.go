package schema

import (
	"github.com/GuiaBolso/darwin"
	"github.com/jmoiron/sqlx"
)

var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add products",
		Script: `CREATE TABLE products (
    product_id   UUID,
    name TEXT,
    cost INT,
    quantity INT,
    date_created TIMESTAMP,
    date_updated TIMESTAMP,
    
    primary key (product_id));`,
	},
	{
		Version:     2,
		Description: "Add sales",
		Script: `CREATE TABLE sales (
			sale_id   UUID,
			product_id UUID,
			quantity UUID,
			paid INT,
			date_created TIMESTAMP,
			
			PRIMARY KEY (sale_id),
			FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE);`,
	},
	{
		Version:     2.1,
		Description: "Update sales",
		Script: `alter table sales
					drop column quantity;
			
				alter table sales
					add quantity INT;`,
	},
}

func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}
