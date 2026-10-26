
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

        stage('Check Go Formatting') {
            steps {
                sh '''
                    if [ -n "$(gofmt -l .)" ]; then
                        echo "The following Go files need formatting:"
                        gofmt -l .
                        exit 1
                    fi

                    echo "All Go files are correctly formatted."
                '''
            }
        }
        stage('Build Todo API') {
            steps {
                sh 'go build -o todo-api ./cmd/server'
            }
        }

	    stage('Archive Build Artifact') {
            steps {
                archiveArtifacts artifacts: 'todo-api', fingerprint: true
            }
        }

        stage('Build Docker Image') {
            steps {
                sh '''
                    docker build -t todo-api:1.0 .
                '''
            }
        }

        stage('Deploy with Docker Compose') {
            steps {
                sh '''
                    cd "$HOME/go/src/github.com/todo"

                    docker compose up -d --no-build api
                '''
            }
        }

        stage('Smoke Test') {
            steps {
                sh '''
                    set -eu

                    for attempt in $(seq 1 15); do
                        if curl --fail --silent http://localhost:8080/info; then
                            echo
                            echo "Todo API smoke test passed."
                            exit 0
                        fi

                        echo "Waiting for Todo API... attempt $attempt/15"
                        sleep 2
                    done

                    echo "Todo API smoke test failed."
                    exit 1
                '''
            }
        }
    }

    post {
        success {
            echo 'Todo API CI/CD pipeline completed successfully!'
        }

        failure {
            echo 'Todo API pipeline failed. Check the stage logs.'
        }

        always {
            echo 'Jenkins pipeline execution finished.'
        }
    }
}
