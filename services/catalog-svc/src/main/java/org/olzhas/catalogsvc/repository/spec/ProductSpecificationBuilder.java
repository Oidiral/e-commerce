package org.olzhas.catalogsvc.repository.spec;

import jakarta.persistence.criteria.Predicate;
import org.olzhas.catalogsvc.dto.ProductFilter;
import org.olzhas.catalogsvc.model.Product;
import org.springframework.data.jpa.domain.Specification;

import java.util.ArrayList;
import java.util.List;

public class ProductSpecificationBuilder {
    public static Specification<Product> build(ProductFilter filter) {
        return ((root, query, criteriaBuilder) -> {
            List<Predicate> predicates = new ArrayList<>();
            if (filter.getCategoryId() != null) {
                predicates.add(criteriaBuilder.equal(root.get("category").get("id"), filter.getCategoryId()));
            }
            if (filter.getSku() != null){
                predicates.add(criteriaBuilder.equal(root.get("sku"), filter.getSku()));
            }
            if (filter.getMinPrice() != null) {
                predicates.add(criteriaBuilder.greaterThanOrEqualTo(root.get("price"), filter.getMinPrice()));
            }
            if (filter.getMaxPrice() != null) {
                predicates.add(criteriaBuilder.lessThanOrEqualTo(root.get("price"), filter.getMaxPrice()));
            }
            if (filter.getType() != null && !filter.getType().isEmpty()) {
                predicates.add(criteriaBuilder.equal(root.get("type"), filter.getType()));
            }
            if (filter.getAvailable() != null) {
                predicates.add(criteriaBuilder.equal(root.get("inStock"), filter.getAvailable()));
            }

            return criteriaBuilder.and(predicates.toArray(new Predicate[0]));
        });
    }
}

