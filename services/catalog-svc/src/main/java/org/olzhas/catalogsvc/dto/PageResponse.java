package org.olzhas.catalogsvc.dto;

import lombok.Data;

import java.util.List;

@Data
public class PageResponse<T> {
    List<T> content;
    PageMetadata page;

    @Data
    public static class PageMetadata {
        int number;
        int size;
        long totalElements;
        int totalPages;
        String sort;
    }
}