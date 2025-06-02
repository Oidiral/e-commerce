package org.olzhas.catalogsvc.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;
import jakarta.validation.constraints.PositiveOrZero;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.util.UUID;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class ProductCreateReq {
    @NotBlank
    String sku;
    @NotBlank String name;
    String description;
    @NotNull
    @Positive
    BigDecimal price;
    @PositiveOrZero
    Integer quantity;
    UUID categoryId;
}