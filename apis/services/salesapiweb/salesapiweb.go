// Package salesapiweb represent the encore application.
package salesapiweb

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"encore.dev"
	"encore.dev/rlog"
	"github.com/ardanlabs/conf/v3"
	homeapp "github.com/ardanlabs/encore/app/core/crud/homeapp"
	"github.com/ardanlabs/encore/app/core/crud/productapp"
	"github.com/ardanlabs/encore/app/core/crud/tranapp"
	"github.com/ardanlabs/encore/app/core/crud/userapp"
	"github.com/ardanlabs/encore/app/core/views/vproductapp"
	"github.com/ardanlabs/encore/business/api/auth"
	"github.com/ardanlabs/encore/business/api/debug"
	"github.com/ardanlabs/encore/business/api/metrics"
	"github.com/ardanlabs/encore/business/core/crud/delegate"
	"github.com/ardanlabs/encore/business/core/crud/homebus"
	"github.com/ardanlabs/encore/business/core/crud/homebus/stores/homedb"
	"github.com/ardanlabs/encore/business/core/crud/productbus"
	"github.com/ardanlabs/encore/business/core/crud/productbus/stores/productdb"
	"github.com/ardanlabs/encore/business/core/crud/userbus"
	"github.com/ardanlabs/encore/business/core/crud/userbus/stores/userdb"
	"github.com/ardanlabs/encore/business/core/views/vproductbus"
	"github.com/ardanlabs/encore/business/core/views/vproductbus/stores/vproductdb"
	"github.com/ardanlabs/encore/business/data/appdb"
	"github.com/ardanlabs/encore/business/data/appdb/migrate"
	"github.com/ardanlabs/encore/business/data/sqldb"
	"github.com/ardanlabs/encore/foundation/keystore"
	"github.com/jmoiron/sqlx"
)

// Service represents the encore service application.
//
//encore:service
type Service struct {
	mtrcs   *metrics.Values
	db      *sqlx.DB
	auth    *auth.Auth
	appCrud appCrud
	appView appView
	busCrud busCrud
	debug   http.Handler
}

// NewService is called to create a new encore Service.
func NewService(db *sqlx.DB, ath *auth.Auth) (*Service, error) {
	delegate := delegate.New()
	userBus := userbus.NewCore(delegate, userdb.NewStore(db))
	productBus := productbus.NewCore(userBus, delegate, productdb.NewStore(db))
	homeBus := homebus.NewCore(userBus, delegate, homedb.NewStore(db))
	vproductBus := vproductbus.NewCore(vproductdb.NewStore(db))

	s := Service{
		mtrcs: newMetrics(),
		db:    db,
		auth:  ath,
		appCrud: appCrud{
			user:    userapp.NewCore(userBus, ath),
			product: productapp.NewCore(productBus),
			home:    homeapp.NewCore(homeBus),
			tran:    tranapp.NewCore(userBus, productBus),
		},
		appView: appView{
			product: vproductapp.NewCore(vproductBus),
		},
		busCrud: busCrud{
			user:    userBus,
			product: productBus,
			home:    homeBus,
		},
		debug: debug.Mux(),
	}

	return &s, nil
}

// Shutdown implements a function that will be called by encore when the service
// is signaled to shutdown.
func (s *Service) Shutdown(force context.Context) {
	defer rlog.Info("shutdown", "status", "shutdown complete")

	rlog.Info("shutdown", "status", "stopping database support")
	s.db.Close()
}

// =============================================================================

// initService is called by Encore to initialize the service.
//
//lint:ignore U1000 "called by encore"
func initService() (*Service, error) {
	db, auth, err := startup()
	if err != nil {
		return nil, err
	}

	return NewService(db, auth)
}

func startup() (*sqlx.DB, *auth.Auth, error) {

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	rlog.Info("initService", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	cfg := struct {
		conf.Version
		Auth struct {
			KeysFolder string `conf:"default:zarf/keys/"`
			ActiveKID  string `conf:"default:54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"`
			Issuer     string `conf:"default:service project"`
		}
		DB struct {
			MaxIdleConns int `conf:"default:2"`
			MaxOpenConns int `conf:"default:0"`
		}
	}{
		Version: conf.Version{
			Build: encore.Meta().Environment.Name,
			Desc:  "Service Project",
		},
	}

	const prefix = "SALES"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil, nil, err
		}
		return nil, nil, fmt.Errorf("parsing config: %w", err)
	}

	// -------------------------------------------------------------------------
	// App Starting

	rlog.Info("initService", "environment", encore.Meta().Environment.Name)

	out, err := conf.String(&cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("generating config for output: %w", err)
	}
	rlog.Info("initService", "config", out)

	// -------------------------------------------------------------------------
	// Database Support

	rlog.Info("initService", "status", "initializing database support")

	db, err := sqldb.Open(sqldb.Config{
		EDB:          appdb.AppDB,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("connecting to db: %w", err)
	}

	// TODO: I don't like this here because it's more of an ops thing, but
	// for now I will leave it as I learn more.
	if err := migrate.Seed(context.Background(), db); err != nil {
		return nil, nil, fmt.Errorf("seeding db: %w", err)
	}

	// -------------------------------------------------------------------------
	// Auth Support

	rlog.Info("initService", "status", "initializing authentication support")

	// Load the private keys files from disk. We can assume some system like
	// Vault has created these files already. How that happens is not our
	// concern.

	ks := keystore.New()
	if err := ks.LoadRSAKeys(os.DirFS(cfg.Auth.KeysFolder)); err != nil {
		return nil, nil, fmt.Errorf("reading keys: %w", err)
	}

	authCfg := auth.Config{
		DB:        db,
		KeyLookup: ks,
	}

	auth, err := auth.New(authCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("constructing auth: %w", err)
	}

	return db, auth, nil
}