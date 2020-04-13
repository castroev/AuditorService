// Organization: Tyler Technologies
// Author: Tyler Cloud Platform Team
pipeline {
    agent any
    options {
        timeout(time: 60, unit: 'MINUTES')
        ansiColor('xterm')
    }

    environment {
        artifactory_repository = "tylertech-corpdev-docker-local.jfrog.io"
        repository = "tylerorg"
        registry = "https://tylertech.jfrog.io"
        registryCredential = "artifactory"
        service = "tcp-auditor-go"
        bootstrap = "tcp-auditor-go-bootstrapper"
        k8_deployment = ""
        epoch = """${sh(
                        returnStdout: true,
                        script: "date +%s | tr -d '\n'"
                    )}"""
        commit = """${sh(
                        returnStdout: true,
                        script: "git log -n 1 --pretty=format:'%H'"
                    )}"""
        tag = "${epoch}_${commit}_${BUILD_NUMBER}"
        deploymentConfig = "auditorDeploy.yaml"
        ciBootstrapperConfig = "ci-bootstrapper-job.yaml"
        qaBootstrapperConfig = "qa-bootstrapper-job.yaml"
        prodBootstrapperConfig = "prod-bootstrapper-job.yaml"
    }

    stages {
        stage('Build Images (Service and Config Bootstrap)') {
            steps {
                echo 'Remember - access credentials are configured on connecting the Jenkins application to the Repo target.'
                echo 'Copy git files to build agent..'
                checkout scm
                echo 'Copy git files complete'
                echo "Docker build project service starting: ${artifactory_repository}/${service}:${tag}"
                echo 'TODO: Should standardize project repos to include a Dockerfile at SAME LEVEL as the JenkinsFile!'
                dir("${service}") {
                  sh "docker build . -t ${artifactory_repository}/${service}:latest"
                }
                echo "Docker build configuration bootstrap starting: ${artifactory_repository}/${bootstrap}:${tag}"
                dir("${bootstrap}") {
                  sh "docker build . -t ${artifactory_repository}/${bootstrap}:latest"
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
            when {
                // if the branch being built is the master branch run this stage, otherwise skip
                branch 'master'
            }
            steps {
                echo 'Pushing to Dockerhub registry....'
                withDockerRegistry([ credentialsId: "dockerhub", url: "" ]) {
                    sh "docker image tag ${artifactory_repository}/${service}:latest ${repository}/${service}:latest"
                    sh "docker push ${repository}/${service}:latest"
                    sh "docker image tag ${repository}/${service}:latest ${repository}/${service}:${tag}"
                    sh "docker push ${repository}/${service}:${tag}"
                }
                echo 'Pushing bootstrapper to Dockerhub registry....'
                withDockerRegistry([ credentialsId: "dockerhub", url: "" ]) {
                    sh "docker image tag ${artifactory_repository}/${bootstrap}:latest ${repository}/${bootstrap}:latest"
                    sh "docker push ${repository}/${bootstrap}:latest"
                    sh "docker image tag ${repository}/${bootstrap}:latest ${repository}/${bootstrap}:${tag}"
                    sh "docker push ${repository}/${bootstrap}:${tag}"
                }
                echo 'Dockerhub push complete'
            }
        }
        stage('Push Images To Artifactory') {
            options {
                timeout(time: 15, unit: 'MINUTES')
            }
            steps {
                echo 'Artifactory push starting'
                dir("${service}") {
                    echo "Connecting to registry: ${registry} and logging into ${artifactory_repository}"
                    sh "docker login ${artifactory_repository}"
                    sh "docker image tag ${artifactory_repository}/${service}:latest ${artifactory_repository}/${service}:${tag}"
                    sh "docker push ${artifactory_repository}/${service}:latest"
                    sh "docker push ${artifactory_repository}/${service}:${tag}"

                    sh "docker image tag ${artifactory_repository}/${bootstrap}:latest ${artifactory_repository}/${bootstrap}:${tag}"
                    sh "docker push ${artifactory_repository}/${bootstrap}:latest"
                    sh "docker push ${artifactory_repository}/${bootstrap}:${tag}"
                }
                echo 'Artifactory push complete'
            }
        }
        stage('Deploy To Kubernetes: TCPCI') {
            steps {
              kubernetesDeploy(configs: "${ciBootstrapperConfig}",
                                 kubeConfig: [path: ''],
                                 kubeconfigId: 'TCP-CI-Cluster',
                                 enableConfigSubstitution: true,
                                 secretName: 'tylerartifactory',
                                 secretNamespace: 'default',
                                 ssh: [sshCredentialsId: '*', sshServer: ''],
                                 textCredentials: [certificateAuthorityData: '', clientCertificateData: '', clientKeyData: '', serverUrl: 'https://'])
              kubernetesDeploy(configs: "${deploymentConfig}",
                                kubeConfig: [path: ''],
                                kubeconfigId: 'TCP-CI-Cluster',
                                enableConfigSubstitution: true,
                                secretName: 'tylerartifactory',
                                secretNamespace: 'default',
                                ssh: [sshCredentialsId: '*', sshServer: ''],
                                textCredentials: [certificateAuthorityData: '', clientCertificateData: '', clientKeyData: '', serverUrl: 'https://'])
            }
        }
        stage('Cleanup Artifactory Tagged Images') {
            steps {
                echo 'Removing built docker images'
                sh "docker rmi ${artifactory_repository}/${service}:${tag}"
                sh "docker rmi ${artifactory_repository}/${service}:latest"
                sh "docker rmi ${artifactory_repository}/${bootstrap}:${tag}"
                sh "docker rmi ${artifactory_repository}/${bootstrap}:latest"
                sh "docker image prune -f"
            }
        }
        stage('Cleanup Dockerhub Tagged Images') {
            when {
                // Cleanup images that were only created because the branch built was master
                branch 'master'
            }
            steps {
                echo 'Removing built docker images'
                sh "docker rmi ${repository}/${service}:${tag}"
                sh "docker rmi ${repository}/${service}:latest"
                sh "docker rmi ${repository}/${bootstrap}:${tag}"
                sh "docker rmi ${repository}/${bootstrap}:latest"
                sh "docker image prune -f"
            }
        }
    }
}
