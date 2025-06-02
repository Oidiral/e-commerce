package org.olzhas.catalogsvc.dto;

import lombok.Builder;
import lombok.Data;

import java.math.BigDecimal;
import java.time.Instant;

@Data
@Builder
public class PriceDto {
    BigDecimal amount;
    String currency;
    Instant createdAt;
}
