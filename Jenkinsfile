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
        )
    ]
) {
    node('event-exporter') {
        stage('Checkout') {
            checkout scm
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
        stage('Deploy') {
            sh('echo NotImplemented')
        }
    }
}
