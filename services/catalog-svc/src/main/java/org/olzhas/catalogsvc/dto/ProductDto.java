package org.olzhas.catalogsvc.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.util.List;
import java.util.UUID;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
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