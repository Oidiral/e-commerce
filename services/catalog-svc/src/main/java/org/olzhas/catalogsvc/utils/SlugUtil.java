package org.olzhas.catalogsvc.utils;

import java.text.Normalizer;
import java.util.Locale;
import java.util.regex.Pattern;

public class SlugUtil {

    private static final Pattern WHITESPACE = Pattern.compile("[\\s]");

    private static final Pattern NON_LATIN = Pattern.compile("[^\\w-]");


    private SlugUtil() {
        throw new UnsupportedOperationException("Utility class");
    }

    public static String toSlug(String input) {
        if (input == null) {
            return "";
        }

        String withHyphens = WHITESPACE.matcher(input).replaceAll("-");

        String normalized = Normalizer.normalize(withHyphens, Normalizer.Form.NFD);

        String slug = NON_LATIN.matcher(normalized).replaceAll("");

        slug = slug.replaceAll("[-]{2,}", "-");

        if (slug.startsWith("-")) {
            slug = slug.substring(1);
        }
        if (slug.endsWith("-")) {
            slug = slug.substring(0, slug.length() - 1);
        }

        return slug.toLowerCase(Locale.ENGLISH);
    }

}
