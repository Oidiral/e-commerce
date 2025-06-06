package org.olzhas.catalogsvc.dto;

import lombok.*;

import java.math.BigDecimal;
import java.time.Instant;
import java.util.UUID;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
public class ProductPriceResponseDto {
    private UUID id;
    private UUID productId;
    private BigDecimal amount;
    private String currency;
    private Instant createdAt;


}