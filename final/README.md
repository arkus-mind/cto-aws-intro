# Cloud Computing, AWS y Go
### Actividad Final

1.- Obtener del api [bitso](https://bitso.com/api_info?shell#ticker) la conversion de criptomoneda [1 BTC a MXN Pesos](https://api.bitso.com/v3/ticker/?book=btc_mxn)   **lo mas cercano al tiempo real** para despues insertar esos datos a dynamodb.

2.- Insertar a una tabla de dynamodb y transformar los valores a RDS Postgresql haciendo que se cree una columna para saber su equivalencia a dolar hongkoné invocando una lambda function. DYNAMODB (STREAMS). LAMBDA FUNCTION

3.- Crear un api que liste los valores de la criptomoneda BTC de una tabla de sql, el api debe permitir filtrar resultados entre dos fechas, que te permita filtrar en dolar usd o/y dolar hongkones con el resultado pagino . API GATEWAY, LAMBDA FUNCTION, RDS POSTGRESQL


## Condiciones
- El codigo de las lambdas debe estar en Go y seguir todas las practicas ya propuestas de los ejercicios anteriores

## Extra
- Usar AWS CDK
- Las lambda function dockerizadas
- Codigo de go con unit test


![image](https://user-images.githubusercontent.com/7213379/152082871-f87a2da3-8f95-401f-bda7-410fd42cee47.png)

