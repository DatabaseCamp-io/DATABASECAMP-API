pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        sh '''cker image build -t backend-main .'''
      }
    }

    stage('Run') {
      steps {
        sh 'docker run -it --pid=host -p 8008:8001 --name back backend-main'
      }
    }

  }
}