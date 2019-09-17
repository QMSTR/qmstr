package org.qmstr.gradle.buildcache;

import java.io.IOException;

import org.gradle.caching.BuildCacheEntryReader;
import org.gradle.caching.BuildCacheEntryWriter;
import org.gradle.caching.BuildCacheException;
import org.gradle.caching.BuildCacheKey;
import org.gradle.caching.BuildCacheService;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class QMSTRBuildCacheWrapperService implements BuildCacheService {

    private final BuildCacheService wrappedSrv;
    private static final Logger logger = LoggerFactory.getLogger(QMSTRBuildCacheWrapperService.class);

    public QMSTRBuildCacheWrapperService(BuildCacheService wrappedSrv) {
        this.wrappedSrv = wrappedSrv;
    }

    @Override
    public void close() throws IOException {
        this.wrappedSrv.close();
    }

    @Override
    public boolean load(BuildCacheKey arg0, BuildCacheEntryReader arg1) throws BuildCacheException {
        logger.warn("Loading from cache\n%s - %s", arg0.getDisplayName(), arg0.getHashCode());
        return this.wrappedSrv.load(arg0, arg1);
    }

    @Override
    public void store(BuildCacheKey arg0, BuildCacheEntryWriter arg1) throws BuildCacheException {
        logger.warn("Storing to cache\n%s - %s", arg0.getDisplayName(), arg0.getHashCode());
        this.wrappedSrv.store(arg0, arg1);
    }

}