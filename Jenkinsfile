pipeline {

    agent {
        docker { image 'endocode/qmstr_buildenv:latest' }
    }

    environment {
        MASTER_CONTAINER_NAME="qmstr-demo-master_${BUILD_NUMBER}"
    }

    stages {

        stage('Build') {
            steps {
                sh "make clients"
            }
        }

        stage('Test') {
            steps {
                sh "make gotest"
            }
        }

        stage('compile curl'){
            steps{
                sh 'cp out/qmstr* /usr/local/bin'
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
