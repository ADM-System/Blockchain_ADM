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
  
![image](https://github.com/user-attachments/assets/a9d27b96-55ff-4928-8a65-8b2254126983)

2) aggiunta del Nounce per dare un numero ai blocchi , ed aggiunto un menu' dal quale possiamo aggiungere un blocco , visualizzare l'intera blockchain everificarne la validità e poi uscire.

![image](https://github.com/user-attachments/assets/6ddf2d95-0de5-4a6b-bb4b-dfec5e121c0d)

3) aggiunta un interfaccia grafica web grazsie alla aggiunta si un server web locale scritto ed implementato nel codice in go

 ![image](https://github.com/user-attachments/assets/e1359ac1-92b6-4923-adef-cfa560bb0da8)

 4) aggiunto un sistema array  che tiene tyraccia delle transazioni all'interno del blocco come da foto
  
  ![image](https://github.com/user-attachments/assets/a3ea5470-fe3d-4dfc-ad15-d67d4126dbf7)

5) aggiunta del meccanismo di mining e aggiunto un bottone per il controllo della validità della blockcahin

![image](https://github.com/user-attachments/assets/ab748ee8-f099-433d-9695-4f45ed08fe0b)

6)aggiunto il meccanismo mempool per memorizzare non solo 1 trasnazione ma molte di piu' in un singolo blocco e poi tramite il tasto eseguire il mining del blocco.
7) aggiunto il meccanismo di ricompensa per il miner che viene specificato nel campo , e la sua ricompensa viene inclusa nel blocco.




  
## Tecnologie utilizzate

- **Go (Golang)**: linguaggio di programmazione utilizzato per la creazione della blockchain.

## Prerequisiti

Per eseguire il progetto, assicurati di avere **Go** installato sul tuo computer. Puoi scaricarlo e installarlo dal sito ufficiale: [https://golang.org/dl/](https://golang.org/dl/)
