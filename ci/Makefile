set-pipeline:
	fly -t kurtmc set-pipeline -c pipeline.yml -p "$$(basename "$$(dirname "$$(pwd)")")"
destroy-pipeline:
	fly -t kurtmc destroy-pipeline -p "$$(basename "$$(dirname "$$(pwd)")")"
