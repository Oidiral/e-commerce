package org.olzhas.catalogsvc.service;


import java.util.UUID;

public interface InventoryService {

    void setQuantity(UUID productId, int qty);
}
