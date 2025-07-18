package org.olzhas.catalogsvc.repository;

import org.olzhas.catalogsvc.model.Category;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface CategoryRepository extends JpaRepository<Category, UUID> {
    boolean existsBySlug(String candidate);

    void deleteBySlug(String slug);
}