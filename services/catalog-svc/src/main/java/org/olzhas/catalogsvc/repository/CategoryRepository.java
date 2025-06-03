package org.olzhas.catalogsvc.repository;

import org.olzhas.catalogsvc.model.Category;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

public interface CategoryRepository extends JpaRepository<Category, UUID> {
}