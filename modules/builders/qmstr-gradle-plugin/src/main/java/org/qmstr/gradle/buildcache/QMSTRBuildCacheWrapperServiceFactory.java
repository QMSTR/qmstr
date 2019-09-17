
package org.qmstr.gradle.buildcache;

import org.gradle.cache.CacheRepository;
import org.gradle.cache.internal.CacheScopeMapping;
import org.gradle.cache.internal.CleanupActionFactory;
import org.gradle.caching.BuildCacheService;
import org.gradle.caching.local.DirectoryBuildCache;
import org.gradle.caching.local.internal.DirectoryBuildCacheFileStoreFactory;
import org.gradle.caching.local.internal.DirectoryBuildCacheServiceFactory;
import org.gradle.internal.file.PathToFileResolver;
import org.gradle.internal.resource.local.FileAccessTimeJournal;

public class QMSTRBuildCacheWrapperServiceFactory extends DirectoryBuildCacheServiceFactory {

    public QMSTRBuildCacheWrapperServiceFactory(CacheRepository cacheRepository, CacheScopeMapping cacheScopeMapping,
            PathToFileResolver resolver, DirectoryBuildCacheFileStoreFactory fileStoreFactory,
            CleanupActionFactory cleanupActionFactory, FileAccessTimeJournal fileAccessTimeJournal) {
        super(cacheRepository, cacheScopeMapping, resolver, fileStoreFactory, cleanupActionFactory,
                fileAccessTimeJournal);
    }

    @Override
    public BuildCacheService createBuildCacheService(DirectoryBuildCache arg0, Describer arg1) {
        BuildCacheService wrappedSrv = super.createBuildCacheService(arg0, arg1);
        return new QMSTRBuildCacheWrapperService(wrappedSrv);
    }
}