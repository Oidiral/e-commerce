package org.olzhas.catalogsvc.service;

import org.olzhas.catalogsvc.dto.PriceCreateReq;
import org.olzhas.catalogsvc.dto.ProductPriceResponseDto;

import java.util.UUID;

public interface PriceService {

    ProductPriceResponseDto addPrice(UUID productId, PriceCreateReq req);
}
