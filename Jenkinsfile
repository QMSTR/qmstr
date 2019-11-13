pipeline {

    agent none

    environment {
        MASTER_CONTAINER_NAME="qmstr-demo-master_${BUILD_NUMBER}"
    }

    stages {

        stage('Build & Test') {
            agent {
                docker { image 'endocode/qmstr_buildenv:latest' }
            }
            steps {
                sh "make clients"
                sh "make gotest"
                stash includes: 'out/qmstr*', name: 'executables' 
            }
            
        }

        stage('compile curl'){
            agent { label 'docker' }
            
            steps{
                unstash 'executables'
                sh 'export PATH=$PATH:$PWD/out/'
                sh 'git submodule update --init'
                sh 'cd demos && make curl'
            }
        }

    }

    post {
        success {
            archiveArtifacts artifacts: 'out/*', fingerprint: true
        }
    }

}
