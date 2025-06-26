# Trabajo Final - Programación Concurrente

El conjunto de datos utilizado en este proyecto está disponible públicamente y puede consultarse en:

https://www.datosabiertos.gob.pe/dataset/encuesta-nacional-de-hogares-enaho-2022-instituto-nacional-de-estadistica-e-informatica

Desde ese enlace se descarga un archivo comprimido `2022.zip`, el cual contiene diversas carpetas. Cada una de estas carpetas corresponde a un módulo diferente de la Encuesta Nacional de Hogares (ENAHO) 2022.

Para el presente trabajo nos interesa específicamente la carpeta **784-Modulo18**, la cual, de acuerdo a su diccionario de datos, representa el módulo de **Equipamiento del Hogar**.

# 🧠 PC4 – Descripción de archivos
modelo_RL.go – Entrenamiento del modelo de regresión logística
Este archivo entrena un modelo de regresión logística en Go de forma concurrente usando goroutines y canales.

- Lee el dataset real Enaho01-2022-612.csv (más de 990 mil registros).

- Extrae las columnas ESTRATO, DOMINIO y P612 (indicador de uso de TIC).

- Aplica gradiente descendente durante 100 épocas para encontrar los pesos del modelo (w0, w1, w2).

- Usa programación concurrente para calcular los gradientes.

- Muestra por consola los pesos finales entrenados del modelo.

Variables usadas: ESTRATO, DOMINIO
Variable objetivo: P612 (1 = usa TIC, 2 = no usa → convertido a 1 y 0)

predecir.go – Funcionamiento de la predicción
Este archivo simula el funcionamiento del modelo ya entrenado.

- Usa los pesos obtenidos del entrenamiento (copiados desde modelo_RL.go).

- Define tres registros de prueba con valores distintos de estrato y dominio.

- Calcula la probabilidad de que cada registro use TIC, aplicando la función sigmoide.

- Imprime por consola la probabilidad de uso de TIC para cada caso de prueba.

Esto demuestra que el modelo funciona con datos nuevos sin necesidad de volver a entrenar.

