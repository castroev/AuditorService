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
                git 'https://github.com/castroev/AuditorService.git'
                echo '(1) COMPLETE'

                sh "docker build . -t $registry:$BUILD_NUMBER"
                sh "docker build . -t $registry:latest"
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
                echo 'TODO: Interation and unit tests go here!'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
                echo '(1) Pushing to registry....'
                withDockerRegistry([ credentialsId: "dockerhub", url: "" ]) {
                    sh "docker push $registry:$BUILD_NUMBER"
                    sh "docker push $registry:latest"
                }
                echo '(1) COMPLETE'
            }
        }
    }
}
