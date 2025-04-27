# Blockchain in Go


## Introduzione

Questo è un esempio di una **blockchain** scritta in **Go**. La blockchain è composta da una sequenza di blocchi, ognuno dei quali contiene:
- Un **indice** univoco,
- Un **timestamp** di quando il blocco è stato creato,
- I **dati** che vengono memorizzati nel blocco,
- L'**hash** del blocco precedente, che lega i blocchi insieme,
- Il **proprio hash** che rappresenta l'identificatore unico di ciascun blocco.

Ogni blocco è creato e collegato alla catena in modo sicuro.

## Funzionalità
1)
- Creazione di blocchi con timestamp.
- Calcolo dell'hash di ogni blocco per garantire l'integrità della blockchain.
- Ogni blocco fa riferimento all'hash del blocco precedente.
2) aggiunta del Nounce per dare un numero ai blocchi , ed aggiunto un menu' dal quale possiamo aggiungere un blocco , visualizzare l'intera blockchain everificarne la validità e poi uscire
  

![image](https://github.com/user-attachments/assets/a9d27b96-55ff-4928-8a65-8b2254126983)

  
## Tecnologie utilizzate

- **Go (Golang)**: linguaggio di programmazione utilizzato per la creazione della blockchain.

## Prerequisiti

Per eseguire il progetto, assicurati di avere **Go** installato sul tuo computer. Puoi scaricarlo e installarlo dal sito ufficiale: [https://golang.org/dl/](https://golang.org/dl/)
