package tests

import (
	"runtime/debug"
	"testing"

	"github.com/ardanlabs/encore/business/data/dbtest"
)

func Test_VProduct(t *testing.T) {
	t.Parallel()

	dbTest := dbtest.NewTest(t, url, "Test_VProduct")
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		dbTest.Teardown()
	}()

	// app := appTest{
	// 	Handler: mux.WebAPI(mux.Config{
	// 		Shutdown: make(chan os.Signal, 1),
	// 		Log:      dbTest.Log,
	// 		Auth:     dbTest.V1.Auth,
	// 		DB:       dbTest.DB,
	// 	}, all.Routes()),
	// 	userToken:  dbTest.TokenV1("user@example.com", "gophers"),
	// 	adminToken: dbTest.TokenV1("admin@example.com", "gophers"),
	// }

	// -------------------------------------------------------------------------

	// sd, err := createProductSeed(dbTest)
	// if err != nil {
	// 	t.Fatalf("Seeding error: %s", err)
	// }

	// -------------------------------------------------------------------------

	//app.test(t, vproductQuery200(sd), "vproduct-query-200")
}
