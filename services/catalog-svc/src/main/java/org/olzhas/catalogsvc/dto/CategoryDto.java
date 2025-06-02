package org.olzhas.catalogsvc.dto;

import lombok.Builder;
import lombok.Data;

import java.util.UUID;

@Data
@Builder
public class CategoryDto {
    UUID id;
    String name;
    String slug;
}