package salesapiweb

import (
	"context"
	"net/http"

	"encore.dev"
	"github.com/ardanlabs/encore/app/core/crud/homeapp"
	"github.com/ardanlabs/encore/app/core/crud/productapp"
	"github.com/ardanlabs/encore/app/core/crud/tranapp"
	"github.com/ardanlabs/encore/app/core/crud/userapp"
	"github.com/ardanlabs/encore/app/core/views/vproductapp"
	"github.com/ardanlabs/encore/business/api/page"
)

// Fallback is called for the debug enpoints.
//
//encore:api public raw path=/!fallback
func (s *Service) Fallback(w http.ResponseWriter, r *http.Request) {

	// If this is a web socket call for statsviz and we are in development.
	if r.URL.String() == "/debug/statsviz/ws" && encore.Meta().Environment.Type == encore.EnvDevelopment {

		// In development the r.Host will be host=127.0.0.1:RandPort while the
		// Origin will be origin=127.0.0.1:4000. These need to match.
		r.Header.Set("Origin", "http://"+r.Host)
	}

	s.debug.ServeHTTP(w, r)
}

// =============================================================================

//lint:ignore U1000 "called by encore"
//encore:api auth method=POST path=/v1/homes tag:metrics tag:authorize tag:as_user_role
func (s *Service) HomeCreate(ctx context.Context, app homeapp.NewHome) (homeapp.Home, error) {
	return s.appCrud.home.Create(ctx, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=PUT path=/v1/homes/:homeID tag:metrics tag:authorize_home
func (s *Service) HomeUpdate(ctx context.Context, homeID string, app homeapp.UpdateHome) (homeapp.Home, error) {
	return s.appCrud.home.Update(ctx, homeID, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=DELETE path=/v1/homes/:homeID tag:metrics tag:authorize_home
func (s *Service) HomeDelete(ctx context.Context, homeID string) error {
	return s.appCrud.home.Delete(ctx, homeID)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/homes tag:metrics tag:authorize tag:as_any_role
func (s *Service) HomeQuery(ctx context.Context, qp homeapp.QueryParams) (page.Document[homeapp.Home], error) {
	return s.appCrud.home.Query(ctx, qp)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/homes/:productID tag:metrics tag:authorize_home
func (s *Service) HomeQueryByID(ctx context.Context, productID string) (homeapp.Home, error) {
	return s.appCrud.home.QueryByID(ctx, productID)
}

// =============================================================================

//lint:ignore U1000 "called by encore"
//encore:api auth method=POST path=/v1/products tag:metrics tag:authorize tag:as_user_role
func (s *Service) ProductCreate(ctx context.Context, app productapp.NewProduct) (productapp.Product, error) {
	return s.appCrud.product.Create(ctx, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=PUT path=/v1/products/:productID tag:metrics tag:authorize_product
func (s *Service) ProductUpdate(ctx context.Context, productID string, app productapp.UpdateProduct) (productapp.Product, error) {
	return s.appCrud.product.Update(ctx, productID, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=DELETE path=/v1/products/:productID tag:metrics tag:authorize_product
func (s *Service) ProductDelete(ctx context.Context, productID string) error {
	return s.appCrud.product.Delete(ctx, productID)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/products tag:metrics tag:authorize tag:as_any_role
func (s *Service) ProductQuery(ctx context.Context, qp productapp.QueryParams) (page.Document[productapp.Product], error) {
	return s.appCrud.product.Query(ctx, qp)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/products/:productID tag:metrics tag:authorize_product
func (s *Service) ProductQueryByID(ctx context.Context, productID string) (productapp.Product, error) {
	return s.appCrud.product.QueryByID(ctx, productID)
}

// =============================================================================

//lint:ignore U1000 "called by encore"
//encore:api auth method=POST path=/v1/tran tag:metrics tag:authorize tag:as_user_role
func (s *Service) TranCreate(ctx context.Context, app tranapp.NewTran) (tranapp.Product, error) {
	return s.appCrud.tran.Create(ctx, app)
}

// =============================================================================

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/token/:kid tag:metrics
func (s *Service) UserToken(ctx context.Context, kid string) (userapp.Token, error) {
	return s.appCrud.user.Token(ctx, kid)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=POST path=/v1/users tag:metrics tag:authorize tag:as_admin_role
func (s *Service) UserCreate(ctx context.Context, app userapp.NewUser) (userapp.User, error) {
	return s.appCrud.user.Create(ctx, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=PUT path=/v1/users/:userID tag:metrics tag:authorize_user
func (s *Service) UserUpdate(ctx context.Context, userID string, app userapp.UpdateUser) (userapp.User, error) {
	return s.appCrud.user.Update(ctx, userID, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=PUT path=/v1/role/:userID tag:metrics tag:authorize_user tag:as_admin_role
func (s *Service) UserUpdateRole(ctx context.Context, userID string, app userapp.UpdateUserRole) (userapp.User, error) {
	return s.appCrud.user.UpdateRole(ctx, userID, app)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=DELETE path=/v1/users/:userID tag:metrics tag:authorize_user
func (s *Service) UserDelete(ctx context.Context, userID string) error {
	return s.appCrud.user.Delete(ctx, userID)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/users tag:metrics tag:authorize tag:as_admin_role
func (s *Service) UserQuery(ctx context.Context, qp userapp.QueryParams) (page.Document[userapp.User], error) {
	return s.appCrud.user.Query(ctx, qp)
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/users/:userID tag:metrics tag:authorize_user
func (s *Service) UserQueryByID(ctx context.Context, userID string) (userapp.User, error) {
	return s.appCrud.user.QueryByID(ctx, userID)
}

// =============================================================================

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/vproducts tag:metrics tag:authorize tag:as_admin_role
func (s *Service) VProductQuery(ctx context.Context, qp vproductapp.QueryParams) (page.Document[vproductapp.Product], error) {
	return s.appView.product.Query(ctx, qp)
}