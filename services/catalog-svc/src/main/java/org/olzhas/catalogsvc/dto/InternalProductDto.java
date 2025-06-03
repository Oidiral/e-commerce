package org.olzhas.catalogsvc.dto;

import lombok.AllArgsConstructor;
import lombok.Data;

import java.math.BigDecimal;
import java.util.UUID;

@Data
@AllArgsConstructor
public class InternalProductDto {
    private UUID id;
    private String name;
    private BigDecimal price;
    private int availableQuantity;
}
