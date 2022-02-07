import { Stack, StackProps,Duration } from 'aws-cdk-lib';
import * as dynamo from 'aws-cdk-lib/aws-dynamodb';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as events from 'aws-cdk-lib/aws-events';
import * as targets from 'aws-cdk-lib/aws-events-targets';
import * as lambdas from 'aws-cdk-lib/aws-lambda';
import * as lamdaEvents from 'aws-cdk-lib/aws-lambda-event-sources';
import { Construct } from 'constructs';

export class BtcPoolingInfraStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);
     
     const  defaultVpc = ec2.Vpc.fromLookup(this,"vpc",{vpcId:"vpc-0e14f401eae69127f"})
     const secGroup = ec2.SecurityGroup.fromLookupByName(this,"secgroup","othergroup",defaultVpc)
     let table = dynamo.Table.fromTableName(this,'LeTable','LeMemeCoins')
     if (!table){
       table = new dynamo.Table(this, 'LeTable', {
         partitionKey: { name: 'id', type: dynamo.AttributeType.STRING },
         tableName:  "LeMemeCoins",
         billingMode: dynamo.BillingMode.PAY_PER_REQUEST,
         stream : dynamo.StreamViewType.NEW_AND_OLD_IMAGES
       });
     }else {
       console.log("dynamo table already exists won't be created")
     }

     const jalacrypto = new lambdas.Function(this, 'jalacrypto', {
       runtime: lambdas.Runtime.GO_1_X,
       handler: 'jalacrypto', 
       code : lambdas.Code.fromAsset("../jalacrypto/jalacrypto.zip"),
       functionName: "jalacryptoCDKV"
     })
     const tuneacrypto = new lambdas.Function(this, 'tuneacrypto', {
       runtime: lambdas.Runtime.GO_1_X,
       handler: 'tuneacrypto', 
       code : lambdas.Code.fromAsset("../tuneacrypto/tuneacrypto.zip"),
       functionName: "tuneacryptoCDKV",
       securityGroups: [secGroup],
       //timeout:Duration.minutes(3),
       vpc : defaultVpc,
     })
     table.grantReadWriteData(jalacrypto)
     const eventRule = new events.Rule(this, 'eachMinuteCheckForShekels', {
      schedule: events.Schedule.cron({}),
     });
     eventRule.addTarget(new targets.LambdaFunction(jalacrypto))
     
     tuneacrypto.addEventSource(new lamdaEvents.DynamoEventSource(table,{
       startingPosition : lambdas.StartingPosition.LATEST
     }))

  }
}
