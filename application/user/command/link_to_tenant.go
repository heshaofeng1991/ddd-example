package command

//
// import (
// 	"context"
//
// 	domainOrder "github.com/heshaofeng1991/ddd-johnny/domain/order"
// 	domainPlatformProduct "github.com/heshaofeng1991/ddd-johnny/domain/platformProduct"
// 	domainStore "github.com/heshaofeng1991/ddd-johnny/domain/store"
// 	"github.com/heshaofeng1991/ddd-johnny/internal"
// 	"github.com/pkg/errors"
// )
//
// type LinkToTenantHandler struct {
// 	storeRepo           domainStore.Repository
// 	orderRepo           domainOrder.Repository
// 	platformProductRepo domainPlatformProduct.Repository
// }
//
// type LinkToTenantCommand struct {
// 	TenantID  int64
// 	StoreCode string
// }
//
// func NewLinkToTenantHandler(storeRepo domainStore.Repository,
// 	orderRepo domainOrder.Repository,
// 	platformProductRepo domainPlatformProduct.Repository,
// ) LinkToTenantHandler {
// 	var handler LinkToTenantHandler
//
// 	if storeRepo == nil || orderRepo == nil || platformProductRepo == nil {
// 		return handler
// 	}
//
// 	return LinkToTenantHandler{
// 		storeRepo:           storeRepo,
// 		orderRepo:           orderRepo,
// 		platformProductRepo: platformProductRepo,
// 	}
// }
//
// func (h LinkToTenantHandler) Handle(ctx context.Context, cmd LinkToTenantCommand) error {
// 	if cmd.StoreCode == "" {
// 		return errors.New("store code is empty")
// 	}
//
// 	store, err := h.storeRepo.GetByStoreCode(ctx, cmd.StoreCode)
// 	if err != nil || store == nil {
// 		return errors.New("store not found")
// 	}
//
// 	if store != nil && store.Edges.Tenant != nil && store.Edges.Tenant.ID == internal.SystemTenantID {
// 		if _, err := h.storeRepo.SetStoreTenantID(ctx, store, cmd.TenantID); err != nil {
// 			return errors.New("failed to set store tenant id")
// 		}
//
// 		if err := h.platformProductRepo.SetPlatformProductTenantID(ctx, store.ID, cmd.TenantID); err != nil {
// 			return errors.New("failed to set platform product tenant id")
// 		}
//
// 		if err := h.orderRepo.SetOrderTenantID(ctx, store.ID, cmd.TenantID); err != nil {
// 			return errors.New("failed to set order tenant id")
// 		}
// 	}
//
// 	return nil
// }
