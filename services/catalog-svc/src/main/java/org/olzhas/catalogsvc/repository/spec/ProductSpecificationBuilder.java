package org.olzhas.catalogsvc.repository.spec;

import jakarta.persistence.criteria.Predicate;
import jakarta.persistence.criteria.Root;
import jakarta.persistence.criteria.Subquery;
import org.olzhas.catalogsvc.dto.ProductFilter;
import org.olzhas.catalogsvc.model.Product;
import org.olzhas.catalogsvc.model.ProductCategory;
import org.olzhas.catalogsvc.model.ProductInventory;
import org.olzhas.catalogsvc.model.ProductPrice;
import org.springframework.data.jpa.domain.Specification;

import java.util.ArrayList;
import java.util.List;

public class ProductSpecificationBuilder {
    public static Specification<Product> build(ProductFilter filter) {
        return ((root, query, criteriaBuilder) -> {
            List<Predicate> predicates = new ArrayList<>();

            if (filter.getSku() != null) {
                predicates.add(criteriaBuilder.equal(root.get("sku"), filter.getSku()));
            }

            if (filter.getName() != null && !filter.getName().isEmpty()) {
                predicates.add(criteriaBuilder.like(
                        criteriaBuilder.lower(root.get("name")),
                        "%" + filter.getName().toLowerCase() + "%"));
            }

            if (filter.getCategoryId() != null) {
                Subquery<java.util.UUID> sub = query.subquery(java.util.UUID.class);
                Root<ProductCategory> pc = sub.from(ProductCategory.class);
                sub.select(pc.get("id").get("productId"))
                        .where(criteriaBuilder.equal(pc.get("id").get("categoryId"), filter.getCategoryId()));
                predicates.add(root.get("id").in(sub));
            }

            if (filter.getMinPrice() != null || filter.getMaxPrice() != null) {
                Subquery<java.util.UUID> sub = query.subquery(java.util.UUID.class);
                Root<ProductPrice> price = sub.from(ProductPrice.class);
                List<Predicate> pricePreds = new ArrayList<>();
                pricePreds.add(criteriaBuilder.equal(price.get("product"), root));
                if (filter.getMinPrice() != null) {
                    pricePreds.add(criteriaBuilder.greaterThanOrEqualTo(price.get("amount"), filter.getMinPrice()));
                }
                if (filter.getMaxPrice() != null) {
                    pricePreds.add(criteriaBuilder.lessThanOrEqualTo(price.get("amount"), filter.getMaxPrice()));
                }
                sub.select(price.get("product").get("id")).where(pricePreds.toArray(new Predicate[0]));
                predicates.add(criteriaBuilder.exists(sub));
            }

            if (filter.getAvailable() != null) {
                Subquery<java.util.UUID> sub = query.subquery(java.util.UUID.class);
                Root<ProductInventory> inv = sub.from(ProductInventory.class);
                List<Predicate> invPreds = new ArrayList<>();
                invPreds.add(criteriaBuilder.equal(inv.get("product"), root));
                if (filter.getAvailable()) {
                    invPreds.add(criteriaBuilder.greaterThan(inv.get("quantity"), 0));
                } else {
                    invPreds.add(criteriaBuilder.equal(inv.get("quantity"), 0));
                }
                sub.select(inv.get("product").get("id")).where(invPreds.toArray(new Predicate[0]));
                predicates.add(criteriaBuilder.exists(sub));
            }

            return criteriaBuilder.and(predicates.toArray(new Predicate[0]));
        });
    }
}