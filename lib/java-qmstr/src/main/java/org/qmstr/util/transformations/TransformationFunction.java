package org.qmstr.util.transformations;

@FunctionalInterface
public interface TransformationFunction<T, R> {
   R apply(T t) throws TransformationException;
}