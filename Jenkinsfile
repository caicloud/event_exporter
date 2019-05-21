podTemplate(
    // 之前配置的 Kubernetes Cloud Provider
    cloud: 'dev-cluster',
    // 这个 pipeline 执行环境名称
    name: 'event-exporter',
    // 运行在带有 always-golang 标签的 Jenkins Slave 上 
    label: 'event-exporter',
    containers: [
        // Kubernetes Pod 的配置, 这个 Pod 包含两个容器
        containerTemplate(
            name: 'jnlp',
            alwaysPullImage: true,
            // Jenkins Slave ， 与 Master 通信进程
            image: 'cargo.caicloud.io/circle/jnlp:2.62',
            command: '',
            args: '${computer.jnlpmac} ${computer.name}',
            resourceRequestCpu: '300m',
            resourceLimitCpu: '5000m',
            resourceRequestMemory: '300Mi',
            resourceLimitMemory: '500Mi',
        ),
        containerTemplate(
            name: 'golang',
            // Jenkins Slave 作业执行环境， 此处为一个 Docker in Docker 环境，用于跑作业
            image: 'cargo.caicloud.io/caicloud/golang:1.7',
            ttyEnabled: true,
            command: '',
            args: '',
            resourceRequestCpu: '1000m',
            resourceLimitCpu: '2000m',
            resourceRequestMemory: '1000Mi',
            resourceLimitMemory: '2000Mi',
        ),
        containerTemplate(
            name: 'docker',
            image: 'docker:stable',
            ttyEnabled: true,
            command: 'cat',
        )
    ],
    volumes: [
        // Please update it with docker secret in k8s cluster
        secretVolume(secretName: 'dockerConfig', mountPath: '/mnt/docker-hub'),
        hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock')
    ]
) {
    def shortSHA = '',
    def dockerImage = 'cargo.caicloud.io/sysinfra/event-exporter'
    def gitBranch = ''
    def dockerTag = dockerImage + ':' + shortSHA
    node('event-exporter') {
        stage('Checkout') {
            varSCM = checkout scm
            shortSHA = varSCM.GIT_COMMIT.substring(0, 7)
            gitBranch = varSCM.GIT_BRANCH
        }
        stage('Test') {
            container('golang') {
                sh('mkdir -p /go/src/github.com/caicloud')
                sh('cp -r . /go/src/github.com/caicloud/event_exporter')
                sh('go test github.com/caicloud/event_exporter')
            }
        }
        stage('Build') {
            container('golang') {
                sh('cd /go/src/github.com/caicloud/event_exporter && make promu')
                // we need a docker client in golang env
                // sh('cd /go/src/github.com/caicloud/event_exporter && make docker')
                // sh('cd /go/src/github.com/caicloud/event_exporter && make push')
            }
        }
        stage('Docker Build and push') {
            container('docker') {
                sh("cp /mnt/docker-hub/.dockercfg ~/.dockercfg")
                sh("docker build -t ${dockerTag} -f Dockerfile .")
                sh("docker push ${dockerTag}")
                if (gitBranch == 'master') {
                        sh("docker tag ${dockerTag} ${dockerImage}:latest")
                        sh("docker push ${dockerImage}:latest")
                    }

            }
        }
        stage('Deploy') {
            sh('echo NotImplemented')
        }
    }
}
