// Package interfaces provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.1 DO NOT EDIT.
package interfaces

const (
	Jwt_header_AuthorizationScopes = "jwt_header_Authorization.Scopes"
)

// BatchQueryOrderInfo defines model for BatchQueryOrderInfo.
type BatchQueryOrderInfo struct {
	// list data
	List []QueryInboundData `json:"list"`

	// MetaData describes the MetaData
	Meta MetaData `json:"meta"`
}

// BatchQueryOrderRsp defines model for BatchQueryOrderRsp.
type BatchQueryOrderRsp struct {
	// code
	Code int64                `json:"code"`
	Data *BatchQueryOrderInfo `json:"data,omitempty"`

	// message
	Message string `json:"message"`
}

// CreateInboundData defines model for CreateInboundData.
type CreateInboundData struct {
	// inbound order number (入库单序列号)
	OrderNumber string `json:"order_number"`
}

// CreateInboundItem defines model for CreateInboundItem.
type CreateInboundItem struct {
	// product id (产品ID)
	ProductId int32 `json:"product_id"`

	// product qty (产品数量)
	ProductQty int32 `json:"product_qty"`
}

// CreateInboundOrderRequestBody defines model for CreateInboundOrderRequestBody.
type CreateInboundOrderRequestBody struct {
	// description (描述）
	Description *string `json:"description,omitempty"`

	// estimated arrival time (预估到达时间)
	EstimatedArrivalAt *string `json:"estimated_arrival_at,omitempty"`

	// items
	Items []CreateInboundItem `json:"items"`

	// ship time (发货时间)
	ShipAt *string `json:"ship_at,omitempty"`

	// tracking number (运单号)
	TrackingNumber *string `json:"tracking_number,omitempty"`

	// warehouse id (仓库ID)
	WarehouseId int32 `json:"warehouse_id"`
}

// CreateInboundRsp defines model for CreateInboundRsp.
type CreateInboundRsp struct {
	// code
	Code int64              `json:"code"`
	Data *CreateInboundData `json:"data,omitempty"`

	// message
	Message string `json:"message"`
}

// Error defines model for Error.
type Error struct {
	// Is the error a server-side fault?
	Fault bool `json:"fault"`

	// ID is a unique identifier for this particular occurrence of the problem.
	Id string `json:"id"`

	// Message is a human-readable explanation specific to this occurrence of the problem.
	Message string `json:"message"`

	// Name is the name of this class of errors.
	Name string `json:"name"`

	// Is the error temporary?
	Temporary bool `json:"temporary"`

	// Is the error a timeout?
	Timeout bool `json:"timeout"`
}

// MetaData describes the MetaData
type MetaData struct {
	// current
	Current int `json:"current"`

	// page_size
	PageSize int `json:"page_size"`

	// total
	Total int64 `json:"total"`
}

// QueryInboundData defines model for QueryInboundData.
type QueryInboundData struct {
	// created time (创建时间)
	CreatedAt string `json:"created_at"`

	// description (描述)
	Description string `json:"description"`

	// estimated arrival time (预估到达时间)
	EstimatedArrivalAt string `json:"estimated_arrival_at"`

	// inbound order id (入库单ID)
	Id int32 `json:"id"`

	// inbound order items (入库单items)
	Items []QueryInboundItem `json:"items"`

	// inbound order number (入库单序列号)
	OrderNumber string `json:"order_number"`

	// ship time (发货时间)
	ShipAt string `json:"ship_at"`

	// inbound order status (入库单状态)
	Status int32 `json:"status"`

	// tracking number (运单号)
	TrackingNumber string `json:"tracking_number"`

	// warehouse name (仓库名称)
	WarehouseName string `json:"warehouse_name"`
}

// QueryInboundItem defines model for QueryInboundItem.
type QueryInboundItem struct {
	// product barcode (产品Barcode)
	ProductBarcode string `json:"product_barcode"`

	// product id (产品ID)
	ProductName string `json:"product_name"`

	// product qty (产品数量)
	ProductQty int32 `json:"product_qty"`

	// product sku (产品SKu)
	ProductSku string `json:"product_sku"`
}

// QueryInboundOrderRequestBody defines model for QueryInboundOrderRequestBody.
type QueryInboundOrderRequestBody struct {
	// current
	Current *int `json:"current,omitempty"`

	// page_size
	PageSize *int `json:"page_size,omitempty"`
}

// QueryInboundRsp defines model for QueryInboundRsp.
type QueryInboundRsp struct {
	// code
	Code int64             `json:"code"`
	Data *QueryInboundData `json:"data,omitempty"`

	// message
	Message string `json:"message"`
}

// InboundOrderBatchQueryInboundOrderParams defines parameters for InboundOrderBatchQueryInboundOrder.
type InboundOrderBatchQueryInboundOrderParams struct {
	// order number
	OrderNumbers *[]string `json:"order_numbers,omitempty"`

	// status
	Status *int `json:"status,omitempty"`

	// current
	Current *int `json:"current,omitempty"`

	// page_size
	PageSize *int `json:"page_size,omitempty"`
}

// InboundOrderCreateInboundOrderJSONBody defines parameters for InboundOrderCreateInboundOrder.
type InboundOrderCreateInboundOrderJSONBody CreateInboundOrderRequestBody

// InboundOrderQueryInboundOrderJSONBody defines parameters for InboundOrderQueryInboundOrder.
type InboundOrderQueryInboundOrderJSONBody QueryInboundOrderRequestBody

// InboundOrderCreateInboundOrderJSONRequestBody defines body for InboundOrderCreateInboundOrder for application/json ContentType.
type InboundOrderCreateInboundOrderJSONRequestBody InboundOrderCreateInboundOrderJSONBody

// InboundOrderQueryInboundOrderJSONRequestBody defines body for InboundOrderQueryInboundOrder for application/json ContentType.
type InboundOrderQueryInboundOrderJSONRequestBody InboundOrderQueryInboundOrderJSONBody
