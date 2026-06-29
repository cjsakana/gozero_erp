package types

const (
	CacheErpSaleSalesDeliveryIdPrefix = "cache:erpSale:salesDelivery:id:%d"

	CacheErpSaleSalesDeliveryDetailIdPrefix = "cache:erpSale:salesDeliveryDetail:id:%d"

	// CacheErpSaleSalesDeliveryDetailIdsByDeliveryId 出库单明细缓存，用list
	CacheErpSaleSalesDeliveryDetailIdsByDeliveryId = "cache:erpSale:salesDeliveryDetail:ids:deliveryId:%d"

	CacheErpSaleSalesOrderIdPrefix = "cache:erpSale:salesOrder:id:%v"

	CacheErpSaleSalesOrderDetailIdPrefix = "cache:erpSale:salesOrderDetail:id:%v"

	CacheErpSaleSalesOrderDetailIdsByOrderId = "cache:erpSale:salesOrderDetail:ids:orderId:%d"
)
