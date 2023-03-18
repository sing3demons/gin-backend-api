pipeline {

    agent {
        label 'main'
    }

    environment {
        image = "sing3demons/gin-backend-apt"
        registry = "docker.io"
    }

    stages {

        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Print Environment') {
            steps {
                sh('ls -al')
                sh('printenv')
            }
        }

        // stage('Ansible prepareations docker ') {
        //     steps{
        //         sh 'ANSIBLE_ROLES_PATH="$PWD/ansible-script/roles" ansible-playbook -vvv ./ansible-script/playbook/web-server/web-server.yml -i ./ansible-script/host -u root -e "state=prepareation tagnumber=${BUILD_NUMBER}"'
        //     }
        // }
        
        stage('Build docker image') {
            steps {
                script {
                    docker.withRegistry('', 'dockerhub') {
                        def slackImage = docker.build("${env.image}:${BUILD_NUMBER}")
                        slackImage.push()
                        slackImage.push('latest')
                    }
                }
            }
        }

        stage('Deployment'){
            steps {
                sh "docker-compose up -d"
            }
            
        }

        // stage('tag docker image') {
        //     steps {
        //        sh "docker tag ${env.image}:${BUILD_NUMBER} ${env.image}:latest"
        //     }
        // }

        // stage('push docker image') {
        //     steps {
        //        sh "docker push ${env.image}:latest"
        //     }
        // }

        // stage('Verify new docker image(s)') {
        //     steps {
        //         sh('docker images')
        //     }
        // }
        
    }
}