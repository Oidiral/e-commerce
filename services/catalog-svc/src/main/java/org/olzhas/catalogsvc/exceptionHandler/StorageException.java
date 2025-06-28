package org.olzhas.catalogsvc.exceptionHandler;

public class StorageException extends RuntimeException {
    public StorageException(String message, Throwable cause) {
        super(message, cause);
    }
}
