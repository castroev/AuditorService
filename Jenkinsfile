pipeline {
    agent any
    options {
        timeout(time: 60, unit: 'MINUTES')
        ansiColor('xterm')
    }

    environment {
        repository = "tylertech-corpdev-docker-local.jfrog.io"
        registry = "https://tylertech.jfrog.io"
        registryCredential = "artifactory"
        service = "tcp-auditor-go"
        epoch = """${sh(
                        returnStdout: true,
                        script: "date +%s | tr -d '\n'"
                    )}"""
        commit = """${sh(
                        returnStdout: true,
                        script: "git log -n 1 --pretty=format:'%H'"
                    )}"""
    }

    stages {
        stage('Build') {
            steps {
                echo 'Building..'
                echo 'Remember - access credentials are configured on connectinig the Jenkins application to the Repo target.'
                echo 'Copy git files to build agent..'
                checkout scm
                echo 'Copy git files complete'
                echo "Docker build starting: ${repository}/${service}:${epoch}_${commit}_${BUILD_NUMBER}"
                echo 'TODO: Should standardize project repos to include a Dockerfile at SAME LEVEL as the JenkinsFile!'
                dir("${service}") {
                  sh "docker build . -t $repository/${service}:${epoch}_${commit}_${BUILD_NUMBER}"
                  sh "docker build . -t $repository/${service}:latest"
                }
                echo 'Docker build complete'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
                echo 'TODO: Interaction and unit tests go here!'
            }
        }
        stage('Push Images to Dockerhub') {
            steps {
                echo 'Pushing to Dockerhub registry....'
                withDockerRegistry([ credentialsId: "dockerhub", url: "" ]) {
                    sh "docker push $registry:$BUILD_NUMBER"
                    sh "docker push $registry:latest"
                }
                echo 'Dockerhub push complete'
            }
        }
        stage('Push Images To Artifactory') {
            steps {
                echo 'Artifactory push starting'
                dir("${service}") {
                    echo "Connecting to: ${registry}"
                    script {
                        def rtServer = Artifactory.server "TylerArtifactory"
                        def rtDocker = Artifactory.docker server: rtServer

                        def buildInfo = rtDocker.push "${repository}/${service}:${epoch}_${commit}_${BUILD_NUMBER}", "${repository}"
                        def buildInfoLatest = rtDocker.push "${repository}/${service}:latest", "${repository}"
                    }
                }
                echo 'Artifactory push complete'
            }
        }
        stage('Cleanup Images') {
            steps {
                echo 'Removing built docker images'
                sh "docker rmi $repository/${service}:${epoch}_${commit}_${BUILD_NUMBER}"
                sh "docker rmi $repository/${service}:latest"
            }
        }
    }
}
