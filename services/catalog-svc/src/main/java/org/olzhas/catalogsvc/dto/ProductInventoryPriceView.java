package org.olzhas.catalogsvc.dto;

import java.math.BigDecimal;
import java.util.UUID;

public interface ProductInventoryPriceView {
    UUID getProductId();
    String getName();
    Integer getQuantity();
    BigDecimal getLatestPrice();
}
