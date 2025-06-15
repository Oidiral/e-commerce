package org.olzhas.catalogsvc.mapper;

import org.mapstruct.*;
import org.olzhas.catalogsvc.dto.CategoryCreateReq;
import org.olzhas.catalogsvc.dto.CategoryDto;
import org.olzhas.catalogsvc.dto.CategoryUpdateReq;
import org.olzhas.catalogsvc.model.Category;

@Mapper(unmappedTargetPolicy = ReportingPolicy.IGNORE, componentModel = MappingConstants.ComponentModel.SPRING)
public interface CategoryMapper {
    Category toEntity(CategoryCreateReq categoryCreateReq);

    CategoryDto toDto(Category category);



    @BeanMapping(nullValuePropertyMappingStrategy = NullValuePropertyMappingStrategy.IGNORE)
    Category partialUpdate(CategoryUpdateReq categoryCreateReq, @MappingTarget Category category);
}