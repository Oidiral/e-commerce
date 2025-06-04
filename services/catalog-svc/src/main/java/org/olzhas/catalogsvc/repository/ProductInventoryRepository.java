package org.olzhas.catalogsvc.repository;

import org.olzhas.catalogsvc.model.ProductInventory;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

public interface ProductInventoryRepository extends JpaRepository<ProductInventory, UUID> {
}