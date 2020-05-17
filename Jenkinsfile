#!/usr/bin/env groovy
pipeline {

  environment {
    BUILD_SCRIPTS_GIT="https://github.com/haidlir/golang-prtg-api-wrapper.git"
    XDG_CACHE_HOME="/tmp/.cache"
  }

  agent any

  stages {
      stage('Unit-Test') {
          agent { docker { image 'golang' } }
          steps {
              cleanWs(disableDeferredWipeout: true, deleteDirs: true)
              sh 'cd ${WORKSPACE}'
              sh 'git clone $BUILD_SCRIPTS_GIT .'
              // Test commit message for flags
              sh(script: "ls -l", returnStatus: true)
              echo 'Get Dependency'
              sh './scripts/getdep'
              echo 'Unit Testing..'
              sh 'go test ./prtg-api/... -covermode=count -coverprofile=cover.out'
          }
      }
  }

  post {
    always {
      cleanWs(disableDeferredWipeout: true, deleteDirs: true)
    }
  }
}
