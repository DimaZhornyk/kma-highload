## VPC
```
aws ec2 create-vpc --cidr-block 10.10.0.0/18 --no-amazon-provided-ipv6-cidr-block --tag-specifications 'ResourceType=vpc,Tags=[{Key=Name,Value=kma-genesis},{Key=Lesson,Value=public-clouds}]' --query Vpc.VpcId --output text
```
➜ export VPC_ID=<vpc_id>

---

### Subnets
```
aws ec2 create-subnet --vpc-id $VPC_ID --availability-zone us-east-1a --cidr-block 10.10.1.0/24
```
➜ export S1_ID=<subnet_id>
```
aws ec2 create-subnet --vpc-id $VPC_ID --availability-zone us-east-1b --cidr-block 10.10.2.0/24
```
➜ export S2_ID=<subnet_id>

```
aws ec2 create-subnet --vpc-id $VPC_ID --availability-zone us-east-1c --cidr-block 10.10.3.0/24
```
➜ export S3_ID=<subnet_id>

---

### IGW
```
aws ec2 create-internet-gateway --query InternetGateway.InternetGatewayId --output text
```
➜ export IG_ID=<ig_id>

```
aws ec2 attach-internet-gateway --internet-gateway-id $IG_ID --vpc-id $VPC_ID
```
---

### Security group
```
aws ec2 create-security-group --group-name kma-highload-sg --description "kma highload sh" --vpc-id $VPC_ID
```

➜ export SG_ID=<sg_id>

```
aws ec2 authorize-security-group-ingress --group-id $SG_ID --protocol tcp --port 22 --cidr  0.0.0.0/0
aws ec2 authorize-security-group-ingress --group-id $SG_ID --protocol tcp --port 80 --cidr  0.0.0.0/0
aws ec2 authorize-security-group-ingress --group-id $SG_ID --protocol tcp --port 443 --cidr  0.0.0.0/0
```

## Instances management

```
aws ec2 create-launch-template --launch-template-name KmaHighloadTemplate --version-description AutoScalingVersion1 --launch-template-data '{"NetworkInterfaces":[{"DeviceIndex":0,"AssociatePublicIpAddress":true,"Groups":["$SG_ID"],"DeleteOnTermination":true}],"ImageId":"ami-0ff8a91507f77f867","InstanceType":"t3.micro","TagSpecifications":[{"ResourceType":"instance","Tags":[{"Key":"Name","Value":"KmaASG"}]}],"BlockDeviceMappings":[{"DeviceName":"/dev/sda1","Ebs":{"VolumeSize":15}}]}' --region us-east-1

aws autoscaling create-auto-scaling-group --auto-scaling-group-name KmaHighloadASG --launch-template "LaunchTemplateName=KmaHighloadTemplate" --min-size 1 --max-size 2 --desired-capacity 1 --vpc-zone-identifier "$S1_ID,$2_ID,$S3_ID" --availability-zones "us-east-1a" "us-east-1b" "us-east-1c"
# ec2 instance will be created by ASG in a couple of minutes

aws elbv2 create-load-balancer --name kma-highload-lb  --subnets $S1_ID $S2_ID $S3_ID --security-groups $SG_ID
```