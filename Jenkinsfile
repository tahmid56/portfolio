pipeline {
    agent any

    options {
        timestamps()
        disableConcurrentBuilds()
    }

    environment {
        // Docker Hub repo to push to, e.g. "tahmid56/portfolio"
        IMAGE_NAME = 'tahmidakter56/portfolio'
        // ID of the "Username with password" credential added in Jenkins
        DOCKERHUB_CREDENTIALS = credentials('dockerhub-credentials')
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build Docker Image') {
            steps {
                sh """
                    docker build \
                      -t ${IMAGE_NAME}:${BUILD_NUMBER} \
                      -t ${IMAGE_NAME}:latest \
                      .
                """
            }
        }

        stage('Login to Docker Hub') {
            steps {
                sh '''
                    echo "$DOCKERHUB_CREDENTIALS_PSW" | docker login -u "$DOCKERHUB_CREDENTIALS_USR" --password-stdin
                '''
            }
        }

        stage('Push Image') {
            steps {
                sh """
                    docker push ${IMAGE_NAME}:${BUILD_NUMBER}
                    docker push ${IMAGE_NAME}:latest
                """
            }
        }
    }

    post {
        always {
            sh 'docker logout || true'
            // Keep the workspace/local images tidy on the Jenkins host
            sh "docker rmi ${IMAGE_NAME}:${BUILD_NUMBER} ${IMAGE_NAME}:latest || true"
        }
        success {
            echo "Pushed ${IMAGE_NAME}:${BUILD_NUMBER} and ${IMAGE_NAME}:latest to Docker Hub"
        }
        failure {
            echo "Build failed — check the console output above."
        }
    }
}
