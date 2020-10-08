allow_k8s_contexts('kind-incluster-leader-election')

# allow_k8s_contexts('abatilo.cloud')
# default_registry("ghcr.io/abatilo")

load('ext://restart_process', 'docker_build_with_restart')

docker_build_with_restart('incluster-leader-election', '.',
  dockerfile='./build/Dockerfile',
  entrypoint='go run cmd/main.go',
  target='builder',
  live_update=[
    sync('./cmd', '/app/cmd'),
  ],
)

k8s_yaml('./deployments/local/leader-election.yml')
k8s_resource('leader-election', port_forwards=['8080'])
