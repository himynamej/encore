package user_test

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	"encore.dev"
	eauth "encore.dev/beta/auth"
	"encore.dev/et"
	authsrv "github.com/ardanlabs/encore/api/services/auth"
	"github.com/ardanlabs/encore/api/services/sales"
	"github.com/ardanlabs/encore/api/services/sales/tests/apitest"
	"github.com/ardanlabs/encore/app/api/auth"
	"github.com/ardanlabs/encore/app/api/mid"
	"github.com/ardanlabs/encore/business/api/dbtest"
)

var url string

func TestMain(m *testing.M) {
	if encore.Meta().Environment.Name == "ci-test" {
		return
	}

	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	url, err = dbtest.StartDB()
	if err != nil {
		return 1, err
	}

	defer func() {
		err = dbtest.StopDB()
	}()

	return m.Run(), nil
}

// =============================================================================

func Test_User(t *testing.T) {
	t.Parallel()

	apitest := startTest(t, url, "Test_User")
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		apitest.DBTest.Teardown()
	}()

	// -------------------------------------------------------------------------

	sd, err := insertSeedData(apitest.DBTest, apitest.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	apitest.Run(t, userQueryOk(sd), "user-query-ok")
	apitest.Run(t, userQueryByIDOk(sd), "user-querybyid-ok")

	apitest.Run(t, userCreateOk(sd), "user-create-ok")
	apitest.Run(t, userCreateAuth(sd), "user-create-auth")
	apitest.Run(t, userCreateBad(sd), "user-create-bad")

	apitest.Run(t, userUpdateOk(sd), "user-update-ok")
	apitest.Run(t, userUpdateAuth(sd), "user-update-auth")
	apitest.Run(t, userUpdateBad(sd), "user-update-bad")

	apitest.Run(t, userDeleteOk(sd), "user-delete-ok")
	apitest.Run(t, userDeleteAuth(sd), "user-delete-auth")
}

func startTest(t *testing.T, url string, testName string) *apitest.AppTest {
	dbTest := dbtest.NewTest(t, url, testName)

	// -------------------------------------------------------------------------

	ath, err := auth.New(auth.Config{
		Log:       dbTest.Log,
		DB:        dbTest.DB,
		KeyLookup: &apitest.KeyStore{},
	})
	if err != nil {
		t.Fatal(err)
	}

	// -------------------------------------------------------------------------

	authService, err := authsrv.NewService(dbTest.Log, dbTest.DB, ath)
	if err != nil {
		t.Fatalf("Auth service init error: %s", err)
	}
	et.MockService("auth", authService)

	salesService, err := sales.NewService(dbTest.Log, dbTest.DB)
	if err != nil {
		t.Fatalf("Sales service init error: %s", err)
	}
	et.MockService("sales", salesService, et.RunMiddleware(true))

	// -------------------------------------------------------------------------

	authHandler := func(ctx context.Context, ap *mid.AuthParams) (eauth.UID, *auth.Claims, error) {
		return mid.BearerBasic(ctx, ath, dbTest.BusDomain.User, ap)
	}

	appTest := apitest.New(dbTest, ath, authHandler)

	return appTest
}