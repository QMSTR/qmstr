pipeline {
    agent { label 'golang' }
    // fixme: this is wrong - the makefile requires golang for some reason. That should not be necessary when building in containers

    environment {
        MASTER_CONTAINER_NAME="qmstr-demo-master_${BUILD_NUMBER}"
    }

    stages {
        stage('Clean') {
            steps {
                sh 'git clean -f -d'
            }
        }
 
        stage('Build master and client images') {
            steps {
                script {
                    sh 'make democontainer'
                    def mastername = sh(script: 'docker create qmstr/master', returnStdout: true)
                    mastername = mastername.trim()
                    sh 'mkdir out'
                    sh "docker cp ${mastername}:/usr/local/bin/qmstr out/qmstr"
                    sh "docker cp ${mastername}:/usr/local/bin/qmstrctl out/qmstrctl"
                    sh "docker rm ${mastername}"
                }
            }
        }
    }

    post {
        success {
            archiveArtifacts artifacts: 'out/*', fingerprint: true
        }
    }

}
