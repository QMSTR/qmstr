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

    }

    post {
        success {
            archiveArtifacts artifacts: 'out/*', fingerprint: true
        }
    }

}
