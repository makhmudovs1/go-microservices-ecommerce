ALTER TABLE cart_item ADD CONSTRAINT cart_item_cart_id_sku_unique UNIQUE (cart_id, sku);
