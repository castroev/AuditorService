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
                    sh "docker image tag $repository/${service}:latest tylerorg/${service}:latest"
                    sh "docker image tag $repository/${service}:latest tylerorg/${service}:${tag}"
                    sh "docker push tylerorg/${service}:${epoch}_${commit}_${BUILD_NUMBER}"
                    sh "docker push tylerorg/${service}:latest"
                }
                echo 'Dockerhub push complete'
            }
        }
        stage('Push Images To Artifactory') {
            options {
                timeout(time: 5, unit: 'MINUTES')
            }
            steps {
                echo 'Artifactory push starting'
                dir("${service}") {
                    sh "docker image tag $repository/${service}:latest $repository/${service}:${tag}"
                    echo "Connecting to registry: ${registry} and logging into ${repository}"
                    sh "docker login ${repository}"
                    sh "docker push ${repository}/${service}:${tag}"
                    sh "docker push ${repository}/${service}:latest"
                }
                echo 'Artifactory push complete'
            }
        }
        stage('Deploy To Kubernetes') {
            steps {
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
        stage('Cleanup Images') {
            steps {
                echo 'Removing built docker images'
                sh "docker rmi $repository/${service}:${tag}"
                sh "docker rmi $repository/${service}:latest"
                sh "docker image prune -f"
            }
        }
    }
}
