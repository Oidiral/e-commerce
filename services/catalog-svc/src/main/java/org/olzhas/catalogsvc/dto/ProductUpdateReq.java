package org.olzhas.catalogsvc.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.util.UUID;

@Data
@NoArgsConstructor
@Builder
@AllArgsConstructor
public class ProductUpdateReq {
    @NotBlank
    @Size(max = 64)
    private String sku;
    @NotBlank
    private String name;
    private String description;
}
