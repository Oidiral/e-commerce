package org.olzhas.catalogsvc.dto;

import jakarta.validation.constraints.Min;
import lombok.Data;
import org.springframework.format.annotation.NumberFormat;

import java.math.BigDecimal;

@Data
public class ProductFilter {
    private Long categoryId;

    @NumberFormat(style = NumberFormat.Style.CURRENCY)
    @Min(value = 0, message = "Min price must be greater than or equal to 0")
    private BigDecimal minPrice;

    @NumberFormat(style = NumberFormat.Style.CURRENCY)
    @Min(value = 0, message = "Max price must be greater than or equal to 0")
    private BigDecimal maxPrice;

    private String sku;

    private String name;

    private Boolean available;
}
