package postgres

import (
	"github.com/go-jet/jet"
	"github.com/go-jet/jet/internal/testutils"
	. "github.com/go-jet/jet/postgres"
	"github.com/go-jet/jet/tests/.gentestdata/jetdb/dvds/enum"
	"github.com/go-jet/jet/tests/.gentestdata/jetdb/dvds/model"
	. "github.com/go-jet/jet/tests/.gentestdata/jetdb/dvds/table"
	"gotest.tools/assert"
	"testing"
)

func TestSelect_ScanToStruct(t *testing.T) {
	expectedSQL := `
SELECT DISTINCT actor.actor_id AS "actor.actor_id",
     actor.first_name AS "actor.first_name",
     actor.last_name AS "actor.last_name",
     actor.last_update AS "actor.last_update"
FROM dvds.actor
WHERE actor.actor_id = 1;
`

	query := Actor.
		SELECT(Actor.AllColumns).
		DISTINCT().
		WHERE(Actor.ActorID.EQ(Int(1)))

	testutils.AssertDebugStatementSql(t, query, expectedSQL, int64(1))

	actor := model.Actor{}
	err := query.Query(db, &actor)

	assert.NilError(t, err)

	expectedActor := model.Actor{
		ActorID:    1,
		FirstName:  "Penelope",
		LastName:   "Guiness",
		LastUpdate: *testutils.TimestampWithoutTimeZone("2013-05-26 14:47:57.62", 2),
	}

	assert.DeepEqual(t, actor, expectedActor)
}

func TestClassicSelect(t *testing.T) {
	expectedSQL := `
SELECT payment.payment_id AS "payment.payment_id",
     payment.customer_id AS "payment.customer_id",
     payment.staff_id AS "payment.staff_id",
     payment.rental_id AS "payment.rental_id",
     payment.amount AS "payment.amount",
     payment.payment_date AS "payment.payment_date",
     customer.customer_id AS "customer.customer_id",
     customer.store_id AS "customer.store_id",
     customer.first_name AS "customer.first_name",
     customer.last_name AS "customer.last_name",
     customer.email AS "customer.email",
     customer.address_id AS "customer.address_id",
     customer.activebool AS "customer.activebool",
     customer.create_date AS "customer.create_date",
     customer.last_update AS "customer.last_update",
     customer.active AS "customer.active"
FROM dvds.payment
     INNER JOIN dvds.customer ON (payment.customer_id = customer.customer_id)
ORDER BY payment.payment_id ASC
LIMIT 30;
`

	query := SELECT(
		Payment.AllColumns,
		Customer.AllColumns,
	).
		FROM(Payment.
			INNER_JOIN(Customer, Payment.CustomerID.EQ(Customer.CustomerID))).
		ORDER_BY(Payment.PaymentID.ASC()).
		LIMIT(30)

	testutils.AssertDebugStatementSql(t, query, expectedSQL, int64(30))

	dest := []model.Payment{}

	err := query.Query(db, &dest)

	assert.NilError(t, err)
	assert.Equal(t, len(dest), 30)
}

func TestSelect_ScanToSlice(t *testing.T) {
	expectedSQL := `
SELECT customer.customer_id AS "customer.customer_id",
     customer.store_id AS "customer.store_id",
     customer.first_name AS "customer.first_name",
     customer.last_name AS "customer.last_name",
     customer.email AS "customer.email",
     customer.address_id AS "customer.address_id",
     customer.activebool AS "customer.activebool",
     customer.create_date AS "customer.create_date",
     customer.last_update AS "customer.last_update",
     customer.active AS "customer.active"
FROM dvds.customer
ORDER BY customer.customer_id ASC;
`
	customers := []model.Customer{}

	query := Customer.SELECT(Customer.AllColumns).ORDER_BY(Customer.CustomerID.ASC())

	testutils.AssertDebugStatementSql(t, query, expectedSQL)

	err := query.Query(db, &customers)
	assert.NilError(t, err)

	assert.Equal(t, len(customers), 599)

	assert.DeepEqual(t, customer0, customers[0])
	assert.DeepEqual(t, customer1, customers[1])
	assert.DeepEqual(t, lastCustomer, customers[598])
}

func TestSelectAndUnionInProjection(t *testing.T) {
	expectedSQL := `
SELECT payment.payment_id AS "payment.payment_id",
     (
          SELECT customer.customer_id AS "customer.customer_id"
          FROM dvds.customer
          LIMIT 1
     ),
     (
          (
               (
                    SELECT payment.payment_id AS "payment.payment_id"
                    FROM dvds.payment
                    LIMIT 1
                    OFFSET 10
               )
               UNION
               (
                    SELECT payment.payment_id AS "payment.payment_id"
                    FROM dvds.payment
                    LIMIT 1
                    OFFSET 2
               )
          )
          LIMIT 1
     )
FROM dvds.payment
LIMIT 12;
`

	query := Payment.
		SELECT(
			Payment.PaymentID,
			Customer.SELECT(Customer.CustomerID).LIMIT(1),
			UNION(
				Payment.SELECT(Payment.PaymentID).LIMIT(1).OFFSET(10),
				Payment.SELECT(Payment.PaymentID).LIMIT(1).OFFSET(2),
			).LIMIT(1),
		).
		LIMIT(12)

	testutils.AssertDebugStatementSql(t, query, expectedSQL, int64(1), int64(1), int64(10), int64(1), int64(2), int64(1), int64(12))
}

func TestJoinQueryStruct(t *testing.T) {

	expectedSQL := `
SELECT film_actor.actor_id AS "film_actor.actor_id",
     film_actor.film_id AS "film_actor.film_id",
     film_actor.last_update AS "film_actor.last_update",
     film.film_id AS "film.film_id",
     film.title AS "film.title",
     film.description AS "film.description",
     film.release_year AS "film.release_year",
     film.language_id AS "film.language_id",
     film.rental_duration AS "film.rental_duration",
     film.rental_rate AS "film.rental_rate",
     film.length AS "film.length",
     film.replacement_cost AS "film.replacement_cost",
     film.rating AS "film.rating",
     film.last_update AS "film.last_update",
     film.special_features AS "film.special_features",
     film.fulltext AS "film.fulltext",
     language.language_id AS "language.language_id",
     language.name AS "language.name",
     language.last_update AS "language.last_update",
     actor.actor_id AS "actor.actor_id",
     actor.first_name AS "actor.first_name",
     actor.last_name AS "actor.last_name",
     actor.last_update AS "actor.last_update",
     inventory.inventory_id AS "inventory.inventory_id",
     inventory.film_id AS "inventory.film_id",
     inventory.store_id AS "inventory.store_id",
     inventory.last_update AS "inventory.last_update",
     rental.rental_id AS "rental.rental_id",
     rental.rental_date AS "rental.rental_date",
     rental.inventory_id AS "rental.inventory_id",
     rental.customer_id AS "rental.customer_id",
     rental.return_date AS "rental.return_date",
     rental.staff_id AS "rental.staff_id",
     rental.last_update AS "rental.last_update"
FROM dvds.film_actor
     INNER JOIN dvds.actor ON (film_actor.actor_id = actor.actor_id)
     INNER JOIN dvds.film ON (film_actor.film_id = film.film_id)
     INNER JOIN dvds.language ON (film.language_id = language.language_id)
     INNER JOIN dvds.inventory ON (inventory.film_id = film.film_id)
     INNER JOIN dvds.rental ON (rental.inventory_id = inventory.inventory_id)
ORDER BY film.film_id ASC
LIMIT 1000;
`
	for i := 0; i < 2; i++ {
		query := FilmActor.
			INNER_JOIN(Actor, FilmActor.ActorID.EQ(Actor.ActorID)).
			INNER_JOIN(Film, FilmActor.FilmID.EQ(Film.FilmID)).
			INNER_JOIN(Language, Film.LanguageID.EQ(Language.LanguageID)).
			INNER_JOIN(Inventory, Inventory.FilmID.EQ(Film.FilmID)).
			INNER_JOIN(Rental, Rental.InventoryID.EQ(Inventory.InventoryID)).
			SELECT(
				FilmActor.AllColumns,
				Film.AllColumns,
				Language.AllColumns,
				Actor.AllColumns,
				Inventory.AllColumns,
				Rental.AllColumns,
			).
			//WHERE(FilmActor.ActorID.GtEqL(1).AND(FilmActor.ActorID.LtEqL(2))).
			ORDER_BY(Film.FilmID.ASC()).
			LIMIT(1000)

		testutils.AssertDebugStatementSql(t, query, expectedSQL, int64(1000))

		var languageActorFilm []struct {
			model.Language

			Films []struct {
				model.Film
				Actors []struct {
					model.Actor
				}

				Inventory []struct {
					model.Inventory

					Rental []model.Rental
				}
			}
		}

		err := query.Query(db, &languageActorFilm)

		assert.NilError(t, err)
		assert.Equal(t, len(languageActorFilm), 1)
		assert.Equal(t, len(languageActorFilm[0].Films), 10)
		assert.Equal(t, len(languageActorFilm[0].Films[0].Actors), 10)
	}
}

func TestJoinQuerySlice(t *testing.T) {
	expectedSQL := `
SELECT language.language_id AS "language.language_id",
     language.name AS "language.name",
     language.last_update AS "language.last_update",
     film.film_id AS "film.film_id",
     film.title AS "film.title",
     film.description AS "film.description",
     film.release_year AS "film.release_year",
     film.language_id AS "film.language_id",
     film.rental_duration AS "film.rental_duration",
     film.rental_rate AS "film.rental_rate",
     film.length AS "film.length",
     film.replacement_cost AS "film.replacement_cost",
     film.rating AS "film.rating",
     film.last_update AS "film.last_update",
     film.special_features AS "film.special_features",
     film.fulltext AS "film.fulltext"
FROM dvds.film
     INNER JOIN dvds.language ON (film.language_id = language.language_id)
WHERE film.rating = 'NC-17'
LIMIT 15;
`

	type FilmsPerLanguage struct {
		Language *model.Language
		Film     []model.Film
	}

	filmsPerLanguage := []FilmsPerLanguage{}
	limit := 15

	query := Film.
		INNER_JOIN(Language, Film.LanguageID.EQ(Language.LanguageID)).
		SELECT(Language.AllColumns, Film.AllColumns).
		WHERE(Film.Rating.EQ(enum.MpaaRating.Nc17)).
		LIMIT(15)

	testutils.AssertDebugStatementSql(t, query, expectedSQL, int64(15))

	err := query.Query(db, &filmsPerLanguage)

	assert.NilError(t, err)
	assert.Equal(t, len(filmsPerLanguage), 1)
	assert.Equal(t, len(filmsPerLanguage[0].Film), limit)

	englishFilms := filmsPerLanguage[0]

	assert.Equal(t, *englishFilms.Film[0].Rating, model.MpaaRating_Nc17)

	filmsPerLanguageWithPtrs := []*FilmsPerLanguage{}
	err = query.Query(db, &filmsPerLanguageWithPtrs)

	assert.NilError(t, err)
	assert.Equal(t, len(filmsPerLanguage), 1)
	assert.Equal(t, len(filmsPerLanguage[0].Film), limit)
}

func TestExecution1(t *testing.T) {
	stmt := City.
		INNER_JOIN(Address, Address.CityID.EQ(City.CityID)).
		INNER_JOIN(Customer, Customer.AddressID.EQ(Address.AddressID)).
		SELECT(
			City.CityID,
			City.City,
			Address.AddressID,
			Address.Address,
			Customer.CustomerID,
			Customer.LastName,
		).
		WHERE(City.City.EQ(String("London")).OR(City.City.EQ(String("York")))).
		ORDER_BY(City.CityID, Address.AddressID, Customer.CustomerID)

	testutils.AssertDebugStatementSql(t, stmt, `
SELECT city.city_id AS "city.city_id",
     city.city AS "city.city",
     address.address_id AS "address.address_id",
     address.address AS "address.address",
     customer.customer_id AS "customer.customer_id",
     customer.last_name AS "customer.last_name"
FROM dvds.city
     INNER JOIN dvds.address ON (address.city_id = city.city_id)
     INNER JOIN dvds.customer ON (customer.address_id = address.address_id)
WHERE (city.city = 'London') OR (city.city = 'York')
ORDER BY city.city_id, address.address_id, customer.customer_id;
`, "London", "York")

	var dest []struct {
		model.City

		Customers []struct {
			model.Customer

			Address model.Address
		}
	}

	err := stmt.Query(db, &dest)

	assert.NilError(t, err)

	assert.Equal(t, len(dest), 2)
	assert.Equal(t, dest[0].City.City, "London")
	assert.Equal(t, dest[1].City.City, "York")
	assert.Equal(t, len(dest[0].Customers), 2)
	assert.Equal(t, dest[0].Customers[0].LastName, "Hoffman")
	assert.Equal(t, dest[0].Customers[1].LastName, "Vines")

}

func TestExecution2(t *testing.T) {

	type MyAddress struct {
		ID          int32 `sql:"primary_key"`
		AddressLine string
	}

	type MyCustomer struct {
		ID       int32 `sql:"primary_key"`
		LastName *string

		Address MyAddress
	}

	type MyCity struct {
		ID   int32 `sql:"primary_key"`
		Name string

		Customers []MyCustomer
	}

	dest := []MyCity{}

	stmt := City.
		INNER_JOIN(Address, Address.CityID.EQ(City.CityID)).
		INNER_JOIN(Customer, Customer.AddressID.EQ(Address.AddressID)).
		SELECT(
			City.CityID.AS("my_city.id"),
			City.City.AS("myCity.Name"),
			Address.AddressID.AS("My_Address.id"),
			Address.Address.AS("my address.address line"),
			Customer.CustomerID.AS("my_customer.id"),
			Customer.LastName.AS("my_customer.last_name"),
		).
		WHERE(City.City.EQ(String("London")).OR(City.City.EQ(String("York")))).
		ORDER_BY(City.CityID, Address.AddressID, Customer.CustomerID)

	testutils.AssertDebugStatementSql(t, stmt, `
SELECT city.city_id AS "my_city.id",
     city.city AS "myCity.Name",
     address.address_id AS "My_Address.id",
     address.address AS "my address.address line",
     customer.customer_id AS "my_customer.id",
     customer.last_name AS "my_customer.last_name"
FROM dvds.city
     INNER JOIN dvds.address ON (address.city_id = city.city_id)
     INNER JOIN dvds.customer ON (customer.address_id = address.address_id)
WHERE (city.city = 'London') OR (city.city = 'York')
ORDER BY city.city_id, address.address_id, customer.customer_id;
`, "London", "York")

	err := stmt.Query(db, &dest)

	assert.NilError(t, err)

	assert.Equal(t, len(dest), 2)
	assert.Equal(t, dest[0].Name, "London")
	assert.Equal(t, dest[1].Name, "York")
	assert.Equal(t, len(dest[0].Customers), 2)
	assert.Equal(t, *dest[0].Customers[0].LastName, "Hoffman")
	assert.Equal(t, *dest[0].Customers[1].LastName, "Vines")

}

func TestExecution3(t *testing.T) {

	var dest []struct {
		CityID   int32 `sql:"primary_key"`
		CityName string

		Customers []struct {
			CustomerID int32 `sql:"primary_key"`
			LastName   *string

			Address struct {
				AddressID   int32 `sql:"primary_key"`
				AddressLine string
			}
		}
	}

	stmt := City.
		INNER_JOIN(Address, Address.CityID.EQ(City.CityID)).
		INNER_JOIN(Customer, Customer.AddressID.EQ(Address.AddressID)).
		SELECT(
			City.CityID.AS("city_id"),
			City.City.AS("city_name"),
			Customer.CustomerID.AS("customer_id"),
			Customer.LastName.AS("last_name"),
			Address.AddressID.AS("address_id"),
			Address.Address.AS("address_line"),
		).
		WHERE(City.City.EQ(String("London")).OR(City.City.EQ(String("York")))).
		ORDER_BY(City.CityID, Address.AddressID, Customer.CustomerID)

	testutils.AssertDebugStatementSql(t, stmt, `
SELECT city.city_id AS "city_id",
     city.city AS "city_name",
     customer.customer_id AS "customer_id",
     customer.last_name AS "last_name",
     address.address_id AS "address_id",
     address.address AS "address_line"
FROM dvds.city
     INNER JOIN dvds.address ON (address.city_id = city.city_id)
     INNER JOIN dvds.customer ON (customer.address_id = address.address_id)
WHERE (city.city = 'London') OR (city.city = 'York')
ORDER BY city.city_id, address.address_id, customer.customer_id;
`, "London", "York")

	err := stmt.Query(db, &dest)

	assert.NilError(t, err)

	assert.Equal(t, len(dest), 2)
	assert.Equal(t, dest[0].CityName, "London")
	assert.Equal(t, dest[1].CityName, "York")
	assert.Equal(t, len(dest[0].Customers), 2)
	assert.Equal(t, *dest[0].Customers[0].LastName, "Hoffman")
	assert.Equal(t, *dest[0].Customers[1].LastName, "Vines")
}

func TestExecution4(t *testing.T) {

	var dest []struct {
		CityID   int32  `sql:"primary_key" alias:"city.city_id"`
		CityName string `alias:"city.city"`

		Customers []struct {
			CustomerID int32   `sql:"primary_key" alias:"customer_id"`
			LastName   *string `alias:"last_name"`

			Address struct {
				AddressID   int32  `sql:"primary_key" alias:"AddressId"`
				AddressLine string `alias:"address.address"`
			} `alias:"address.*"`
		} `alias:"customer"`
	}

	stmt := City.
		INNER_JOIN(Address, Address.CityID.EQ(City.CityID)).
		INNER_JOIN(Customer, Customer.AddressID.EQ(Address.AddressID)).
		SELECT(
			City.CityID,
			City.City,
			Customer.CustomerID,
			Customer.LastName,
			Address.AddressID,
			Address.Address,
		).
		WHERE(City.City.EQ(String("London")).OR(City.City.EQ(String("York")))).
		ORDER_BY(City.CityID, Address.AddressID, Customer.CustomerID)

	testutils.AssertDebugStatementSql(t, stmt, `
SELECT city.city_id AS "city.city_id",
     city.city AS "city.city",
     customer.customer_id AS "customer.customer_id",
     customer.last_name AS "customer.last_name",
     address.address_id AS "address.address_id",
     address.address AS "address.address"
FROM dvds.city
     INNER JOIN dvds.address ON (address.city_id = city.city_id)
     INNER JOIN dvds.customer ON (customer.address_id = address.address_id)
WHERE (city.city = 'London') OR (city.city = 'York')
ORDER BY city.city_id, address.address_id, customer.customer_id;
`, "London", "York")

	err := stmt.Query(db, &dest)

	assert.NilError(t, err)
	assert.Equal(t, len(dest), 2)
	testutils.AssertJSON(t, `
[
	{
		"CityID": 312,
		"CityName": "London",
		"Customers": [
			{
				"CustomerID": 252,
				"LastName": "Hoffman",
				"Address": {
					"AddressID": 256,
					"AddressLine": "1497 Yuzhou Drive"
				}
			},
			{
				"CustomerID": 512,
				"LastName": "Vines",
				"Address": {
					"AddressID": 517,
					"AddressLine": "548 Uruapan Street"
				}
			}
		]
	},
	{
		"CityID": 589,
		"CityName": "York",
		"Customers": [
			{
				"CustomerID": 497,
				"LastName": "Sledge",
				"Address": {
					"AddressID": 502,
					"AddressLine": "1515 Korla Way"
				}
			}
		]
	}
]
`, dest)
}

func TestJoinQuerySliceWithPtrs(t *testing.T) {
	type FilmsPerLanguage struct {
		Language model.Language
		Film     *[]*model.Film
	}

	limit := int64(3)

	query := Film.INNER_JOIN(Language, Film.LanguageID.EQ(Language.LanguageID)).
		SELECT(Language.AllColumns, Film.AllColumns).
		LIMIT(limit)

	filmsPerLanguageWithPtrs := []*FilmsPerLanguage{}
	err := query.Query(db, &filmsPerLanguageWithPtrs)

	assert.NilError(t, err)
	assert.Equal(t, len(filmsPerLanguageWithPtrs), 1)
	assert.Equal(t, len(*filmsPerLanguageWithPtrs[0].Film), int(limit))
}

func TestSelect_WithoutUniqueColumnSelected(t *testing.T) {
	query := Customer.SELECT(Customer.FirstName, Customer.LastName, Customer.Email)

	customers := []model.Customer{}

	err := query.Query(db, &customers)

	assert.NilError(t, err)

	//spew.Dump(customers)

	assert.Equal(t, len(customers), 599)
}

func TestSelectOrderByAscDesc(t *testing.T) {
	customersAsc := []model.Customer{}

	err := Customer.SELECT(Customer.CustomerID, Customer.FirstName, Customer.LastName).
		ORDER_BY(Customer.FirstName.ASC()).
		Query(db, &customersAsc)

	assert.NilError(t, err)

	firstCustomerAsc := customersAsc[0]
	lastCustomerAsc := customersAsc[len(customersAsc)-1]

	customersDesc := []model.Customer{}
	err = Customer.SELECT(Customer.CustomerID, Customer.FirstName, Customer.LastName).
		ORDER_BY(Customer.FirstName.DESC()).
		Query(db, &customersDesc)

	assert.NilError(t, err)

	firstCustomerDesc := customersDesc[0]
	lastCustomerDesc := customersDesc[len(customersAsc)-1]

	assert.DeepEqual(t, firstCustomerAsc, lastCustomerDesc)
	assert.DeepEqual(t, lastCustomerAsc, firstCustomerDesc)

	customersAscDesc := []model.Customer{}
	err = Customer.SELECT(Customer.CustomerID, Customer.FirstName, Customer.LastName).
		ORDER_BY(Customer.FirstName.ASC(), Customer.LastName.DESC()).
		Query(db, &customersAscDesc)

	assert.NilError(t, err)

	customerAscDesc326 := model.Customer{
		CustomerID: 67,
		FirstName:  "Kelly",
		LastName:   "Torres",
	}

	customerAscDesc327 := model.Customer{
		CustomerID: 546,
		FirstName:  "Kelly",
		LastName:   "Knott",
	}

	assert.DeepEqual(t, customerAscDesc326, customersAscDesc[326])
	assert.DeepEqual(t, customerAscDesc327, customersAscDesc[327])
}

func TestSelectFullJoin(t *testing.T) {
	expectedSQL := `
SELECT customer.customer_id AS "customer.customer_id",
     customer.store_id AS "customer.store_id",
     customer.first_name AS "customer.first_name",
     customer.last_name AS "customer.last_name",
     customer.email AS "customer.email",
     customer.address_id AS "customer.address_id",
     customer.activebool AS "customer.activebool",
     customer.create_date AS "customer.create_date",
     customer.last_update AS "customer.last_update",
     customer.active AS "customer.active",
     address.address_id AS "address.address_id",
     address.address AS "address.address",
     address.address2 AS "address.address2",
     address.district AS "address.district",
     address.city_id AS "address.city_id",
     address.postal_code AS "address.postal_code",
     address.phone AS "address.phone",
     address.last_update AS "address.last_update"
FROM dvds.customer
     FULL JOIN dvds.address ON (customer.address_id = address.address_id)
ORDER BY customer.customer_id ASC;
`
	query := Customer.
		FULL_JOIN(Address, Customer.AddressID.EQ(Address.AddressID)).
		SELECT(Customer.AllColumns, Address.AllColumns).
		ORDER_BY(Customer.CustomerID.ASC())

	testutils.AssertDebugStatementSql(t, query, expectedSQL)

	allCustomersAndAddress := []struct {
		Address  *model.Address
		Customer *model.Customer
	}{}

	err := query.Query(db, &allCustomersAndAddress)

	assert.NilError(t, err)
	assert.Equal(t, len(allCustomersAndAddress), 603)

	assert.DeepEqual(t, allCustomersAndAddress[0].Customer, &customer0)
	assert.Assert(t, allCustomersAndAddress[0].Address != nil)

	lastCustomerAddress := allCustomersAndAddress[len(allCustomersAndAddress)-1]

	assert.Assert(t, lastCustomerAddress.Customer == nil)
	assert.Assert(t, lastCustomerAddress.Address != nil)

}

func TestSelectFullCrossJoin(t *testing.T) {
	expectedSQL := `
SELECT customer.customer_id AS "customer.customer_id",
     customer.store_id AS "customer.store_id",
     customer.first_name AS "customer.first_name",
     customer.last_name AS "customer.last_name",
     customer.email AS "customer.email",
     customer.address_id AS "customer.address_id",
     customer.activebool AS "customer.activebool",
     customer.create_date AS "customer.create_date",
     customer.last_update AS "customer.last_update",
     customer.active AS "customer.active",
     address.address_id AS "address.address_id",
     address.address AS "address.address",
     address.address2 AS "address.address2",
     address.district AS "address.district",
     address.city_id AS "address.city_id",
     address.postal_code AS "address.postal_code",
     address.phone AS "address.phone",
     address.last_update AS "address.last_update"
FROM dvds.customer
     CROSS JOIN dvds.address
ORDER BY customer.customer_id ASC
LIMIT 1000;
`
	query := Customer.
		CROSS_JOIN(Address).
		SELECT(Customer.AllColumns, Address.AllColumns).
		ORDER_BY(Customer.CustomerID.ASC()).
		LIMIT(1000)

	testutils.AssertDebugStatementSql(t, query, expectedSQL, int64(1000))

	var customerAddresCrosJoined []struct {
		model.Customer
		model.Address
	}

	err := query.Query(db, &customerAddresCrosJoined)

	assert.Equal(t, len(customerAddresCrosJoined), 1000)

	assert.NilError(t, err)
}

func TestSelectSelfJoin(t *testing.T) {
	expectedSQL := `
SELECT f1.film_id AS "f1.film_id",
     f1.title AS "f1.title",
     f1.description AS "f1.description",
     f1.release_year AS "f1.release_year",
     f1.language_id AS "f1.language_id",
     f1.rental_duration AS "f1.rental_duration",
     f1.rental_rate AS "f1.rental_rate",
     f1.length AS "f1.length",
     f1.replacement_cost AS "f1.replacement_cost",
     f1.rating AS "f1.rating",
     f1.last_update AS "f1.last_update",
     f1.special_features AS "f1.special_features",
     f1.fulltext AS "f1.fulltext",
     f2.film_id AS "f2.film_id",
     f2.title AS "f2.title",
     f2.description AS "f2.description",
     f2.release_year AS "f2.release_year",
     f2.language_id AS "f2.language_id",
     f2.rental_duration AS "f2.rental_duration",
     f2.rental_rate AS "f2.rental_rate",
     f2.length AS "f2.length",
     f2.replacement_cost AS "f2.replacement_cost",
     f2.rating AS "f2.rating",
     f2.last_update AS "f2.last_update",
     f2.special_features AS "f2.special_features",
     f2.fulltext AS "f2.fulltext"
FROM dvds.film AS f1
     INNER JOIN dvds.film AS f2 ON ((f1.film_id < f2.film_id) AND (f1.length = f2.length))
ORDER BY f1.film_id ASC;
`
	f1 := Film.AS("f1")

	f2 := Film.AS("f2")

	query := f1.
		INNER_JOIN(f2, f1.FilmID.LT(f2.FilmID).AND(f1.Length.EQ(f2.Length))).
		SELECT(f1.AllColumns, f2.AllColumns).
		ORDER_BY(f1.FilmID.ASC())

	testutils.AssertDebugStatementSql(t, query, expectedSQL)

	type F1 model.Film
	type F2 model.Film

	theSameLengthFilms := []struct {
		F1 F1
		F2 F2
	}{}

	err := query.Query(db, &theSameLengthFilms)

	assert.NilError(t, err)

	//spew.Dump(theSameLengthFilms)

	//assert.Equal(t, len(theSameLengthFilms), 100)
}

func TestSelectAliasColumn(t *testing.T) {
	expectedSQL := `
SELECT f1.title AS "thesame_length_films.title1",
     f2.title AS "thesame_length_films.title2",
     f1.length AS "thesame_length_films.length"
FROM dvds.film AS f1
     INNER JOIN dvds.film AS f2 ON ((f1.film_id != f2.film_id) AND (f1.length = f2.length))
ORDER BY f1.length ASC, f1.title ASC, f2.title ASC
LIMIT 1000;
`
	f1 := Film.AS("f1")
	f2 := Film.AS("f2")

	f1.FilmID.EQ(Int(11))

	query := f1.
		INNER_JOIN(f2, f1.FilmID.NOT_EQ(f2.FilmID).AND(f1.Length.EQ(f2.Length))).
		SELECT(f1.Title.AS("thesame_length_films.title1"),
			f2.Title.AS("thesame_length_films.title2"),
			f1.Length.AS("thesame_length_films.length")).
		ORDER_BY(f1.Length.ASC(), f1.Title.ASC(), f2.Title.ASC()).
		LIMIT(1000)

	testutils.AssertDebugStatementSql(t, query, expectedSQL, int64(1000))

	type thesameLengthFilms struct {
		Title1 string
		Title2 string
		Length int16
	}
	films := []thesameLengthFilms{}

	err := query.Query(db, &films)

	assert.NilError(t, err)

	//spew.Dump(films)

	assert.Equal(t, len(films), 1000)
	assert.DeepEqual(t, films[0], thesameLengthFilms{"Alien Center", "Iron Moon", 46})
}

func TestSubQuery(t *testing.T) {
	expectedQuery := `
SELECT actor.actor_id AS "actor.actor_id",
     actor.first_name AS "actor.first_name",
     actor.last_name AS "actor.last_name",
     actor.last_update AS "actor.last_update",
     film_actor.actor_id AS "film_actor.actor_id",
     film_actor.film_id AS "film_actor.film_id",
     film_actor.last_update AS "film_actor.last_update",
     "rFilms"."film.film_id" AS "film.film_id",
     "rFilms"."film.title" AS "film.title",
     "rFilms"."film.rating" AS "film.rating"
FROM dvds.actor
     INNER JOIN dvds.film_actor ON (actor.actor_id = film_actor.film_id)
     INNER JOIN (
          SELECT film.film_id AS "film.film_id",
               film.title AS "film.title",
               film.rating AS "film.rating"
          FROM dvds.film
          WHERE film.rating = 'R'
     ) AS "rFilms" ON (film_actor.film_id = "rFilms"."film.film_id");
`

	rRatingFilms := Film.
		SELECT(
			Film.FilmID,
			Film.Title,
			Film.Rating,
		).
		WHERE(Film.Rating.EQ(enum.MpaaRating.R)).
		AsTable("rFilms")

	rFilmID := Film.FilmID.From(rRatingFilms)

	query := Actor.
		INNER_JOIN(FilmActor, Actor.ActorID.EQ(FilmActor.FilmID)).
		INNER_JOIN(rRatingFilms, FilmActor.FilmID.EQ(rFilmID)).
		SELECT(
			Actor.AllColumns,
			FilmActor.AllColumns,
			rRatingFilms.AllColumns(),
		)

	testutils.AssertDebugStatementSql(t, query, expectedQuery)

	dest := []model.Actor{}

	err := query.Query(db, &dest)

	assert.NilError(t, err)
}

func TestSelectFunctions(t *testing.T) {
	expectedQuery := `
SELECT MAX(film.rental_rate) AS "max_film_rate"
FROM dvds.film;
`
	query := Film.SELECT(
		MAXf(Film.RentalRate).AS("max_film_rate"),
	)

	testutils.AssertDebugStatementSql(t, query, expectedQuery)

	ret := struct {
		MaxFilmRate float64
	}{}

	err := query.Query(db, &ret)

	assert.NilError(t, err)
	assert.Equal(t, ret.MaxFilmRate, 4.99)
}

func TestSelectQueryScalar(t *testing.T) {
	expectedSQL := `
SELECT film.film_id AS "film.film_id",
     film.title AS "film.title",
     film.description AS "film.description",
     film.release_year AS "film.release_year",
     film.language_id AS "film.language_id",
     film.rental_duration AS "film.rental_duration",
     film.rental_rate AS "film.rental_rate",
     film.length AS "film.length",
     film.replacement_cost AS "film.replacement_cost",
     film.rating AS "film.rating",
     film.last_update AS "film.last_update",
     film.special_features AS "film.special_features",
     film.fulltext AS "film.fulltext"
FROM dvds.film
WHERE film.rental_rate = (
          SELECT MAX(film.rental_rate)
          FROM dvds.film
     )
ORDER BY film.film_id ASC;
`

	maxFilmRentalRate := FloatExp(
		Film.
			SELECT(MAXf(Film.RentalRate)),
	)

	query := Film.
		SELECT(Film.AllColumns).
		WHERE(Film.RentalRate.EQ(maxFilmRentalRate)).
		ORDER_BY(Film.FilmID.ASC())

	testutils.AssertDebugStatementSql(t, query, expectedSQL)

	maxRentalRateFilms := []model.Film{}
	err := query.Query(db, &maxRentalRateFilms)

	assert.NilError(t, err)

	assert.Equal(t, len(maxRentalRateFilms), 336)

	gRating := model.MpaaRating_G

	assert.DeepEqual(t, maxRentalRateFilms[0], model.Film{
		FilmID:          2,
		Title:           "Ace Goldfinger",
		Description:     StringPtr("A Astounding Epistle of a Database Administrator And a Explorer who must Find a Car in Ancient China"),
		ReleaseYear:     Int32Ptr(2006),
		LanguageID:      1,
		RentalRate:      4.99,
		Length:          Int16Ptr(48),
		ReplacementCost: 12.99,
		Rating:          &gRating,
		RentalDuration:  3,
		LastUpdate:      *testutils.TimestampWithoutTimeZone("2013-05-26 14:50:58.951", 3),
		SpecialFeatures: StringPtr("{Trailers,\"Deleted Scenes\"}"),
		Fulltext:        "'ace':1 'administr':9 'ancient':19 'astound':4 'car':17 'china':20 'databas':8 'epistl':5 'explor':12 'find':15 'goldfing':2 'must':14",
	})
}

func TestSelectGroupByHaving(t *testing.T) {
	expectedSQL := `
SELECT customer.customer_id AS "customer.customer_id",
     customer.store_id AS "customer.store_id",
     customer.first_name AS "customer.first_name",
     customer.last_name AS "customer.last_name",
     customer.email AS "customer.email",
     customer.address_id AS "customer.address_id",
     customer.activebool AS "customer.activebool",
     customer.create_date AS "customer.create_date",
     customer.last_update AS "customer.last_update",
     customer.active AS "customer.active",
     SUM(payment.amount) AS "amount.sum",
     AVG(payment.amount) AS "amount.avg",
     MAX(payment.amount) AS "amount.max",
     MIN(payment.amount) AS "amount.min",
     COUNT(payment.amount) AS "amount.count"
FROM dvds.payment
     INNER JOIN dvds.customer ON (customer.customer_id = payment.customer_id)
GROUP BY customer.customer_id
HAVING SUM(payment.amount) > 125.6
ORDER BY customer.customer_id, SUM(payment.amount) ASC;
`
	query := Payment.
		INNER_JOIN(Customer, Customer.CustomerID.EQ(Payment.CustomerID)).
		SELECT(
			Customer.AllColumns,

			SUMf(Payment.Amount).AS("amount.sum"),
			AVG(Payment.Amount).AS("amount.avg"),
			MAXf(Payment.Amount).AS("amount.max"),
			MINf(Payment.Amount).AS("amount.min"),
			COUNT(Payment.Amount).AS("amount.count"),
		).
		GROUP_BY(Customer.CustomerID).
		HAVING(
			SUMf(Payment.Amount).GT(Float(125.6)),
		).
		ORDER_BY(
			Customer.CustomerID, SUMf(Payment.Amount).ASC(),
		)

	//fmt.Println(query.DebugSql())

	testutils.AssertDebugStatementSql(t, query, expectedSQL, float64(125.6))

	var dest []struct {
		model.Customer

		Amount struct {
			Sum   float64
			Avg   float64
			Max   float64
			Min   float64
			Count int64
		} `alias:"amount"`
	}

	err := query.Query(db, &dest)

	assert.NilError(t, err)

	//testutils.JsonPrint(dest)

	assert.Equal(t, len(dest), 104)

	//testutils.JsonSave(dest, "postgres/testdata/customer_payment_sum.json")
	testutils.AssertJSONFile(t, dest, "postgres/testdata/customer_payment_sum.json")
}

func TestSelectGroupBy2(t *testing.T) {
	expectedSQL := `
SELECT customer.customer_id AS "customer.customer_id",
     customer.store_id AS "customer.store_id",
     customer.first_name AS "customer.first_name",
     customer.last_name AS "customer.last_name",
     customer.email AS "customer.email",
     customer.address_id AS "customer.address_id",
     customer.activebool AS "customer.activebool",
     customer.create_date AS "customer.create_date",
     customer.last_update AS "customer.last_update",
     customer.active AS "customer.active",
     customer_payment_sum."amount_sum" AS "CustomerWithAmounts.AmountSum"
FROM dvds.customer
     INNER JOIN (
          SELECT payment.customer_id AS "payment.customer_id",
               SUM(payment.amount) AS "amount_sum"
          FROM dvds.payment
          GROUP BY payment.customer_id
     ) AS customer_payment_sum ON (customer.customer_id = customer_payment_sum."payment.customer_id")
ORDER BY customer_payment_sum."amount_sum" ASC;
`

	customersPayments := Payment.
		SELECT(
			Payment.CustomerID,
			SUMf(Payment.Amount).AS("amount_sum"),
		).
		GROUP_BY(Payment.CustomerID).
		AsTable("customer_payment_sum")

	customerID := Payment.CustomerID.From(customersPayments)
	amountSum := FloatColumn("amount_sum").From(customersPayments)

	query := Customer.
		INNER_JOIN(customersPayments, Customer.CustomerID.EQ(customerID)).
		SELECT(
			Customer.AllColumns,
			amountSum.AS("CustomerWithAmounts.AmountSum"),
		).
		ORDER_BY(amountSum.ASC())

	testutils.AssertDebugStatementSql(t, query, expectedSQL)

	type CustomerWithAmounts struct {
		Customer  *model.Customer
		AmountSum float64
	}
	customersWithAmounts := []CustomerWithAmounts{}

	err := query.Query(db, &customersWithAmounts)
	assert.NilError(t, err)
	assert.Equal(t, len(customersWithAmounts), 599)

	assert.DeepEqual(t, customersWithAmounts[0].Customer, &model.Customer{
		CustomerID: 318,
		StoreID:    1,
		FirstName:  "Brian",
		LastName:   "Wyman",
		AddressID:  323,
		Email:      StringPtr("brian.wyman@sakilacustomer.org"),
		Activebool: true,
		CreateDate: *testutils.TimestampWithoutTimeZone("2006-02-14 00:00:00", 0),
		LastUpdate: testutils.TimestampWithoutTimeZone("2013-05-26 14:49:45.738", 3),
		Active:     Int32Ptr(1),
	})

	assert.Equal(t, customersWithAmounts[0].AmountSum, 27.93)
}

func TestSelectStaff(t *testing.T) {
	staffs := []model.Staff{}

	err := Staff.SELECT(Staff.AllColumns).Query(db, &staffs)

	assert.NilError(t, err)

	testutils.AssertJSON(t, `
[
	{
		"StaffID": 1,
		"FirstName": "Mike",
		"LastName": "Hillyer",
		"AddressID": 3,
		"Email": "Mike.Hillyer@sakilastaff.com",
		"StoreID": 1,
		"Active": true,
		"Username": "Mike",
		"Password": "8cb2237d0679ca88db6464eac60da96345513964",
		"LastUpdate": "2006-05-16T16:13:11.79328Z",
		"Picture": "iVBORw0KWgo="
	},
	{
		"StaffID": 2,
		"FirstName": "Jon",
		"LastName": "Stephens",
		"AddressID": 4,
		"Email": "Jon.Stephens@sakilastaff.com",
		"StoreID": 2,
		"Active": true,
		"Username": "Jon",
		"Password": "8cb2237d0679ca88db6464eac60da96345513964",
		"LastUpdate": "2006-05-16T16:13:11.79328Z",
		"Picture": null
	}
]
`, staffs)
}

func TestSelectTimeColumns(t *testing.T) {

	expectedSQL := `
SELECT payment.payment_id AS "payment.payment_id",
     payment.customer_id AS "payment.customer_id",
     payment.staff_id AS "payment.staff_id",
     payment.rental_id AS "payment.rental_id",
     payment.amount AS "payment.amount",
     payment.payment_date AS "payment.payment_date"
FROM dvds.payment
WHERE payment.payment_date < '2007-02-14 22:16:01.000'::timestamp without time zone
ORDER BY payment.payment_date ASC;
`

	query := Payment.SELECT(Payment.AllColumns).
		WHERE(Payment.PaymentDate.LT(Timestamp(2007, 02, 14, 22, 16, 01, 0))).
		ORDER_BY(Payment.PaymentDate.ASC())

	testutils.AssertDebugStatementSql(t, query, expectedSQL, "2007-02-14 22:16:01.000")

	payments := []model.Payment{}

	err := query.Query(db, &payments)

	assert.NilError(t, err)

	//spew.Dump(payments)

	assert.Equal(t, len(payments), 9)
	assert.DeepEqual(t, payments[0], model.Payment{
		PaymentID:   17793,
		CustomerID:  416,
		StaffID:     2,
		RentalID:    1158,
		Amount:      2.99,
		PaymentDate: *testutils.TimestampWithoutTimeZone("2007-02-14 21:21:59.996577", 6),
	})
}

func TestUnion(t *testing.T) {
	expectedQuery := `
(
     (
          SELECT payment.payment_id AS "payment.payment_id",
               payment.amount AS "payment.amount"
          FROM dvds.payment
          WHERE payment.amount <= 100
     )
     UNION ALL
     (
          SELECT payment.payment_id AS "payment.payment_id",
               payment.amount AS "payment.amount"
          FROM dvds.payment
          WHERE payment.amount >= 200
     )
)
ORDER BY "payment.payment_id" ASC, "payment.amount" DESC
LIMIT 10
OFFSET 20;
`
	query := UNION_ALL(
		Payment.
			SELECT(Payment.PaymentID.AS("payment.payment_id"), Payment.Amount).
			WHERE(Payment.Amount.LT_EQ(Float(100))),
		Payment.
			SELECT(Payment.PaymentID, Payment.Amount).
			WHERE(Payment.Amount.GT_EQ(Float(200))),
	).
		ORDER_BY(IntegerColumn("payment.payment_id").ASC(), Payment.Amount.DESC()).
		LIMIT(10).
		OFFSET(20)

	testutils.AssertDebugStatementSql(t, query, expectedQuery, float64(100), float64(200), int64(10), int64(20))

	dest := []model.Payment{}

	err := query.Query(db, &dest)

	assert.NilError(t, err)
	assert.Equal(t, len(dest), 10)
	assert.DeepEqual(t, dest[0], model.Payment{
		PaymentID: 17523,
		Amount:    4.99,
	})
	assert.DeepEqual(t, dest[1], model.Payment{
		PaymentID: 17524,
		Amount:    0.99,
	})
	assert.DeepEqual(t, dest[9], model.Payment{
		PaymentID: 17532,
		Amount:    8.99,
	})
}

func TestAllSetOperators(t *testing.T) {

	select1 := Payment.SELECT(Payment.AllColumns).WHERE(Payment.PaymentID.GT_EQ(Int(17600)).AND(Payment.PaymentID.LT(Int(17610))))
	select2 := Payment.SELECT(Payment.AllColumns).WHERE(Payment.PaymentID.GT_EQ(Int(17620)).AND(Payment.PaymentID.LT(Int(17630))))

	type setOperator func(lhs, rhs jet.SelectStatement, selects ...jet.SelectStatement) jet.SelectStatement
	operators := []setOperator{
		UNION,
		UNION_ALL,
		INTERSECT,
		INTERSECT_ALL,
		//EXCEPT,
		//EXCEPT_ALL,
	}

	expectedDestLen := []int{
		20,
		20,
		0,
		0,
		10,
		10,
	}

	for i, operator := range operators {
		query := operator(select1, select2)

		dest := []model.Payment{}
		err := query.Query(db, &dest)

		assert.NilError(t, err)
		assert.Equal(t, len(dest), expectedDestLen[i])
	}
}

func TestSelectWithCase(t *testing.T) {
	expectedQuery := `
SELECT (CASE payment.staff_id WHEN 1 THEN 'ONE' WHEN 2 THEN 'TWO' WHEN 3 THEN 'THREE' ELSE 'OTHER' END) AS "staff_id_num"
FROM dvds.payment
ORDER BY payment.payment_id ASC
LIMIT 20;
`
	query := Payment.SELECT(
		CASE(Payment.StaffID).
			WHEN(Int(1)).THEN(String("ONE")).
			WHEN(Int(2)).THEN(String("TWO")).
			WHEN(Int(3)).THEN(String("THREE")).
			ELSE(String("OTHER")).AS("staff_id_num"),
	).
		ORDER_BY(Payment.PaymentID.ASC()).
		LIMIT(20)

	testutils.AssertDebugStatementSql(t, query, expectedQuery, int64(1), "ONE", int64(2), "TWO", int64(3), "THREE", "OTHER", int64(20))

	dest := []struct {
		StaffIDNum string
	}{}

	err := query.Query(db, &dest)

	assert.NilError(t, err)
	assert.Equal(t, len(dest), 20)
	assert.Equal(t, dest[0].StaffIDNum, "TWO")
	assert.Equal(t, dest[1].StaffIDNum, "ONE")
}

func getRowLockTestData() map[SelectLock]string {
	return map[SelectLock]string{
		UPDATE():        "UPDATE",
		NO_KEY_UPDATE(): "NO KEY UPDATE",
		SHARE():         "SHARE",
		KEY_SHARE():     "KEY SHARE",
	}
}

func TestRowLock(t *testing.T) {
	expectedSQL := `
SELECT *
FROM dvds.address
LIMIT 3
FOR`
	query := Address.
		SELECT(STAR).
		LIMIT(3)

	for lockType, lockTypeStr := range getRowLockTestData() {
		query.FOR(lockType)

		testutils.AssertDebugStatementSql(t, query, expectedSQL+" "+lockTypeStr+";\n", int64(3))

		tx, _ := db.Begin()

		res, err := query.Exec(tx)
		assert.NilError(t, err)
		rowsAffected, _ := res.RowsAffected()
		assert.Equal(t, rowsAffected, int64(3))

		err = tx.Rollback()
		assert.NilError(t, err)
	}

	for lockType, lockTypeStr := range getRowLockTestData() {
		query.FOR(lockType.NOWAIT())

		testutils.AssertDebugStatementSql(t, query, expectedSQL+" "+lockTypeStr+" NOWAIT;\n", int64(3))

		tx, _ := db.Begin()

		res, err := query.Exec(tx)
		assert.NilError(t, err)
		rowsAffected, _ := res.RowsAffected()
		assert.Equal(t, rowsAffected, int64(3))

		err = tx.Rollback()
		assert.NilError(t, err)
	}

	for lockType, lockTypeStr := range getRowLockTestData() {
		query.FOR(lockType.SKIP_LOCKED())

		testutils.AssertDebugStatementSql(t, query, expectedSQL+" "+lockTypeStr+" SKIP LOCKED;\n", int64(3))

		tx, _ := db.Begin()

		res, err := query.Exec(tx)
		assert.NilError(t, err)
		rowsAffected, _ := res.RowsAffected()
		assert.Equal(t, rowsAffected, int64(3))

		err = tx.Rollback()
		assert.NilError(t, err)
	}
}

func TestQuickStart(t *testing.T) {

	var expectedSQL = `
SELECT actor.actor_id AS "actor.actor_id",
     actor.first_name AS "actor.first_name",
     actor.last_name AS "actor.last_name",
     actor.last_update AS "actor.last_update",
     film.film_id AS "film.film_id",
     film.title AS "film.title",
     film.description AS "film.description",
     film.release_year AS "film.release_year",
     film.language_id AS "film.language_id",
     film.rental_duration AS "film.rental_duration",
     film.rental_rate AS "film.rental_rate",
     film.length AS "film.length",
     film.replacement_cost AS "film.replacement_cost",
     film.rating AS "film.rating",
     film.last_update AS "film.last_update",
     film.special_features AS "film.special_features",
     film.fulltext AS "film.fulltext",
     language.language_id AS "language.language_id",
     language.name AS "language.name",
     language.last_update AS "language.last_update",
     category.category_id AS "category.category_id",
     category.name AS "category.name",
     category.last_update AS "category.last_update"
FROM dvds.actor
     INNER JOIN dvds.film_actor ON (actor.actor_id = film_actor.actor_id)
     INNER JOIN dvds.film ON (film.film_id = film_actor.film_id)
     INNER JOIN dvds.language ON (language.language_id = film.language_id)
     INNER JOIN dvds.film_category ON (film_category.film_id = film.film_id)
     INNER JOIN dvds.category ON (category.category_id = film_category.category_id)
WHERE ((language.name = 'English') AND (category.name != 'Action')) AND (film.length > 180)
ORDER BY actor.actor_id ASC, film.film_id ASC;
`

	stmt := SELECT(
		Actor.ActorID, Actor.FirstName, Actor.LastName, Actor.LastUpdate, // list of all actor columns (equivalent to Actor.AllColumns)
		Film.AllColumns, // list of all film columns (equivalent to Film.FilmID, Film.Title, ...)
		Language.AllColumns,
		Category.AllColumns,
	).FROM(
		Actor.
			INNER_JOIN(FilmActor, Actor.ActorID.EQ(FilmActor.ActorID)). // INNER JOIN Actor with FilmActor on condition Actor.ActorID = FilmActor.ActorID
			INNER_JOIN(Film, Film.FilmID.EQ(FilmActor.FilmID)).         // then with Film, Language, FilmCategory and Category.
			INNER_JOIN(Language, Language.LanguageID.EQ(Film.LanguageID)).
			INNER_JOIN(FilmCategory, FilmCategory.FilmID.EQ(Film.FilmID)).
			INNER_JOIN(Category, Category.CategoryID.EQ(FilmCategory.CategoryID)),
	).WHERE(
		Language.Name.EQ(String("English")). // note that every column has type.
							AND(Category.Name.NOT_EQ(String("Action"))). // String column Language.Name and Category.Name can be compared only with string expression
							AND(Film.Length.GT(Int(180))),               // Film.Length is integer column and can be compared only with integer expression
	).ORDER_BY(
		Actor.ActorID.ASC(),
		Film.FilmID.ASC(),
	)

	testutils.AssertDebugStatementSql(t, stmt, expectedSQL, "English", "Action", int64(180))

	var dest []struct {
		model.Actor

		Films []struct {
			model.Film

			Language model.Language

			Categories []model.Category
		}
	}

	err := stmt.Query(db, &dest)
	assert.NilError(t, err)

	//jsonSave("./testdata/quick-start-dest.json", dest)
	testutils.AssertJSONFile(t, dest, "./postgres/testdata/quick-start-dest.json")

	var dest2 []struct {
		model.Category

		Films  []model.Film
		Actors []model.Actor
	}

	err = stmt.Query(db, &dest2)
	assert.NilError(t, err)

	//jsonSave("./testdata/quick-start-dest2.json", dest2)
	testutils.AssertJSONFile(t, dest2, "./postgres/testdata/quick-start-dest2.json")
}

func TestQuickStartWithSubQueries(t *testing.T) {

	filmLogerThan180 := Film.
		SELECT(Film.AllColumns).
		WHERE(Film.Length.GT(Int(180))).
		AsTable("films")

	filmID := Film.FilmID.From(filmLogerThan180)
	filmLanguageID := Film.LanguageID.From(filmLogerThan180)

	categoriesNotAction := Category.
		SELECT(Category.AllColumns).
		WHERE(Category.Name.NOT_EQ(String("Action"))).
		AsTable("categories")

	categoryID := Category.CategoryID.From(categoriesNotAction)

	stmt := Actor.
		INNER_JOIN(FilmActor, Actor.ActorID.EQ(FilmActor.ActorID)).
		INNER_JOIN(filmLogerThan180, filmID.EQ(FilmActor.FilmID)).
		INNER_JOIN(Language, Language.LanguageID.EQ(filmLanguageID)).
		INNER_JOIN(FilmCategory, FilmCategory.FilmID.EQ(filmID)).
		INNER_JOIN(categoriesNotAction, categoryID.EQ(FilmCategory.CategoryID)).
		SELECT(
			Actor.AllColumns,
			filmLogerThan180.AllColumns(),
			Language.AllColumns,
			categoriesNotAction.AllColumns(),
		).ORDER_BY(
		Actor.ActorID.ASC(),
		filmID.ASC(),
	)

	var dest []struct {
		model.Actor

		Films []struct {
			model.Film

			Language model.Language

			Categories []model.Category
		}
	}

	err := stmt.Query(db, &dest)
	assert.NilError(t, err)

	//jsonSave("./testdata/quick-start-dest.json", dest)
	testutils.AssertJSONFile(t, dest, "./postgres/testdata/quick-start-dest.json")

	var dest2 []struct {
		model.Category

		Films  []model.Film
		Actors []model.Actor
	}

	err = stmt.Query(db, &dest2)
	assert.NilError(t, err)

	//jsonSave("./testdata/quick-start-dest2.json", dest2)
	testutils.AssertJSONFile(t, dest2, "./postgres/testdata/quick-start-dest2.json")
}