pipeline {

    agent none

    environment {
        PATH = "/tmp:$PATH"
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
                archiveArtifacts artifacts: 'out/*', fingerprint: true
            }
            
        }

        stage('Compile with QMSTR'){

            parallel{

                stage('compile curl'){

                    agent { label 'docker' }

                    steps {
                        unstash 'executables'
                        sh 'export PATH=$PATH:$PWD/out/'
                        sh 'make container'
                        sh 'git submodule update --init'
                        sh 'echo $PATH'
                        sh "cd demos && make curl"
                       
                    }
                }

                stage('compile openssl'){

                    agent { label 'docker' }

                    steps {
                        unstash 'executables'
                        sh 'export PATH=$PATH:$PWD/out/'
                        sh 'make container'
                        sh 'git submodule update --init'
                        sh 'echo $PATH'
                        sh "cd demos && make openssl"
                    }
                }
            }
        }

    }

}
