
pipeline {
    agent any

    environment {
        GO111MODULE = 'on'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Check Go Version') {
            steps {
                sh 'go version'
            }
        }

        stage('Download Dependencies') {
            steps {
                sh 'go mod download'
            }
        }

        stage('Compile') {
            steps {
                sh 'go build ./...'
            }
        }

        stage('Unit Tests') {
            steps {
                sh 'go test ./...'
            }
        }

        stage('Static Analysis') {
            steps {
                sh 'go vet ./...'
            }
        }

        stage('Build Todo API') {
            steps {
                sh 'go build -o todo-api ./cmd/server'
            }
        }
    }

    post {
        success {
            echo 'Todo API CI pipeline completed successfully!'
        }

        failure {
            echo 'Todo API CI pipeline failed. Check the stage logs.'
        }

        always {
            echo 'Jenkins pipeline execution finished.'
        }
    }
}
