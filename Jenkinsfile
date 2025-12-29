#!/usr/bin/env groovy

/**
 * goossify CI/CD Pipeline
 *
 * Builds and tests the Go OSS project boilerplate generator.
 */

pipeline {
    agent any

    tools {
        go 'Go 1.22'
    }

    environment {
        GO111MODULE = 'on'
        CGO_ENABLED = '0'
        GOFLAGS = '-mod=readonly'
    }

    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
        timestamps()
        timeout(time: 15, unit: 'MINUTES')
        disableConcurrentBuilds()
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
                echo "Building branch: ${GIT_BRANCH}"
            }
        }

        stage('Download Dependencies') {
            steps {
                sh 'go mod download'
                sh 'go mod verify'
            }
        }

        stage('Lint') {
            steps {
                sh '''
                    if command -v golangci-lint &> /dev/null; then
                        golangci-lint run ./...
                    else
                        echo "golangci-lint not installed, running go vet instead..."
                        go vet ./...
                    fi
                '''
            }
        }

        stage('Test') {
            steps {
                sh 'go test -v -race -cover -coverprofile=coverage.out ./...'
            }
            post {
                always {
                    script {
                        if (fileExists('coverage.out')) {
                            sh 'go tool cover -func=coverage.out'
                        }
                    }
                }
            }
        }

        stage('Build') {
            steps {
                sh '''
                    mkdir -p bin
                    go build -ldflags="-s -w" -o bin/goossify .
                '''
            }
        }

        stage('Verify Build') {
            steps {
                sh './bin/goossify --version || ./bin/goossify --help | head -5'
            }
        }

        stage('Archive Artifacts') {
            steps {
                archiveArtifacts artifacts: 'bin/*',
                                 fingerprint: true,
                                 allowEmptyArchive: true
            }
        }
    }

    post {
        success {
            echo 'Build completed successfully!'
        }
        failure {
            echo 'Build failed!'
        }
        always {
            cleanWs(
                deleteDirs: true,
                notFailBuild: true,
                patterns: [
                    [pattern: 'bin', type: 'INCLUDE'],
                    [pattern: 'coverage.out', type: 'INCLUDE']
                ]
            )
        }
    }
}
