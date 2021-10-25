pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        sh '''docker rm back &&
docker image build -t backend-main .'''
      }
    }

    stage('Run') {
      steps {
        sh 'docker run -it --pid=host -p 8008:8080 --name back backend-main'
      }
    }

  }
}