package org.olzhas.catalogsvc.mapper;

import org.mapstruct.*;
import org.olzhas.catalogsvc.dto.ProductCreateReq;
import org.olzhas.catalogsvc.dto.ProductDto;
import org.olzhas.catalogsvc.dto.ProductPatchReq;
import org.olzhas.catalogsvc.dto.ProductUpdateReq;
import org.olzhas.catalogsvc.model.Product;

@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE, componentModel = MappingConstants.ComponentModel.SPRING)
public interface ProductMapper {
    Product toEntity(ProductDto productDto);

    Product toEntity(ProductCreateReq productCreateReq);

    @Mapping(target = "id", ignore = true)
    @Mapping(target = "createdAt", ignore = true)
    @Mapping(target = "updatedAt", ignore = true)
    Product toEntity(ProductUpdateReq request);

    ProductDto toDto(Product product);


    @BeanMapping(nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.IGNORE)
    Product partialUpdate(ProductPatchReq request, @MappingTarget Product product);


    @BeanMapping(nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.IGNORE)
    Product partialUpdate(ProductDto productDto, @MappingTarget Product product);
}