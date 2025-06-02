package org.olzhas.catalogsvc.dto;

import lombok.Builder;
import lombok.Data;

import java.math.BigDecimal;
import java.util.List;
import java.util.UUID;

@Data
@Builder
public class ProductDto {
    UUID id;
    String sku;
    String name;
    String description;
    BigDecimal price;
    String currency;
    Integer quantity;
    List<String> images;
    List<CategoryDto> categories;
}