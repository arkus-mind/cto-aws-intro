deploy:
	echo "building jalacrypto"
	cd jalacrypto;go build 
	echo deploying jalacrypto
	cd jalacrypto;zip jalacrypto.zip jalacrypto; 
	echo "building tuneacrypto"
	cd tuneacrypto; go build; zip tuneacrypto.zip tuneacrypto
	cd apicrypto; go build; zip apicrypto.zip apicrypto
	echo "deploying infraestructure"
	cd btc-pooling-infra; cdk deploy;
