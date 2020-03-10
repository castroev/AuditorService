pipeline {
    agent any
    options {
        timeout(time: 60, unit: 'MINUTES')
        ansiColor('xterm')
    }

    environment {
        registry = "joenunnelley/docker-test"
        registryCredential = "dockerhub"
    }

    stages {
        stage('Build') {
            steps {
                echo 'Building..'
                echo 'Remember - access credentials are configured on connectinig the Jenkins application to the Repo target.'
                echo '(1) Copy git files to build agent..'
                checkout scm
                echo '(1) COMPLETE'

                echo '(2) Doocker build'
                echo 'TODO: Should standardize project repos to include Dockerfile at SAME LEVEL as JenksnisFile!'
                dir("tcp-auditor-go") {
                  sh "ls"
                  sh "docker build . -t $registry:$BUILD_NUMBER"
                  sh "docker build . -t $registry:latest"
                }
                echo '(2) COMPLETE'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
                echo 'TODO: Interation and unit tests go here!'
            }
        }
        stage('Push To Artifactory') {
            steps {
                echo 'Pushing to Artifactory....'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
                echo '(1) Pushing to registry....'
                dir("tcp-auditor-go") {
                  withDockerRegistry([ credentialsId: "dockerhub", url: "" ]) {
                    sh "docker push $registry:$BUILD_NUMBER"
                    sh "docker push $registry:latest"
                  }
                }
                echo '(1) COMPLETE'
            }
        }
    }
}
