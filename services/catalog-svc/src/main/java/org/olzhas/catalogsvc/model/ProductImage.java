package org.olzhas.catalogsvc.model;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;
import org.hibernate.annotations.ColumnDefault;
import org.hibernate.annotations.OnDelete;
import org.hibernate.annotations.OnDeleteAction;

import java.time.Instant;
import java.util.UUID;

@Getter
@Setter
@Entity
@Table(name = "product_image")
public class ProductImage {
    @Id
    @ColumnDefault("uuid_generate_v4()")
    @Column(name = "id", nullable = false)
    private UUID id;

    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @OnDelete(action = OnDeleteAction.CASCADE)
    @JoinColumn(name = "product_id", nullable = false)
    private Product product;

    @Column(name = "s3_key", nullable = false, length = 512)
    private String s3Key;

    @Column(name = "url", nullable = false, length = 512)
    private String url;

    @ColumnDefault("false")
    @Column(name = "is_primary", nullable = false)
    private Boolean isPrimary = false;

    @ColumnDefault("now()")
    @Column(name = "created_at")
    private Instant createdAt;

}