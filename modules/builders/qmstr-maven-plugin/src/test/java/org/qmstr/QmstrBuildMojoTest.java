package org.qmstr;


import org.apache.maven.plugin.testing.MojoRule;
import org.apache.maven.plugin.testing.WithoutMojo;

import org.junit.Rule;
import static org.junit.Assert.*;
import org.junit.Test;
import java.io.File;

public class QmstrBuildMojoTest
{
    @Rule
    public MojoRule rule = new MojoRule()
    {
        @Override
        protected void before() throws Throwable 
        {
        }

        @Override
        protected void after()
        {
        }
    };

    /**
     * @throws Exception if any
     */
    @Test
    public void testSomething()
            throws Exception
    {
        File pom = new File( "target/test-classes/project-to-test/" );
        assertNotNull( pom );
        assertTrue( pom.exists() );

        QmstrBuildMojo myMojo = ( QmstrBuildMojo ) rule.lookupConfiguredMojo( pom, "qmstrbuild" );
        assertNotNull( myMojo );
        myMojo.execute();

        String qmstrAddress = ( String ) rule.getVariableValueFromObject( myMojo, "qmstrAddress" );
        assertNotNull( qmstrAddress );
        assertTrue( qmstrAddress.equals("localhost") );
    }

    /** Do not need the MojoRule. */
    @WithoutMojo
    @Test
    public void testSomethingWhichDoesNotNeedTheMojoAndProbablyShouldBeExtractedIntoANewClassOfItsOwn()
    {
        assertTrue( true );
    }

}

