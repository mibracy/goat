pipeline {
    agent any

    parameters {
        // NOTE: The 'when' directive is not valid here.
        // The Extended Choice Parameter plugin is required for multiChoice.
        choice(name: 'ARTIFACTS_TO_BUILD', choices: ['Web', 'Batch', 'Database'], description: 'Select artifact type(s) to build')
        choice(name: 'WEB_PROJECTS', choices: ['Frontend', 'Backend', 'Messaging'], description: 'Select Web Project(s) to build')
        choice(name: 'BATCH_JOBS', choices: ['Copy', 'Archive', 'Purge'], description: 'Select Batch Job(s) to build')
        string(name: 'RELEASE_NUMBER', defaultValue: '1.0', description: 'Enter the release number')
    }

    environment {
        // Use a relative path for cross-platform compatibility
        TARGET_DIR = 'target'
        GRADLE_HOME = tool 'Gradle 7.6.2'
        GITHUB_REPOS = [
            'WebFrontend': 'your-organization/webfrontend',
            'WebBackend': 'your-organization/webbackend',
            'WebMessaging': 'your-organization/webmessaging',
            'BatchCopy': 'your-organization/batchcopy',
            'BatchArchive': 'your-organization/batarchive',
            'BatchPurge': 'your-organization/batchpurge',
            'Database': 'your-organization/database',
            'Deploy': 'your-organization/deploy'
        ]
    }

    stages {
        stage('Checkout & Build') {
            steps {
                script {
                    // Always clone the deploy repo
                    def deployRepoUrl = GITHUB_REPOS['Deploy']
                    // FIX: 'git clone' is a shell command and must be run with 'bat' on Windows.
                    // The 'git' step is for declarative checkouts, not arbitrary commands.
                    bat "git clone https://github.com/${deployRepoUrl}.git --branch release/${RELEASE_NUMBER}"

                    // Build selected Web projects
                    if (params.ARTIFACTS_TO_BUILD.contains('Web') && params.WEB_PROJECTS) {
                        params.WEB_PROJECTS.split(',').each { project ->
                            def repoName = "Web${project}"
                            def repoUrl = GITHUB_REPOS[repoName]
                            echo "Building ${repoName} from ${repoUrl}"
                            // FIX: There is no 'gradle' pipeline step. It must be called via a shell step.
                            bat "gradle clean build --github-repo ${repoUrl} --branch release/${RELEASE_NUMBER}"
                        }
                    }

                    // Build selected Batch jobs
                    if (params.ARTIFACTS_TO_BUILD.contains('Batch') && params.BATCH_JOBS) {
                        params.BATCH_JOBS.split(',').each { job ->
                            def repoName = "Batch${job}"
                            def repoUrl = GITHUB_REPOS[repoName]
                            echo "Building ${repoName} from ${repoUrl}"
                            bat "gradle clean build --github-repo ${repoUrl} --branch release/${RELEASE_NUMBER}"
                        }
                    }

                    // Clone Database repo
                    if (params.ARTIFACTS_TO_BUILD.contains('Database')) {
                        def repoUrl = GITHUB_REPOS['Database']
                        echo "Cloning Database repo"
                        bat "git clone https://github.com/${repoUrl}.git --branch release/${RELEASE_NUMBER}"
                    }
                }
            }
        }

        stage('Package Artifacts') {
            steps {
                script {
                    // This stage will collect all build outputs into staging folders and then zip them.
                    if (params.ARTIFACTS_TO_BUILD.contains('Web') && params.WEB_PROJECTS) {
                        def stagingDir = "${TARGET_DIR}\\Web"
                        bat "if not exist ${stagingDir} mkdir ${stagingDir}"
                        
                        params.WEB_PROJECTS.split(',').each { project ->
                            def artifactDir = "Web${project}"
                            def sourcePath = "${artifactDir}\\build\\libs\\*.*" // Copy all build artifacts
                            // FIX: 'mv' is a Linux command. The Windows equivalent is 'move' or 'copy', run with 'bat'.
                            echo "Copying from ${sourcePath} to ${stagingDir}"
                            bat "copy ${sourcePath} ${stagingDir}"
                        }
                        
                        def zipFileName = "Web-${RELEASE_NUMBER}.zip"
                        echo "Archiving ${stagingDir} to ${zipFileName}"
                        powershell "Compress-Archive -Path '${stagingDir}\\*' -DestinationPath '${zipFileName}' -Force"
                    }

                    if (params.ARTIFACTS_TO_BUILD.contains('Batch') && params.BATCH_JOBS) {
                        def stagingDir = "${TARGET_DIR}\\Batch"
                        bat "if not exist ${stagingDir} mkdir ${stagingDir}"

                        params.BATCH_JOBS.split(',').each { job ->
                            def artifactDir = "Batch${job}"
                            def sourcePath = "${artifactDir}\\build\\libs\\*.*"
                            echo "Copying from ${sourcePath} to ${stagingDir}"
                            bat "copy ${sourcePath} ${stagingDir}"
                        }

                        def zipFileName = "Batch-${RELEASE_NUMBER}.zip"
                        echo "Archiving ${stagingDir} to ${zipFileName}"
                        // FIX: Removed redundant 'zip' step. 'powershell' is sufficient.
                        powershell "Compress-Archive -Path '${stagingDir}\\*' -DestinationPath '${zipFileName}' -Force"
                    }
                }
            }
        }

        stage('Upload') {
            steps {
                // This will archive any zip file created in the previous stage.
                archiveArtifacts artifacts: '*.zip', followSymlinks: false
                echo "Artifacts ending in .zip have been archived by Jenkins."
                // For Artifactory, you would use a command like this:
                // bat "jf rt upload *.zip your-repo/"
            }
        }
    }
}

