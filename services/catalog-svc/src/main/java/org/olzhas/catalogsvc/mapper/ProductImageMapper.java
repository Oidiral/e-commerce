package org.olzhas.catalogsvc.mapper;

import org.mapstruct.*;
import org.olzhas.catalogsvc.dto.ProductImageDto;
import org.olzhas.catalogsvc.model.ProductImage;

@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE, componentModel = MappingConstants.ComponentModel.SPRING)
public interface ProductImageMapper {
    ProductImage toEntity(ProductImageDto productImageDto);

    ProductImageDto toDto(ProductImage productImage);

    @BeanMapping(nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.IGNORE)
    ProductImage partialUpdate(ProductImageDto productImageDto, @MappingTarget ProductImage productImage);
}