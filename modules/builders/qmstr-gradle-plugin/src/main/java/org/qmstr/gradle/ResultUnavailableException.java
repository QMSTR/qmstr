package org.qmstr.gradle;

public class ResultUnavailableException extends Exception {
    public ResultUnavailableException() {
        super();
    }

    public ResultUnavailableException(String message) {
        super(message);
    }

    public ResultUnavailableException(String message, Throwable cause) {
        super(message, cause);
    }

    public ResultUnavailableException(Throwable cause) {
        super(cause);
    }

}