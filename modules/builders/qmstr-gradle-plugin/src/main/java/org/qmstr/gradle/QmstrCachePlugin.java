package org.qmstr.gradle;

import org.gradle.api.Plugin;
import org.gradle.api.initialization.Settings;
import org.gradle.caching.configuration.BuildCacheConfiguration;
import org.qmstr.gradle.buildcache.QMSTRBuildCacheWrapper;
import org.qmstr.gradle.buildcache.QMSTRBuildCacheWrapperServiceFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class QmstrCachePlugin implements Plugin<Settings> {
    private static final Logger logger = LoggerFactory.getLogger(QmstrCachePlugin.class);
  
    @Override
    public void apply(Settings settings) {
        logger.warn("Registering QMSTR BuildCache Wrapper");    
        BuildCacheConfiguration bcc = settings.getBuildCache();
        bcc.registerBuildCacheService(QMSTRBuildCacheWrapper.class, QMSTRBuildCacheWrapperServiceFactory.class);
    }
}