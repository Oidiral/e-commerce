package org.olzhas.catalogsvc.utils;

import org.olzhas.catalogsvc.dto.PageResponse;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Sort;

import java.util.List;
import java.util.stream.Collectors;


public final class PageConverter {

    private PageConverter() {
        throw new UnsupportedOperationException("Utility class");
    }

    public static <DTO, E> PageResponse<DTO> toPageResponse(List<DTO> dtos, Page<E> page) {
        PageResponse<DTO> response = new PageResponse<>();

        response.setContent(dtos);

        PageResponse.PageMetadata meta = new PageResponse.PageMetadata();
        meta.setNumber(page.getNumber());
        meta.setSize(page.getSize());
        meta.setTotalElements(page.getTotalElements());
        meta.setTotalPages(page.getTotalPages());
        meta.setSort(formatSort(page.getSort()));

        response.setPage(meta);

        return response;
    }

    private static String formatSort(Sort sort) {
        if (sort == null || sort.isUnsorted()) {
            return "";
        }
        return sort.stream()
                .map(order -> order.getProperty() + "," + order.getDirection())
                .collect(Collectors.joining(";"));
    }
}
