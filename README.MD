# Install Docker
[DockerDownloadURL](https://docs.docker.com/engine/install/)

Windows Docker Desktop install URL

[DockerDesktopForWindowsDownloadURL](https://docs.docker.com/desktop/setup/install/windows-install/)

# Build the Docker image
In the same directory as Dockerfile, run the following commands:

```sh
docker build -t blockchain-app .
```
# Run the container

```sh
docker run -p 8080:8080 blockchain-app
```
# Installing GO

If you haven’t already installed Go, download it from Go Downloads URL shown below.
[GODownloadURL](https://go.dev/doc/install)

# Run the Go application
In the same directory as main.go, run the following commands:

```sh
go run main.go
```

This will start the server on http://localhost:8080.

In order to access the Endpoint open your browser or use curl to check the response:
Once the Blockchain Application is hosted in ECS Fargate (through terraform) it will accessible publicly from the public IP address. As per the commands mentioned below replace the localhost with the public IP of the ECS Fargate.

Open any browser:

```sh
GetBlockNumber: http://localhost:8080/blockNumber
GetBlockByNumber: http://localhost:8080/block/<blockNumber>
```
Open terminal for CLI command:
```sh
curl -i -X GET "http://localhost:8080/blockNumber"
curl -i -X GET "http://localhost:8080/block/<blockNumber>"
```
# Unit Test Cases

```sh
go test
```

handler_test.go has the unit test cases written

# Install Terraform

```sh
sudo yum install -y yum-utils shadow-utils
sudo yum-config-manager --add-repo https://rpm.releases.hashicorp.com/AmazonLinux/hashicorp.repo
sudo yum -y install terraform
```
# Intialise Terrafrom
Update the AWS IAM Access Key and Scecret Key in the main.tf file along with the AWS region.
Before terraform apply do push the created docket image to either DockerHub or ECR repository in AWS, please get the image URL copied and update in the main.tf file for image name in container_definitions section for ECS Task Definition.

```sh
terraform init
terraform plan
terraform apply
```