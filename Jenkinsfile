pipeline {
    agent any
    stages {
        stage('Clone'){
            steps {
                git 'https://github.com/dinhcanh303/go-microservices.git'
            }
        }
    }
    post {
        always {
            mail bcc: '', body: 'Please check server ?', cc: 'dinhcanhng303@gmail.com', from: '', replyTo: '', subject: 'Jenkins Clone Code', to: 'rapidstars97@gmail.com'
        }
    }
}