package org.olzhas.catalogsvc.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;
import java.util.UUID;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class ProductImageDto implements Serializable {
    UUID id;
    String url;
    Boolean isPrimary;
}