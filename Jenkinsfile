#!groovy

ansiColor('xterm') {
    node('executor') {
        checkout scm

        def authorName = sh(returnStdout: true, script: 'git --no-pager show --format="%an" --no-patch')
        def isMain = env.BRANCH_NAME == "main"
        def serviceName = env.JOB_NAME.tokenize("/")[1]

        def commitHash = sh(returnStdout: true, script: 'git rev-parse HEAD | cut -c-7').trim()
        def imageTag = "${env.BUILD_NUMBER}-${commitHash}"

        try {
            stage("Run Tests") {
                sh "make test-ci"
            }
        } catch (e) {
            slackSend(color: '#b20000', message: "FAILED: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL}) by ${authorName}")
            throw e
        } finally {
            stage("Clean Up") {
                sh "IMAGE_TAG=${imageTag} make clean-ci"
            }
        }
        slackSend(color: '#006600', message: "SUCCESSFUL: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL}) by ${authorName}")
    }
}
