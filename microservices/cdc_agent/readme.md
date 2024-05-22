#Docker
From parent folder
docker build -t static-cdc -f ./cdc_agent/Dockerfile .
docker run static-cdc
docker tag static-cdc antrad1978/static-cdc:latest
docker push antrad1978/static-cdc:latest