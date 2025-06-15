package org.olzhas.catalogsvc.dto;

import jakarta.validation.constraints.Size;

public class ProductPatchReq {
    @Size(max = 64)
    private String sku;
    private String name;
    private String description;
}