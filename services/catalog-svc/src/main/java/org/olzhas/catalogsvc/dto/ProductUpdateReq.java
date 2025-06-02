package org.olzhas.catalogsvc.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.util.UUID;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class ProductUpdateReq {
    String name;
    String description;
    BigDecimal price;
    Integer quantity;
    UUID categoryId;
}
